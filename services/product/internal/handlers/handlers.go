package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
)

type RouterHandler struct{}

func StartRouterHandlers(r *chi.Mux) error {
	//&RouterHandler{}

	r.Get("/test", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello"))
	})

	return nil
}
