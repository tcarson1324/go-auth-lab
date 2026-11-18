package middleware

import (
	"go-auth-lab/internal/testutil"
	
	"net/http"
	"testing"
	"time"
	"net/http/httptest"
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
	`, time.Now().Add(-1 * time.Hour), tokenHash)
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