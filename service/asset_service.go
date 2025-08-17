package service

import (
	"context"
	"database/sql"
	"equipmentManager/domain"
	"equipmentManager/internal/database/assets"
	model "equipmentManager/internal/database/model/tables"
	"equipmentManager/utils"
	"fmt"
	"strconv"
	"strings"
	"time"
	"log"
)

type AssetService struct {
	DB *sql.DB
}

func NewAssetService(db *sql.DB) *AssetService {
	return &AssetService{DB: db}
}

// --- ドメイン変換共通関数群 ---
func ModelToDomainAsset(model model.Asset) domain.Asset {
	return domain.Asset{
		ID:              model.ID,
		AssetMasterID:   model.AssetMasterID,
		Quantity:        model.Quantity,
		StatusID:        model.StatusID,
		SerialNumber:    utils.NullStringToPtr(model.SerialNumber),
		PurchaseDate:    utils.NullTimeToPtr(model.PurchaseDate),
		Owner:           utils.NullStringToPtr(model.Owner),
		Location:        utils.NullStringToPtr(model.Location),
		DefaultLocation: utils.NullStringToPtr(model.DefaultLocation),
		LastCheckDate:   utils.NullTimeToPtr(model.LastCheckDate),
		LastChecker:     utils.NullStringToPtr(model.LastChecker),
		Notes:           utils.NullStringToPtr(model.Notes),
	}
}

func toDomainAssetMaster(m model.AssetsMaster) domain.AssetMaster {
	return domain.AssetMaster{
		ID:                   m.ID,
		Name:                 m.Name,
		ManagementNumber:     m.ManagementNumber,
		ManagementCategoryID: m.ManagementCategoryID,
		GenreID:              nullInt64ToPtr(m.GenreID),
		Manufacturer:         nullStringToPtr(m.Manufacturer),
		ModelNumber:          nullStringToPtr(m.ModelNumber),
	}
}

// --- 変換処理 ---
func convertAssetsModelListToDomain(src []model.Asset) []domain.Asset {
	result := make([]domain.Asset, 0, len(src))
	for i := 0; i < len(src); i++ {
		result = append(result, ModelToDomainAsset(src[i]))
	}
	return result
}

func convertAssetMasterModelListToDomain(src []model.AssetsMaster) []domain.AssetMaster {
	result := make([]domain.AssetMaster, 0, len(src))
	for i := 0; i < len(src); i++ {
		result = append(result, toDomainAssetMaster(src[i]))
	}
	return result
}

func convertDomainToModel(domainAsset domain.CreateAssetRequest) (model.AssetsMaster, model.Asset) {
	master := model.AssetsMaster{
		Name:                 utils.PtrStringToString(domainAsset.Name),
		ManagementCategoryID: utils.PtrInt64ToInt64(domainAsset.ManagementCategoryID),
		GenreID:              utils.Int64ToNullInt64(domainAsset.GenreID),
		Manufacturer:         utils.StringToNullString(domainAsset.Manufacturer),
		ModelNumber:          utils.StringToNullString(domainAsset.ModelNumber),
	}

	asset := model.Asset{
		AssetMasterID: utils.PtrInt64ToInt64(domainAsset.AssetMasterID),
		Quantity:      domainAsset.Quantity,
		StatusID:      domainAsset.StatusID,
		SerialNumber:  utils.StringToNullString(domainAsset.SerialNumber),
		PurchaseDate:  utils.StringToNullTime(domainAsset.PurchaseDate),
		Owner:         utils.StringToNullString(domainAsset.Owner),
		Location:      utils.StringToNullString(domainAsset.DefaultLocation),
		LastCheckDate: utils.StringToNullTime(domainAsset.LastCheckDate),
		LastChecker:   utils.StringToNullString(domainAsset.LastChecker),
		Notes:         utils.StringToNullString(domainAsset.Notes),
	}

	return master, asset
}

func convertToDomainEditAsset(editedAsset domain.EditAssetRequest) model.Asset {
	return model.Asset{
		ID:            editedAsset.AssetID,
		Quantity:      utils.PtrIntToInt(editedAsset.Quantity),
		SerialNumber:  utils.StringToNullString(editedAsset.SerialNumber),
		StatusID:      utils.PtrInt64ToInt64(editedAsset.StatusID),
		PurchaseDate:  utils.StringToNullTime(editedAsset.PurchaseDate),
		Owner:         utils.StringToNullString(editedAsset.Owner),
		Location:      utils.StringToNullString(editedAsset.Location),
		LastCheckDate: utils.StringToNullTime(editedAsset.LastCheckDate),
		LastChecker:   utils.StringToNullString(editedAsset.LastChecker),
		Notes:         utils.StringToNullString(editedAsset.Notes),
	}
}

func nullStringToPtr(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	s := v.String
	return &s
}
func nullInt64ToPtr(v sql.NullInt64) *int64 {
	if !v.Valid {
		return nil
	}
	x := v.Int64
	return &x
}

