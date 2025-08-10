package disposal

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	model "equipmentManager/internal/database/model/tables"
)

func RegisterDisposal(db *sql.DB, req model.AssetsDisposal) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		log.Println("資産廃棄：初期化エラー:", err)
		return false, err
	}
	defer func() { _ = tx.Rollback() }()

	// 減算量の決定（個別管理は常に1）
	decQty := req.Quantity
	if req.IsIndividual {
		decQty = 1
	}
	if decQty <= 0 {
		return false, fmt.Errorf("廃棄数量が不正です: %d", decQty)
	}

	// 在庫をロックして確認
	var curQty int64
	if err := tx.QueryRowContext(ctx,
		`SELECT quantity FROM assets WHERE id = ? FOR UPDATE`, req.AssetID,
	).Scan(&curQty); err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("asset_id %d の資産が存在しません", req.AssetID)
		}
		return false, err
	}
	if curQty < int64(decQty) {
		return false, fmt.Errorf("在庫不足: 現在=%d, 減算=%d", curQty, decQty)
	}

	// 廃棄レコードを追加（is_individual も保存推奨）
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO asset_disposals
		  (asset_id, quantity, disposal_date, reason, processed_by, is_individual)
		VALUES (?, ?, ?, ?, ?, ?)`,
		req.AssetID, decQty, req.DisposalDate, req.Reason, req.ProcessedBy, req.IsIndividual,
	); err != nil {
		log.Println("資産廃棄：登録エラー:", err)
		return false, err
	}

	// 在庫を減算
	res, err := tx.ExecContext(ctx,
		`UPDATE assets SET quantity = quantity - ? WHERE id = ?`, decQty, req.AssetID,
	)
	if err != nil {
		log.Println("資産廃棄：数量更新(減算)エラー:", err)
		return false, err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return false, fmt.Errorf("資産廃棄：数量更新対象なし (asset_id=%d)", req.AssetID)
	}

	// 在庫がゼロなら廃棄済みに
	if _, err := tx.ExecContext(ctx,
		`UPDATE assets SET status_id = 5 WHERE id = ? AND quantity = 0`, req.AssetID,
	); err != nil {
		log.Println("資産廃棄：ステータス更新(廃棄)エラー:", err)
		return false, err
	}

	if err := tx.Commit(); err != nil {
		log.Println("資産廃棄：コミットエラー:", err)
		return false, err
	}
	log.Println("資産廃棄：成功")
	return true, nil
}
