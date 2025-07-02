package model

type Developer struct{
	ID    int64          `db:"id"`    // 主キー (AUTO_INCREMENT)
	StudentNumber  string         `db:"student_number"`  // 学籍番号（NOT NULL, UNIQUE, VARCHAR(20)）
	PasswordHash string         `db:"password_hash"`    // bcryptハッシュ（NOT NULL, VARCHAR(100)）
}