var categoryNameMap = map[int]string{
	1: "IND",
	2: "OFS",
	3: "FAC",
	4: "EMB",
	5: "ADV",
}

// --- サービス ---
// POST /assets
func (e *AssetService) CreateAssetWithMaster(newAsset domain.CreateAssetRequest) (string, error) {
	genrePrefix, ok := categoryNameMap[int(*newAsset.GenreID)]
	if !ok {
		return "", fmt.Errorf("未知のジャンルID: %d", newAsset.GenreID)
	}
	loc, _ := time.LoadLocation("Asia/Tokyo")
	dateStr := time.Now().In(loc).Format("20060102")

	model_master, model_asset := convertDomainToModel(newAsset)

	// 管理番号は登録時に生成するので適当な文字列を設定
	model_master.ManagementNumber = fmt.Sprintf("%s-%s", "tmp", utils.GenerateUUID())
	model_asset.StatusID = 1

	switch model_master.ManagementCategoryID {
	case 1:
		return assets.CrateAssetIndivisual(e.DB, model_master, model_asset, genrePrefix, dateStr)
	case 2:
		return assets.CreateAssetCollective(e.DB, model_master, model_asset, genrePrefix, dateStr)
	default:
		return "", fmt.Errorf("未知の管理カテゴリID: %d", model_master.ManagementCategoryID)
	}
}

// GET /assets/all
func (e *AssetService) GetAssetAll() ([]domain.Asset, error) {
	dbAssets, err := assets.FetchAssetsAll(e.DB)
	if err != nil {
		return nil, err
	}
	return convertAssetsModelListToDomain(dbAssets), nil
}

// GET /assets/:id
func (e *AssetService) GetAssetById(id int) (domain.Asset, error) {
	dbAsset, err := assets.FetchAssetsByID(e.DB, int64(id))
	if err != nil {
		return domain.Asset{}, err
	}
	return ModelToDomainAsset(*dbAsset), nil
}

// PUT /assets/edit/:id
func (e *AssetService) PutAssetsEdit(editedAsset domain.EditAssetRequest, id int64) (bool, error) {

	return assets.UpdateAsset(e.DB, convertToDomainEditAsset(editedAsset), id)
}

// GET /assets/master/all
func (e *AssetService) GetAssetMasterAll() ([]domain.AssetMaster, error) {
	dbAssetMasters, err := assets.FetchAllAssetsMaster(e.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch asset masters: %w", err)
	}
	return convertAssetMasterModelListToDomain(dbAssetMasters), nil
}

// GET /assets/master/:id
func (e *AssetService) GetAssetMasterById(id int) (domain.AssetMaster, error) {
	dbAssetMaster, err := assets.FetchAssetsMasterByID(e.DB, int64(id))
	if err != nil {
		return domain.AssetMaster{}, fmt.Errorf("failed to fetch asset master by ID: %w", err)
	}
	return toDomainAssetMaster(*dbAssetMaster), nil
}

// DELETE /assets/master/:id 主キーはmasterの主キー
// 将来的に責任分解する
func (e *AssetService) DeleteAssetMasterById(id int) (bool, error) {
	status, err := assets.DeleteAssetMasterByID(e.DB, int64(id))
	if err != nil {
		return false, fmt.Errorf("failed to delete asset master by ID: %w", err)
	}
	return status, nil
}

func (e *AssetService) Search(ctx context.Context, q string, page, size int) ([]domain.AssetMaster, int64, error) {
	// 学内コード/日本語名/数値ID から genre_id 候補を推定
	genreIDs := guessGenreIDs(q)
	log.Printf("Search: q=%q, genreIDs=%v\n", q, genreIDs)
	return assets.SearchMasters(e.DB, ctx, q, genreIDs, page, size)
}

// 要件に出てきた5ジャンルを素直にマッピング
func guessGenreIDs(q string) []int64 {
	if q == "" {
		return nil
	}
	normalized := strings.ToUpper(strings.TrimSpace(q))
	jp := strings.TrimSpace(q)

	// 数値ID指定（"3" など）
	if id, err := strconv.ParseInt(normalized, 10, 64); err == nil && id >= 1 && id <= 9999 {
		return []int64{id}
	}

	type g struct {
		id         int64
		code, name string
	}
	genres := []g{
		{1, "IND", "個人"},
		{2, "OFS", "事務"},
		{3, "FAC", "ファシリティ"},
		{4, "EMB", "組込みシステム"},
		{5, "ADV", "高度情報演習"},
	}

	out := make([]int64, 0, 2)
	for i := 0; i < len(genres); i++ {
		if normalized == genres[i].code || strings.Contains(genres[i].name, jp) {
			out = append(out, genres[i].id)
		}
	}
	return out
}

func (e *AssetService) GetAssetSummary() (domain.AssetSummary, error) {
	dbSummary, err := assets.FetchAssetSummary(e.DB)
	if err != nil {
		return domain.AssetSummary{}, fmt.Errorf("failed to fetch asset summary: %w", err)
	}
	return dbSummary, nil
}
