package handlers

import (
	"database/sql"
	"go-auth-lab/middleware"
	"go-auth-lab/notes"
	"net/http"
	"strings"
)

func HandlerCreateNote(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	content := r.FormValue("content")
	if strings.TrimSpace(content) == "" {
		http.Error(w, "note cannot be empty", http.StatusBadRequest)
		return
	}
	err := notes.CreateNote(db, userID, content)
	if err != nil {
		http.Error(w, "could not create note", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
