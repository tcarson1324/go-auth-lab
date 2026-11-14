package handlers

import (
	"database/sql"
	"go-auth-lab/middleware"
	"go-auth-lab/notes"
	"net/http"
	"strings"
	"strconv"
	"errors"
)

func HandlerUpdateNote(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
	ID := r.FormValue("note_id")
	noteID, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		http.Error(w, "noteID must be an integer", http.StatusBadRequest)
		return
	}
	if noteID <= 0 {
		http.Error(w, "id must be greater than 0", 400)
		return
	}
	err = notes.UpdateNote(db, noteID, userID, content)
	if err != nil {
		if errors.Is(err, notes.ErrNoteNotFound) {
			http.Error(w, "note not found", http.StatusNotFound)
			return
		}
		http.Error(w, "could not update note", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
