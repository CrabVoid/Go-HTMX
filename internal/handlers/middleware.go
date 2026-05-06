/*
+--------------------------------------------------------------------------------+
| PACKAGE DECLARATION                                                            |
| PURPOSE: Global middleware for the Internship Manager.                         |
| INFO: Implements authentication checks and context injection.                  |
+--------------------------------------------------------------------------------+
*/
package handlers

/*
+--------------------------------------------------------------------------------+
| EXTERNAL & INTERNAL IMPORTS                                                    |
+| PURPOSE: Import standard library and internal auth package.                    |
+| INFO: Includes UUID for session parsing.                                       |
++--------------------------------------------------------------------------------+
+*/
import (
	"context"
	"net/http"

	"internship-manager/internal/auth"

	"github.com/google/uuid"
)

/*
+--------------------------------------------------------------------------------+
| MIDDLEWARE: AUTH MIDDLEWARE                                                    |
+| PURPOSE: Protects routes by verifying session cookies.                         |
+| INFO: Injects the UserID into the request context if authenticated.            |
++--------------------------------------------------------------------------------+
+*/
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userID, err := uuid.Parse(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), auth.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
