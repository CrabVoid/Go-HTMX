package handlers

import (
	"context"
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
