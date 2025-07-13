package domain

/*要修正*/
// CreateDisposalRequest は廃棄登録リクエストの構造体POST /disposal/:id
type CreateDisposalRequest struct {
	Reason      string `json:"reason"`
	ProcessedBy string `json:"processed_by"` // 処理者の学籍番号
	Quantity    int    `json:"quantity"`     
	IsIndividual bool   `json:"is_individual"` // 個別管理かどうか
}

type DisposalResponse struct {
	ID           int64  `json:"id"`           // 主キー
	AssetID      int64  `json:"asset_id"`      // 対象資産 ID
	Quantity     int    `json:"quantity"`     // 廃棄数量
	DisposalDate string `json:"disposal_date"` // 廃棄日 (YYYY-MM-DD)
	Reason       string `json:"reason"`       // 廃棄理由
	ProcessedBy  string `json:"processed_by"` // 処理担当者
}