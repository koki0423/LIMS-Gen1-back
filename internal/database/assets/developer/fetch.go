package developer

import (
	"database/sql"

	model "equipmentManager/internal/database/model/tables"
	_ "github.com/go-sql-driver/mysql"
)

// developer/repository.go
func GetDeveloperByStudentNumber(db *sql.DB, studentNumber string) (*model.Developer, error) {
	query := `
		SELECT id, student_number, password_hash
		FROM developers
		WHERE student_number = ?
	`

	var dev model.Developer
	err := db.QueryRow(query, studentNumber).Scan(
		&dev.ID,
		&dev.StudentNumber,
		&dev.PasswordHash,
	)
	if err == sql.ErrNoRows {
		return nil, nil // 開発者でない
	}
	if err != nil {
		return nil, err
	}
	return &dev, nil
}
