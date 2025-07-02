package assets

import (
	"database/sql"
	model "equipmentManager/internal/database/model/tables"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// 個別管理テスト
func RegistrationEquipmentForIndivisual(db *sql.DB, asset model.AllDataOfAsset) (bool, error) {
	tx, err := db.Begin()
	if err != nil {
		return false, err
	}

	query := "INSERT INTO assets_masters ( name, management_category_id, genre_id, manufacturer, model_number)  VALUES (?, ?, ?, ?, ?)"
	res, err := tx.Exec(query, asset.Master.Name, asset.Master.ManagementCategoryID, asset.Master.GenreID.Int64, asset.Master.Manufacturer, asset.Master.ModelNumber)
	if err != nil {
		log.Println("(個別管理)機器のマスタ登録：エラー:", err)
		tx.Rollback()
		return false, err

	}
	log.Println("(個別管理)機器のマスタ登録：成功")

	asset.Asset.ItemMasterID, _ = res.LastInsertId()

	query = `INSERT INTO assets 
	         (asset_master_id, quantity, serial_number, status_id, purchase_date, owner, location, last_check_date, last_checker, notes)
	         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err = tx.Exec(query, asset.Asset.ItemMasterID, asset.Asset.Quantity, asset.Asset.SerialNumber, asset.Asset.StatusID, asset.Asset.PurchaseDate, asset.Asset.Owner, asset.Asset.Location, asset.Asset.LastCheckDate, asset.Asset.LastChecker, asset.Asset.Notes)
	if err != nil {
		log.Println("(個別管理)機器の資産登録：エラー:", err)
		tx.Rollback()
		return false, err
	}
	log.Println("(個別管理)機器の資産登録：成功")
	err = tx.Commit()
	if err != nil {
		log.Println("(個別管理)トランザクションコミット失敗:", err)
		return false, err
	}
	return true, nil
}

// 全体管理テスト
func RegistrationEquipmentForGeneral(db *sql.DB, asset model.AllDataOfAsset) (bool, error) {
	tx, err := db.Begin()
	if err != nil {
		return false, err
	}

	query := `INSERT INTO assets_masters (name, management_category_id, genre_id, manufacturer, model_number) 
	          VALUES (?, ?, ?, ?, ?)`
	res, err := tx.Exec(query,
		asset.Master.Name,
		asset.Master.ManagementCategoryID,
		asset.Master.GenreID.Int64,
		asset.Master.Manufacturer,
		asset.Master.ModelNumber,
	)
	if err != nil {
		log.Println("(全体管理)機器のマスタ登録：エラー:", err)
		tx.Rollback()
		return false, err
	}
	log.Println("(全体管理)機器のマスタ登録：成功")

	asset.Asset.ItemMasterID, _ = res.LastInsertId()

	query = `INSERT INTO assets 
	         (asset_master_id, quantity, serial_number, status_id, purchase_date, owner, location, last_check_date, last_checker, notes)
	         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.Exec(query,
		asset.Asset.ItemMasterID,
		asset.Asset.Quantity,
		asset.Asset.SerialNumber,
		asset.Asset.StatusID,
		asset.Asset.PurchaseDate,
		asset.Asset.Owner,
		asset.Asset.Location,
		asset.Asset.LastCheckDate,
		asset.Asset.LastChecker,
		asset.Asset.Notes,
	)
	if err != nil {
		log.Println("(全体管理)機器の資産登録：エラー:", err)
		tx.Rollback()
		return false, err
	}
	log.Println("(全体管理)機器の資産登録：成功")

	err = tx.Commit()
	if err != nil {
		log.Println("(全体管理)トランザクションコミット失敗:", err)
		return false, err
	}

	return true, nil
}
