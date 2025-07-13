package lends

import (
	"database/sql"
	model "equipmentManager/internal/database/model/tables"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// 全体管理のみに適用
// 個別管理はassets.ownerが借りている人として扱うため
func RegisterLend(db *sql.DB, asset model.AssetsLend) (bool, int64, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Println("資産貸出：初期化エラー:", err)
		return false, -1, err
	}
	defer tx.Rollback() // 関数終了時にコミットされていなければ自動でロールバック

	// 在庫確認
	var stock int
	query := `SELECT quantity FROM assets WHERE id = ?`
	err = tx.QueryRow(query, asset.AssetID).Scan(&stock)
	if err != nil {
		log.Println("資産貸出：在庫数量取得エラー:", err)
		return false, -1, err
	}

	if stock < asset.Quantity {
		log.Printf("資産貸出：在庫不足（在庫 %d < 要求 %d）\n", stock, asset.Quantity)
		return false, -1, fmt.Errorf("在庫不足：在庫 %d < 要求 %d", stock, asset.Quantity)
	}

	// 1. asset_lends テーブルに貸出記録をINSERT
	query = `
		INSERT INTO asset_lends 
			(asset_id, borrower, quantity, lend_date, expected_return_date, actual_return_date, notes) 
		VALUES (?, ?, ?, ?, ?, ?, ?);`

	res, err := tx.Exec(query, asset.AssetID, asset.Borrower, asset.Quantity, asset.LendDate, asset.ExpectedReturnDate, asset.ActualReturnDate, asset.Notes)
	if err != nil {
		log.Println("資産貸出：登録エラー:", err)
		return false, -1, err
	}

	// 2. assets テーブルの location と owner を borrower でUPDATE
	updateQuery := `UPDATE assets SET location = ?, owner = ? WHERE id = ?`
	_, err = tx.Exec(updateQuery, asset.Borrower, asset.Borrower, asset.AssetID)
	if err != nil {
		log.Println("資産貸出：資産場所更新エラー:", err)
		return false, -1, err
	}

	rendId, err := res.LastInsertId()
	if err != nil {
		log.Println("資産貸出：ID取得エラー:", err)
		return false, -1, err
	}

	log.Println("資産貸出：登録成功")
	// 全ての処理が成功したのでコミット
	return true, rendId, tx.Commit()
}
