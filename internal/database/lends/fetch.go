package lends

import (
	"database/sql"
	model "equipmentManager/internal/database/model/tables"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// 全貸出情報を取得
func FetchLendsAll(db *sql.DB) ([]model.AssetsLend, error) {
	query := `SELECT id, asset_id, borrower, quantity, lend_date, expected_return_date, actual_return_date, notes FROM asset_lends ORDER BY lend_date DESC, id DESC`

	rows, err := db.Query(query)
	if err != nil {
		log.Println("貸出一覧取得：クエリエラー:", err)
		return nil, err
	}
	defer rows.Close()

	var lends []model.AssetsLend
	for rows.Next() {
		var lend model.AssetsLend
		err := rows.Scan(
			&lend.ID,
			&lend.AssetID,
			&lend.Borrower,
			&lend.Quantity,
			&lend.LendDate,
			&lend.ExpectedReturnDate,
			&lend.ActualReturnDate,
			&lend.Notes,
		)
		if err != nil {
			log.Println("貸出一覧取得：スキャンエラー:", err)
			return nil, err
		}
		lends = append(lends, lend)
	}
	return lends, nil
}

// 特定資産IDの貸出情報を取得
func FetchLendsByAssetID(db *sql.DB, assetID int64) ([]model.AssetsLend, error) {
	query := `SELECT id, asset_id, borrower, quantity, lend_date, expected_return_date, actual_return_date, notes FROM asset_lends WHERE asset_id = ? ORDER BY lend_date DESC, id DESC`

	rows, err := db.Query(query, assetID)
	if err != nil {
		log.Println("貸出情報取得：クエリエラー:", err)
		return nil, err
	}
	defer rows.Close()

	var lends []model.AssetsLend
	for rows.Next() {
		var lend model.AssetsLend
		err := rows.Scan(
			&lend.ID,
			&lend.AssetID,
			&lend.Borrower,
			&lend.Quantity,
			&lend.LendDate,
			&lend.ExpectedReturnDate,
			&lend.ActualReturnDate,
			&lend.Notes,
		)
		if err != nil {
			return nil, err
		}
		lends = append(lends, lend)
	}

	return lends, nil
}
