package lends

import (
	"database/sql"
	"time"
)

// DBモデル（テーブル構造に極力1:1）
type Lend struct {
	LendID         uint64
	LendULID       string
	AssetMasterID  uint64
	ManagementNumber string
	Quantity       uint
	BorrowerID     string
	DueOn          sql.NullString // DATEを文字列で扱う（"2006-01-02"）
	LentByID       sql.NullString
	LentAt         time.Time
	Note           sql.NullString
}

type Return struct {
	ReturnID       uint64
	ReturnULID     string
	LendID         uint64
	Quantity       uint
	ProcessedByID  sql.NullString
	ReturnedAt     time.Time
	Note           sql.NullString
}
