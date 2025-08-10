package disposal

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func UndoRegisterDisposal(db *sql.DB, assetID int64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		log.Println("資産廃棄取り消し：初期化エラー:", err)
		return false, err
	}
	defer func() { _ = tx.Rollback() }()

	// 最新の廃棄記録を取得（日時同値の競合に備え id も降順）
	var disposalID int64
	var qty int64
	err = tx.QueryRowContext(ctx, `
		SELECT id, quantity
		FROM asset_disposals
		WHERE asset_id = ?
		ORDER BY disposal_date DESC, id DESC
		LIMIT 1`, assetID,
	).Scan(&disposalID, &qty)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("資産廃棄取り消し：対象なし")
		} else {
			log.Println("資産廃棄取り消し：記録取得エラー:", err)
		}
		return false, err
	}

	// 対象資産の行をロック（在庫とステータス更新の一貫性を保つ）
	var curQty int64
	if err := tx.QueryRowContext(ctx,
		`SELECT quantity FROM assets WHERE id = ? FOR UPDATE`, assetID,
	).Scan(&curQty); err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("asset_id %d が存在しません", assetID)
		}
		return false, err
	}

	// 廃棄記録を削除
	if _, err := tx.ExecContext(ctx,
		`DELETE FROM asset_disposals WHERE id = ?`, disposalID,
	); err != nil {
		log.Println("資産廃棄取り消し：削除エラー:", err)
		return false, err
	}

	// 数量を戻す（加算）＋在庫が正なら status を「正常(=1)」へ
	// 一発更新でステータスも整合させる（quantity は式右辺の元値で計算される）
	if _, err := tx.ExecContext(ctx, `
		UPDATE assets
		SET quantity = quantity + ?,
		    status_id = CASE WHEN quantity + ? > 0 THEN 1 ELSE status_id END
		WHERE id = ?`,
		qty, qty, assetID,
	); err != nil {
		log.Println("資産廃棄取り消し：数量/ステータス更新エラー:", err)
		return false, err
	}

	if err := tx.Commit(); err != nil {
		log.Println("資産廃棄取り消し：コミットエラー:", err)
		return false, err
	}
	log.Println("資産廃棄取り消し：成功")
	return true, nil
}
