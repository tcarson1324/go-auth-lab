package main

import (
	"database/sql"
	"go-auth-lab/handlers"
	"go-auth-lab/middleware"
	"log"
	"net/http"

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

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlerHomePage(w, r, db)
	})
	mux.Handle("POST /notes",
		middleware.RequireAuth(db,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlers.HandlerCreateNote(w, r, db)
			}),
		),
	)
	mux.Handle("POST /notes/delete",
		middleware.RequireAuth(db,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlers.HandlerDelteNote(w, r, db)
			})))
	mux.Handle("GET /profile",
		middleware.RequireAuth(db,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlers.HandlerProfile(w, r, db)
			}),
		),
	)

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
