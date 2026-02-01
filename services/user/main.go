package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	// User service MVP

	// Keep is simple for now:
	// Schemas:
	// User: id, name, email, role, refresh_token (hash), refresh_created_at
	// Access Token: token (hash), created_at, expires_at, fk => refresh_token
	// Endpoints:
	// Register
	// Login
	// Logout
	// Authentication

	// Flow => Register => Login => Create Requests (Authentication guard) => Logout

	r := chi.NewRouter()

	r.Get("/tester", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello"))
	})

	http.ListenAndServe(":8080", r)
}
