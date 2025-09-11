package disposals

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Store struct{ db *sql.DB }

func NewStore(db *sql.DB) *Store { return &Store{db: db} }

// --- assets_master / assets 参照 ---

func (s *Store) ResolveMasterID(ctx context.Context, managementNumber string) (uint64, error) {
	const q = `SELECT asset_master_id FROM assets_master WHERE management_number = ?`
	var id uint64
	if err := s.db.QueryRowContext(ctx, q, managementNumber).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrNotFound("assets_master not found")
		}
		return 0, err
	}
	return id, nil
}

func (s *Store) LockAssetRow(ctx context.Context, tx *sql.Tx, masterID uint64) (assetID uint64, quantity uint, err error) {
	const q = `SELECT asset_id, quantity FROM assets WHERE asset_master_id = ? LIMIT 1 FOR UPDATE`
	row := tx.QueryRowContext(ctx, q, masterID)
	if err = row.Scan(&assetID, &quantity); err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, ErrNotFound("asset row not found")
		}
		return 0, 0, err
	}
	return assetID, quantity, nil
}

func (s *Store) UpdateAssetQuantity(ctx context.Context, tx *sql.Tx, assetID uint64, delta int) error {
	const q = `UPDATE assets SET quantity = quantity + ? WHERE asset_id = ?`
	res, err := tx.ExecContext(ctx, q, delta, assetID)
	if err != nil {
		return err
	}
	if aff, _ := res.RowsAffected(); aff != 1 {
		return ErrInternal("failed to update assets.quantity")
	}
	return nil
}

func (s *Store) UpdateAssetStatus(ctx context.Context, tx *sql.Tx, assetID uint64, statusCode int) error {
	const q = `UPDATE assets SET status_id = ? WHERE asset_id = ?`
	res, err := tx.ExecContext(ctx, q, statusCode, assetID)
	if err != nil {
		return err
	}
	if aff, _ := res.RowsAffected(); aff != 1 {
		return ErrInternal("failed to update assets.status")
	}
	return nil
}

// --- disposals ---

func (s *Store) InsertDisposal(ctx context.Context, tx *sql.Tx, m *Disposal) (uint64, error) {
	const q = `
	INSERT INTO disposals
	(disposal_ulid, management_number, quantity, disposed_at, reason, processed_by_id)
	VALUES
	(?, ?, ?, CURRENT_TIMESTAMP, ?, ?)`
	res, err := tx.ExecContext(ctx, q,
		m.DisposalULID, m.ManagementNumber, m.Quantity,
		nullStrOrNil(m.Reason), nullStrOrNil(m.ProcessedByID),
	)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return uint64(id), nil
}

func (s *Store) GetByULID(ctx context.Context, ul string) (*Disposal, error) {
	const q = `
	SELECT disposal_id, disposal_ulid, management_number, quantity, disposed_at, reason, processed_by_id
	FROM disposals WHERE disposal_ulid = ?`
	var m Disposal
	if err := s.db.QueryRowContext(ctx, q, ul).Scan(
		&m.DisposalID, &m.DisposalULID, &m.ManagementNumber, &m.Quantity,
		&m.DisposedAt, &m.Reason, &m.ProcessedByID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound("disposal not found")
		}
		return nil, err
	}
	return &m, nil
}

func (s *Store) List(ctx context.Context, f DisposalFilter, p Page) ([]Disposal, int64, error) {
	sb := strings.Builder{}
	sb.WriteString(`
	SELECT disposal_id, disposal_ulid, management_number, quantity, disposed_at, reason, processed_by_id
	FROM disposals WHERE 1=1`)

	args := []any{}
	if f.ManagementNumber != nil {
		sb.WriteString(` AND management_number = ?`)
		args = append(args, *f.ManagementNumber)
	}
	if f.ProcessedByID != nil {
		sb.WriteString(` AND processed_by_id = ?`)
		args = append(args, *f.ProcessedByID)
	}
	if f.From != nil {
		sb.WriteString(` AND disposed_at >= ?`)
		args = append(args, *f.From)
	}
	if f.To != nil {
		sb.WriteString(` AND disposed_at < ?`)
		args = append(args, *f.To)
	}

	order := "DESC"
	if strings.ToLower(p.Order) == "asc" {
		order = "ASC"
	}
	if p.Limit <= 0 {
		p.Limit = 50
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
	sb.WriteString(fmt.Sprintf(` ORDER BY disposed_at %s LIMIT ? OFFSET ?`, order))
	args = append(args, p.Limit, p.Offset)

	rows, err := s.db.QueryContext(ctx, sb.String(), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []Disposal
	for rows.Next() {
		var m Disposal
		if err := rows.Scan(&m.DisposalID, &m.DisposalULID, &m.ManagementNumber, &m.Quantity, &m.DisposedAt, &m.Reason, &m.ProcessedByID); err != nil {
			return nil, 0, err
		}
		items = append(items, m)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	// count
	cb := strings.Builder{}
	cb.WriteString(`SELECT COUNT(*) FROM disposals WHERE 1=1`)
	argsC := []any{}
	if f.ManagementNumber != nil {
		cb.WriteString(` AND management_number = ?`)
		argsC = append(argsC, *f.ManagementNumber)
	}
	if f.ProcessedByID != nil {
		cb.WriteString(` AND processed_by_id = ?`)
		argsC = append(argsC, *f.ProcessedByID)
	}
	if f.From != nil {
		cb.WriteString(` AND disposed_at >= ?`)
		argsC = append(argsC, *f.From)
	}
	if f.To != nil {
		cb.WriteString(` AND disposed_at < ?`)
		argsC = append(argsC, *f.To)
	}
	var total int64
	if err := s.db.QueryRowContext(ctx, cb.String(), argsC...).Scan(&total); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func nullStrOrNil(ns sql.NullString) any {
	if ns.Valid {
		return ns.String
	}
	return nil
}
