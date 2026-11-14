package testutil

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open(
		"pgx",
		"host=/var/run/postgresql dbname=go_auth_lab_test user=taylor",
	)
	if err != nil {
		t.Fatalf("could not open test database: %v", err)
	}

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("could not close database: %v", err)
		}
	})

	if err := db.Ping(); err != nil {
		t.Fatalf("could not connect to test database: %v", err)
	}
	return db
}
