package service

import (
	"database/sql"
	"equipmentManager/domain"
	"equipmentManager/internal/database/disposal"
	model "equipmentManager/internal/database/model/tables"
	"fmt"
	time "time"
)

type DisposalService struct {
	DB *sql.DB
}

const timeFormat = "2006-01-02"

func NewDisposalService(db *sql.DB) *DisposalService {
	return &DisposalService{DB: db}
}

func (ds *DisposalService) CreateDisposal(req domain.CreateDisposalRequest, assetID int64) error {
	model_req := model.AssetsDisposal{
		AssetID:      assetID,
		Quantity:     req.Quantity,
		DisposalDate: time.Now(),
		Reason:       sql.NullString{String: req.Reason, Valid: req.Reason != ""},
		ProcessedBy:  sql.NullString{String: req.ProcessedBy, Valid: req.ProcessedBy != ""},
		IsIndividual: req.IsIndividual,
	}

	_, err := disposal.RegisterDisposal(ds.DB, model_req)
	if err != nil {
		return fmt.Errorf("廃棄登録に失敗しました: %w", err)
	}
	return nil
}

func (ds *DisposalService) GetDisposalAll() ([]domain.DisposalResponse, error) {
	allDisposals, err := disposal.FetchDisposalAll(ds.DB)
	if err != nil {
		return nil, fmt.Errorf("廃棄情報の取得に失敗しました: %w", err)
	}

	// モデルからドメインへの変換
	domainDisposals := make([]domain.DisposalResponse, len(allDisposals))
	for i, d := range allDisposals {
		domainDisposals[i] = ModelToDomainDisposal(d)
	}

	return domainDisposals, nil
}

func (ds *DisposalService) GetDisposalByID(id int64) (domain.DisposalResponse, error) {
	disposalData, err := disposal.FetchDisposalByID(ds.DB, id)
	if err != nil {
		return domain.DisposalResponse{}, fmt.Errorf("廃棄情報の取得に失敗しました: %w", err)
	}

	// モデルからドメインへの変換
	return ModelToDomainDisposal(disposalData), nil
}

func ModelToDomainDisposal(disposal model.AssetsDisposal) domain.DisposalResponse {
	return domain.DisposalResponse{
		ID:			disposal.ID,
		AssetID:     disposal.AssetID,
		Quantity:    disposal.Quantity,
		DisposalDate: disposal.DisposalDate.Format(timeFormat),
		Reason:      disposal.Reason.String,
		ProcessedBy: disposal.ProcessedBy.String,
	}
}
