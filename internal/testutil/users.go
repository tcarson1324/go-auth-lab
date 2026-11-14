package testutil

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"testing"
)

func CreateTestUser(t *testing.T, db *sql.DB) int64 {
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
