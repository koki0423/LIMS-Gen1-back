package disposal

import (
	"context"
	"database/sql"
	"time"

	model "equipmentManager/internal/database/model/tables"
)

// GET /disposals/all
func FetchDisposalAll(db *sql.DB) ([]model.AssetsDisposal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	const query = `
SELECT id, asset_id, quantity, disposal_date, reason, processed_by, is_individual
FROM asset_disposals
ORDER BY id ASC;
`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	disposals := make([]model.AssetsDisposal, 0, 128)
	for rows.Next() {
		var d model.AssetsDisposal
		if err := rows.Scan(
			&d.ID,
			&d.AssetID,
			&d.Quantity,
			&d.DisposalDate,
			&d.Reason,
			&d.ProcessedBy,
			&d.IsIndividual,
		); err != nil {
			return nil, err
		}
		disposals = append(disposals, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return disposals, nil
}

// GET /disposals/:id
func FetchDisposalByID(db *sql.DB, id int64) (model.AssetsDisposal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	const query = `
SELECT id, asset_id, quantity, disposal_date, reason, processed_by, is_individual
FROM asset_disposals WHERE id = ?;
`
	var d model.AssetsDisposal
	err := db.QueryRowContext(ctx, query, id).Scan(
		&d.ID,
		&d.AssetID,
		&d.Quantity,
		&d.DisposalDate,
		&d.Reason,
		&d.ProcessedBy,
		&d.IsIndividual,
	)
	if err != nil {
		return model.AssetsDisposal{}, err // ErrNoRows もここでそのまま返す
	}
	return d, nil
}
