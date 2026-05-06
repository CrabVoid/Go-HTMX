/*
+--------------------------------------------------------------------------------+
| PACKAGE DECLARATION                                                            |
| PURPOSE: Handles request logic for the Dashboard view.                         |
| INFO: Aggregates metrics and pipeline information for the main view.           |
+--------------------------------------------------------------------------------+
*/
package handlers

/*
+--------------------------------------------------------------------------------+
| EXTERNAL & INTERNAL IMPORTS                                                    |
| PURPOSE: Import standard library and project dependencies.                     |
| INFO: Includes templ components and database access objects.                   |
+--------------------------------------------------------------------------------+
*/
import (
	"net/http"

	"internship-manager/components"
	"internship-manager/internal/auth"
	"internship-manager/internal/db"
)

/*
+--------------------------------------------------------------------------------+
| DASHBOARD HANDLER STRUCT                                                       |
| PURPOSE: State management for dashboard-related operations.                    |
| INFO: Holds a reference to the generated database queries.                      |
+--------------------------------------------------------------------------------+
*/
type DashboardHandler struct {
	Queries *db.Queries
}

/*
+--------------------------------------------------------------------------------+
| CONSTRUCTOR: NEW DASHBOARD HANDLER                                             |
| PURPOSE: Initialize a new DashboardHandler instance.                           |
| INFO: Accepts a pointer to db.Queries for data access.                         |
+--------------------------------------------------------------------------------+
*/
func NewDashboardHandler(queries *db.Queries) *DashboardHandler {
	return &DashboardHandler{Queries: queries}
}

/*
+--------------------------------------------------------------------------------+
| HANDLER: GET COUNTS                                                            |
| PURPOSE: Fetches high-level metrics for the dashboard header.                  |
| INFO: Returns an HTMX fragment containing total, applied, and interview counts.|
+--------------------------------------------------------------------------------+
*/
func (h *DashboardHandler) GetCounts(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r.Context())
	stats, err := h.Queries.GetDashboardStats(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	component := components.DashboardCounts(stats)
	component.Render(r.Context(), w)
}

/*
+--------------------------------------------------------------------------------+
| HANDLER: GET PIPELINE                                                          |
| PURPOSE: Fetches application pipeline status data.                             |
| INFO: Returns an HTMX fragment visualising the flow of applications.           |
+--------------------------------------------------------------------------------+
*/
func (h *DashboardHandler) GetPipeline(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r.Context())
	stats, err := h.Queries.GetDashboardStats(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	component := components.DashboardPipeline(stats)
	component.Render(r.Context(), w)
}
