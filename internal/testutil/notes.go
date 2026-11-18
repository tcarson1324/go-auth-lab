package testutil

import (
	"database/sql"
	"testing"
)

func CreateTestNote(t *testing.T, db *sql.DB, userID int64) int64 {
	t.Helper()
	var noteID int64

	err := db.QueryRow(`
		INSERT INTO notes (user_id, content)
		VALUES ($1, $2)
		RETURNING id
	`, userID, "This is a test note").Scan(&noteID)
	if err != nil {
		t.Fatalf("could not create test note: %v", err)
	}

	// Registered after user cleanup, so the note is removed first.
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM notes WHERE id = $1", noteID)
		if err != nil {
			t.Errorf("could not clean up test note: %v", err)
		}
	})

	return noteID
}

func EditTestNote(t *testing.T, db *sql.DB, noteID, userID int64) error {
	t.Helper()
	result, err := db.Exec(`
	UPDATE notes
	SET content = "this test note has been updated"
	WHERE id = $2
	AND user_id = $3
	`, noteID, userID)
	rows, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("get rows affected: %v", err)
	}

	if rows != 1 {
		t.Errorf("rows affected = %d, want 1", rows)
	}
	return nil
}
