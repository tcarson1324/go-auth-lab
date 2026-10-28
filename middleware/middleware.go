package middleware

import (
	"database/sql"
	"errors"
	"net/http"
	"context"
	"go-auth-lab/sessions"
)
type contextKey string	

const userIDKey contextKey = "userID"
func RequireAuth(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Retrieve the session cookie from the request.
		sessionCookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate the session against PostgreSQL.
		session, err := sessions.GetSession(db, sessionCookie.Value)
		if err != nil {
			if errors.Is(err, sessions.ErrSessionNotFound) {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, session.UserID)
		r = r.WithContext(ctx)
		// Authentication succeeded. Execute the protected handler.
		next.ServeHTTP(w, r)
	})
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
    userID, ok := ctx.Value(userIDKey).(int64)
    return userID, ok
}