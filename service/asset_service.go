package service

import (
	"database/sql"
	"equipmentManager/domain"
	"equipmentManager/internal/database/assets"
	model "equipmentManager/internal/database/model/tables"
	"equipmentManager/utils"
	"fmt"
	"time"
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
		ItemMasterID:    model.ItemMasterID,
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
	var ret_value domain.AssetMaster
	ret_value = domain.AssetMaster{
		ID:           m.ID,
		Name:         m.Name,
		Manufacturer: utils.NullStringToPtr(m.Manufacturer),
		ModelNumber:  utils.NullStringToPtr(m.ModelNumber),
	}
	return ret_value
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
		ItemMasterID:  utils.PtrInt64ToInt64(domainAsset.AssetMasterID),
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
	model_master.ManagementNumber = "abc"

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
