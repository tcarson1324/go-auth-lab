package handlers

import (
	"database/sql"
	"go-auth-lab/middleware"
	"go-auth-lab/notes"
	"net/http"
	"strconv"
	"strings"
	"errors"
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
func HandlerDelteNote(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	noteID := r.FormValue("note_id")
	if strings.TrimSpace(noteID) == "" {
		http.Error(w, "note id cannot be empty", http.StatusBadRequest)
		return
	}

	ID, err := strconv.ParseInt(noteID, 10, 64)
	if err != nil {
		http.Error(w, "id must be a valid number", http.StatusBadRequest)
		return
	}
	if ID <= 0 {
		http.Error(w, "id must be greater than 0", http.StatusBadRequest)
		return
	}
	err = notes.DeleteNoteByID(db, ID, userID)
	if err != nil {
		if errors.Is(err, notes.ErrNoteNotFound) {
			http.Error(w, "note not found", http.StatusNotFound)
			return
		}
		http.Error(w, "could not delete note", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
