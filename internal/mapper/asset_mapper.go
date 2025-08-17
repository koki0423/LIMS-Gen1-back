// internal/mapper/asset_mapper.go
package mapper

import (
	"database/sql"
	"time"

	model "equipmentManager/internal/database/model/tables"
	"equipmentManager/domain"
)

// --- Null系 → *T 変換ユーティリティ ---

func ptrInt64(v sql.NullInt64) *int64 {
	if !v.Valid {
		return nil
	}
	x := v.Int64
	return &x
}
func ptrString(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	s := v.String
	return &s
}
func ptrTime(v sql.NullTime) *time.Time {
	if !v.Valid {
		return nil
	}
	t := v.Time
	return &t
}

// --- 単体変換 ---

func ToDomainAssetMaster(m model.AssetsMaster) domain.AssetMaster {
	return domain.AssetMaster{
		ID:                   m.ID,
		Name:                 m.Name,
		ManagementNumber:     m.ManagementNumber,          // 既存値をそのまま
		ManagementCategoryID: m.ManagementCategoryID,
		GenreID:              ptrInt64(m.GenreID),
		Manufacturer:         ptrString(m.Manufacturer),
		ModelNumber:          ptrString(m.ModelNumber),
	}
}

func ToDomainAsset(a model.Asset) domain.Asset {
	return domain.Asset{
		ID:               a.ID,
		AssetMasterID:    a.AssetMasterID,
		Quantity:         a.Quantity,
		StatusID:         a.StatusID,
		SerialNumber:     ptrString(a.SerialNumber),
		PurchaseDate:     ptrTime(a.PurchaseDate),
		Owner:            ptrString(a.Owner),
		Location:         ptrString(a.Location),
		DefaultLocation:  ptrString(a.DefaultLocation),
		LastCheckDate:    ptrTime(a.LastCheckDate),
		LastChecker:      ptrString(a.LastChecker),
		Notes:            ptrString(a.Notes),
	}
}

// --- 複合（AllDataOfAsset） ---

// ハンドラに渡しやすい形でまとめたい場合
type AllDataOfAsset struct {
	Master domain.AssetMaster `json:"master"`
	Asset  domain.Asset       `json:"asset"`
}

func ToDomainAll(m model.AllDataOfAsset) AllDataOfAsset {
	return AllDataOfAsset{
		Master: ToDomainAssetMaster(m.Master),
		Asset:  ToDomainAsset(m.Asset),
	}
}

// --- スライス変換 ---

func ToDomainAssetMasters(src []model.AssetsMaster) []domain.AssetMaster {
	out := make([]domain.AssetMaster, 0, len(src))
	for i := range src {
		out = append(out, ToDomainAssetMaster(src[i]))
	}
	return out
}

func ToDomainAssets(src []model.Asset) []domain.Asset {
	out := make([]domain.Asset, 0, len(src))
	for i := range src {
		out = append(out, ToDomainAsset(src[i]))
	}
	return out
}

func ToDomainAllSlice(src []model.AllDataOfAsset) []AllDataOfAsset {
	out := make([]AllDataOfAsset, 0, len(src))
	for i := range src {
		out = append(out, ToDomainAll(src[i]))
	}
	return out
}
