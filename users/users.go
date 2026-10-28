package users

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"go-auth-lab/passwords"

	"github.com/jackc/pgx/v5/pgconn"
)

var ErrUsernameTaken = errors.New("username already exists")
var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID           int64
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}

func AddUser(username, password string, db *sql.DB) error {

	hashedPassword, err := passwords.Hash(password)
	if err != nil {
		return err
	}
	query := `
	INSERT INTO users(username, password_hash)
	VALUES($1, $2)
	`
	_, err = db.Exec(query, username, hashedPassword)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("could not register user: %w", ErrUsernameTaken)
			}
		}
		return err
	}
	return nil
}
func GetUserByName(name string, db *sql.DB) (User, error) {

	query := `
	SELECT id, username, password_hash, created_at
	FROM users
	WHERE username = $1
	`
	currentUser := User{}
	row := db.QueryRow(query, name)
	err := row.Scan(
		&currentUser.ID,
		&currentUser.Username,
		&currentUser.PasswordHash,
		&currentUser.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return currentUser, ErrUserNotFound
		}
		return currentUser, err
	}
	return currentUser, nil
}

