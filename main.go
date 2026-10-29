package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"go-auth-lab/handlers"
	"go-auth-lab/middleware"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	db, err := sql.Open(
		"pgx",
		"host=/var/run/postgresql dbname=go_auth_lab user=taylor",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	profileHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.UserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "You are logged in as user %d\n", userID)
	})
	mux.Handle("GET /profile", middleware.RequireAuth(db, profileHandler))

	mux.HandleFunc("GET /register", handlers.HandlerRegisterPage)
	mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlerAddUser(w, r, db)
	})
	mux.HandleFunc("GET /login", handlers.HandlerLoginPage)
	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlerLogin(w, r, db)
	})
	mux.HandleFunc("POST /logout", func(w http.ResponseWriter, r *http.Request) {
    handlers.HandlerLogout(w, r, db)
})
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}

}
