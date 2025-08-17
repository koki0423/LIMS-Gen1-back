package assets

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	model "equipmentManager/internal/database/model/tables"
)

func UpdateAsset(db *sql.DB, updated model.Asset, id int64) (bool, error) {
	updated.ID = id

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer func() { _ = tx.Rollback() }()

	const query = `
	UPDATE assets
	SET quantity = ?, serial_number = ?, status_id = ?, purchase_date = ?,
		owner = ?, location = ?, last_check_date = ?, last_checker = ?, notes = ?
	WHERE id = ?`

	if _, err := tx.ExecContext(ctx, query,
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
	); err != nil {
		log.Println("資産情報更新：エラー:", err)
		return false, err
	}

	if err := tx.Commit(); err != nil {
		log.Println("資産情報更新：コミットエラー:", err)
		return false, err
	}
	log.Println("資産情報更新：成功")
	return true, nil
}

// すでに外でTxを張っている想定のユーティリティ
func ResetAssetStatus(tx *sql.Tx, assetID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	const query = `UPDATE assets SET owner = NULL, location = default_location WHERE id = ?`
	res, err := tx.ExecContext(ctx, query, assetID)
	if err != nil {
		log.Println("資産状態のリセットエラー:", err)
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("asset_id %d に該当する資産が見つかりません", assetID)
	}
	return nil
}
