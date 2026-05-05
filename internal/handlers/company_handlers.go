package handlers

import (
	"context"
	"net/http"

	"internship-manager/components"
	"internship-manager/internal/db"
)

type CompanyHandler struct {
	Queries *db.Queries
}

func NewCompanyHandler(queries *db.Queries) *CompanyHandler {
	return &CompanyHandler{Queries: queries}
}

func (h *CompanyHandler) ListCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := h.Queries.ListCompanies(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	component := components.CompanyList(companies)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	industry := r.FormValue("industry")
	website := r.FormValue("website")

	var industryPtr *string
	if industry != "" {
		industryPtr = &industry
	}
	var websitePtr *string
	if website != "" {
		websitePtr = &website
	}

	company, err := h.Queries.CreateCompany(context.Background(), db.CreateCompanyParams{
		Name:     name,
		Industry: industryPtr,
		Website:  websitePtr,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render only the new company card (HTMX will append it)
	component := components.CompanyCard(company)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
