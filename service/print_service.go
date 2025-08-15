package service

import (
	"context"
	"database/sql"

	"equipmentManager/domain"
	"equipmentManager/internal/print"
	"equipmentManager/internal/print/model"
)

type PrintService struct {
	DB *sql.DB
}

func NewPrintService(db *sql.DB) *PrintService {
	return &PrintService{DB: db}
}

func (ps *PrintService) Create(ctx context.Context, inputDomain domain.PrintRequest) (bool, error) {
	rows := []model.PrintRow{{Checked: inputDomain.Label.Checked,
		ColB: inputDomain.Label.ColB,
		ColC: inputDomain.Label.ColC,
		ColD: inputDomain.Label.ColD,
		ColE: inputDomain.Label.ColE,
	}}

	params := model.PrintParams{
		TemplateWidthMM:     inputDomain.Width,
		BarcodeType:         inputDomain.Type,
		UseHalfcut:          inputDomain.Config.UseHalfcut,
		ConfirmTapeWidthDlg: inputDomain.Config.ConfirmTapeWidth,
		EnablePrintLog:      inputDomain.Config.EnablePrintLog,
		PrinterName:         "",
	}

	if err := print.PrintLabels(rows, params); err != nil {
		if err == print.ErrTemplateNotFound {
			return false, err
		}
		return false, err
	}

	return true, nil
}
