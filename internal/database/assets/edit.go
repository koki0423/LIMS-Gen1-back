package assets

import (
	"database/sql"

	model "equipmentManager/internal/database/model/tables"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"fmt"
)

func UpdateAsset(db *sql.DB, updated model.Asset, id int64) (bool, error) {
	updated.ID = id
	// log.Printf("	資産情報更新：ID=%d, 数量=%d, シリアル番号=%s, ステータスID=%d, 購入日=%s, 所有者=%s, 配置場所=%s, 最終チェック日=%s, 最終チェック者=%s, メモ=%s",
	// 	updated.ID, updated.Quantity, updated.SerialNumber, updated.StatusID, updated.PurchaseDate, updated.Owner, updated.Location, updated.LastCheckDate, updated.LastChecker, updated.Notes)

	tx, err := db.Begin()
	if err != nil {
		return false, err
	}

	query := `UPDATE assets
		SET quantity = ?, serial_number = ?, status_id = ?, purchase_date = ?, 
		    owner = ?, location = ?, last_check_date = ?, last_checker = ?, notes = ?
		WHERE id = ?`

	_, err = tx.Exec(query,
		updated.Quantity,
		updated.SerialNumber,
		updated.StatusID,
		updated.PurchaseDate,
		updated.Owner,
		updated.Location,
		updated.LastCheckDate,
		updated.LastChecker,
		updated.Notes,
		updated.ID,
	)

	if err != nil {
		log.Println("資産情報更新：エラー:", err)
		tx.Rollback()
		return false, err
	}

	if err = tx.Commit(); err != nil {
		log.Println("資産情報更新：コミットエラー:", err)
		return false, err
	}

	log.Println("資産情報更新：成功")
	return true, nil
}

func ResetAssetStatus(tx *sql.Tx, assetID int64) error {
	query := `UPDATE assets SET owner = NULL, location = default_location WHERE id = ?`
	res, err := tx.Exec(query, assetID)
	if err != nil {
		log.Println("資産状態のリセットエラー:", err)
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("asset_id %d に該当する資産が見つかりません", assetID)
	}
	return nil
}