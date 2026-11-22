package middleware

import (
	"go-auth-lab/internal/testutil"

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRequireAuthExpiredSession(t *testing.T) {
	funcExecuted := false
	db := testutil.SetupTestDB(t)
	userID := testutil.CreateTestUser(t, db)
	token, tokenHash := testutil.CreateTestSession(t, db, userID)
	_, err := db.Exec(`
	UPDATE sessions
	SET expires_at = $1
	WHERE token_hash = $2
	`, time.Now().Add(-1*time.Hour), tokenHash)
	if err != nil {
		t.Fatalf("database error: %v", err)
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		funcExecuted = true
	})
	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: token,
	})

	recorder := httptest.NewRecorder()

	protectedHandler := RequireAuth(db, handler)
	protectedHandler.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("expected error code 401, got %d", recorder.Code)
	}
	if funcExecuted == true {
		t.Error("handler not expected to execute")
	}

}
func TestRequireCSRFMissingToken(t *testing.T) {
	funcExecuted := false
	db := testutil.SetupTestDB(t)
	userID := testutil.CreateTestUser(t, db)
	token, _ := testutil.CreateTestSession(t, db, userID)
	req := httptest.NewRequest(
		http.MethodPost,
		"/notes/delete",
		strings.NewReader("note_id=1"),
	)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: token,
	})
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler := RequireCSRF(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		funcExecuted = true
	}))
	recorder := httptest.NewRecorder()
	protectedHandler := RequireAuth(db, handler)
	protectedHandler.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusForbidden {
		t.Errorf("expected error code 403, got %v", recorder.Code)
	}
	if funcExecuted {
		t.Errorf("handler not expected to execute")
	}

}

func TestRequireCSRFWrongToken(t *testing.T) {
	wrongToken := "this_is_an_incorrect_token_value"
	funcExecuted := false
	db := testutil.SetupTestDB(t)
	userID := testutil.CreateTestUser(t, db)

	token, _ := testutil.CreateTestSession(t, db, userID)
	req := httptest.NewRequest(
		http.MethodPost,
		"/notes/delete",
		strings.NewReader("note_id=1&csrf_token="+wrongToken),
	)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: token,
	})
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler := RequireCSRF(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		funcExecuted = true
	}))
	recorder := httptest.NewRecorder()
	protectedHandler := RequireAuth(db, handler)
	protectedHandler.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusForbidden {
		t.Errorf("expected error code 403, got %v", recorder.Code)
	}
	if funcExecuted {
		t.Errorf("handler not expected to execute")
	}

}

func TestRequireCSRF(t *testing.T) {
	var csrfToken string
	funcExecuted := false
	db := testutil.SetupTestDB(t)
	userID := testutil.CreateTestUser(t, db)
	token, tokenHash := testutil.CreateTestSession(t, db, userID)
	query := `
	SELECT csrf_token
	FROM sessions
	WHERE token_hash = $1
	`
	result := db.QueryRow(query, tokenHash)
	err := result.Scan(&csrfToken)
	if err != nil {
		t.Fatal("internal database error")
	}
	req := httptest.NewRequest(
		http.MethodPost,
		"/notes/delete",
		strings.NewReader("note_id=1&csrf_token="+csrfToken),
	)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: token,
	})
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler := RequireCSRF(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		funcExecuted = true
	}))
	recorder := httptest.NewRecorder()
	protectedHandler := RequireAuth(db, handler)
	protectedHandler.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %v", recorder.Code)
	}
	if !funcExecuted {
		t.Error("expected handler function to execute")
	}
}
