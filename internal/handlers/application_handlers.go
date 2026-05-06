/*
+--------------------------------------------------------------------------------+
| PACKAGE DECLARATION                                                            |
| PURPOSE: Request logic for Job Applications.                                   |
| INFO: Manages the lifecycle of applications and associated interviews.          |
+--------------------------------------------------------------------------------+
*/
package handlers

/*
+--------------------------------------------------------------------------------+
| EXTERNAL & INTERNAL IMPORTS                                                    |
| PURPOSE: Import standard library and project dependencies.                     |
| INFO: Includes time for interview scheduling and Chi for routing.               |
+--------------------------------------------------------------------------------+
*/
import (
	"context"
	"fmt"
	"net/http"
	"time"

	"internship-manager/components"
	"internship-manager/internal/auth"
	"internship-manager/internal/db"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

/*
+--------------------------------------------------------------------------------+
| APPLICATION HANDLER STRUCT                                                     |
| PURPOSE: State management for application-related logic.                       |
| INFO: Holds a reference to the database access layer.                           |
+--------------------------------------------------------------------------------+
*/
type ApplicationHandler struct {
	Queries *db.Queries
}

/*
+--------------------------------------------------------------------------------+
| CONSTRUCTOR: NEW APPLICATION HANDLER                                           |
| PURPOSE: Initialize a new ApplicationHandler instance.                         |
| INFO: Injects the database queries object.                                     |
+--------------------------------------------------------------------------------+
*/
func NewApplicationHandler(queries *db.Queries) *ApplicationHandler {
	return &ApplicationHandler{Queries: queries}
}

/*
+--------------------------------------------------------------------------------+
| HANDLER: LIST APPLICATIONS                                                     |
| PURPOSE: Displays all applications for the current user.                       |
| INFO: Renders the application list component.                                  |
+--------------------------------------------------------------------------------+
*/
func (h *ApplicationHandler) ListApplications(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r.Context())
	apps, err := h.Queries.ListApplications(context.Background(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	component := components.ApplicationList(apps)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/*
+--------------------------------------------------------------------------------+
| HANDLER: DELETE APPLICATION                                                    |
| PURPOSE: Removes an application from the database.                             |
| INFO: Validates ownership before deletion.                                     |
+--------------------------------------------------------------------------------+
*/
func (h *ApplicationHandler) DeleteApplication(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Security: check if application belongs to user
	app, err := h.Queries.GetApplication(context.Background(), id)
	if err != nil {
		http.Error(w, "Application not found", http.StatusNotFound)
		return
	}
	if app.UserID != auth.GetUserID(r.Context()) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	err = h.Queries.DeleteApplication(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

/*
+--------------------------------------------------------------------------------+
| HANDLER: CREATE APPLICATION                                                    |
| PURPOSE: Creates a new application for a specific position.                   |
| INFO: Defaults status to 'Applied' and redirects to detail view.               |
+--------------------------------------------------------------------------------+
*/
func (h *ApplicationHandler) CreateApplication(w http.ResponseWriter, r *http.Request) {
	positionIDStr := r.URL.Query().Get("position_id")
	if positionIDStr == "" {
		http.Error(w, "Missing position_id", http.StatusBadRequest)
		return
	}

	positionID, err := uuid.Parse(positionIDStr)
	if err != nil {
		http.Error(w, "Invalid position ID", http.StatusBadRequest)
		return
	}

	userID := auth.GetUserID(r.Context())
	app, err := h.Queries.CreateApplication(context.Background(), db.CreateApplicationParams{
		PositionID: positionID,
		UserID:     userID,
		Status:     "Applied", // Default status
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to the new application detail page
	http.Redirect(w, r, fmt.Sprintf("/applications/%s", app.ID.String()), http.StatusSeeOther)
}

/*
+--------------------------------------------------------------------------------+
| HANDLER: GET APPLICATION                                                       |
| PURPOSE: Renders the detail view for a specific application.                   |
| INFO: Includes position info, status history, and interviews.                  |
+--------------------------------------------------------------------------------+
*/
func (h *ApplicationHandler) GetApplication(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid application ID", http.StatusBadRequest)
		return
	}

	app, err := h.Queries.GetApplication(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if app.UserID != auth.GetUserID(r.Context()) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	component := components.ApplicationDetail(app)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/*
+--------------------------------------------------------------------------------+
| HANDLER: CREATE INTERVIEW                                                      |
| PURPOSE: Schedules a new interview for an application.                         |
| INFO: Parses scheduled time and optional notes.                                |
+--------------------------------------------------------------------------------+
*/
func (h *ApplicationHandler) CreateInterview(w http.ResponseWriter, r *http.Request) {
	appIDStr := chi.URLParam(r, "id")
	appID, err := uuid.Parse(appIDStr)
	if err != nil {
		http.Error(w, "Invalid application ID", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stageName := r.FormValue("stage_name")
	scheduledAtStr := r.FormValue("scheduled_at")
	notes := r.FormValue("notes")

	scheduledAt, err := time.Parse("2006-01-02T15:04", scheduledAtStr)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	var notesPtr *string
	if notes != "" {
		notesPtr = &notes
	}

	interview, err := h.Queries.CreateInterview(context.Background(), db.CreateInterviewParams{
		ApplicationID: appID,
		StageName:     stageName,
		ScheduledAt:   scheduledAt,
		Notes:         notesPtr,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	component := components.InterviewItem(interview)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
