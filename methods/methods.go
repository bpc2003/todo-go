package methods

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"todo/todoDB"
)

type formContent struct {
	Todo     string
	Email    string
	Password string
}

type todoList struct {
	List []todoDB.List
}

var t *template.Template
var db *sql.DB
var q *todoDB.Queries
var ctx = context.Background()

//go:embed schema.sql
var ddl string

func init() {
	var err error
	t, err = template.ParseGlob("public/views/*.html")
	if err != nil {
		fmt.Fprintf(os.Stderr, "todo: %v\n", err)
		os.Exit(1)
	}
	db, err = sql.Open("sqlite3", "list.db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "todo: %v\n", err)
		os.Exit(1)
	}
	if _, err = db.ExecContext(ctx, ddl); err != nil {
		fmt.Fprintf(os.Stderr, "todo: %v\n", err)
		os.Exit(1)
	}
	q = todoDB.New(db)
}

func GetIndex(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("auth")
	if c == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	uid, _ := strconv.Atoi(c.Value)
	list, err := q.GetList(ctx, int64(uid))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "todo: %v\n", err)
		return
	}
	tl := todoList{List: list}
	t.ExecuteTemplate(w, "index", tl)
}

func PostIndex(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("auth")
	if c == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "todo: %v\n", err)
		return
	}

	uid, _ := strconv.Atoi(c.Value)
	f := todoDB.PostListParams{Todo: r.PostForm.Get("todo"), Userid: int64(uid)}
	if f.Todo != "" {
		if _, err := q.PostList(ctx, f); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "todo: %v\n", err)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteEntry(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("auth")
	if c == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	uid, _ := strconv.Atoi(c.Value)
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "invalid id of", r.URL.Path)
		return
	}

	f := todoDB.DeleteEntryParams{ID: int64(id), Userid: int64(uid)}
	if err := q.DeleteEntry(ctx, f); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "item with id of", id, "doesn't exist")
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
