package developer

import (
	"database/sql"
	"log"
	model "equipmentManager/internal/database/model/tables"
)

func RegisterDeveloper(db *sql.DB, user model.Developer) (int64, error) {
	query := `
		INSERT INTO developers (student_number, password_hash)
		VALUES (?, ?)
	`
	result, err := db.Exec(query, user.StudentNumber, user.PasswordHash)
	if err != nil {
		log.Println("開発者登録：DBエラー", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("開発者登録：ID取得エラー", err)
		return 0, err
	}

	log.Println("開発者登録：成功、ID=", id)
	return id, nil
}