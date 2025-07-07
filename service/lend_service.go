package service

import (
	"database/sql"
	"equipmentManager/utils"
	"equipmentManager/domain"
	"equipmentManager/internal/database/lends"
	model "equipmentManager/internal/database/model/tables"

	"equipmentManager/internal/database/returns"
	"fmt"
	"time"
)

type LendService struct {
	DB *sql.DB
}

func NewLendService(db *sql.DB) *LendService {
	return &LendService{DB: db}
}

// --- ドメイン変換共通関数群 ---
func convertAssetsLendModelListToDomain(src []model.AssetsLend) []domain.AssetsLend {
	result := make([]domain.AssetsLend, 0, len(src))
	for i := 0; i < len(src); i++ {
		result = append(result, toDomainAssetsLend(src[i]))
	}
	return result
}

func toDomainAssetsLend(m model.AssetsLend) domain.AssetsLend {
	return domain.AssetsLend{
		ID:                 m.ID,
		AssetID:            m.AssetID,
		Borrower:           m.Borrower,
		Quantity:           m.Quantity,
		LendDate:           m.LendDate,
		ExpectedReturnDate: utils.NullTimeToPtr(m.ExpectedReturnDate),
		ActualReturnDate:   utils.NullTimeToPtr(m.ActualReturnDate),
		Notes:              utils.NullStringToPtr(m.Notes),
	}
}

// GET /lend/all
func (e *LendService) GetLends() ([]domain.AssetsLend, error) {
	lendsData, err := lends.FetchLendsAll(e.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch lends: %w", err)
	}
	return convertAssetsLendModelListToDomain(lendsData), nil
}

// POST /lend/:id
func (e *LendService) PostLend(req domain.LendAssetRequest) (bool, error) {
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

// POST /lend/return/:id
func (e *LendService) PostReturn(req domain.ReturnAssetRequest) (bool, error) {
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

	parsedDate, err := time.Parse("2006-01-02", *req.ActualReturnDate)
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

func (e *LendService) GetLendById(id int64) ([]domain.AssetsLend, error) {
	lendDataList, err := lends.FetchLendsByAssetID(e.DB, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch lends by ID: %w", err)
	}
	return convertAssetsLendModelListToDomain(lendDataList), nil
}