package notes

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupTestDB(t *testing.T) *sql.DB {
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
func createTestUser(t *testing.T, db *sql.DB) int64 {
	t.Helper()
	bytes := make([]byte, 4)

	_, err := rand.Read(bytes)
	if err != nil {
		t.Fatalf("could not generate random bytes: %v", err)
	}
	userName := "test_user_" + hex.EncodeToString(bytes)
	var userID int64

	err = db.QueryRow(`
		INSERT INTO users (username, password_hash)
		VALUES ($1, 'test_hash')
		RETURNING id
	`, userName).Scan(&userID)

	if err != nil {
		t.Fatalf("could not create test user: %v", err)
	}

	t.Cleanup(func() {
		_, err := db.Exec(
			"DELETE FROM users WHERE id = $1",
			userID,
		)
		if err != nil {
			t.Errorf("could not clean up user: %v", err)
		}
	})
	return userID
}

func createTestNote(t *testing.T, db *sql.DB, userID int64) int64 {
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

func TestDeleteNoteNotFound(t *testing.T) {
	// Arrange
	db := setupTestDB(t)

	userID := createTestUser(t, db)

	// Act
	err := DeleteNoteByID(db, -1, userID)

	// Assert
	if !errors.Is(err, ErrNoteNotFound) {
		t.Errorf("expected ErrNoteNotFound, got %v", err)
	}
}

func TestDeleteNoteSuccess(t *testing.T) {
	// Arrange
	db := setupTestDB(t)

	userID := createTestUser(t, db)
	noteID := createTestNote(t, db, userID)

	// Act
	err := DeleteNoteByID(db, noteID, userID)

	// Assert
	if err != nil {
		t.Fatalf("expected successful deletion, got error: %v", err)
	}

	// Verify that the note was deleted.
	var count int

	err = db.QueryRow(`
		SELECT COUNT(*)
		FROM notes
		WHERE id = $1 AND user_id = $2
	`, noteID, userID).Scan(&count)

	if err != nil {
		t.Fatalf("could not verify note deletion: %v", err)
	}

	if count != 0 {
		t.Errorf("expected 0 notes after deletion, found %d", count)
	}
}

func TestDeleteNoteWrongUser(t *testing.T) {
	// Arrange
	db := setupTestDB(t)

	userAID := createTestUser(t, db)
	userBID := createTestUser(t, db)
	noteID := createTestNote(t, db, userAID)

	// Act
	err := DeleteNoteByID(db, noteID, userBID)

	// Assert
	if !errors.Is(err, ErrNoteNotFound) {
		t.Errorf("expected ErrNoteNotFound, got %v", err)
	}

	// Verify that the unauthorized attempt did not delete the note.
	var count int

	err = db.QueryRow(`
		SELECT COUNT(*)
		FROM notes
		WHERE id = $1 AND user_id = $2
	`, noteID, userAID).Scan(&count)

	if err != nil {
		t.Fatalf("could not verify note existence: %v", err)
	}

	if count != 1 {
		t.Errorf("expected note to still exist, found %d matching notes", count)
	}
}
