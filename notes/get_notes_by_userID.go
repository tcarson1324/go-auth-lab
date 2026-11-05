package notes

import (
	"database/sql"


)



func GetNotesByUserID(db *sql.DB, userID int64) ([]Note, error) {
	query := `
	SELECT id, content
	FROM notes
	WHERE user_id = $1
	`
	rows, err := db.Query(query, userID)
	if err != nil {
		return []Note{}, err
	}
	defer rows.Close()
	var notes []Note
	for rows.Next() {
		var note Note
		err := rows.Scan(
			&note.ID,
			&note.Content,
		)
		if err != nil {
			return []Note{}, err
		}
		notes = append(notes, note)
	}
	if err = rows.Err(); err != nil {
		return notes, err
	}
	return notes, nil
}
