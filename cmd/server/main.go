package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	log.Println("starting server at :8080")
	http.ListenAndServe("127.0.0.1:8080", nil)
}
