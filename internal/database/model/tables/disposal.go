package model

import (
	"database/sql"
	"time"
)

// disposals テーブル
type AssetsDisposal struct {
	ID           int64          `db:"id"`            // 主キー (AUTO_INCREMENT)
	AssetID      int64          `db:"asset_id"`      // 対象資産 ID（NOT NULL）
	Quantity     int            `db:"quantity"`      // 廃棄数量（NOT NULL）
	DisposalDate time.Time      `db:"disposal_date"` // 廃棄日（NOT NULL / DATE）
	Reason       sql.NullString `db:"reason"`        // 廃棄理由（NULL 可 / TEXT）
	ProcessedBy  sql.NullString `db:"processed_by"`  // 処理担当者（NULL 可 / VARCHAR(100)）
	IsIndividual bool           `db:"-"` // 個別管理かどうか（NOT NULL / BOOLEAN）データベースには保存しないので、`-`タグを使用
}