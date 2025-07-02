package model

import (
	"database/sql"
)

type AllDataOfAsset struct {
	Master AssetsMaster
	Asset  Asset
}

// assets_masters テーブル
type AssetsMaster struct {
	ID                   int64          `db:"id"`                     // 主キー (AUTO_INCREMENT)
	Name                 string         `db:"name"`                   // 備品の正式名称
	ManagementCategoryID int64          `db:"management_category_id"` // 管理区分 ID（NOT NULL）
	GenreID              sql.NullInt64  `db:"genre_id"`               // ジャンル ID（NULL 可）
	Manufacturer         sql.NullString `db:"manufacturer"`           // メーカー名（NULL 可）
	ModelNumber          sql.NullString `db:"model_number"`           // 型番（NULL 可）
}

// assets テーブル（個別資産）
type Asset struct {
	ID            int64          `db:"id"`              // 主キー (AUTO_INCREMENT)
	ItemMasterID  int64          `db:"item_master_id"`  // どの種類の備品か（FK, NOT NULL）
	Quantity      int            `db:"quantity"`        // 個数（NOT NULL, デフォルト 1）
	SerialNumber  sql.NullString `db:"serial_number"`   // 個別管理用識別子（NULL 可, UNIQUE）
	StatusID      int64          `db:"status_id"`       // 状態 ID（FK, NOT NULL）
	PurchaseDate  sql.NullTime   `db:"purchase_date"`   // 購入日（NULL 可 / DATE）
	Owner         sql.NullString `db:"owner"`           // 所有者（NULL 可 / VARCHAR(100)）
	Location      sql.NullString `db:"location"`        // 保管場所（NULL 可 / VARCHAR(255)）
	LastCheckDate sql.NullTime   `db:"last_check_date"` // 最終確認日（NULL 可 / DATE）
	LastChecker   sql.NullString `db:"last_checker"`    // 最終確認者（NULL 可 / VARCHAR(100)）
	Notes         sql.NullString `db:"notes"`           // 備考欄（NULL 可 / TEXT）
}
