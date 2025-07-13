package lends

import (
	"database/sql"
	model "equipmentManager/internal/database/model/tables"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
	"errors"
)

func UpdateLend(db *sql.DB, asset model.AssetsLend) (bool, error) {

	tx, err := db.Begin()
	if err != nil {
		log.Println("資産貸出：更新初期化エラー:", err)
		return false, err
	}

	query := `
		UPDATE asset_lends
		SET 
			borrower = ?,
			quantity = ?,
			lend_date = ?,
			expected_return_date = ?,
			actual_return_date = ?,
			notes = ?
		WHERE id = ?;`

	res, err := tx.Exec(query, asset.Borrower, asset.Quantity, asset.LendDate, asset.ExpectedReturnDate, asset.ActualReturnDate, asset.Notes, asset.ID)
	if err != nil {
		log.Println("資産貸出：更新エラー:", err)
		tx.Rollback()
		return false, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("資産貸出：更新行数取得エラー:", err)
		tx.Rollback()
		return false, err
	}
	log.Printf("資産貸出：更新行数: %d\n", rowsAffected)

	log.Println("資産貸出：更新成功")
	return true, tx.Commit()
}

func UpdateReturnDateForAssetlist(tx *sql.Tx, lendID int64, returnDate time.Time, notes sql.NullString) (bool, error) {
	log.Printf("資産貸出：返却登録: %v", returnDate)
	query := `UPDATE asset_lends SET actual_return_date = ?, notes = ? WHERE id = ?`
	res, err := tx.Exec(query, returnDate, notes, lendID)
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

	return rows > 0, nil
}
