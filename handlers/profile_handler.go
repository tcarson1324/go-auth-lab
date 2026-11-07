package handlers

import (
	"database/sql"
	"errors"
	"go-auth-lab/middleware"
	"go-auth-lab/notes"
	"go-auth-lab/users"
	"html/template"
	"net/http"
)

type ProfilePageData struct {
	UserID   int64
	Username string
	Notes    []notes.Note
}

func HandlerProfile(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	user, err := users.GetUserByID(userID, db)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	notes, err := notes.GetNotesByUserID(db, userID)
	if err != nil {
		http.Error(w, "couldn't retieve notes", http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles("templates/profile.html")
	if err != nil {
		http.Error(w, "couldn't parse template", http.StatusInternalServerError)
		return
	}
	profileData := ProfilePageData{
		UserID:   userID,
		Username: user.Username,
		Notes:    notes,
	}
	err = tmpl.Execute(w, profileData)
	if err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}
