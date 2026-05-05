package handlers

import (
	"context"
	"net/http"

	"internship-manager/components"
	"internship-manager/internal/db"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PositionHandler struct {
	Queries *db.Queries
}

func NewPositionHandler(queries *db.Queries) *PositionHandler {
	return &PositionHandler{Queries: queries}
}

func (h *PositionHandler) ListPositions(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	userID := GetUserID(r.Context())
	
	positions, err := h.Queries.ListPositions(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter by skill if provided
	skillFilter := r.URL.Query().Get("skill")
	if skillFilter != "" {
		filtered := []db.PositionWithDetails{}
		for _, p := range positions {
			for _, s := range p.Skills {
				if s.Name == skillFilter {
					filtered = append(filtered, p)
					break
				}
			}
		}
		positions = filtered
	}

	// If it's an HTMX request for the grid, only render the grid
	if r.Header.Get("HX-Request") == "true" {
		component := components.PositionGrid(positions)
		err = component.Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Otherwise render the full page
	skills, err := h.Queries.ListSkills(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	companies, err := h.Queries.ListCompanies(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	component := components.PositionList(positions, skills, companies)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PositionHandler) CreatePosition(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	companyIDStr := r.FormValue("company_id")
	title := r.FormValue("title")
	location := r.FormValue("location")
	workMode := r.FormValue("work_mode")
	salaryRange := r.FormValue("salary_range")
	postURL := r.FormValue("post_url")

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		http.Error(w, "Invalid company ID", http.StatusBadRequest)
		return
	}

	userID := GetUserID(r.Context())

	var locationPtr *string
	if location != "" {
		locationPtr = &location
	}
	var salaryRangePtr *string
	if salaryRange != "" {
		salaryRangePtr = &salaryRange
	}
	var postURLPtr *string
	if postURL != "" {
		postURLPtr = &postURL
	}

	pos, err := h.Queries.CreatePosition(context.Background(), db.CreatePositionParams{
		CompanyID:   companyID,
		UserID:      userID,
		Title:       title,
		Location:    locationPtr,
		WorkMode:    workMode,
		SalaryRange: salaryRangePtr,
		PostURL:     postURLPtr,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch details for the card (like company name)
	posDetails, err := h.Queries.GetPosition(context.Background(), pos.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	component := components.PositionCard(posDetails)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PositionHandler) DeletePosition(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Security: check if position belongs to user
	pos, err := h.Queries.GetPosition(context.Background(), id)
	if err != nil {
		http.Error(w, "Position not found", http.StatusNotFound)
		return
	}
	if pos.UserID != GetUserID(r.Context()) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	err = h.Queries.DeletePosition(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

