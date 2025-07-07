package domain
import "time"


// AssetsLend は資産の貸出履歴を表すドメインモデルです。
type AssetsLend struct {
	ID                 int64      `json:"id"`                   // 主キー
	AssetID            int64      `json:"asset_id"`             // 貸出対象資産ID
	Borrower           string     `json:"borrower"`             // 借用者
	Quantity           int        `json:"quantity"`             // 貸出数
	LendDate           time.Time     `json:"lend_date"`            // 貸出日 ("YYYY-MM-DD")
	ExpectedReturnDate *time.Time    `json:"expected_return_date"` // 返却予定日 ("YYYY-MM-DD" or null)
	ActualReturnDate   *time.Time    `json:"actual_return_date"`   // 実際の返却日 ("YYYY-MM-DD" or null)
	Notes              *string    `json:"notes"`                // 備考 (null可)
}

// LendAssetRequest は貸出登録リクエスト用の構造体
type LendAssetRequest struct {
	AssetID            int64   `json:"asset_id" binding:"required"`
	Borrower           string  `json:"borrower" binding:"required"`
	Quantity           int     `json:"quantity" binding:"required"`
	LendDate           string  `json:"lend_date" binding:"required"` // "YYYY-MM-DD" 形式
	ExpectedReturnDate *string `json:"expected_return_date"`         // NULL 可
	ActualReturnDate   *string `json:"actual_return_date"`           // NULL 可（通常登録時はNULL）
	Notes              *string `json:"notes"`                        // NULL 可
}

type ReturnAssetRequest struct {
	LendID           int64   `json:"lend_id" binding:"required"` // 貸出ID
	Quantity		int     `json:"quantity" binding:"required"` // 返却数量
	ActualReturnDate *string `json:"actual_return_date" binding:"required"`
	Notes            *string `json:"notes"` // NULL 可
}
