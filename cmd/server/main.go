package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"internship-manager/components"
	"internship-manager/internal/db"
	"internship-manager/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/internship_manager?sslmode=disable"
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse DATABASE_URL: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	queries := db.New(pool)
	authHandler := handlers.NewAuthHandler(queries)
	companyHandler := handlers.NewCompanyHandler(queries)
	positionHandler := handlers.NewPositionHandler(queries)
	applicationHandler := handlers.NewApplicationHandler(queries)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Static files
	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Public routes
	r.Get("/login", authHandler.GetLogin)
	r.Post("/login", authHandler.PostLogin)
	r.Get("/register", authHandler.GetRegister)
	r.Post("/register", authHandler.PostRegister)
	r.Post("/logout", authHandler.PostLogout)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(handlers.AuthMiddleware)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			component := components.Dashboard()
			err := component.Render(r.Context(), w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
		r.Get("/companies", companyHandler.ListCompanies)
		r.Post("/companies", companyHandler.CreateCompany)
		r.Get("/companies/{id}/edit", companyHandler.GetCompanyForm)
		r.Get("/companies/{id}/card", companyHandler.GetCompanyCard)
		r.Put("/companies/{id}", companyHandler.UpdateCompany)
		r.Delete("/companies/{id}", companyHandler.DeleteCompany)
		r.Get("/positions", positionHandler.ListPositions)
		r.Post("/positions", positionHandler.CreatePosition)
		r.Delete("/positions/{id}", positionHandler.DeletePosition)
		r.Get("/applications", applicationHandler.ListApplications)
		r.Post("/applications", applicationHandler.CreateApplication)
		r.Get("/applications/{id}", applicationHandler.GetApplication)
		r.Delete("/applications/{id}", applicationHandler.DeleteApplication)
		r.Post("/applications/{id}/interviews", applicationHandler.CreateInterview)
	})

	fmt.Printf("Server starting on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
