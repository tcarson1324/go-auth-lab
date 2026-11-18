package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-auth-lab/internal/testutil"
	"go-auth-lab/middleware"
)

func TestEditNoteWrongUser(t *testing.T) {
	db := testutil.SetupTestDB(t)
	userAID := testutil.CreateTestUser(t, db)
	userBID := testutil.CreateTestUser(t, db)
	noteID := testutil.CreateTestNote(t, db, userAID)
	// Request user A's note.
	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/notes/edit?id=%d", noteID),
		nil,
	)

	// Make the request as user B.
	ctx := middleware.ContextWithUserID(req.Context(), userBID)
	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()
	HandlerEditNote(recorder, req, db)
	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected error code 404, got %v", recorder.Code)
	}
}

func TestEditNoteSuccess(t *testing.T) {
	t.Chdir("..")
	db := testutil.SetupTestDB(t)
	userID := testutil.CreateTestUser(t, db)
	noteID := testutil.CreateTestNote(t, db, userID)
	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/notes/edit?id=%d", noteID),
		nil,
	)
	ctx := middleware.ContextWithUserID(req.Context(), userID)
	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()
	HandlerEditNote(recorder, req, db)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	if !strings.Contains(recorder.Body.String(), "This is a test note") {
		t.Error("expected response body to contain original note content")
	}
}

func TestHandlerUpdateNoteSuccess(t *testing.T) {
	db := testutil.SetupTestDB(t)
	userID := testutil.CreateTestUser(t, db)
	noteID := testutil.CreateTestNote(t, db, userID)
	content := "This note has been updated"

	req := httptest.NewRequest(
		http.MethodPost,
		"/notes/edit",
		strings.NewReader(fmt.Sprintf("note_id=%d&content=%s", noteID, content)),
	)
	ctx := middleware.ContextWithUserID(req.Context(), userID)

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recorder := httptest.NewRecorder()
	HandlerUpdateNote(recorder, req, db)
	if recorder.Code != http.StatusSeeOther {
		t.Fatalf("expected status 303, got %d: %s", recorder.Code, recorder.Body.String())
	}
	location := recorder.Header().Get("Location")
	if location != "/profile" {
		t.Fatalf("expected redirect to '/profile', got %s", location)
	}
	var updatedContent string
	note := db.QueryRow(`
	SELECT content
	FROM notes
	WHERE id = $1
	AND user_id = $2
	`, noteID, userID)
	err := note.Scan(&updatedContent)
	if err != nil {
		t.Fatalf("database error: %v", err)
	}
	if updatedContent != content {
		t.Errorf("expected note to be updated, got %s", updatedContent)
	}

}

func TestHandlerUpdateNoteWrongUser(t *testing.T) {
	db := testutil.SetupTestDB(t)
	userAID := testutil.CreateTestUser(t, db)
	userBID := testutil.CreateTestUser(t, db)
	noteID := testutil.CreateTestNote(t, db, userAID)
	content := "This note has been updated"

	req := httptest.NewRequest(
		http.MethodPost,
		"/notes/edit",
		strings.NewReader(fmt.Sprintf("note_id=%d&content=%s", noteID, content)),
	)
	ctx := middleware.ContextWithUserID(req.Context(), userBID)

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recorder := httptest.NewRecorder()
	HandlerUpdateNote(recorder, req, db)
	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected status code 404, got %v", recorder.Code)
	}
	var updatedContent string
	note := db.QueryRow(`
	SELECT content
	FROM notes
	WHERE id = $1
	AND user_id = $2
	`, noteID, userAID)
	err := note.Scan(&updatedContent)
	if err != nil {
		t.Fatalf("database error: %v", err)
	}
	if updatedContent != "This is a test note" {
		t.Errorf("expected note to contain original content, got %s", updatedContent)
	}

}
func TestUpdateNoteEmptyContent(t *testing.T) {
	db := testutil.SetupTestDB(t)
	userID := testutil.CreateTestUser(t, db)
	noteID := testutil.CreateTestNote(t, db, userID)
	content := ""

	req := httptest.NewRequest(
		http.MethodPost,
		"/notes/edit",
		strings.NewReader(fmt.Sprintf("note_id=%d&content=%s", noteID, content)),
	)
	ctx := middleware.ContextWithUserID(req.Context(), userID)

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recorder := httptest.NewRecorder()
	HandlerUpdateNote(recorder, req, db)
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("expected error code 400, got %v", recorder.Code)
	}
	var updatedContent string
	note := db.QueryRow(`
	SELECT content
	FROM notes
	WHERE id = $1
	AND user_id = $2
	`, noteID, userID)
	err := note.Scan(&updatedContent)
	if err != nil {
		t.Fatalf("database error: %v", err)
	}
	if updatedContent != "This is a test note" {
		t.Errorf("expected original note content, got %v", updatedContent)
	}
}