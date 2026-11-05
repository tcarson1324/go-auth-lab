package handlers

import (
	"database/sql"
	"errors"
	"go-auth-lab/sessions"
	"html/template"
	"net/http"
)

type HomePageData struct {
	LoggedIn bool
}

func HandlerHomePage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data := HomePageData{
		LoggedIn: false,
	}

	sessionCookie, err := r.Cookie("session")

	if err == nil {
		_, err = sessions.GetSession(db, sessionCookie.Value)

		if err == nil {
			data.LoggedIn = true
		} else if !errors.Is(err, sessions.ErrSessionNotFound) {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "couldn't parse template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}
