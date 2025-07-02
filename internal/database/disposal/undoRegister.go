package disposal
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func UndoRegisterDisposal(db *sql.DB, assetId int) (bool, error) {
	var disposalID int
	var qty int

	tx, err := db.Begin()
	if err != nil {
		log.Println("資産廃棄取り消し：初期化エラー:", err)
		return false, err
	}

	// 最新の廃棄記録を取得（1件）
	query := `
		SELECT id, quantity 
		FROM asset_disposals 
		WHERE asset_id = ? 
		ORDER BY disposal_date DESC 
		LIMIT 1`
	err = tx.QueryRow(query, assetId).Scan(&disposalID, &qty)
	if err != nil {
		log.Println("資産廃棄取り消し：記録取得エラー or 対象なし")
		tx.Rollback()
		return false, err
	}

	// 廃棄記録削除
	_, err = tx.Exec("DELETE FROM asset_disposals WHERE id = ?", disposalID)
	if err != nil {
		log.Println("資産廃棄取り消し：削除エラー:", err)
		tx.Rollback()
		return false, err
	}

	// 数量戻す（加算）
	query = "UPDATE assets SET quantity = quantity + ? WHERE id = ?"
	_, err = tx.Exec(query, qty, assetId)
	if err != nil {
		log.Println("資産廃棄取り消し：数量更新(加算)エラー:", err)
		tx.Rollback()
		return false, err
	}

	// ステータス戻す（0個状態 → 1個以上になったら正常に）
	query = "UPDATE assets SET status_id = 1 WHERE id = ? AND quantity > 0"
	_, err = tx.Exec(query, assetId)
	if err != nil {
		log.Println("資産廃棄取り消し：ステータス更新(正常)エラー:", err)
		tx.Rollback()
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		log.Println("資産廃棄取り消し：コミットエラー:", err)
		return false, err
	}

	log.Println("資産廃棄取り消し：成功")
	return true, nil
}
