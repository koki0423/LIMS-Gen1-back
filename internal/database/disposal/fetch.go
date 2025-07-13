package disposal

import (
	"database/sql"
	model "equipmentManager/internal/database/model/tables"
)

func FetchDisposalAll(db *sql.DB) ([]model.AssetsDisposal, error) {
	query := "SELECT * FROM asset_disposals"

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query(query)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()

	var disposals []model.AssetsDisposal
	for rows.Next() {
		var disposal model.AssetsDisposal
		err := rows.Scan(
			&disposal.ID,
			&disposal.AssetID,
			&disposal.Quantity,
			&disposal.DisposalDate,
			&disposal.Reason,
			&disposal.ProcessedBy,
			&disposal.IsIndividual,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		disposals = append(disposals, disposal)
	}
	return disposals, tx.Commit()
}

func FetchDisposalByID(db *sql.DB, id int64) (model.AssetsDisposal, error) {
	tx,err:= db.Begin()
	if err != nil {
		return model.AssetsDisposal{}, err
	}
	query := "SELECT * FROM asset_disposals WHERE id = ?"

	rows, err := tx.Query(query, id)
	if err != nil {
		tx.Rollback()
		return model.AssetsDisposal{}, err
	}
	var disposal model.AssetsDisposal
	if rows.Next() {
		err := rows.Scan(
			&disposal.ID,
			&disposal.AssetID,
			&disposal.Quantity,
			&disposal.DisposalDate,
			&disposal.Reason,
			&disposal.ProcessedBy,
			&disposal.IsIndividual,
		)
		if err != nil {
			tx.Rollback()
			return model.AssetsDisposal{}, err
		}
	} else {
		tx.Rollback()
		return model.AssetsDisposal{}, sql.ErrNoRows
	}
	return disposal, tx.Commit()
}
