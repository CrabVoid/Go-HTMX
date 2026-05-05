# Dashboard Overhaul Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement a dynamic, HTMX-powered dashboard with live statistics and a modern visual theme.

**Architecture:** 
- Modular loading: The dashboard shell loads fragments via HTMX (`/dashboard/counts` and `/dashboard/pipeline`).
- Data isolation: All stats are scoped to the authenticated user.
- Component-based UI: Shared Templ components for stat cards and status lists.

**Tech Stack:** Go (Chi), Templ, HTMX, Vanilla CSS.

---

### Task 1: Database Stats Logic

**Files:**
- Modify: `internal/db/db.go`
- Modify: `internal/db/models.go`

- [ ] **Step 1: Define DashboardStats struct**
Add the following to `internal/db/models.go`:
```go
type DashboardStats struct {
	TotalApplications int
	InterviewingCount int
	OfferCount        int
	StatusBreakdown   map[string]int
}
```

- [ ] **Step 2: Implement GetDashboardStats in db.go**
Add this method to `internal/db/db.go`:
```go
func (q *Queries) GetDashboardStats(ctx context.Context, userID uuid.UUID) (DashboardStats, error) {
	var stats DashboardStats
	stats.StatusBreakdown = make(map[string]int)

	// Total Apps
	err := q.db.QueryRow(ctx, "SELECT COUNT(*) FROM applications WHERE user_id = $1", userID).Scan(&stats.TotalApplications)
	if err != nil {
		return stats, err
	}

	// Interviewing
	err = q.db.QueryRow(ctx, "SELECT COUNT(*) FROM applications WHERE user_id = $1 AND status = 'Interviewing'", userID).Scan(&stats.InterviewingCount)
	if err != nil {
		return stats, err
	}

	// Offers
	err = q.db.QueryRow(ctx, "SELECT COUNT(*) FROM applications WHERE user_id = $1 AND status = 'Offer'", userID).Scan(&stats.OfferCount)
	if err != nil {
		return stats, err
	}

	// Status Breakdown
	rows, err := q.db.Query(ctx, "SELECT status, COUNT(*) FROM applications WHERE user_id = $1 GROUP BY status", userID)
	if err != nil {
		return stats, err
	}
	defer rows.Close()
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return stats, err
		}
		stats.StatusBreakdown[status] = count
	}

	return stats, nil
}
```

- [ ] **Step 3: Commit**
```bash
git add internal/db/db.go internal/db/models.go
git commit -m "feat(db): add GetDashboardStats method"
```

---

### Task 2: Dashboard Fragment Components

**Files:**
- Create: `components/dashboard_fragments.templ`

- [ ] **Step 1: Create DashboardCounts component**
```templ
package components

import (
    "internship-manager/internal/db"
    "fmt"
)

templ DashboardCounts(stats db.DashboardStats) {
	<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
		<div class="bg-blue-50 p-6 rounded-xl border border-blue-100 shadow-sm">
			<h3 class="text-blue-800 font-semibold text-sm uppercase tracking-wider mb-1">Total Applications</h3>
			<p class="text-4xl font-bold text-blue-900">{ fmt.Sprint(stats.TotalApplications) }</p>
		</div>
		<div class="bg-indigo-50 p-6 rounded-xl border border-indigo-100 shadow-sm">
			<h3 class="text-indigo-800 font-semibold text-sm uppercase tracking-wider mb-1">Interviewing</h3>
			<p class="text-4xl font-bold text-indigo-900">{ fmt.Sprint(stats.InterviewingCount) }</p>
		</div>
		<div class="bg-emerald-50 p-6 rounded-xl border border-emerald-100 shadow-sm">
			<h3 class="text-emerald-800 font-semibold text-sm uppercase tracking-wider mb-1">Offers</h3>
			<p class="text-4xl font-bold text-emerald-900">{ fmt.Sprint(stats.OfferCount) }</p>
		</div>
	</div>
}
```

- [ ] **Step 2: Create DashboardPipeline component**
```templ
templ DashboardPipeline(stats db.DashboardStats) {
	<div class="bg-white p-6 rounded-xl border border-gray-200 shadow-sm">
		<h3 class="text-xl font-bold text-gray-800 mb-6">Application Pipeline</h3>
		<div class="space-y-4">
			if len(stats.StatusBreakdown) == 0 {
				<p class="text-gray-500 italic">No applications yet. Start tracking to see your pipeline!</p>
			}
			for status, count := range stats.StatusBreakdown {
				<div>
					<div class="flex justify-between items-center mb-1">
						<span class="font-medium text-gray-700">{ status }</span>
						<span class="text-sm font-bold text-gray-900">{ fmt.Sprint(count) }</span>
					</div>
					<div class="w-full bg-gray-100 rounded-full h-2.5">
						<div class="bg-indigo-600 h-2.5 rounded-full" style={ fmt.Sprintf("width: %d%%", (count * 100 / max(stats.TotalApplications, 1))) }></div>
					</div>
				</div>
			}
		</div>
	</div>
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
```

