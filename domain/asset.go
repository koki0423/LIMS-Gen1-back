package domain

import "time"

type AssetMaster struct {
	ID                   int64   `json:"id"`
	Name                 string  `json:"name"`
	ManagementCategoryID int64   `json:"management_category_id"`
	GenreID              *int64  `json:"genre_id,omitempty"`
	Manufacturer         *string `json:"manufacturer,omitempty"`
	ModelNumber          *string `json:"model_number,omitempty"`
}

// AssetはDBモデルのNULL許容フィールドをポインタで表現
type Asset struct {
	ID           int64 `json:"id"`
	ItemMasterID int64 `json:"item_master_id"`
	Quantity     int   `json:"quantity"`
	StatusID     int64 `json:"status_id"`

	SerialNumber    *string    `json:"serial_number,omitempty"`
	PurchaseDate    *time.Time `json:"purchase_date,omitempty"`
	Owner           *string    `json:"owner,omitempty"`
	Location        *string    `json:"location,omitempty"`
	DefaultLocation *string    `json:"default_location,omitempty"`
	LastCheckDate   *time.Time `json:"last_check_date,omitempty"`
	LastChecker     *string    `json:"last_checker,omitempty"`
	Notes           *string    `json:"notes,omitempty"`
}

// CreateAssetRequest は資産の新規登録リクエスト
type CreateAssetRequest struct {
	// --- マスター情報（任意） ---
	AssetMasterID        *int64  `json:"asset_master_id"`
	Name                 *string `json:"name"`
	ManagementCategoryID *int64  `json:"management_category_id"`
	GenreID              *int64  `json:"genre_id"`
	Manufacturer         *string `json:"manufacturer"`
	ModelNumber          *string `json:"model_number"`

	// --- 個別資産情報 ---
	Quantity        int     `json:"quantity" binding:"required"`
	StatusID        int64   `json:"status_id" binding:"required"`
	SerialNumber    *string `json:"serial_number"`
	PurchaseDate    *string `json:"purchase_date"` //  (形式: "YYYY-MM-DD")
	Owner           *string `json:"owner"`
	DefaultLocation *string `json:"default_location"`
	LastCheckDate   *string `json:"last_check_date"`
	LastChecker     *string `json:"last_checker"`
	Notes           *string `json:"notes"`
}

// EditAssetRequest は資産情報更新リクエスト
// CreateAssetRequestと型の一貫性が保たれる
type EditAssetRequest struct {
	AssetID       int64   `json:"asset_id"` // 主キー
	Quantity      *int    `json:"quantity"` // 更新時は任意
	SerialNumber  *string `json:"serial_number"`
	StatusID      *int64  `json:"status_id"` // 更新時は任意, IDなのでint64
	PurchaseDate  *string `json:"purchase_date"`
	Owner         *string `json:"owner"`
	Location      *string `json:"location"`
	LastCheckDate *string `json:"last_check_date"`
	LastChecker   *string `json:"last_checker"`
	Notes         *string `json:"notes"`
}

// // UpdateAssetRequest は資産情報更新リクエストの構造体 PUT /assets/edit/:id
type UpdateAssetRequest struct {
	StatusID *int    `json:"status_id"`
	Owner    *string `json:"owner"`
	Location *string `json:"location"`
	Notes    *string `json:"notes"`
}

// // CheckAssetRequest は資産点検リクエストの構造体 POST /assets/check/:id
type CheckAssetRequest struct {
	LastChecker string `json:"last_checker"` // 点検者の学籍番号
	StatusID    int    `json:"status_id"`    // 点検後の状態
	Notes       string `json:"notes"`        // 点検に関する備考
}

// // CreateLendRequest は貸出登録リクエストの構造体 POST /lend/:id
type CreateLendRequest struct {
	Borrower           string `json:"borrower"`             // 借主の学籍番号
	ExpectedReturnDate string `json:"expected_return_date"` // YYYY-MM-DD
	Notes              string `json:"notes"`
	Quantity           int    `json:"quantity"` // 通常は1
}

// // UpdateLendRequest は貸出情報更新リクエストの構造体 PUT /lend/edit/:id
type UpdateLendRequest struct {
	ExpectedReturnDate *string `json:"expected_return_date"`
	Notes              *string `json:"notes"`
}

// // CreateReturnRequest は返却登録リクエストの構造体POST /lend/return/:id
type CreateReturnRequest struct {
	Notes            string `json:"notes"`
	ReturnedQuantity int    `json:"returned_quantity"` // 通常は1
}
