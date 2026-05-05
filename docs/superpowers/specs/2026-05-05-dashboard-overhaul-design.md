# Spec: Dashboard Overhaul & Visual Improvement

This document outlines the design for making the Internship Manager dashboard functional using a modular HTMX approach and improving the overall visual aesthetic of the application.

## 1. Goal
- Transform the static dashboard into a dynamic, data-driven summary of the user's progress.
- Improve the visual polish of the entire site (layout, colors, spacing).
- Implement modular loading for dashboard components using HTMX.

## 2. Architecture

### Dashboard Shell
The main dashboard page (`/`) will serve as a container that triggers the loading of its sub-components.

- **Component:** `Dashboard()`
- **HTMX Triggers:**
  - `<div hx-get="/dashboard/counts" hx-trigger="load"></div>`
  - `<div hx-get="/dashboard/pipeline" hx-trigger="load"></div>`

### API Endpoints (Internal)
1.  `GET /dashboard/counts`: Returns the summary stat cards.
2.  `GET /dashboard/pipeline`: Returns the status breakdown (Applied, Interviewing, etc.).

## 3. Data Model & Database

### New DB Method
`internal/db/db.go` will include a `GetDashboardStats` method:

```go
type DashboardStats struct {
    TotalApplications int
    InterviewingCount int
    OfferCount        int
    StatusBreakdown   map[string]int
}
```

### Query Logic
- Count applications where `user_id = $1`.
- Count applications with status 'Interviewing' where `user_id = $1`.
- Count applications with status 'Offer' where `user_id = $1`.
- Group by status to get the breakdown.

## 4. UI/UX Improvements

### Visual Theme
- **Background**: Light slate gray (`#f8fafc`).
- **Cards**: White with subtle borders and shadows.
- **Typography**: Clean sans-serif (Inter/system).
- **Primary Color**: Indigo/Blue for actions.

### Stat Cards
Each card in the `counts` component will feature:
- A clear label (e.g., "Total Applications").
- A large, bold number.
- A background color hint (Blue for total, Purple for interviews, Green for offers).

### Pipeline View
A list or grid showing:
- Status name.
- Count of applications in that status.
- A visual "bar" indicating relative volume.

## 5. Implementation Plan (High Level)
1.  Update `db.go` with stats logic.
2.  Create new handlers for `/dashboard/counts` and `/dashboard/pipeline`.
3.  Implement new Templ components for these fragments.
4.  Update `Dashboard()` templ to use HTMX triggers.
5.  Refactor `styles.css` for a more modern look.
6.  Update `Layout` for better navigation and container spacing.

## 6. Success Criteria
- Dashboard loads stats automatically without full page refresh.
- Stats are accurate and specific to the logged-in user.
- The UI feels modern, clean, and responsive.
