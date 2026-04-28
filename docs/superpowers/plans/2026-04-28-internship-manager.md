# Internship Manager Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a comprehensive internship application tracker with a normalized PostgreSQL backend, Go server, and HTMX-powered frontend using Templ for type-safe components.

**Architecture:** A Go web server using `chi` for routing, `sqlc` for type-safe database access, and `templ` for rendering server-side components. HTMX will handle all dynamic UI updates (filtering, form submissions, modals) without full page reloads.

**Tech Stack:** Go 1.22+, PostgreSQL, HTMX, Templ, Chi, sqlc, pgx/v5.

---

### Task 1: Project Initialization & Tooling

**Files:**
- Create: `go.mod`
- Create: `Makefile`
- Create: `sqlc.yaml`
- Create: `tools.go`

- [ ] **Step 1: Initialize Go module**
Run: `go mod init internship-manager`

- [ ] **Step 2: Create tools.go to track dev dependencies**
```go
// +build tools
package tools

import (
	_ "github.com/a-h/templ/cmd/templ"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
)
```

- [ ] **Step 3: Create Makefile for common tasks**
```makefile
.PHONY: generate run test db-up db-down

generate:
	templ generate
	sqlc generate

run: generate
	go run cmd/server/main.go

test:
	go test ./... -v
```

- [ ] **Step 4: Configure sqlc.yaml**
```yaml
version: "2"
sql:
  - schema: "sql/schema.sql"
    queries: "sql/queries.sql"
    engine: "postgresql"
    gen:
      go:
        package: "db"
        out: "internal/db"
        sql_package: "pgx/v5"
```

- [ ] **Step 5: Commit**
```bash
git add .
git commit -m "chore: initialize project and tooling"
```

---

### Task 2: Database Schema & Migration

**Files:**
- Create: `sql/schema.sql`
- Create: `sql/queries.sql`

- [ ] **Step 1: Write SQL Schema**
```sql
CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    website TEXT,
    industry TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE positions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    location TEXT,
    work_mode TEXT NOT NULL, -- Remote, Hybrid, Onsite
    salary_range TEXT,
    post_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE skills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE position_skills (
    position_id UUID NOT NULL REFERENCES positions(id) ON DELETE CASCADE,
    skill_id UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    PRIMARY KEY (position_id, skill_id)
);

CREATE TABLE applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    position_id UUID NOT NULL REFERENCES positions(id) ON DELETE CASCADE,
    status TEXT NOT NULL, -- Pending, Applied, Interviewing, Rejected, Offer, Ghosted
    source TEXT,
    applied_at TIMESTAMP NOT NULL DEFAULT NOW(),
    notes TEXT
);

CREATE TABLE interviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    stage_name TEXT NOT NULL,
    scheduled_at TIMESTAMP NOT NULL,
    notes TEXT,
    feedback TEXT
);

CREATE TABLE contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    role TEXT,
    email TEXT,
    linkedin_url TEXT
);
```

- [ ] **Step 2: Generate DB code with sqlc**
Run: `sqlc generate` (Note: ensure sqlc is installed or use `go run github.com/sqlc-dev/sqlc/cmd/sqlc generate`)

- [ ] **Step 3: Commit**
```bash
git add sql/ internal/db/
git commit -m "feat: add database schema and generated code"
```

---

### Task 3: Base UI & Layout with Templ

**Files:**
- Create: `components/layout.templ`
- Create: `static/css/styles.css`

- [ ] **Step 1: Create base Layout component**
```templ
package components

templ Layout(title string) {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="UTF-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
			<title>{ title } - Internship Manager</title>
			<script src="https://unpkg.com/htmx.org@1.9.10"></script>
			<link rel="stylesheet" href="/static/css/styles.css"/>
		</head>
		<body class="bg-gray-100 font-sans">
			<nav class="bg-blue-600 text-white p-4 shadow-md">
				<div class="container mx-auto flex justify-between items-center">
					<h1 class="text-xl font-bold">Internship Manager</h1>
					<div class="space-x-4">
						<a href="/" class="hover:underline">Dashboard</a>
						<a href="/companies" class="hover:underline">Companies</a>
						<a href="/positions" class="hover:underline">Positions</a>
					</div>
				</div>
			</nav>
			<main class="container mx-auto p-6">
				{ children... }
			</main>
		</body>
	</html>
}
```

- [ ] **Step 2: Add basic CSS**
```css
/* Minimalist styling */
body { margin: 0; background: #f4f7f6; }
.container { max-width: 1200px; margin: 0 auto; }
/* Add more utility-like styles here */
```

- [ ] **Step 3: Generate Templ code**
Run: `templ generate`

- [ ] **Step 4: Commit**
```bash
git add components/ static/
git commit -m "feat: add base layout and styles"
```

---

### Task 4: Server Setup & Company CRUD

**Files:**
- Create: `cmd/server/main.go`
- Create: `internal/handlers/company_handlers.go`
- Modify: `sql/queries.sql`

- [ ] **Step 1: Add Company queries to queries.sql**
```sql
-- name: ListCompanies :many
SELECT * FROM companies ORDER BY name;

-- name: CreateCompany :one
INSERT INTO companies (name, website, industry)
VALUES ($1, $2, $3)
RETURNING *;
```

- [ ] **Step 2: Implement Company handlers**
Implement a handler to list companies and a handler to create one.

- [ ] **Step 3: Setup Chi router and server in main.go**
Connect to Postgres using `pgxpool` and wire up the routes.

- [ ] **Step 4: Run the app and verify the Companies list page**
Run: `make run`

- [ ] **Step 5: Commit**
```bash
git add .
git commit -m "feat: implement basic company listing"
```

---

### Task 5: Positions & N-to-M Skills

**Files:**
- Create: `components/position_list.templ`
- Modify: `sql/queries.sql`
- Create: `internal/handlers/position_handlers.go`

- [ ] **Step 1: Add complex queries for Positions and Skills**
Need a query that joins positions with their skills and company names.

- [ ] **Step 2: Create Templ components for Position list cards**
Show tags for skills.

- [ ] **Step 3: Implement filtering by skill using HTMX**
Clicking a skill tag should trigger `hx-get="/positions?skill=Go"` and replace the list.

- [ ] **Step 4: Commit**
```bash
git add .
git commit -m "feat: add position list with skill filtering"
```

---

### Task 6: Applications & Timeline View

**Files:**
- Create: `components/application_detail.templ`
- Create: `internal/handlers/application_handlers.go`

- [ ] **Step 1: Implement Application detail page**
Show application status and a timeline of interviews.

- [ ] **Step 2: Implement "Add Interview" inline form**
Use HTMX to append a new interview stage to the timeline without reload.

- [ ] **Step 3: Commit**
```bash
git add .
git commit -m "feat: add application details and interview timeline"
```
