/*
+--------------------------------------------------------------------------------+
| PACKAGE DECLARATION                                                            |
| PURPOSE: CRUD operations for Company entities.                                 |
| INFO: Handles listing, creating, updating, and deleting companies.             |
+--------------------------------------------------------------------------------+
*/
package handlers

/*
+--------------------------------------------------------------------------------+
| EXTERNAL & INTERNAL IMPORTS                                                    |
| PURPOSE: Import standard library and third-party dependencies.                 |
| INFO: Includes Chi for URL params and UUID for entity identification.         |
+--------------------------------------------------------------------------------+
*/
import (
	"context"
	"net/http"

	"internship-manager/components"
	"internship-manager/internal/db"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

/*
+--------------------------------------------------------------------------------+
| COMPANY HANDLER STRUCT                                                         |
| PURPOSE: State management for company-related operations.                      |
| INFO: Holds a reference to the database queries object.                         |
+--------------------------------------------------------------------------------+
*/
type CompanyHandler struct {
	Queries *db.Queries
}

/*
+--------------------------------------------------------------------------------+
| CONSTRUCTOR: NEW COMPANY HANDLER                                               |
| PURPOSE: Initialize a new CompanyHandler instance.                              |
| INFO: Connects the handler to the data access layer.                           |
+--------------------------------------------------------------------------------+
*/
func NewCompanyHandler(queries *db.Queries) *CompanyHandler {
	return &CompanyHandler{Queries: queries}
}

/*
+--------------------------------------------------------------------------------+
| HANDLER: LIST COMPANIES                                                        |
| PURPOSE: Fetches and displays all companies.                                   |
| INFO: Renders the company list component.                                      |
+--------------------------------------------------------------------------------+
*/
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

/*
+--------------------------------------------------------------------------------+
| HANDLER: CREATE COMPANY                                                        |
| PURPOSE: Processes the creation of a new company.                              |
| INFO: Accepts form data and returns the newly created company card.            |
+--------------------------------------------------------------------------------+
*/
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

/*
+--------------------------------------------------------------------------------+
| HANDLER: GET COMPANY FORM                                                      |
| PURPOSE: Renders the edit form for a specific company.                         |
| INFO: Fetches company data by UUID and returns the form component.             |
+--------------------------------------------------------------------------------+
*/
func (h *CompanyHandler) GetCompanyForm(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	company, err := h.Queries.GetCompany(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	component := components.EditCompanyForm(company)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/*
+--------------------------------------------------------------------------------+
| HANDLER: GET COMPANY CARD                                                      |
| PURPOSE: Renders a single company card.                                        |
| INFO: Used for HTMX updates to refresh specific entries in the list.           |
+--------------------------------------------------------------------------------+
*/
func (h *CompanyHandler) GetCompanyCard(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	company, err := h.Queries.GetCompany(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	component := components.CompanyCard(company)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/*
+--------------------------------------------------------------------------------+
| HANDLER: UPDATE COMPANY                                                        |
| PURPOSE: Processes updates to an existing company record.                       |
| INFO: Validates the UUID and updates the database record.                      |
+--------------------------------------------------------------------------------+
*/
func (h *CompanyHandler) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
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

	company, err := h.Queries.UpdateCompany(context.Background(), db.UpdateCompanyParams{
		ID:       id,
		Name:     name,
		Industry: industryPtr,
		Website:  websitePtr,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	component := components.CompanyCard(company)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/*
+--------------------------------------------------------------------------------+
| HANDLER: DELETE COMPANY                                                        |
| PURPOSE: Removes a company record from the database.                           |
| INFO: Deletes the record by UUID and returns a 200 OK status.                  |
+--------------------------------------------------------------------------------+
*/
func (h *CompanyHandler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Queries.DeleteCompany(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
