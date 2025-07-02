package disposal


import (
	"database/sql"
	model "equipmentManager/internal/database/model/tables"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func RegisterDisposal(db *sql.DB, asset model.AssetsDisposal) (bool, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Println("資産廃棄：初期化エラー:", err)
		return false, err
	}

	if asset.IsIndividual {
		// 個別管理：1個ずつ登録
		query := `
			INSERT INTO asset_disposals 
				(asset_id, quantity, disposal_date, reason, processed_by) 
			VALUES (?, ?, ?, ?, ?);`
		_, err = tx.Exec(query,asset.AssetID, 1, asset.DisposalDate, asset.Reason, asset.ProcessedBy)
		if err != nil {
			log.Println("資産廃棄：登録エラー（個別）:", err)
			tx.Rollback()
			return false, err
		}
	} else {
		// 全体管理：任意の数量（通常は1）
		query := `
			INSERT INTO asset_disposals 
				(asset_id, quantity, disposal_date, reason, processed_by) 
			VALUES (?, ?, ?, ?, ?);`
		_, err = tx.Exec(query, asset.AssetID, asset.Quantity, asset.DisposalDate, asset.Reason, asset.ProcessedBy)
		if err != nil {
			log.Println("資産廃棄：登録エラー（全体）:", err)
			tx.Rollback()
			return false, err
		}
	}

	// 数量更新（減算）
	query := "UPDATE assets SET quantity = quantity - ? WHERE id = ?"
	_, err = tx.Exec(query, asset.Quantity, asset.AssetID)
	if err != nil {
		log.Println("資産廃棄：数量更新(減算)エラー:", err)
		tx.Rollback()
		return false, err
	}

	// ステータス更新：数量がゼロになったら廃棄済みに
	query = "UPDATE assets SET status_id = 5 WHERE id = ? AND quantity = 0"
	_, err = tx.Exec(query, asset.AssetID)
	if err != nil {
		log.Println("資産廃棄：ステータス更新(廃棄)エラー:", err)
		tx.Rollback()
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		log.Println("資産廃棄：コミットエラー:", err)
		return false, err
	}

	log.Println("資産廃棄：成功")
	return true, nil
}