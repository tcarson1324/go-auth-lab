package handlers

import (
	"database/sql"
	"errors"
	"go-auth-lab/passwords"
	"go-auth-lab/sessions"
	"go-auth-lab/users"
	"html/template"
	"net/http"
)

func HandlerLoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, "couldn't parse template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}

func HandlerLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	name := r.FormValue("username")
	password := r.FormValue("password")

	user, err := users.GetUserByName(name, db)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			http.Error(w, "invalid username or password", http.StatusUnauthorized)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	found, err := passwords.Verify(password, user.PasswordHash)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if !found {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}
	token, err := sessions.GenerateToken()
	if err != nil {
		http.Error(w, "could not create session", http.StatusInternalServerError)
		return 
	}

	tokenHash := sessions.HashToken(token)
	expiresAt, err := sessions.CreateSession(db, user.ID, tokenHash)
	if err != nil {
		http.Error(w, "could not create session", http.StatusInternalServerError)
		return
	}
	sessions.SetCookie(w, token, expiresAt)

}
