package handlers

import (
	"fmt"
	"go-auth-lab/internal/testutil"
	"go-auth-lab/middleware"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDeleteNoteInvalidID(t *testing.T) {
	db := testutil.SetupTestDB(t)
	handler := func(w http.ResponseWriter, r *http.Request) {
		HandlerDelteNote(w, r, db)
	}
	req := httptest.NewRequest(
		http.MethodPost,
		"/notes/delete",
		strings.NewReader("note_id=abc"),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := middleware.ContextWithUserID(req.Context(), int64(1))
	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()
	handler(recorder, req)
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", recorder.Code)
	}
}

func TestDeleteNoteNotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)

	// Create a user and a note.
	userID := testutil.CreateTestUser(t, db)
	noteID := testutil.CreateTestNote(t, db, userID)

	// Delete the note directly from PostgreSQL.
	_, err := db.Exec("DELETE FROM notes WHERE id = $1", noteID)
	if err != nil {
		t.Fatalf("could not delete test note: %v", err)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		HandlerDelteNote(w, r, db)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/notes/delete",
		strings.NewReader(fmt.Sprintf("note_id=%d", noteID)),
	)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ctx := middleware.ContextWithUserID(req.Context(), userID)
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()

	handler(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", recorder.Code)
	}
}

func TestDeleteNoteWrongUser(t *testing.T) {
	db := testutil.SetupTestDB(t)
	userID1 := testutil.CreateTestUser(t, db)
	noteID := testutil.CreateTestNote(t, db, userID1)
	userID2 := testutil.CreateTestUser(t, db)
	handler := func(w http.ResponseWriter, r *http.Request) {
		HandlerDelteNote(w, r, db)
	}
	req := httptest.NewRequest(
		http.MethodPost,
		"/notes/delete",
		strings.NewReader(fmt.Sprintf("note_id=%d", noteID)),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := middleware.ContextWithUserID(req.Context(), userID2)
	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()
	handler(recorder, req)
	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected status 400, got %d", recorder.Code)
	}
	var count int
	err := db.QueryRow(`
	SELECT COUNT(*) FROM notes 
	WHERE id = $1
	`, noteID).Scan(&count)
	if err != nil {
		t.Fatalf("failed to query database: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 note after failed deletion, got %d", count)
	}

}
func TestDeleteNoteSuccess(t *testing.T) {
	db := testutil.SetupTestDB(t)
	userID := testutil.CreateTestUser(t, db)
	noteID := testutil.CreateTestNote(t, db, userID)
	handler := func(w http.ResponseWriter, r *http.Request) {
		HandlerDelteNote(w, r, db)
	}
	req := httptest.NewRequest(
		http.MethodPost,
		"/notes/delete",
		strings.NewReader(fmt.Sprintf("note_id=%d", noteID)),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := middleware.ContextWithUserID(req.Context(), userID)
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()

	handler(recorder, req)
	if recorder.Code != http.StatusSeeOther {
		t.Errorf("expected status 303, got %d", recorder.Code)
	}
	var count int

	err := db.QueryRow(
		"SELECT COUNT(*) FROM notes WHERE id = $1",
		noteID,
	).Scan(&count)
	if err != nil {
		t.Fatalf("failed to query database: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 notes after deletion, got %d", count)
	}
}
