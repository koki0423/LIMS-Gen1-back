package assets

import (
	"database/sql"
	model "equipmentManager/internal/database/model/tables"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func CrateAssetIndivisual(db *sql.DB, asset model.AllDataOfAsset) (bool, error)  {
	tx,err:= db.Begin()
	if err != nil {
		log.Println("(個別管理)トランザクション開始失敗:", err)
		return false, err
	}
	masterId,err:=CreateMaster(tx, asset.Master)
	if err != nil {
		tx.Rollback()
		log.Println("(個別管理)マスタ登録失敗:", err)
		return false, err
	}
	asset.Asset.ItemMasterID = masterId
	err = CreateAsset(tx, asset.Asset)
	if err != nil {
		tx.Rollback()
		log.Println("(個別管理)資産登録失敗:", err)
		return false, err
	}
	err = tx.Commit()
	if err != nil {
		log.Println("(個別管理)トランザクションコミット失敗:", err)
		return false, err
	}
	return true, nil
}

// 個別管理v2
func CreateMaster(tx *sql.Tx, master model.AssetsMaster) (int64, error) {
	query := "INSERT INTO assets_masters (name, management_category_id, genre_id, manufacturer, model_number) VALUES (?, ?, ?, ?, ?)"
	res, err := tx.Exec(query, master.Name, master.ManagementCategoryID, master.GenreID.Int64, master.Manufacturer, master.ModelNumber)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// 個別管理v2とセットで使う
func CreateAsset(tx *sql.Tx, asset model.Asset) error {
	query := `INSERT INTO assets (asset_master_id, quantity, serial_number, status_id, purchase_date, owner, location, last_check_date, last_checker, notes)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := tx.Exec(query, asset.ItemMasterID, asset.Quantity, asset.SerialNumber, asset.StatusID, asset.PurchaseDate, asset.Owner, asset.Location, asset.LastCheckDate, asset.LastChecker, asset.Notes)
	return err
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
