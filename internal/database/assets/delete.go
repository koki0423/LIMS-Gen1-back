package assets

import (
	"database/sql"
	"log"
)

func DeleteAssetMasterByID(db *sql.DB, assetMasterID int64) (bool, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Println("資産マスター削除：トランザクション開始エラー:", err)
		return false, err
	}

	// 1. assets テーブルから asset_id を取得
	assetRows, err := tx.Query(`SELECT id FROM assets WHERE asset_master_id = ?`, assetMasterID)
	if err != nil {
		log.Println("資産マスター削除：資産ID取得クエリエラー:", err)
		tx.Rollback()
		return false, err
	}
	defer assetRows.Close()

	var assetIDs []int64
	for assetRows.Next() {
		var id int64
		if err := assetRows.Scan(&id); err != nil {
			log.Println("資産マスター削除：資産IDスキャンエラー:", err)
			tx.Rollback()
			return false, err
		}
		assetIDs = append(assetIDs, id)
	}

	// 2. 各 asset_id ごとに処理
	for _, assetID := range assetIDs {
		// 2-1. asset_lends.id を取得
		lendRows, err := tx.Query(`SELECT id FROM asset_lends WHERE asset_id = ?`, assetID)
		if err != nil {
			log.Println("資産マスター削除：貸出ID取得エラー:", err)
			tx.Rollback()
			return false, err
		}

		var lendIDs []int64
		for lendRows.Next() {
			var lid int64
			if err := lendRows.Scan(&lid); err != nil {
				log.Println("資産マスター削除：貸出IDスキャンエラー:", err)
				tx.Rollback()
				return false, err
			}
			lendIDs = append(lendIDs, lid)
		}
		lendRows.Close()

		// 2-2. asset_returns を削除
		for _, lendID := range lendIDs {
			_, err := tx.Exec(`DELETE FROM asset_returns WHERE lend_id = ?`, lendID)
			if err != nil {
				log.Println("資産マスター削除：返却履歴削除エラー:", err)
				tx.Rollback()
				return false, err
			}
		}

		// 2-3. asset_lends を削除
		_, err = tx.Exec(`DELETE FROM asset_lends WHERE asset_id = ?`, assetID)
		if err != nil {
			log.Println("資産マスター削除：貸出削除エラー:", err)
			tx.Rollback()
			return false, err
		}
	}

	// 3. assets を削除
	_, err = tx.Exec(`DELETE FROM assets WHERE asset_master_id = ?`, assetMasterID)
	if err != nil {
		log.Println("資産マスター削除：資産削除エラー:", err)
		tx.Rollback()
		return false, err
	}

	// 4. assets_masters を削除
	_, err = tx.Exec(`DELETE FROM assets_masters WHERE id = ?`, assetMasterID)
	if err != nil {
		log.Println("資産マスター削除：マスター削除エラー:", err)
		tx.Rollback()
		return false, err
	}

	// コミット
	if err = tx.Commit(); err != nil {
		log.Println("資産マスター削除：コミットエラー:", err)
		tx.Rollback()
		return false, err
	}

	log.Println("資産マスター削除：成功")
	return true, nil
}
