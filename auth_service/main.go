package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("auth service OK"))
	})

	policyServiceURI := os.Getenv("POLICY_SERVICE_URL")
	if policyServiceURI == "" {
		fmt.Println("POLICY_SERVICE_URL environment variable not set")
		return
	} else {
		fmt.Printf("connected to %s\n", policyServiceURI)
	}

	handler := &AuthHandler{policyServiceURI: policyServiceURI}

	r.Mount("/auth", AuthRoutes(handler))
	http.ListenAndServe(":3001", r)
}

func AuthRoutes(authHandler *AuthHandler) chi.Router {

	r := chi.NewRouter()
	fmt.Println("Got here")
	r.Get("/", authHandler.PerformAuthorization)

	return r
}
