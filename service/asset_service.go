package service

import (
	"database/sql"
	"equipmentManager/domain"
	"equipmentManager/internal/database/assets"
	"equipmentManager/internal/database/lends"
	model "equipmentManager/internal/database/model/tables"
	"equipmentManager/internal/database/returns"
	"equipmentManager/utils"
	"fmt"
	// "log"

	"time"
)

type AssetService struct {
	DB *sql.DB
}

func NewAssetService(db *sql.DB) *AssetService {
	return &AssetService{DB: db}
}

// --- ドメイン変換共通関数群 ---
func toDomainAsset(m model.Asset) domain.Asset {
	return domain.Asset{
		ID:            m.ID,
		ItemMasterID:  m.ItemMasterID,
		Quantity:      m.Quantity,
		SerialNumber:  utils.NullStringToString(m.SerialNumber),
		StatusID:      m.StatusID,
		PurchaseDate:  utils.NullTimeToPtr(m.PurchaseDate),
		Owner:         utils.NullStringToString(m.Owner),
		Location:      utils.NullStringToString(m.Location),
		LastCheckDate: utils.NullTimeToPtr(m.LastCheckDate),
		LastChecker:   utils.NullStringToString(m.LastChecker),
		Notes:         utils.NullStringToString(m.Notes),
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

func toDomainAssetsLend(m model.AssetsLend) domain.AssetsLend {
	return domain.AssetsLend{
		ID:                 m.ID,
		AssetID:            m.AssetID,
		Borrower:           m.Borrower,
		Quantity:           m.Quantity,
		LendDate:           m.LendDate.Format("2006-01-02"),
		ExpectedReturnDate: utils.NullTimeToPtrString(m.ExpectedReturnDate),
		ActualReturnDate:   utils.NullTimeToPtrString(m.ActualReturnDate),
		Notes:              utils.NullStringToPtr(m.Notes),
	}
}

// --- 変換処理 ---
func convertAssetsModelListToDomain(src []model.Asset) []domain.Asset {
	result := make([]domain.Asset, 0, len(src))
	for i := 0; i < len(src); i++ {
		result = append(result, toDomainAsset(src[i]))
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

func convertAssetsLendModelListToDomain(src []model.AssetsLend) []domain.AssetsLend {
	result := make([]domain.AssetsLend, 0, len(src))
	for i := 0; i < len(src); i++ {
		result = append(result, toDomainAssetsLend(src[i]))
	}
	return result
}

func convertToDomainAssetModel(domainAsset domain.CreateAssetRequest) model.AllDataOfAsset {
	return model.AllDataOfAsset{
		Asset: model.Asset{
			ItemMasterID:  domainAsset.Asset.ItemMasterID,
			Quantity:      domainAsset.Asset.Quantity,
			SerialNumber:  utils.StringToNullString(domainAsset.Asset.SerialNumber),
			StatusID:      domainAsset.Asset.StatusID,
			PurchaseDate:  utils.StringToNullTime(domainAsset.Asset.PurchaseDate),
			Owner:         utils.StringToNullString(domainAsset.Asset.Owner),
			Location:      utils.StringToNullString(domainAsset.Asset.Location),
			LastCheckDate: utils.StringToNullTime(domainAsset.Asset.LastCheckDate),
			LastChecker:   utils.StringToNullString(domainAsset.Asset.LastChecker),
			Notes:         utils.StringToNullString(domainAsset.Asset.Notes),
		},
		Master: model.AssetsMaster{
			Name:                 domainAsset.AssetMaster.Name,
			ManagementCategoryID: domainAsset.AssetMaster.ManagementCategoryID,
			GenreID:              utils.Int64ToNullInt64(domainAsset.AssetMaster.GenreID),
			Manufacturer:         utils.StringToNullString(domainAsset.AssetMaster.Manufacturer),
			ModelNumber:          utils.StringToNullString(domainAsset.AssetMaster.ModelNumber),
		},
	}
}

func convertToDomainEditAsset(editedAsset domain.EditAssetRequest) model.Asset {
	return model.Asset{
		ID:            editedAsset.AssetID,
		Quantity:      editedAsset.Quantity,
		SerialNumber:  utils.StringToNullString(editedAsset.SerialNumber),
		StatusID:      editedAsset.StatusID,
		PurchaseDate:  utils.StringToNullTime(editedAsset.PurchaseDate),
		Owner:         utils.StringToNullString(editedAsset.Owner),
		Location:      utils.StringToNullString(editedAsset.Location),
		LastCheckDate: utils.StringToNullTime(editedAsset.LastCheckDate),
		LastChecker:   utils.StringToNullString(editedAsset.LastChecker),
		Notes:         utils.StringToNullString(editedAsset.Notes),
	}
}

// --- サービス ---
func (e *AssetService) CreateAssetWithMaster(newAsset domain.CreateAssetRequest) (bool, error) {
	data := convertToDomainAssetModel(newAsset)
	switch data.Master.ManagementCategoryID {
	case 1:
		return assets.RegistrationEquipmentForIndivisual(e.DB, data)
	case 2:
		return assets.RegistrationEquipmentForGeneral(e.DB, data)
	default:
		return false, fmt.Errorf("未知の管理カテゴリID: %d", data.Master.ManagementCategoryID)
	}
}


func (e *AssetService) GetAssetAll() ([]domain.Asset, error) {
	dbAssets, err := assets.FetchAssetsAll(e.DB)
	if err != nil {
		return nil, err
	}
	return convertAssetsModelListToDomain(dbAssets), nil
}

func (e *AssetService) GetAssetMasterAll() ([]domain.AssetMaster, error) {
	dbAssetMasters, err := assets.FetchAllAssetsMaster(e.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch asset masters: %w", err)
	}
	return convertAssetMasterModelListToDomain(dbAssetMasters), nil
}

func (e *AssetService) GetAssetMasterById(id int) (domain.AssetMaster, error) {
	dbAssetMaster, err := assets.FetchAssetsMasterByID(e.DB, int64(id))
	if err != nil {
		return domain.AssetMaster{}, fmt.Errorf("failed to fetch asset master by ID: %w", err)
	}
	return toDomainAssetMaster(*dbAssetMaster), nil
}

func (e *AssetService) DeleteAssetMasterById(id int) (bool, error) {
	return assets.DeleteAssetMasterByID(e.DB, int64(id))
}

func (e *AssetService) GetAssetById(id int) (domain.Asset, error) {
	dbAsset, err := assets.FetchAssetsByID(e.DB, int64(id))
	if err != nil {
		return domain.Asset{}, err
	}
	return toDomainAsset(*dbAsset), nil
}

func (e *AssetService) PutAssetsEdit(editedAsset domain.EditAssetRequest) (bool, error) {
	return assets.UpdateAsset(e.DB, convertToDomainEditAsset(editedAsset))
}

func (e *AssetService) PostLend(req domain.LendAssetRequest) (bool, error) {
	modelLend := model.AssetsLend{
		AssetID:            req.AssetID,
		Borrower:           req.Borrower,
		Quantity:           req.Quantity,
		LendDate:           utils.MustParseDate(req.LendDate),
		ExpectedReturnDate: utils.StringToNullTime(req.ExpectedReturnDate),
		ActualReturnDate:   utils.StringToNullTime(req.ActualReturnDate),
		Notes:              utils.StringToNullString(req.Notes),
	}

	//返り値はbool, int64, errorだがint64の主キーは何に使うか決めてないので _ にしてある
	status, _, err := lends.RegisterLend(e.DB, modelLend)
	if err != nil {
		return false, fmt.Errorf("failed to register lend: %w", err)
	}

	return status, nil
}

func (e *AssetService) PostAssetReturnHistory(req domain.AssetReturnRequest) (bool, error) {
	tx, err := e.DB.Begin()
	if err != nil {
		return false, fmt.Errorf("トランザクション開始失敗: %w", err)
	}
	// panic対策としてrollbackを保証
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
		}
	}()

	parsedDate, err := time.Parse("2006-01-02", req.ReturnDate)
	if err != nil {
		_ = tx.Rollback()
		return false, fmt.Errorf("日付フォーマット不正: %w", err)
	}

	parsedNotes := utils.StringToNullString(req.Notes)
	returnData := model.AssetReturn{
		LendID:           req.LendID,
		ReturnedQuantity: req.Quantity,
		ReturnDate:       parsedDate,
		Notes:            parsedNotes,
	}

	// --- asset_lendsの更新 ---
	_, err = lends.UpdateReturnDateForAssetlist(tx, req.LendID, parsedDate, parsedNotes)

	if err != nil {
		_ = tx.Rollback()
		return false, fmt.Errorf("貸出情報の更新に失敗: %w", err)
	}

	// --- asset_returnsの新規登録 ---
	_, err = returns.RegisterAssetReturn(tx, returnData)
	if err != nil {
		_ = tx.Rollback()
		return false, fmt.Errorf("返却履歴の登録に失敗: %w", err)
	}

	// --- commit ---
	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("コミットに失敗: %w", err)
	}

	return true, nil
}

func (e *AssetService) GetLends() ([]domain.AssetsLend, error) {
	lendsData, err := lends.FetchLendsAll(e.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch lends: %w", err)
	}
	return convertAssetsLendModelListToDomain(lendsData), nil
}
