package lends

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	model "equipmentManager/internal/database/model/tables"
)

// 資産貸出の更新
func UpdateLend(db *sql.DB, asset model.AssetsLend) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("資産貸出：更新初期化エラー:", err)
		return false, err
	}
	defer func() { _ = tx.Rollback() }()

	const query = `
UPDATE asset_lends
SET
  borrower = ?,
  quantity = ?,
  lend_date = ?,
  expected_return_date = ?,
  actual_return_date = ?,
  notes = ?
WHERE id = ?;
`
	res, err := tx.ExecContext(ctx, query,
		asset.Borrower,
		asset.Quantity,
		asset.LendDate,
		asset.ExpectedReturnDate,
		asset.ActualReturnDate,
		asset.Notes,
		asset.ID,
	)
	if err != nil {
		log.Println("資産貸出：更新エラー:", err)
		return false, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("資産貸出：更新行数取得エラー:", err)
		return false, err
	}
	if rowsAffected == 0 {
		log.Println("資産貸出：更新対象なし")
		// 存在しないIDなど
		return false, sql.ErrNoRows
	}

	if err := tx.Commit(); err != nil {
		log.Println("資産貸出：コミットエラー:", err)
		return false, err
	}
	log.Printf("資産貸出：更新成功（%d件）\n", rowsAffected)
	return true, nil
}

// 返却登録（Tx外でまとめ処理している想定のユーティリティ）
func UpdateReturnDateForAssetlist(tx *sql.Tx, lendID int64, returnDate time.Time, notes sql.NullString) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	const query = `UPDATE asset_lends SET actual_return_date = ?, notes = ? WHERE id = ?`
	res, err := tx.ExecContext(ctx, query, returnDate, notes, lendID)
	if err != nil {
		log.Println("資産貸出：返却日更新エラー:", err)
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("資産貸出：返却日更新行数取得エラー:", err)
		return false, err
	}
	if rows == 0 {
		log.Println("返却日更新: 該当レコードなし")
		return false, errors.New("返却対象の貸出データが存在しません")
	}
	return true, nil
}
