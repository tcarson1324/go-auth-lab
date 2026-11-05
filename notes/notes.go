package notes

import (
	"time"
	"database/sql"

)

type Note struct {
	ID int64
	UserId int64
	Content string
	CreatedAt time.Time
}

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