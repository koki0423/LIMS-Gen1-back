package lends

import (
	"context"
	"database/sql"
	"log"
	"time"
	"fmt"

	model "equipmentManager/internal/database/model/tables"
)

// 全貸出情報を取得
func FetchLendsAll(db *sql.DB) ([]model.AssetsLend, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	const query = `
SELECT id, asset_id, borrower, quantity, lend_date, expected_return_date, actual_return_date, notes
FROM asset_lends
ORDER BY lend_date DESC, id DESC;
`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Println("貸出一覧取得：クエリエラー:", err)
		return nil, err
	}
	defer rows.Close()

	lends := make([]model.AssetsLend, 0, 128)
	for rows.Next() {
		var lend model.AssetsLend
		if err := rows.Scan(
			&lend.ID,
			&lend.AssetID,
			&lend.Borrower,
			&lend.Quantity,
			&lend.LendDate,
			&lend.ExpectedReturnDate,
			&lend.ActualReturnDate,
			&lend.Notes,
		); err != nil {
			log.Println("貸出一覧取得：スキャンエラー:", err)
			return nil, err
		}
		lends = append(lends, lend)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lends, nil
}

// 特定資産IDの貸出情報を取得
func FetchLendsByAssetID(db *sql.DB, assetID int64) ([]model.AssetsLend, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	const query = `
SELECT id, asset_id, borrower, quantity, lend_date, expected_return_date, actual_return_date, notes
FROM asset_lends
WHERE asset_id = ?
ORDER BY lend_date DESC, id DESC;
`
	rows, err := db.QueryContext(ctx, query, assetID)
	if err != nil {
		log.Println("貸出情報取得：クエリエラー:", err)
		return nil, err
	}
	defer rows.Close()

	lends := make([]model.AssetsLend, 0, 64)
	for rows.Next() {
		var lend model.AssetsLend
		if err := rows.Scan(
			&lend.ID,
			&lend.AssetID,
			&lend.Borrower,
			&lend.Quantity,
			&lend.LendDate,
			&lend.ExpectedReturnDate,
			&lend.ActualReturnDate,
			&lend.Notes,
		); err != nil {
			return nil, err
		}
		lends = append(lends, lend)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lends, nil
}

// 全貸出情報を資産名とともに取得（未返却のみ）
func FetchAllLendingDetails(db *sql.DB) ([]model.LendingDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	const query = `
SELECT
  al.id, al.borrower, al.quantity, al.lend_date, al.expected_return_date, al.notes,
  am.name, am.manufacturer, am.model_number
FROM asset_lends AS al
JOIN assets AS a ON al.asset_id = a.id
JOIN assets_masters AS am ON a.asset_master_id = am.id
WHERE al.actual_return_date IS NULL
ORDER BY al.lend_date DESC, al.id DESC;
`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Println("貸出情報取得エラー:", err)
		return nil, err
	}
	defer rows.Close()

	details := make([]model.LendingDetail, 0, 128)
	for rows.Next() {
		var d model.LendingDetail
		if err := rows.Scan(
			&d.ID,
			&d.Borrower,
			&d.Quantity,
			&d.LendDate,
			&d.ExpectedReturnDate,
			&d.Notes,
			&d.Name,
			&d.Manufacturer,
			&d.ModelNumber,
		); err != nil {
			log.Println("スキャンエラー:", err)
			return nil, err
		}
		details = append(details, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return details, nil
}

// Tx内で lend_id → asset_id を引く（関数内で短いタイムアウトを持つ）
func GetAssetIDByLendID(tx *sql.Tx, lendID int64) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	const query = `SELECT asset_id FROM asset_lends WHERE id = ?;`
	var assetID int64
	if err := tx.QueryRowContext(ctx, query, lendID).Scan(&assetID); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("lend_id %d に該当する貸出記録が見つかりません", lendID)
		}
		return 0, err
	}
	return assetID, nil
}
