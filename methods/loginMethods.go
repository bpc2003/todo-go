package methods

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"strconv"

	"todo/todoDB"
)

func GetLogin(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "login", nil)
}

func PostLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "todo: %v\n", err)
		return
	}

	f := todoDB.LoginUserParams{
		Email:    r.PostForm.Get("email"),
		Password: fmt.Sprintf("%x", sha1.Sum([]byte(r.PostForm.Get("password")))),
	}
	if f.Email != "" && f.Password != "" {
		u, err := q.LoginUser(ctx, f)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintln(w, "Invalid username or password")
			return
		}
		c := http.Cookie{
			Name:   "auth",
			Value:  strconv.Itoa(int(u.ID)),
			Path:   "/",
			MaxAge: 3600,
		}
		http.SetCookie(w, &c)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Email and Password are required")
	}
}

func PostRegister(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "todo: %v\n", err)
		return
	}

	f := todoDB.CreateUserParams{
		Email:    r.PostForm.Get("email"),
		Password: fmt.Sprintf("%x", sha1.Sum([]byte(r.PostForm.Get("password")))),
	}
	if f.Email != "" && f.Password != "" {		
		u, err := q.CreateUser(ctx, f)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "todo: %v\n", err)
			return
		}
		c := http.Cookie{
			Name:   "auth",
			Value:  strconv.Itoa(int(u.ID)),
			Path:   "/",
			MaxAge: 3600,
		}
		http.SetCookie(w, &c)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Email and Password are required")
	}
}
