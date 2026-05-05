package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"internship-manager/components"
	"internship-manager/internal/db"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ApplicationHandler struct {
	Queries *db.Queries
}

func NewApplicationHandler(queries *db.Queries) *ApplicationHandler {
	return &ApplicationHandler{Queries: queries}
}

func (h *ApplicationHandler) ListApplications(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r.Context())
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
	if app.UserID != GetUserID(r.Context()) {
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

	userID := GetUserID(r.Context())
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

	if app.UserID != GetUserID(r.Context()) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	component := components.ApplicationDetail(app)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


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
