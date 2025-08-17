package assets

import (
	"context"
	"database/sql"
	"time"

	"equipmentManager/domain"
)

func FetchAssetSummary(db *sql.DB) (domain.AssetSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// SQLクエリを実行して集計データを取得
	query := `
		SELECT (SELECT COUNT(*) FROM assets_masters) AS total_assets,
		COALESCE((SELECT SUM(CASE WHEN status_id = 4 THEN 1 ELSE 0 END) FROM assets), 0) AS lending_assets,
		COALESCE((SELECT SUM(CASE WHEN status_id = 2 THEN 1 ELSE 0 END) FROM assets), 0) AS breakdown_assets,
		COALESCE((SELECT SUM(CASE WHEN status_id = 5 THEN 1 ELSE 0 END) FROM assets), 0) AS dispose_assets,
		(SELECT COUNT(*) FROM assets_masters WHERE genre_id = 1) AS ind_assets,
		(SELECT COUNT(*) FROM assets_masters WHERE genre_id = 2) AS ofs_assets,
		(SELECT COUNT(*) FROM assets_masters WHERE genre_id = 3) AS fac_assets,
		(SELECT COUNT(*) FROM assets_masters WHERE genre_id = 4) AS emb_assets,
		(SELECT COUNT(*) FROM assets_masters WHERE genre_id = 5) AS adv_assets;
	`
	var s domain.AssetSummary
	err := db.QueryRowContext(ctx, query).Scan(
		&s.TotalAssets,
		&s.LendingAssets,
		&s.BreakdownAssets,
		&s.DisposeAssets,
		&s.IND_Assets,
		&s.OFS_Assets,
		&s.FAC_Assets,
		&s.EMB_Assets,
		&s.ADV_Assets,
	)
	return s, err
}
