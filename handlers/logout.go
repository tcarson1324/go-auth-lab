package handlers

import (
	"database/sql"
	"errors"
	"go-auth-lab/sessions"
	"net/http"
)

func HandlerLogout(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		if !errors.Is(err, http.ErrNoCookie) {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		sessions.ClearCookie(w)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = sessions.DeleteSession(db, sessionCookie.Value)
	if err != nil {
		if !errors.Is(err, sessions.ErrSessionNotFound) {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	sessions.ClearCookie(w)
	w.WriteHeader(http.StatusNoContent)
}
