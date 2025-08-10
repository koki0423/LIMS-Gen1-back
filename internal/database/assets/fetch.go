package assets

import (
	"context"
	"database/sql"
	"log"
	"time"

	model "equipmentManager/internal/database/model/tables"
)

// GET /assets/all
func FetchAssetsAll(db *sql.DB) ([]model.Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 列を明示（スキーマ変更に強くする）
	const query = `
SELECT id, item_master_id, quantity, serial_number, status_id,
       purchase_date, owner, location, default_location,
       last_check_date, last_checker, notes
FROM assets
ORDER BY id ASC;
`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Println("資産一覧取得：クエリエラー:", err)
		return nil, err
	}
	defer rows.Close()

	assets := make([]model.Asset, 0, 128)
	for rows.Next() {
		var a model.Asset
		if err := rows.Scan(
			&a.ID,
			&a.ItemMasterID,
			&a.Quantity,
			&a.SerialNumber,
			&a.StatusID,
			&a.PurchaseDate,
			&a.Owner,
			&a.Location,
			&a.DefaultLocation,
			&a.LastCheckDate,
			&a.LastChecker,
			&a.Notes,
		); err != nil {
			log.Println("資産一覧取得：スキャンエラー:", err)
			return nil, err
		}
		assets = append(assets, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return assets, nil
}

// GET /assets/:id
func FetchAssetsByID(db *sql.DB, assetID int64) (*model.Asset, error) {
	if assetID <= 0 {
		log.Println("資産情報取得：無効な資産ID")
		return nil, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	const query = `
SELECT id, item_master_id, quantity, serial_number, status_id,
       purchase_date, owner, location, default_location,
       last_check_date, last_checker, notes
FROM assets
WHERE id = ?;
`
	row := db.QueryRowContext(ctx, query, assetID)

	var a model.Asset
	if err := row.Scan(
		&a.ID,
		&a.ItemMasterID,
		&a.Quantity,
		&a.SerialNumber,
		&a.StatusID,
		&a.PurchaseDate,
		&a.Owner,
		&a.Location,
		&a.DefaultLocation,
		&a.LastCheckDate,
		&a.LastChecker,
		&a.Notes,
	); err != nil {
		if err == sql.ErrNoRows {
			log.Println("資産情報取得：該当資産が存在しません")
			return nil, nil
		}
		log.Println("資産情報取得：スキャンエラー:", err)
		return nil, err
	}
	return &a, nil
}

// GET /assets/master/all
func FetchAllAssetsMaster(db *sql.DB) ([]model.AssetsMaster, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// ※ ManagementNumber を含め、列順をテーブル定義に合わせて明示
	const query = `
SELECT id, management_number, name, management_category_id, genre_id, manufacturer, model_number
FROM assets_masters
ORDER BY id ASC;
`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Println("資産マスター一覧取得：クエリエラー:", err)
		return nil, err
	}
	defer rows.Close()

	list := make([]model.AssetsMaster, 0, 128)
	for rows.Next() {
		var m model.AssetsMaster
		if err := rows.Scan(
			&m.ID,
			&m.ManagementNumber,
			&m.Name,
			&m.ManagementCategoryID,
			&m.GenreID,
			&m.Manufacturer,
			&m.ModelNumber,
		); err != nil {
			log.Println("資産マスター一覧取得：スキャンエラー:", err)
			return nil, err
		}
		list = append(list, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

// GET /assets/master/:id
func FetchAssetsMasterByID(db *sql.DB, assetMasterID int64) (*model.AssetsMaster, error) {
	if assetMasterID <= 0 {
		log.Println("資産マスター情報取得：無効な資産マスターID")
		return nil, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// ここも列を明示。元コードは Scan の列順がズレていたので注意
	const query = `
SELECT id, management_number, name, management_category_id, genre_id, manufacturer, model_number
FROM assets_masters WHERE id = ?;
`
	row := db.QueryRowContext(ctx, query, assetMasterID)

	var m model.AssetsMaster
	if err := row.Scan(
		&m.ID,
		&m.ManagementNumber,
		&m.Name,
		&m.ManagementCategoryID,
		&m.GenreID,
		&m.Manufacturer,
		&m.ModelNumber,
	); err != nil {
		if err == sql.ErrNoRows {
			log.Println("資産マスター情報取得：該当資産マスターが存在しません")
			return nil, nil
		}
		log.Println("資産マスター情報取得：スキャンエラー:", err)
		return nil, err
	}
	return &m, nil
}
