package printLabels

import (
	"context"
	"errors"
)

type Service struct {
}

func NewService() *Service { return &Service{} }

func (s *Service) PrintLabels(ctx context.Context, input PrintRequest) (*PrintResponse, error) {
	rows := []PrintRow{{Checked: input.Label.Checked,
		ColB: input.Label.ColB,
		ColC: input.Label.ColC,
		ColD: input.Label.ColD,
		ColE: input.Label.ColE,
	}}

	params := PrintParams{
		TemplateWidthMM:     input.Width,
		BarcodeType:         input.Type,
		UseHalfcut:          input.Config.UseHalfcut,
		ConfirmTapeWidthDlg: input.Config.ConfirmTapeWidth,
		EnablePrintLog:      input.Config.EnablePrintLog,
		PrinterName:         "",
	}

	if err := PrintLabels(rows, params); err != nil {
		// store.goから返る各エラーをここでハンドリングする
		if errors.Is(err, ErrTapeSizeNotMatched) {
			// テープ幅の不一致は「クライアントからの要求とサーバーの状態の競合」として409 Conflictを返す
			return nil, ErrConflict(err.Error())
		}
		if errors.Is(err, ErrTemplateNotFound) {
			// テンプレートが見つからないのは404 Not Found
			return nil, ErrNotFound(err.Error())
		}
		if errors.Is(err, ErrNoPrintableSelected) {
			// 印刷対象が選択されていないのは「クライアントのリクエストが不正」として400 Bad Request
			return nil, ErrInvalid(err.Error())
		}
		if errors.Is(err, ErrSPC10NotFound) {
			// SPC10.exeが見つからないのはサーバー内部の問題として500 Internal
			// ただし、メッセージは具体的で分かりやすいものにする
			return nil, ErrInternal(err.Error())
		}

		// その他の予期せぬエラーも500 Internal
		return nil, ErrInternal(err.Error())
	}

	// 成功時は空のレスポンスとnil errorを返す
	return &PrintResponse{}, nil
}
