package main

import (
	"log"
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
	if _, err := w.Write([]byte("Hello World")); err != nil {
		log.Printf("%v", err)
	}

}

func protoHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte(r.Proto)); err != nil {
		log.Printf("%v", err)
	}
}
