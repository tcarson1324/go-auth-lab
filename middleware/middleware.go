package middleware

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"errors"
	"go-auth-lab/sessions"
	"net/http"
)

type contextKey string

const userIDKey contextKey = "userID"
const csrfIDKEY contextKey = "csrfKey"

func RequireCSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedToken, ok := CSRFTokenFromContext(r.Context())

		if !ok {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		submittedToken := r.FormValue("csrf_token")
		if expectedToken == "" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if submittedToken == "" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		result := subtle.ConstantTimeCompare(
			[]byte(expectedToken),
			[]byte(submittedToken),
		)

		if result != 1 {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

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
		ctx = context.WithValue(ctx, csrfIDKEY, session.CSRFToken)
		r = r.WithContext(ctx)
		// Authentication succeeded. Execute the protected handler.
		next.ServeHTTP(w, r)
	})
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDKey).(int64)
	return userID, ok
}
func CSRFTokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(csrfIDKEY).(string)
	return token, ok
}

func ContextWithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}
