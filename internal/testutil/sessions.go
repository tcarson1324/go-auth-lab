package testutil

import (
	"database/sql"
	"go-auth-lab/sessions"
	"testing"
)

func CreateTestSession(t *testing.T, db *sql.DB, userID int64) (token, tokenHash string) {
	token, err := sessions.GenerateToken()
	if err != nil {
		t.Fatalf("could not generate token: %v", err)
	}
	tokenHash = sessions.HashToken(token)
	_, err = sessions.CreateSession(db, userID, tokenHash)
	if err != nil {
		t.Fatalf("database error: %v", err)
	}
	t.Cleanup(func() {
		_, err := db.Exec(`
        DELETE FROM sessions
        WHERE token_hash = $1
    `, tokenHash)
		if err != nil {
			t.Errorf("session cleanup failed: %v", err)
		}
	})
	return token, tokenHash
}
