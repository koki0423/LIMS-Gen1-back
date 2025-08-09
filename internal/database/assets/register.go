package assets

import (
	"database/sql"
	model "equipmentManager/internal/database/model/tables"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func CrateAssetIndivisual(db *sql.DB, master model.AssetsMaster, asset model.Asset, genrePrefix string, dateStr string) (string, error) {
	//個別管理はQuantityを1に固定する
	//フロントからは1が送られるはずだが一応再登録しておく
	asset.Quantity = 1

	tx, err := db.Begin()
	if err != nil {
		log.Println("(個別管理)トランザクション開始失敗:", err)
		return "", err
	}
	masterId, err := createMaster(tx, master)
	if err != nil {
		tx.Rollback()
		log.Println("(個別管理)マスタ登録失敗:", err)
		return "", err
	}
	asset.ItemMasterID = masterId

	managementNumber := fmt.Sprintf("%s-%s-%04d", genrePrefix, dateStr, masterId)

	err = insertManagementNumber(tx, masterId, managementNumber)
	if err != nil {
		tx.Rollback()
		log.Println("(個別管理)管理番号登録失敗:", err)
		return "", err
	}

	err = createAsset(tx, asset)
	if err != nil {
		tx.Rollback()
		log.Println("(個別管理)資産登録失敗:", err)
		return "", err
	}
	err = tx.Commit()
	if err != nil {
		log.Println("(個別管理)トランザクションコミット失敗:", err)
		return "", err
	}
	return managementNumber, nil
}

func CreateAssetCollective(db *sql.DB, master model.AssetsMaster, asset model.Asset, genrePrefix string, dateStr string) (string, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Println("(全体管理)トランザクション開始失敗:", err)
		return "", err
	}
	masterID, err := createMaster(tx, master)
	if err != nil {
		tx.Rollback()
		log.Println("(全体管理)マスタ登録失敗:", err)
		return "", err
	}

	asset.ItemMasterID = masterID
	managementNumber := fmt.Sprintf("%s-%s-%04d", genrePrefix, dateStr, masterID)

	err = insertManagementNumber(tx, masterID, managementNumber)
	if err != nil {
		tx.Rollback()
		log.Println("(全体管理)管理番号登録失敗:", err)
		return "", err
	}

	err = createAsset(tx, asset)
	if err != nil {
		tx.Rollback()
		log.Println("(全体管理)資産登録失敗:", err)
		return "", err
	}
	err = tx.Commit()
	if err != nil {
		log.Println("(全体管理)トランザクションコミット失敗:", err)
		return "", err
	}
	return managementNumber, nil
}

// createMaster は資産マスタをデータベースに登録し、登録されたマスタのIDを返す
func createMaster(tx *sql.Tx, master model.AssetsMaster) (int64, error) {
	query := "INSERT INTO assets_masters (management_number,name, management_category_id, genre_id, manufacturer, model_number) VALUES (?, ?, ?, ?, ?, ?)"
	res, err := tx.Exec(query, master.ManagementNumber, master.Name, master.ManagementCategoryID, master.GenreID.Int64, master.Manufacturer, master.ModelNumber)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func insertManagementNumber(tx *sql.Tx, masterID int64, managementNumber string) error {
	query := "UPDATE assets_masters SET management_number = ? WHERE id = ?"
	_, err := tx.Exec(query, managementNumber, masterID)
	return err
}

// createAsset は資産をデータベースに登録する
// 登録直後は場所＝デフォルト保管場所とする → 貸出時に場所を借りている人に変更するのでクエリパラメータのdefault_locationはlocationを入れる
// 貸出時にownerが決まるので、ownerはnilのまま
func createAsset(tx *sql.Tx, asset model.Asset) error {
	query := `INSERT INTO assets (asset_master_id, quantity, serial_number, status_id, purchase_date, location,default_location ,last_check_date, last_checker, notes)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := tx.Exec(query, asset.ItemMasterID, asset.Quantity, asset.SerialNumber, asset.StatusID, asset.PurchaseDate, asset.Location, asset.Location, asset.LastCheckDate, asset.LastChecker, asset.Notes)
	return err
}
