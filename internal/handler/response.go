package handler

import "equipmentManager/domain"

/* --- 汎用レスポンス (前回定義したものと同じ) --- */
// SuccessResponse は汎用的な成功レスポンスです
type SuccessResponse struct {
	Message string `json:"message" example:"処理に成功しました"`
}

// ErrorResponse は汎用的なエラーレスポンスです
type ErrorResponse struct {
	Error string `json:"error" example:"不正なリクエストです"`
}

/* --- Asset関連のレスポンス --- */
// CreateAssetResponse は資産作成成功時のレスポンスです
type CreateAssetResponse struct {
	Message       string `json:"message" example:"Asset created successfully"`
	AssetMasterID int64  `json:"asset_master_id" example:"101"`
}

// AssetListResponse は資産リストのレスポンスです
type AssetListResponse struct {
	Message string         `json:"message" example:"Assets fetched successfully"`
	Assets  []domain.Asset `json:"assets"`
}

// AssetResponse は資産単体のレスポンスです
type AssetResponse struct {
	Message string       `json:"message" example:"Asset fetched successfully"`
	Asset   domain.Asset `json:"asset"`
}

// AssetMasterListResponse は資産マスターリストのレスポンスです
type AssetMasterListResponse struct {
	Message string               `json:"message" example:"Asset masters fetched successfully"`
	Masters []domain.AssetMaster `json:"masters"`
}

// AssetMasterResponse は資産マスター単体のレスポンスです
type AssetMasterResponse struct {
	Message string             `json:"message" example:"Asset master fetched successfully"`
	Master  domain.AssetMaster `json:"master"`
}

/* --- Lend関連のレスポンス --- */
// 貸出情報リストのレスポンス（AssetsLendの配列を内包）
type LendListResponse struct {
	Lends []domain.AssetsLend `json:"lends"`
}

// 貸出情報単体のレスポンス
type LendResponse struct {
	Message string            `json:"message" example:"Fetch completed"`
	Lend    domain.AssetsLend `json:"lend"`
}

// 名称付き貸出情報リストのレスポンス（LendingDetailの配列を内包）
type LendingDetailListResponse struct {
	Lends []domain.LendingDetail `json:"lends"`
}

/* --- Disposal (廃棄) 関連レスポンス --- */
// DisposalListResponse は廃棄情報リストのレスポンスです
type DisposalListResponse struct {
	Message   string                   `json:"message" example:"Disposals fetched successfully"`
	Disposals []domain.DisposalResponse `json:"disposals"`
}

// DisposalResponseWrapper は廃棄情報単体のレスポンスです
type DisposalResponseWrapper struct {
	Message  string                  `json:"message" example:"Disposal fetched successfully"`
	Disposal domain.DisposalResponse `json:"disposal"`
}

/* --- System関連のレスポンス --- */
type PingResponse struct {
	Message string `json:"message" example:"pong"`
}

type PingErrorResponse struct {
	Error string `json:"error" example:"Internal Server Error"`
}