package sessions

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"net/http"
	"time"
)

var ErrSessionNotFound = errors.New("session not found")

type Session struct {
	TokenHash string
	UserID    int64
	ExpiresAt time.Time
}

func SetCookie(w http.ResponseWriter, token string, expiresAt time.Time) {
	cookie := &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expiresAt,
	}
	http.SetCookie(w, cookie)
}
func GenerateToken() (string, error) {
	token := make([]byte, 32)

	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	hashString := hex.EncodeToString(hash[:])
	return hashString
}

func CreateSession(db *sql.DB, userID int64, tokenHash string) (time.Time, error) {
	expiresAt := time.Now().Add(24 * time.Hour)

	query := `
        INSERT INTO sessions (token_hash, user_id, expires_at)
        VALUES ($1, $2, $3)
    `

	_, err := db.Exec(query, tokenHash, userID, expiresAt)
	if err != nil {
		return time.Time{}, err
	}

	return expiresAt, nil
}
func GetSession(db *sql.DB, token string) (Session, error) {
	currentSession := Session{}
	tokenHash := HashToken(token)
	query := `
	SELECT token_hash, user_id, expires_at
	FROM sessions
	WHERE token_hash = $1
	AND expires_at > NOW();
	`
	sessionData := db.QueryRow(query, tokenHash)
	err := sessionData.Scan(
		&currentSession.TokenHash,
		&currentSession.UserID,
		&currentSession.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Session{}, ErrSessionNotFound
		}
		return Session{}, err
	}

	return currentSession, nil
}
