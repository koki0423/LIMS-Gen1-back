package assets

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	model "equipmentManager/internal/database/model/tables"
)

// 個別管理（Quantityは常に1）
func CrateAssetIndivisual(db *sql.DB, master model.AssetsMaster, asset model.Asset, genrePrefix string, dateStr string) (string, error) {
	asset.Quantity = 1

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("(個別管理)トランザクション開始失敗:", err)
		return "", err
	}
	defer func() { _ = tx.Rollback() }()

	// 1) マスタ挿入（管理番号は後で付与）
	// log.Println("(個別管理)マスタ登録開始 管理番号:", master.ManagementNumber)
	masterID, err := createMaster(ctx, tx, master)
	if err != nil {
		log.Println("(個別管理)マスタ登録失敗:", err)
		return "", err
	}

	// 2) 管理番号生成 → 付番
	managementNumber := fmt.Sprintf("%s-%s-%04d", genrePrefix, dateStr, masterID)
	if err := insertManagementNumber(ctx, tx, masterID, managementNumber); err != nil {
		log.Println("(個別管理)管理番号登録失敗:", err)
		return "", err
	}

	// 3) 資産挿入
	asset.AssetMasterID = masterID
	if err := createAsset(ctx, tx, asset); err != nil {
		log.Println("(個別管理)資産登録失敗:", err)
		return "", err
	}

	if err := tx.Commit(); err != nil {
		log.Println("(個別管理)トランザクションコミット失敗:", err)
		return "", err
	}
	return managementNumber, nil
}

// 一括管理
func CreateAssetCollective(db *sql.DB, master model.AssetsMaster, asset model.Asset, genrePrefix string, dateStr string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("(全体管理)トランザクション開始失敗:", err)
		return "", err
	}
	defer func() { _ = tx.Rollback() }()

	masterID, err := createMaster(ctx, tx, master)
	if err != nil {
		log.Println("(全体管理)マスタ登録失敗:", err)
		return "", err
	}

	managementNumber := fmt.Sprintf("%s-%s-%04d", genrePrefix, dateStr, masterID)
	if err := insertManagementNumber(ctx, tx, masterID, managementNumber); err != nil {
		log.Println("(全体管理)管理番号登録失敗:", err)
		return "", err
	}

	asset.AssetMasterID = masterID
	if err := createAsset(ctx, tx, asset); err != nil {
		log.Println("(全体管理)資産登録失敗:", err)
		return "", err
	}

	if err := tx.Commit(); err != nil {
		log.Println("(全体管理)トランザクションコミット失敗:", err)
		return "", err
	}
	return managementNumber, nil
}

// --- 内部ユーティリティ ---

// 管理番号なしでマスタを追加し、採番用のIDを返す
func createMaster(ctx context.Context, tx *sql.Tx, master model.AssetsMaster) (int64, error) {
	// management_numberは後更新にする
	const q = `
INSERT INTO assets_masters (name, management_category_id, genre_id, manufacturer, model_number)
VALUES (?, ?, ?, ?, ?)`
	res, err := tx.ExecContext(ctx, q,
		master.Name,
		master.ManagementCategoryID,
		master.GenreID.Int64, // ここは型に合わせて（sql.NullInt64 なら .Int64）
		master.Manufacturer,
		master.ModelNumber,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func insertManagementNumber(ctx context.Context, tx *sql.Tx, masterID int64, managementNumber string) error {
	const q = `UPDATE assets_masters SET management_number = ? WHERE id = ?`
	_, err := tx.ExecContext(ctx, q, managementNumber, masterID)
	return err
}

func createAsset(ctx context.Context, tx *sql.Tx, asset model.Asset) error {
	// default_location に location を初期設定
	const q = `
INSERT INTO assets (asset_master_id, quantity, serial_number, status_id, purchase_date,
                    location, default_location, last_check_date, last_checker, notes)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := tx.ExecContext(ctx, q,
		asset.AssetMasterID,
		asset.Quantity,
		asset.SerialNumber,
		asset.StatusID,
		asset.PurchaseDate,
		asset.Location,
		asset.Location,
		asset.LastCheckDate,
		asset.LastChecker,
		asset.Notes,
	)
	return err
}