- [ ] **Step 3: Generate Templ code**
```bash
go run github.com/a-h/templ/cmd/templ generate
```

- [ ] **Step 4: Commit**
```bash
git add components/dashboard_fragments.templ components/dashboard_fragments_templ.go
git commit -m "feat(ui): add dashboard fragment components"
```

---

### Task 3: Dashboard Fragment Handlers

**Files:**
- Create: `internal/handlers/dashboard_handlers.go`
- Modify: `cmd/server/main.go`

- [ ] **Step 1: Create DashboardHandler**
```go
package handlers

import (
	"context"
	"net/http"

	"internship-manager/components"
	"internship-manager/internal/db"
)

type DashboardHandler struct {
	Queries *db.Queries
}

func NewDashboardHandler(queries *db.Queries) *DashboardHandler {
	return &DashboardHandler{Queries: queries}
}

func (h *DashboardHandler) GetCounts(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r.Context())
	stats, err := h.Queries.GetDashboardStats(context.Background(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	component := components.DashboardCounts(stats)
	component.Render(r.Context(), w)
}

func (h *DashboardHandler) GetPipeline(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r.Context())
	stats, err := h.Queries.GetDashboardStats(context.Background(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	component := components.DashboardPipeline(stats)
	component.Render(r.Context(), w)
}
```

- [ ] **Step 2: Register routes in main.go**
Add to `cmd/server/main.go` (inside the protected route group):
```go
dashboardHandler := handlers.NewDashboardHandler(queries)
// ... inside r.Group(func(r chi.Router) { ...
r.Get("/dashboard/counts", dashboardHandler.GetCounts)
r.Get("/dashboard/pipeline", dashboardHandler.GetPipeline)
```

- [ ] **Step 3: Commit**
```bash
git add internal/handlers/dashboard_handlers.go cmd/server/main.go
git commit -m "feat(api): add dashboard fragment handlers and routes"
```

---

### Task 4: Update Dashboard Shell

**Files:**
- Modify: `components/dashboard.templ`

- [ ] **Step 1: Update Dashboard shell to use HTMX**
```templ
package components

templ Dashboard() {
	@Layout("Dashboard") {
		<div class="space-y-8">
			<div>
				<h2 class="text-3xl font-bold text-gray-900 mb-2">Welcome Back</h2>
				<p class="text-gray-600">Here's a summary of your internship search.</p>
			</div>

			<div hx-get="/dashboard/counts" hx-trigger="load" class="min-h-[120px]">
				<div class="flex justify-center items-center h-full">
					<div class="animate-pulse text-gray-400">Loading stats...</div>
				</div>
			</div>

			<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
				<div hx-get="/dashboard/pipeline" hx-trigger="load" class="lg:col-span-2 min-h-[300px]">
					<div class="bg-white p-6 rounded-xl border border-gray-200 shadow-sm animate-pulse">
						<div class="h-6 w-48 bg-gray-200 rounded mb-6"></div>
						<div class="space-y-4">
							<div class="h-4 w-full bg-gray-100 rounded"></div>
							<div class="h-4 w-full bg-gray-100 rounded"></div>
							<div class="h-4 w-full bg-gray-100 rounded"></div>
						</div>
					</div>
				</div>

				<div class="space-y-6">
					<div class="bg-indigo-900 p-6 rounded-xl text-white shadow-lg">
						<h3 class="text-lg font-bold mb-2">Quick Actions</h3>
						<div class="space-y-3">
							<a href="/positions" class="block w-full text-center py-2 px-4 bg-white text-indigo-900 rounded-lg font-semibold hover:bg-indigo-50 transition-colors">
								Add Position
							</a>
							<a href="/companies" class="block w-full text-center py-2 px-4 border border-indigo-400 text-indigo-100 rounded-lg font-semibold hover:bg-indigo-800 transition-colors">
								Manage Companies
							</a>
						</div>
					</div>
				</div>
			</div>
		</div>
	}
}
```

- [ ] **Step 2: Generate Templ code**
```bash
go run github.com/a-h/templ/cmd/templ generate
```

