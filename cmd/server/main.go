/*
+--------------------------------------------------------------------------------+
| PACKAGE DECLARATION                                                            |
| PURPOSE: Entry point for the Internship Manager server.                        |
| INFO: Initializes the database, handlers, and defines routing.                 |
+--------------------------------------------------------------------------------+
*/
package main

/*
+--------------------------------------------------------------------------------+
| EXTERNAL & INTERNAL IMPORTS                                                    |
| PURPOSE: Import standard library and project dependencies.                     |
| INFO: Includes chi router, pgx database pool, and internal handlers.           |
+--------------------------------------------------------------------------------+
*/
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

/*
+--------------------------------------------------------------------------------+
| MAIN FUNCTION                                                                  |
| PURPOSE: Main execution block for the server application.                      |
| INFO: Orchestrates the startup sequence of the entire application.             |
+--------------------------------------------------------------------------------+
*/
func main() {
	/*
	+----------------------------------------------------------------------------+
	| ENVIRONMENT CONFIGURATION                                                  |
	| PURPOSE: Load configuration from environment variables.                    |
	| INFO: Configures the listening port and database connection URL.           |
	+----------------------------------------------------------------------------+
	*/
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/internship_manager?sslmode=disable"
	}

	/*
	+----------------------------------------------------------------------------+
	| DATABASE INITIALIZATION                                                    |
	| PURPOSE: Setup connection pool to PostgreSQL.                              |
	| INFO: Uses pgxpool for efficient connection management.                    |
	+----------------------------------------------------------------------------+
	*/
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse DATABASE_URL: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	/*
	+----------------------------------------------------------------------------+
	| HANDLER INITIALIZATION                                                     |
	| PURPOSE: Instantiate all logic handlers for different domain entities.      |
	| INFO: Each handler receives the SQLC queries object.                       |
	+----------------------------------------------------------------------------+
	*/
	queries := db.New(pool)
	authHandler := handlers.NewAuthHandler(queries)
	companyHandler := handlers.NewCompanyHandler(queries)
	positionHandler := handlers.NewPositionHandler(queries)
	applicationHandler := handlers.NewApplicationHandler(queries)
	dashboardHandler := handlers.NewDashboardHandler(queries)

	/*
	+----------------------------------------------------------------------------+
	| ROUTER & MIDDLEWARE SETUP                                                  |
	| PURPOSE: Initialize the Chi router and global middlewares.                 |
	| INFO: Includes standard request logging and recovery from panics.          |
	+----------------------------------------------------------------------------+
	*/
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	/*
	+----------------------------------------------------------------------------+
	| STATIC FILE SERVING                                                        |
	| PURPOSE: Serve CSS, JS, and image assets from the static directory.        |
	| INFO: Mounted at /static/* path.                                           |
	+----------------------------------------------------------------------------+
	*/
	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	/*
	+----------------------------------------------------------------------------+
	| PUBLIC ROUTES                                                              |
	| PURPOSE: Routes accessible without authentication.                         |
	| INFO: Handles user login, registration, and logout.                        |
	+----------------------------------------------------------------------------+
	*/
	r.Get("/login", authHandler.GetLogin)
	r.Post("/login", authHandler.PostLogin)
	r.Get("/register", authHandler.GetRegister)
	r.Post("/register", authHandler.PostRegister)
	r.Post("/logout", authHandler.PostLogout)

	/*
	+----------------------------------------------------------------------------+
	| PROTECTED ROUTES                                                           |
	| PURPOSE: Routes requiring valid user authentication.                       |
	| INFO: Wrapped in AuthMiddleware to ensure session validity.                |
	+----------------------------------------------------------------------------+
	*/
	r.Group(func(r chi.Router) {
		r.Use(handlers.AuthMiddleware)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			component := components.Dashboard()
			err := component.Render(r.Context(), w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
		r.Get("/dashboard/counts", dashboardHandler.GetCounts)
		r.Get("/dashboard/pipeline", dashboardHandler.GetPipeline)
		r.Get("/companies", companyHandler.ListCompanies)
		r.Post("/companies", companyHandler.CreateCompany)
		r.Get("/companies/{id}/edit", companyHandler.GetCompanyForm)
		r.Get("/companies/{id}/card", companyHandler.GetCompanyCard)
		r.Put("/companies/{id}", companyHandler.UpdateCompany)
		r.Delete("/companies/{id}", companyHandler.DeleteCompany)
		r.Get("/positions", positionHandler.ListPositions)
		r.Post("/positions", positionHandler.CreatePosition)
		r.Get("/positions/{id}", positionHandler.GetPosition)
		r.Delete("/positions/{id}", positionHandler.DeletePosition)
		r.Get("/applications", applicationHandler.ListApplications)
		r.Post("/applications", applicationHandler.CreateApplication)
		r.Get("/applications/{id}", applicationHandler.GetApplication)
		r.Delete("/applications/{id}", applicationHandler.DeleteApplication)
		r.Post("/applications/{id}/interviews", applicationHandler.CreateInterview)
	})

	/*
	+----------------------------------------------------------------------------+
	| SERVER LIFECYCLE                                                           |
	| PURPOSE: Start the HTTP server.                                            |
	| INFO: Listens on the configured port and blocks until termination.         |
	+----------------------------------------------------------------------------+
	*/
	fmt.Printf("Server starting on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
