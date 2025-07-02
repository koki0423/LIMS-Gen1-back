package assets

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	model "equipmentManager/internal/database/model/tables"
	"log"
)

func UpdateAsset(db *sql.DB, updated model.Asset) (bool, error) {
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


// func EditLocation(db *sql.DB, assetId int64, newLocation string) (bool, error) {
// 	//テストデータ
// 	assetId = int64(4) // 更新対象の資産ID
// 	newLocation = "seminar room 2"

// 	tx, err := db.Begin()
// 	if err != nil {
// 		return false, err
// 	}

// 	query := "UPDATE assets SET location = ? WHERE id = ?"
// 	_, err = tx.Exec(query, newLocation, assetId)
// 	if err != nil {
// 		log.Println("資産の配置場所更新：エラー:", err)
// 		tx.Rollback()
// 		return false, err
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		log.Println("資産の配置場所更新：コミットエラー:", err)
// 		return false, err
// 	}

// 	log.Println("資産の配置場所更新：成功")
// 	return true, nil
// }

// func EditOwner(db *sql.DB, assetId int64, newOwner string) (bool, error) {
// 	//テストデータ
// 	assetId = int64(4) // 更新対象の資産ID
// 	newOwner = "Arona"

// 	tx, err := db.Begin()
// 	if err != nil {
// 		return false, err
// 	}

// 	query := "UPDATE assets SET owner = ? WHERE id = ?"
// 	_, err = tx.Exec(query, newOwner, assetId)
// 	if err != nil {
// 		log.Println("資産の所有者更新：エラー:", err)
// 		tx.Rollback()
// 		return false, err
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		log.Println("資産の所有者更新：コミットエラー:", err)
// 		return false, err
// 	}

// 	log.Println("資産の所有者更新：成功")
// 	return true, nil
// }

// // 全体管理の場合のみ実行
// func EditQuantity(db *sql.DB, assetId int64, quantity int) (bool, error) {
// 	// 事前に管理区分を確認
// 	var categoryID int
// 	err := db.QueryRow(`
// 		SELECT m.management_category_id
// 		FROM assets AS a
// 		JOIN assets_masters AS m ON a.asset_master_id = m.id
// 		WHERE a.id = ?`, assetId).Scan(&categoryID)
// 	if err != nil {
// 		return false, fmt.Errorf("資産の管理区分取得エラー: %w", err)
// 	}

// 	if categoryID != 2 {
// 		return false, fmt.Errorf("全体管理でない資産の数量は更新できません（id=%d）", assetId)
// 	}

// 	tx, err := db.Begin()
// 	if err != nil {
// 		return false, err
// 	}

// 	query := `UPDATE assets SET quantity = ? WHERE id = ?`
// 	res, err := tx.Exec(query, quantity, assetId)
// 	if err != nil {
// 		log.Println("資産の数量更新：エラー:", err)
// 		tx.Rollback()
// 		return false, err
// 	}

// 	affected, err := res.RowsAffected()
// 	if err != nil {
// 		log.Println("資産の数量更新：件数取得エラー:", err)
// 		tx.Rollback()
// 		return false, err
// 	}
// 	if affected == 0 {
// 		tx.Rollback()
// 		return false, fmt.Errorf("更新対象が見つかりませんでした（id=%d）", assetId)
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		log.Println("資産の数量更新：コミットエラー:", err)
// 		return false, err
// 	}

// 	log.Printf("資産の数量更新：成功（id=%d, 新数量=%d）", assetId, quantity)
// 	return true, nil
// }


// func EditLastCheck(db *sql.DB, assetId int64, lastCheckDate string, lastChecker string) (bool, error) {
// 	//テストデータ
// 	assetId = int64(1) // 更新対象の資産ID
// 	lastCheckDate = "2023-10-15"
// 	lastChecker = "Shiroko Sunaookami"

// 	tx, err := db.Begin()
// 	if err != nil {
// 		return false, err
// 	}

// 	query := "UPDATE assets SET last_check_date = ?, last_checker = ? WHERE assetId = ?"
// 	_, err = tx.Exec(query, lastCheckDate, lastChecker, assetId)
// 	if err != nil {
// 		log.Println("資産の最終チェック更新：エラー:", err)
// 		tx.Rollback()
// 		return false, err
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		log.Println("資産の最終チェック更新：コミットエラー:", err)
// 		return false, err
// 	}

// 	log.Println("資産の最終チェック更新：成功")
// 	return true, nil
// }

// func EditNotes(db *sql.DB, assetId int64, notes string) (bool, error) {
// 	//テストデータ
// 	notes = "バッテリー交換済み"
// 	assetId = int64(1)

// 	tx, err := db.Begin()
// 	if err != nil {
// 		return false, err
// 	}

// 	query := "UPDATE assets SET notes = ? WHERE id = ?"
// 	_, err = tx.Exec(query, notes, assetId)
// 	if err != nil {
// 		log.Println("資産のメモ更新：エラー:", err)
// 		tx.Rollback()
// 		return false, err
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		log.Println("資産のメモ更新：コミットエラー:", err)
// 		return false, err
// 	}

// 	log.Println("資産のメモ更新：成功")
// 	return true, nil
// }

// func EditStatus(db *sql.DB, statusId, assetId int) (bool, error) {
// 	//テストデータ
// 	statusId = 2     // 1: 正常, 2: 故障, 3: 修理中, 4: 貸出中, 5: 廃棄済み
// 	assetId = 1 // 更新対象の資産ID

// 	tx, err := db.Begin()
// 	if err != nil {
// 		return false, err
// 	}

// 	query := "UPDATE assets SET status_id = ? WHERE id = ?"
// 	_, err = tx.Exec(query, statusId, assetId)
// 	if err != nil {
// 		log.Println("資産のステータス更新：エラー:", err)
// 		tx.Rollback()
// 		return false, err
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		log.Println("資産のステータス更新：コミットエラー:", err)
// 		return false, err
// 	}

// 	log.Println("資産のステータス更新：成功")
// 	return true, nil
// }
