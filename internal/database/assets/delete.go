package assets

import (
	"context"
	"database/sql"
	"log"
	"time"
)

func DeleteAssetMasterByID(db *sql.DB, assetMasterID int64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("資産マスター削除：トランザクション開始エラー:", err)
		return false, err
	}
	// 以降どこでreturnしても安全
	defer func() { _ = tx.Rollback() }()

	// 1) assets.id を取得
	assetRows, err := tx.QueryContext(ctx, `SELECT id FROM assets WHERE asset_master_id = ?`, assetMasterID)
	if err != nil {
		log.Println("資産マスター削除：資産ID取得クエリエラー:", err)
		return false, err
	}
	defer assetRows.Close()

	var assetIDs []int64
	for assetRows.Next() {
		var id int64
		if err := assetRows.Scan(&id); err != nil {
			log.Println("資産マスター削除：資産IDスキャンエラー:", err)
			return false, err
		}
		assetIDs = append(assetIDs, id)
	}
	if err := assetRows.Err(); err != nil {
		return false, err
	}

	// 2) 各 asset_id に対して関連削除
	for _, assetID := range assetIDs {
		// 2-1) 貸出IDを取得（必ずCloseされるスコープに包む）
		lendIDs, err := func() ([]int64, error) {
			rows, err := tx.QueryContext(ctx, `SELECT id FROM asset_lends WHERE asset_id = ?`, assetID)
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			var ids []int64
			for rows.Next() {
				var lid int64
				if err := rows.Scan(&lid); err != nil {
					return nil, err
				}
				ids = append(ids, lid)
			}
			return ids, rows.Err()
		}()
		if err != nil {
			log.Println("資産マスター削除：貸出ID取得エラー:", err)
			return false, err
		}

		// 2-2) 返却履歴を削除
		for _, lendID := range lendIDs {
			if _, err := tx.ExecContext(ctx, `DELETE FROM asset_returns WHERE lend_id = ?`, lendID); err != nil {
				log.Println("資産マスター削除：返却履歴削除エラー:", err)
				return false, err
			}
		}

		// 2-3) 貸出レコードを削除
		if _, err := tx.ExecContext(ctx, `DELETE FROM asset_lends WHERE asset_id = ?`, assetID); err != nil {
			log.Println("資産マスター削除：貸出削除エラー:", err)
			return false, err
		}
	}

	// 3) assets を削除
	if _, err := tx.ExecContext(ctx, `DELETE FROM assets WHERE asset_master_id = ?`, assetMasterID); err != nil {
		log.Println("資産マスター削除：資産削除エラー:", err)
		return false, err
	}

	// 4) assets_masters を削除
	if _, err := tx.ExecContext(ctx, `DELETE FROM assets_masters WHERE id = ?`, assetMasterID); err != nil {
		log.Println("資産マスター削除：マスター削除エラー:", err)
		return false, err
	}

	if err := tx.Commit(); err != nil {
		log.Println("資産マスター削除：コミットエラー:", err)
		return false, err
	}
	log.Println("資産マスター削除：成功")
	return true, nil
}
