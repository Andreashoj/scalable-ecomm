package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()

	r.Get("/tester", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello"))
	})

	http.ListenAndServe(":8080", r)
}
