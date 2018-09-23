package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", indexHandler)
	r.Get("/proto", protoHandler)
	return r
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func protoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.Proto))
}
