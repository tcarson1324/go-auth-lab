package handlers

import (
	"database/sql"
	"errors"
	"go-auth-lab/users"
	"html/template"
	"net/http"
	"unicode/utf8"
)

func HandlerRegisterPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/register.html")
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

func HandlerAddUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" {
		http.Error(w, "missing username or password", 400)
		return
	}
	if usernameLength := utf8.RuneCountInString(username); usernameLength < 3 || usernameLength > 30 {
		http.Error(w, "username must be between 3 and 30 characters", 400)
		return
	}
	for _, char := range username {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_') {

			http.Error(w, "username contains invalid characters", http.StatusBadRequest)
			return
		}
	}
	if password == "" {
		http.Error(w, "missing username or password", http.StatusBadRequest)
		return
	}
	if utf8.RuneCountInString(password) < 6 {
		http.Error(w, "password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	if len(password) > 72 {
		http.Error(w, "password exceeds 72 bytes", http.StatusBadRequest)
		return
	}
	err := users.AddUser(username, password, db)
	if err != nil {
		if errors.Is(err, users.ErrUsernameTaken) {
			http.Error(w, "username already exists", http.StatusConflict)
			return
		}
		http.Error(w, "couldn't create new user", http.StatusInternalServerError)
		return
	}
}
