package routes

import (
	"net/http"

	"github.com/amrohan/expenso-go/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func LoadRoutes(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World!"))
		})
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		})

		r.Post("/login", handlers.LoginUser)
		r.Post("/register", handlers.RegisterUser)
		r.Post("/logout", handlers.LogoutUser)
	})

	r.With(handlers.AuthMiddleware).Route("/api/transaction", func(r chi.Router) {
		r.Post("/", handlers.CreateTransaction)
		r.Get("/", handlers.GetAllTransaction)
		r.Get("/{id}", handlers.GetTransactionById)
		r.Get("/{month}-{year}", handlers.GetTransactionByMonthAndYear)
		r.Get("/u/{month}-{year}-{userId}", handlers.GetTransactionByMonthAndYearByUserId)
		r.Get("/user/{id}", handlers.GetTransactionByUserId)
		r.Get("/category/{id}", handlers.GetTransactionByCategoryId)
		r.Get("/account/{id}", handlers.GetTransactionByAccountId)
		r.Put("/", handlers.UpdateTransaction)
		r.Delete("/{id}", handlers.DeleteTransaction)
	})

	r.With(handlers.AuthMiddleware).Route("/api/category", func(r chi.Router) {
		r.Post("/", handlers.CreateCategory)
		r.Get("/", handlers.GetAllCategory)
		r.Get("/{id}", handlers.GetCategoryById)
		r.Get("/user/{id}", handlers.GetCategoryByUserId)
		r.Put("/", handlers.UpdateCategory)
		r.Delete("/{id}", handlers.DeleteCategory)
	})

	r.With(handlers.AuthMiddleware).Route("/api/account", func(r chi.Router) {
		r.Post("/", handlers.CreateAccount)
		r.Get("/", handlers.GetAllAccount)
		r.Get("/{id}", handlers.GetAccountById)
		r.Get("/user/{id}", handlers.GetAccountsByUserId)
		r.Put("/", handlers.UpdateAccount)
		r.Delete("/{id}", handlers.DeleteAccount)
	})

	r.With(handlers.AuthMiddleware).Route("/api/user", func(r chi.Router) {
		r.Post("/", handlers.CreateUser)
		r.Get("/", handlers.GetAllUsers)
		r.Get("/{id}", handlers.GetUserById)
		r.Get("/email/{email}", handlers.GetUserByEmail)
		r.Get("/username/{username}", handlers.GetUserByUsername)
		r.Get("/du/", handlers.GetAllDeletedUser)
		r.Post("/du/{id}", handlers.RestoreUser)
		r.Put("/", handlers.UpdateUser)
		r.Delete("/{id}", handlers.DeleteUser)
	})

}
