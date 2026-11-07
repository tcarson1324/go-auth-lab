package notes

import (
	"time"
	"database/sql"
	"errors"
)

type Note struct {
	ID int64
	UserId int64
	Content string
	CreatedAt time.Time
}

var ErrNoteNotFound = errors.New("note not found")

func CreateNote(db *sql.DB, userID int64, content string) error {
	query := `
	INSERT INTO notes (user_id, content)
	VALUES ($1, $2)
	`
	_, err := db.Exec(query, userID, content)
	if err != nil {
		return err
	}
	
	return nil
}

func DeleteNoteByID(db *sql.DB, noteID int64, userID int64) error{
	query := `
	DELETE FROM notes
	WHERE id = $1 AND user_id = $2
	`
	result, err := db.Exec(query, noteID, userID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNoteNotFound
	}
	
	return nil
}