/*
+--------------------------------------------------------------------------------+
| PACKAGE DECLARATION                                                            |
| PURPOSE: Authentication logic for the Internship Manager.                      |
| INFO: Handles registration, login, logout and session management.              |
+--------------------------------------------------------------------------------+
*/
package handlers

/*
+--------------------------------------------------------------------------------+
| EXTERNAL & INTERNAL IMPORTS                                                    |
| PURPOSE: Import standard library and third-party dependencies.                 |
| INFO: Includes bcrypt for security and uuid for session tracking.              |
+--------------------------------------------------------------------------------+
*/
import (
	"context"
	"net/http"
	"time"

	"internship-manager/components"
	"internship-manager/internal/db"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

/*
+--------------------------------------------------------------------------------+
| AUTH HANDLER STRUCT                                                            |
| PURPOSE: State management for authentication operations.                       |
| INFO: Holds a reference to the database queries object.                         |
+--------------------------------------------------------------------------------+
*/
type AuthHandler struct {
	Queries *db.Queries
}

/*
+--------------------------------------------------------------------------------+
| CONSTRUCTOR: NEW AUTH HANDLER                                                  |
| PURPOSE: Initialize a new AuthHandler instance.                                |
| INFO: Connects the handler to the data access layer.                           |
+--------------------------------------------------------------------------------+
*/
func NewAuthHandler(queries *db.Queries) *AuthHandler {
	return &AuthHandler{Queries: queries}
}

/*
+--------------------------------------------------------------------------------+
| HANDLER: GET LOGIN                                                             |
| PURPOSE: Renders the login page.                                               |
| INFO: Returns the full HTML page for user authentication.                      |
+--------------------------------------------------------------------------------+
*/
func (h *AuthHandler) GetLogin(w http.ResponseWriter, r *http.Request) {
	component := components.LoginPage()
	component.Render(r.Context(), w)
}

/*
+--------------------------------------------------------------------------------+
| HANDLER: POST LOGIN                                                            |
| PURPOSE: Processes the login form submission.                                  |
| INFO: Validates credentials and establishes a session cookie.                  |
+--------------------------------------------------------------------------------+
*/
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

/*
+--------------------------------------------------------------------------------+
| HANDLER: GET REGISTER                                                          |
| PURPOSE: Renders the registration page.                                        |
| INFO: Allows new users to create an account.                                   |
+--------------------------------------------------------------------------------+
*/
func (h *AuthHandler) GetRegister(w http.ResponseWriter, r *http.Request) {
	component := components.RegisterPage()
	component.Render(r.Context(), w)
}

/*
+--------------------------------------------------------------------------------+
| HANDLER: POST REGISTER                                                         |
| PURPOSE: Processes account creation requests.                                   |
| INFO: Hashes passwords and stores new user records in the database.            |
+--------------------------------------------------------------------------------+
*/
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

/*
+--------------------------------------------------------------------------------+
| HANDLER: POST LOGOUT                                                           |
| PURPOSE: Terminates the current user session.                                  |
| INFO: Invalidates the session cookie and redirects to login.                   |
+--------------------------------------------------------------------------------+
*/
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

/*
+--------------------------------------------------------------------------------+
| PRIVATE HELPER: SET SESSION                                                    |
| PURPOSE: Sets an HTTP-only cookie for session persistence.                     |
| INFO: Stores the user's UUID in the browser for authentication identification.|
+--------------------------------------------------------------------------------+
*/
func (h *AuthHandler) setSession(w http.ResponseWriter, userID uuid.UUID) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    userID.String(),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})
}
