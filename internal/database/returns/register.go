package returns

import (
	"database/sql"
	model "equipmentManager/internal/database/model/tables"
	"log"
)

func RegisterAssetReturn(db *sql.Tx, r model.AssetReturn) (bool, error) {
	query := `
        INSERT INTO asset_returns (lend_id, returned_quantity, return_date, notes)
        VALUES (?, ?, ?, ?)
    `
	result, err := db.Exec(query, r.LendID, r.ReturnedQuantity, r.ReturnDate, r.Notes)
	if err != nil {
		log.Println("返却履歴の登録エラー:", err)
		return false, err
	}

	rows, _ := result.RowsAffected()
	log.Printf("返却履歴：%d 行を登録", rows)
	return true, nil
}
