package handlers

import (
	"context"
	"net/http"

	"internship-manager/components"
	"internship-manager/internal/db"
)

type PositionHandler struct {
	Queries *db.Queries
}

func NewPositionHandler(queries *db.Queries) *PositionHandler {
	return &PositionHandler{Queries: queries}
}

func (h *PositionHandler) ListPositions(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	
	positions, err := h.Queries.ListPositions(ctx)
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

	component := components.PositionList(positions, skills)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
