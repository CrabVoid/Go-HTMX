# Design Spec: Internship Manager (Go + HTMX + Templ)

## 1. Overview
A comprehensive web-based internship application tracker built with Go, HTMX, and Templ. The system uses a normalized PostgreSQL database with complex relationships to track companies, positions, skills, applications, and interview stages.

## 2. Technical Stack
- **Backend:** Go 1.22+
- **Routing:** `go-chi/chi`
- **Templates:** `templ` (Type-safe HTML components)
- **Database:** PostgreSQL
- **DB Driver:** `pgx/v5`
- **DB Tooling:** `sqlc` (Type-safe SQL generator)
- **Frontend:** `htmx` (for dynamic updates)
- **Styling:** Vanilla CSS

## 3. Data Schema (PostgreSQL)

### 3.1 Tables
- **`companies`**
    - `id`: UUID (Primary Key)
    - `name`: TEXT (Unique)
    - `website`: TEXT
    - `industry`: TEXT
    - `created_at`: TIMESTAMP
- **`positions`**
    - `id`: UUID (Primary Key)
    - `company_id`: UUID (FK -> companies.id)
    - `title`: TEXT
    - `location`: TEXT (e.g., "Remote", "City, ST")
    - `work_mode`: TEXT (Enum: Remote, Hybrid, Onsite)
    - `salary_range`: TEXT
    - `post_url`: TEXT
- **`skills`**
    - `id`: UUID (Primary Key)
    - `name`: TEXT (Unique)
- **`position_skills` (N-to-M)**
    - `position_id`: UUID (FK -> positions.id)
    - `skill_id`: UUID (FK -> skills.id)
    - Primary Key: (position_id, skill_id)
- **`applications`**
    - `id`: UUID (Primary Key)
    - `position_id`: UUID (FK -> positions.id)
    - `status`: TEXT (Enum: Pending, Applied, Interviewing, Rejected, Offer, Ghosted)
    - `source`: TEXT (e.g., "LinkedIn", "Referral")
    - `applied_at`: TIMESTAMP
    - `notes`: TEXT
- **`interviews`**
    - `id`: UUID (Primary Key)
    - `application_id`: UUID (FK -> applications.id)
    - `stage_name`: TEXT (e.g., "Technical", "Behavioral")
    - `scheduled_at`: TIMESTAMP
    - `notes`: TEXT
    - `feedback`: TEXT
- **`contacts`**
    - `id`: UUID (Primary Key)
    - `company_id`: UUID (FK -> companies.id)
    - `name`: TEXT
    - `role`: TEXT
    - `email`: TEXT
    - `linkedin_url`: TEXT

## 4. UI/UX & Components (HTMX)
- **Dashboard:** Overview cards for total applications, active interviews, and offers.
- **Position List:** Searchable table with skill tags. Clicking a tag filters the list using HTMX (`hx-get`).
- **Application Detail View:** A drill-down page showing the timeline of interviews and linked contacts.
- **Dynamic Forms:**
    - Adding an interview stage updates the timeline without a full page reload.
    - Skill selection using a dynamic tag-input system.
- **Modals:** Used for quick-adding companies or contacts.

## 5. Directory Structure
```
.
├── cmd/
│   └── server/          # Entry point
├── internal/
│   ├── db/              # sqlc generated code & migrations
│   ├── handlers/        # HTTP route handlers
│   ├── models/          # Domain models (if needed beyond sqlc)
│   └── services/        # Business logic
├── components/          # Templ files (*.templ)
├── static/              # CSS, JS, Images
├── sql/
│   ├── schema.sql       # Database schema
│   └── queries.sql      # sqlc queries
├── sqlc.yaml            # sqlc configuration
└── Makefile             # Task automation (templ generate, sqlc generate)
```

## 6. Success Criteria
- [ ] Users can create a company and link multiple positions to it.
- [ ] Users can assign multiple skills to a position (N-to-M).
- [ ] Users can track the full lifecycle of an application through multiple interview stages.
- [ ] All UI updates are performed via HTMX fragments for a smooth experience.
- [ ] The system is fully type-safe from SQL to HTML using `sqlc` and `templ`.
