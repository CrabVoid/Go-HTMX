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
	stats, err := h.Queries.GetDashboardStats(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	component := components.DashboardCounts(stats)
	component.Render(r.Context(), w)
}

func (h *DashboardHandler) GetPipeline(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r.Context())
	stats, err := h.Queries.GetDashboardStats(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	component := components.DashboardPipeline(stats)
	component.Render(r.Context(), w)
}
