package lends

import (
	"database/sql"
	model "equipmentManager/internal/database/model/tables"
	"fmt"
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

// 全貸出情報を資産名とともに取得
func FetchAllLendingDetails(db *sql.DB) ([]model.LendingDetail, error) {
	query := `
        SELECT
            al.id, al.borrower, al.quantity, al.lend_date, al.expected_return_date, al.notes,
            am.name, am.manufacturer, am.model_number
        FROM
            asset_lends AS al
        JOIN
            assets AS a ON al.asset_id = a.id
        JOIN
            assets_masters AS am ON a.asset_master_id = am.id
        WHERE
            al.actual_return_date IS NULL` // 返却されていないものに絞る

	rows, err := db.Query(query)
	if err != nil {
		log.Println("貸出情報取得エラー:", err)
		return nil, err
	}
	defer rows.Close()

	var lendingDetails []model.LendingDetail

	for rows.Next() {
		var detail model.LendingDetail
		err := rows.Scan(
			&detail.ID,
			&detail.Borrower,
			&detail.Quantity,
			&detail.LendDate,
			&detail.ExpectedReturnDate,
			&detail.Notes,
			&detail.Name,
			&detail.Manufacturer,
			&detail.ModelNumber,
		)
		if err != nil {
			log.Println("スキャンエラー:", err)
			return nil, err
		}
		lendingDetails = append(lendingDetails, detail)
	}

	return lendingDetails, nil
}

func GetAssetIDByLendID(tx *sql.Tx, lendID int64) (int64, error) {
	var assetID int64
	query := `SELECT asset_id FROM asset_lends WHERE id = ?`
	err := tx.QueryRow(query, lendID).Scan(&assetID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("lend_id %d に該当する貸出記録が見つかりません", lendID)
		}
		return 0, err
	}
	return assetID, nil
}