- [ ] **Step 3: Commit**
```bash
git add components/dashboard.templ components/dashboard_templ.go
git commit -m "feat(ui): update dashboard shell with HTMX triggers"
```

---

### Task 5: Visual Polish (CSS & Layout)

**Files:**
- Modify: `static/css/styles.css`
- Modify: `components/layout.templ`

- [ ] **Step 1: Refactor styles.css**
```css
:root {
    --primary: #4f46e5;
    --primary-hover: #4338ca;
    --bg: #f8fafc;
    --text-main: #0f172a;
    --text-muted: #64748b;
    --border: #e2e8f0;
}

body {
    margin: 0;
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
    background-color: var(--bg);
    color: var(--text-main);
    line-height: 1.5;
}

.container {
    max-width: 1100px;
    margin: 0 auto;
    padding: 0 1.5rem;
}

nav {
    background-color: white;
    border-bottom: 1px solid var(--border);
    padding: 1rem 0;
    margin-bottom: 2.5rem;
    position: sticky;
    top: 0;
    z-index: 50;
}

nav .container {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.nav-links {
    display: flex;
    gap: 1.5rem;
    align-items: center;
}

nav a {
    color: var(--text-muted);
    text-decoration: none;
    font-weight: 500;
    font-size: 0.9375rem;
    transition: color 0.2s;
}

nav a:hover, nav a.active {
    color: var(--primary);
}

.btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 0.5rem 1rem;
    border-radius: 0.5rem;
    font-weight: 600;
    transition: all 0.2s;
    cursor: pointer;
    text-decoration: none;
}

.btn-primary {
    background-color: var(--primary);
    color: white;
}

.btn-primary:hover {
    background-color: var(--primary-hover);
}

/* Tailwind-like utility classes for ease of use in Templ */
.shadow-sm { box-shadow: 0 1px 2px 0 rgb(0 0 0 / 0.05); }
.rounded-xl { border-radius: 0.75rem; }
```

- [ ] **Step 2: Update Layout component**
```templ
package components

import "internship-manager/internal/handlers"

templ Layout(title string) {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="UTF-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
			<title>{ title } - Internship Manager</title>
			<script src="https://unpkg.com/htmx.org@1.9.10"></script>
			<script src="https://cdn.tailwindcss.com"></script>
			<link rel="stylesheet" href="/static/css/styles.css"/>
			<link rel="preconnect" href="https://fonts.googleapis.com"/>
			<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin/>
			<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap" rel="stylesheet"/>
		</head>
		<body class="bg-slate-50">
			<nav class="bg-white border-b border-slate-200 py-4 mb-8 sticky top-0 z-50">
				<div class="container mx-auto px-4 flex justify-between items-center">
					<div class="flex items-center gap-8">
						<h1 class="text-xl font-bold text-indigo-600 tracking-tight">InternshipMgr</h1>
						<div class="hidden md:flex gap-6">
							<a href="/" class="text-slate-600 hover:text-indigo-600 font-medium transition-colors">Dashboard</a>
							<a href="/companies" class="text-slate-600 hover:text-indigo-600 font-medium transition-colors">Companies</a>
							<a href="/positions" class="text-slate-600 hover:text-indigo-600 font-medium transition-colors">Positions</a>
						</div>
					</div>
					<div class="flex items-center gap-4">
						if handlers.GetUserID(ctx).String() != "00000000-0000-0000-0000-000000000000" {
							<form hx-post="/logout" hx-target="body" hx-push-url="true">
								<button type="submit" class="text-sm font-semibold text-slate-500 hover:text-red-600 transition-colors">Logout</button>
							</form>
						} else {
							<a href="/login" class="text-sm font-semibold text-indigo-600 hover:text-indigo-700">Login</a>
						}
					</div>
				</div>
			</nav>
			<main class="container mx-auto px-4 pb-12">
				{ children... }
			</main>
		</body>
	</html>
}
```

- [ ] **Step 3: Generate Templ code**
```bash
go run github.com/a-h/templ/cmd/templ generate
```

- [ ] **Step 4: Commit**
```bash
git add static/css/styles.css components/layout.templ components/layout_templ.go
git commit -m "style: modernize visuals and improve layout"
```

---

### Task 6: Final Verification

- [ ] **Step 1: Build and Run**
```bash
go build -o server.exe cmd/server/main.go
# Assuming server is already running in background, restart it or just run it to check for errors
./server.exe
```

- [ ] **Step 2: Manual Check**
1. Login to the application.
2. Navigate to Dashboard.
3. Verify that stats load automatically.
4. Add a new application and check if stats update upon refresh.
5. Verify the visual style matches the modern design.
