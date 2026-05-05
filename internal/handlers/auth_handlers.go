package handlers

import (
	"context"
	"net/http"
	"time"

	"internship-manager/components"
	"internship-manager/internal/db"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Queries *db.Queries
}

func NewAuthHandler(queries *db.Queries) *AuthHandler {
	return &AuthHandler{Queries: queries}
}

func (h *AuthHandler) GetLogin(w http.ResponseWriter, r *http.Request) {
	component := components.LoginPage()
	component.Render(r.Context(), w)
}

func (h *AuthHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := h.Queries.GetUserByEmail(context.Background(), email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	h.setSession(w, user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) GetRegister(w http.ResponseWriter, r *http.Request) {
	component := components.RegisterPage()
	component.Render(r.Context(), w)
}

func (h *AuthHandler) PostRegister(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := h.Queries.CreateUser(context.Background(), db.CreateUserParams{
		Email:        email,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	h.setSession(w, user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) PostLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *AuthHandler) setSession(w http.ResponseWriter, userID uuid.UUID) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    userID.String(),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})
}
