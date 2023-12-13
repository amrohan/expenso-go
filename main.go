package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/amrohan/expenso-go/api/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	routes.LoadRoutes(r)

	fmt.Println("Server is running on port " + port)
	http.ListenAndServe(":3000", r)

}
