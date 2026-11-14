package notes

import (
	"errors"
	"testing"

	"go-auth-lab/internal/testutil"
)

func TestDeleteNoteNotFound(t *testing.T) {
	// Arrange
	db := testutil.SetupTestDB(t)

	userID := testutil.CreateTestUser(t, db)

	// Act
	err := DeleteNoteByID(db, -1, userID)

	// Assert
	if !errors.Is(err, ErrNoteNotFound) {
		t.Errorf("expected ErrNoteNotFound, got %v", err)
	}
}

func TestDeleteNoteSuccess(t *testing.T) {
	// Arrange
	db := testutil.SetupTestDB(t)

	userID := testutil.CreateTestUser(t, db)
	noteID := testutil.CreateTestNote(t, db, userID)

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
	db := testutil.SetupTestDB(t)

	userAID := testutil.CreateTestUser(t, db)
	userBID := testutil.CreateTestUser(t, db)
	noteID := testutil.CreateTestNote(t, db, userAID)

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
