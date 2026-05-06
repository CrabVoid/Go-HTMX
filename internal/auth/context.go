/*
+--------------------------------------------------------------------------------+
| PACKAGE DECLARATION                                                            |
| PURPOSE: Context-based authentication helpers.                                 |
| INFO: Provides keys and getters for UserID in request contexts.                |
+--------------------------------------------------------------------------------+
*/
package auth

/*
+--------------------------------------------------------------------------------+
| EXTERNAL IMPORTS                                                               |
| PURPOSE: Standard library dependencies.                                        |
| INFO: Includes context for KV storage and UUID for identity.                   |
+--------------------------------------------------------------------------------+
*/
import (
	"context"

	"github.com/google/uuid"
)

/*
+--------------------------------------------------------------------------------+
| CONTEXT TYPES & KEYS                                                           |
| PURPOSE: Type-safe keys for context values.                                    |
| INFO: Uses a private string type to avoid collisions.                          |
+--------------------------------------------------------------------------------+
*/
type contextKey string

const UserIDKey contextKey = "userID"

/*
+--------------------------------------------------------------------------------+
| HELPER: GET USER ID                                                            |
| PURPOSE: Extract the UserID from a context.                                    |
| INFO: Returns uuid.Nil if the key is not found or invalid.                     |
+--------------------------------------------------------------------------------+
*/
func GetUserID(ctx context.Context) uuid.UUID {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return userID
}
