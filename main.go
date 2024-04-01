package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"todo/methods"
)

func main() {
	r := chi.NewRouter()
	fs := http.FileServer(http.Dir("public/"))

	r.Use(middleware.Logger)

	r.Handle("/public/*", http.StripPrefix("/public/", fs))
	r.Get("/", methods.GetIndex)
	r.Post("/", methods.PostIndex)

	r.Get("/login", methods.GetLogin)
	r.Post("/login", methods.PostLogin)

	r.Get("/register", methods.GetLogin)
	r.Post("/register", methods.PostRegister)

	r.Delete("/{id}", methods.DeleteEntry)

	log.Fatal(http.ListenAndServe(":8080", r))
}
