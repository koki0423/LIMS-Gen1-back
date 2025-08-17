package assets

import (
	"context"
	"database/sql"
	"strings"
	// "log"

	"equipmentManager/domain"
)

func SearchMasters(db *sql.DB,
	ctx context.Context,
	q string,
	genreIDs []int64,
	page, size int,
) ([]domain.AssetMaster, int64, error) {

	like := "%" + escapeLike(q) + "%"

	// ★ ESCAPE は Go 文字列で "\\\\", SQL 上では '\\' になる
	orParts := []string{
		"`management_number` LIKE ? ESCAPE '\\\\'",
		"`name` LIKE ? ESCAPE '\\\\'",
	}

	args := []any{like, like}

	if len(genreIDs) > 0 {
		orParts = append(orParts, "genre_id IN ("+placeholders(len(genreIDs))+")")
		for i := 0; i < len(genreIDs); i++ {
			args = append(args, genreIDs[i])
		}
	}

	where := "(" + strings.Join(orParts, " OR ") + ")"

	// 総件数
	countSQL := "SELECT COUNT(*) FROM assets_masters WHERE " + where
	var total int64
	if err := db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// 本体
	offset := (page - 1) * size
	listSQL := "SELECT id, management_number, name, management_category_id, " +
		"genre_id, manufacturer, model_number " +
		"FROM assets_masters WHERE " + where + " " +
		"ORDER BY id ASC LIMIT ? OFFSET ?"

	args2 := append(append([]any{}, args...), size, offset)

	rows, err := db.QueryContext(ctx, listSQL, args2...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	out := make([]domain.AssetMaster, 0, size)
	for rows.Next() {
		var (
			id                int64
			managementNumber  sql.NullString
			name              string
			catID             int64
			genreID           sql.NullInt64
			manufacturer      sql.NullString
			modelNumber       sql.NullString
		)
		if err := rows.Scan(&id, &managementNumber, &name, &catID, &genreID, &manufacturer, &modelNumber); err != nil {
			return nil, 0, err
		}
		out = append(out, domain.AssetMaster{
			ID:                   id,
			Name:                 name,
			ManagementNumber:     nzStr(managementNumber),
			ManagementCategoryID: catID,
			GenreID:              ptrI64(genreID),
			Manufacturer:         ptrStr(manufacturer),
			ModelNumber:          ptrStr(modelNumber),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

// --- 小物ヘルパ ---
func escapeLike(s string) string {
	// % と _ と \ をエスケープ
	repl := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return repl.Replace(s)
}
func placeholders(n int) string {
	// n=3 → "?,?,?"
	if n <= 0 {
		return ""
	}
	b := strings.Builder{}
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteByte('?')
	}
	return b.String()
}
func nzStr(v sql.NullString) string {
	if !v.Valid {
		return ""
	}
	return v.String
}
func ptrStr(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	s := v.String
	return &s
}
func ptrI64(v sql.NullInt64) *int64 {
	if !v.Valid {
		return nil
	}
	x := v.Int64
	return &x
}
