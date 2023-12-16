package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/amrohan/expenso-go/api/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	r := chi.NewRouter()

	r.Use(
		middleware.Logger,
		middleware.CleanPath,
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://localhost:4200", "http://localhost:4200"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}),
	)

	routes.LoadRoutes(r)

	fmt.Println("Server is running on port " + port)
	http.ListenAndServe(":3000", r)

}
