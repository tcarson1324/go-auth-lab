package notes

import (
	"database/sql"
	"errors"
	"time"
)

type Note struct {
	ID        int64
	UserID    int64
	Content   string
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

func DeleteNoteByID(db *sql.DB, noteID int64, userID int64) error {
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

func GetNoteByID(db *sql.DB, noteID, userID int64) (Note, error) {
	query := `
	SELECT id, user_id, content FROM notes
	WHERE id = $1 
	AND user_id = $2
	`
	note := Note{}
	data := db.QueryRow(query, noteID, userID)
	err := data.Scan(&note.ID, &note.UserID, &note.Content)
	if err != nil {
		if err == sql.ErrNoRows {
			return Note{}, ErrNoteNotFound
		}
		return Note{}, err
	}
	return note, nil
}

func UpdateNote(db *sql.DB, noteID, userID int64, content string) error {
	query := `
	UPDATE notes
	SET content = $1
	WHERE id = $2
	AND user_id = $3
	`
	result, err := db.Exec(query, content, noteID, userID)
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