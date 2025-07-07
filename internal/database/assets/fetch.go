package assets

import (
	"database/sql"
	model "equipmentManager/internal/database/model/tables"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// GET /assets/all
func FetchAssetsAll(db *sql.DB) ([]model.Asset, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Println("資産一覧取得：トランザクション開始エラー:", err)
		return nil, err
	}

	query := `SELECT * FROM assets ORDER BY id ASC;`
	rows, err := tx.Query(query)
	if err != nil {
		log.Println("資産一覧取得：クエリエラー:", err)
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()

	var assets []model.Asset
	for rows.Next() {
		var asset model.Asset
		err := rows.Scan(
			&asset.ID,
			&asset.ItemMasterID,
			&asset.Quantity,
			&asset.SerialNumber,
			&asset.StatusID,
			&asset.PurchaseDate,
			&asset.Owner,
			&asset.Location,
			&asset.LastCheckDate,
			&asset.LastChecker,
			&asset.Notes,
		)
		if err != nil {
			log.Println("資産一覧取得：スキャンエラー:", err)
			tx.Rollback()
			return nil, err
		}
		assets = append(assets, asset)
	}
	return assets, nil
}

// GET /assets/:id
func FetchAssetsByID(db *sql.DB, assetID int64) (*model.Asset, error) {
	if assetID <= 0 {
		log.Println("資産情報取得：無効な資産ID")
		return nil, nil
	}

	query := `SELECT * FROM assets WHERE id = ?`
	row := db.QueryRow(query, assetID)

	var asset model.Asset
	err := row.Scan(
		&asset.ID,
		&asset.ItemMasterID,
		&asset.Quantity,
		&asset.SerialNumber,
		&asset.StatusID,
		&asset.PurchaseDate,
		&asset.Owner,
		&asset.Location,
		&asset.LastCheckDate,
		&asset.LastChecker,
		&asset.Notes,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("資産情報取得：該当資産が存在しません")
			return nil, nil
		}
		log.Println("資産情報取得：スキャンエラー:", err)
		return nil, err
	}
	return &asset, nil
}

// GET /assets/master/all
func FetchAllAssetsMaster(db *sql.DB) ([]model.AssetsMaster, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Println("資産マスター一覧取得：トランザクション開始エラー:", err)
		return nil, err
	}
	query := `SELECT * FROM assets_masters ORDER BY id DESC;`
	rows, err := tx.Query(query)
	if err != nil {
		log.Println("資産マスター一覧取得：クエリエラー:", err)
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()

	var assetMasters []model.AssetsMaster
	for rows.Next() {
		var assetMaster model.AssetsMaster
		err := rows.Scan(
			&assetMaster.ID,
			&assetMaster.Name,
			&assetMaster.ManagementCategoryID,
			&assetMaster.GenreID,
			&assetMaster.Manufacturer,
			&assetMaster.ModelNumber,
		)
		if err != nil {
			log.Println("資産マスター一覧取得：スキャンエラー:", err)
			tx.Rollback()
			return nil, err
		}
		assetMasters = append(assetMasters, assetMaster)
	}
	return assetMasters, nil
}

// GET /assets/master/:id
func FetchAssetsMasterByID(db *sql.DB, assetMasterID int64) (*model.AssetsMaster, error) {
	if assetMasterID <= 0 {
		log.Println("資産マスター情報取得：無効な資産マスターID")
		return nil, nil
	}

	query := `SELECT * FROM assets_masters WHERE id = ?`
	row := db.QueryRow(query, assetMasterID)

	var assetMaster model.AssetsMaster
	err := row.Scan(
		&assetMaster.ID,
		&assetMaster.Name,
		&assetMaster.ManagementCategoryID,
		&assetMaster.GenreID,
		&assetMaster.Manufacturer,
		&assetMaster.ModelNumber,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("資産マスター情報取得：該当資産マスターが存在しません")
			return nil, nil
		}
		log.Println("資産マスター情報取得：スキャンエラー:", err)
		return nil, err
	}
	return &assetMaster, nil
}
