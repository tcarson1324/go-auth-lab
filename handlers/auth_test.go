package handlers

import (
	"go-auth-lab/internal/testutil"
	"go-auth-lab/middleware"
	"go-auth-lab/sessions"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireAuthNoCookie(t *testing.T) {
	db := testutil.SetupTestDB(t)
	funcExecuted := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		funcExecuted = true
	})
	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	protectedHandler := middleware.RequireAuth(db, handler)
	recorder := httptest.NewRecorder()
	protectedHandler.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("expected code 401, got %d", recorder.Code)
	}
	if funcExecuted {
		t.Error("protected handler executed without authentication")
	}
}

func TestRequireAuthSuccess(t *testing.T) {
	db := testutil.SetupTestDB(t)
	funcExecuted := false
	userID := testutil.CreateTestUser(t, db)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		funcExecuted = true

		authenticatedID, ok := middleware.UserIDFromContext(r.Context())

		if !ok {
			t.Error("user ID missing from context")
			return
		}

		if authenticatedID != userID {
			t.Errorf("expected user ID %d, got %d", userID, authenticatedID)
		}
	})

	recorder := httptest.NewRecorder()
	token, err := sessions.GenerateToken()
	if err != nil {
		t.Fatalf("could not generate token: %v", err)
	}
	hashedToken := sessions.HashToken(token)
	_, err = sessions.CreateSession(db, userID, hashedToken)
	
	if err != nil {
		t.Fatalf("could not create session: %v", err)
	}
	t.Cleanup(func() {
    _, err := db.Exec(
        "DELETE FROM sessions WHERE token_hash = $1",
        hashedToken,
    )
    if err != nil {
        t.Errorf("could not clean up session: %v", err)
    }
})
	req := httptest.NewRequest(http.MethodGet, "/profile", nil)

	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: token,
	})
	protectedHandler := middleware.RequireAuth(db, handler)
	protectedHandler.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Errorf("expected code 200, got %d", recorder.Code)
	}
	if !funcExecuted {
		t.Error("protected handler did not execute")
	}

}

func TestRequireAuthInvalidSession(t *testing.T) {
		db := testutil.SetupTestDB(t)
	funcExecuted := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		funcExecuted = true
	})

	recorder := httptest.NewRecorder()
	
	req := httptest.NewRequest(http.MethodGet, "/profile", nil)

	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: "this token does not exist",
	})
	protectedHandler := middleware.RequireAuth(db, handler)
	protectedHandler.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("expected code 401, got %d", recorder.Code)
	}
	if funcExecuted {
		t.Error("protected handler should not execute")
	}

	
}

