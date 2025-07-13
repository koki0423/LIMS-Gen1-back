package model

import (
	"database/sql"
	"time"
)

// lends テーブル（貸出履歴）
type AssetsLend struct {
	ID                 int64          `db:"id"`                   // 主キー (AUTO_INCREMENT)
	AssetID            int64          `db:"asset_id"`             // 貸出対象資産 ID（NOT NULL）
	Borrower           string         `db:"borrower"`             // 借用者（NOT NULL, VARCHAR(100)）
	Quantity           int            `db:"quantity"`             // 貸出数（NOT NULL）
	LendDate           time.Time      `db:"lend_date"`            // 貸出日（NOT NULL / DATE）
	ExpectedReturnDate sql.NullTime   `db:"expected_return_date"` // 返却予定日（NULL 可 / DATE）
	ActualReturnDate   sql.NullTime   `db:"actual_return_date"`   // 実際の返却日（NULL 可 / DATE）
	Notes              sql.NullString `db:"notes"`                // 備考（NULL 可 / TEXT）
}

// LendingDetail は、貸出情報に備品マスター情報を加えたレスポンス用の構造体
type LendingDetail struct {
	ID                 int64          `db:"id"`
	Borrower           string         `db:"borrower"`
	Quantity           int            `db:"quantity"`
	LendDate           time.Time      `db:"lend_date"`
	ExpectedReturnDate sql.NullTime   `db:"expected_return_date"`
	Notes              sql.NullString `db:"notes"`
	Name               string         `db:"name"`
	Manufacturer       sql.NullString `db:"manufacturer"`
	ModelNumber        sql.NullString `db:"model_number"`
}
