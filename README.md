# Internship Manager Tracker

A modern web application to track internship applications, companies, and positions. Built with **Go (Chi/PGX)**, **HTMX**, and **Templ** for a fast, interactive experience with minimal client-side JavaScript.

## Tech Stack
- **Backend**: Go 1.23+ with Chi router.
- **Frontend**: Templ components + HTMX for dynamic updates.
- **Database**: PostgreSQL (managed via Docker).
- **Styling**: Vanilla CSS.

## Getting Started

### 1. Prerequisites
- [Go](https://go.dev/dl/) installed.
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) installed and running.
- [Templ](https://templ.guide/quick-start/installation) (for code generation).

### 2. Database Setup
Start the PostgreSQL database using Docker:
```powershell
docker-compose up -d
```
Initialize the schema (runs inside the container):
```powershell
Get-Content sql/schema.sql | docker exec -i gohtmx-db-1 psql -U postgres -d internship_manager
```

### 3. Running the Application
Set the environment variables and run:
```powershell
$env:DATABASE_URL="postgres://postgres:postgres@localhost:5432/internship_manager?sslmode=disable"
go run cmd/server/main.go
```
The server will be available at: **http://localhost:8080**

## Project Structure
- `cmd/server/`: Application entry point.
- `components/`: Templ UI components.
- `internal/handlers/`: HTTP request handlers.
- `internal/db/`: Database models and query interfaces.
- `sql/`: SQL schema and queries.
- `static/`: Static assets (CSS).

## Current Status
The project is currently in a prototype phase using manual database stubs to bypass resource-heavy code generation during CLI setup.
