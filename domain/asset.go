package domain

import "time"

type AllDataOfAsset struct {
	Master AssetMaster
	Asset  Asset
}

type AssetMaster struct {
	ID                   int64   `json:"id"`
	Name                 string  `json:"name"`
	ManagementCategoryID int64   `json:"management_category_id"`
	GenreID              *int64  `json:"genre_id,omitempty"`
	Manufacturer         *string `json:"manufacturer,omitempty"`
	ModelNumber          *string `json:"model_number,omitempty"`
}

// assets テーブル（個別資産）
type Asset struct {
	ID            int64      `json:"id"`
	ItemMasterID  int64      `json:"item_master_id"`
	Quantity      int        `json:"quantity"`
	SerialNumber  string     `json:"serial_number,omitempty"` // omitemptyは、値が空ならJSONに含めないというGoのタグ
	StatusID      int64      `json:"status_id"`
	PurchaseDate  *time.Time `json:"purchase_date,omitempty"`
	Owner         string     `json:"owner,omitempty"`
	Location      string     `json:"location,omitempty"`
	LastCheckDate *time.Time `json:"last_check_date,omitempty"`
	LastChecker   string     `json:"last_checker,omitempty"`
	Notes         string     `json:"notes,omitempty"`
}

type CreateAssetRequest struct {
	// assets_masters テーブル相当（同時登録の場合）
	AssetMaster struct {
		Name                 string  `json:"name" binding:"required"`
		ManagementCategoryID int64   `json:"management_category_id" binding:"required"`
		GenreID              *int64  `json:"genre_id"`     // NULL許容
		Manufacturer         *string `json:"manufacturer"` // NULL許容
		ModelNumber          *string `json:"model_number"` // NULL許容
	} `json:"asset_master" binding:"required"`

	// assets テーブル相当（個別資産）
	Asset struct {
		ItemMasterID  int64   `json:"item_master_id"` //マスタ登録後に判明するのでフロントからはNULL許容
		Quantity      int     `json:"quantity" binding:"required"`
		SerialNumber  *string `json:"serial_number"`
		StatusID      int64   `json:"status_id" binding:"required"`
		PurchaseDate  *string `json:"purchase_date"` // "2024-01-01" 形式
		Owner         *string `json:"owner"`
		Location      *string `json:"location"`
		LastCheckDate *string `json:"last_check_date"` // "2024-01-01"
		LastChecker   *string `json:"last_checker"`
		Notes         *string `json:"notes"`
	} `json:"asset" binding:"required"`
}

type EditAssetRequest struct {
	AssetID       int64   `json:"asset_id"` //主キー
	Quantity      int     `json:"quantity" binding:"required"`
	SerialNumber  *string `json:"serial_number"`
	StatusID      int64   `json:"status_id" binding:"required"`
	PurchaseDate  *string `json:"purchase_date"` // "2024-01-01" 形式
	Owner         *string `json:"owner"`
	Location      *string `json:"location"`
	LastCheckDate *string `json:"last_check_date"` // "2024-01-01"
	LastChecker   *string `json:"last_checker"`
	Notes         *string `json:"notes"`
}

/*将来的にこっちへ移行*/
// // CreateAssetRequest は資産の新規登録リクエストの構造体 POST /assets
// type CreateAssetRequest struct {
// 	// 既存の備品マスタID (これを指定した場合、以下のマスタ情報は無視される)
// 	AssetMasterID *int `json:"asset_master_id"`
	
// 	// --- 新規マスタ登録用の情報 ---
// 	Name                 *string `json:"name"`
// 	ManagementCategoryID *int    `json:"management_category_id"`
// 	GenreID              *int    `json:"genre_id"`
// 	Manufacturer         *string `json:"manufacturer"`
// 	ModelNumber          *string `json:"model_number"`
	
// 	// --- 個別資産の情報 ---
// 	SerialNumber string `json:"serial_number"`
// 	StatusID     int    `json:"status_id"`
// 	PurchaseDate string `json:"purchase_date"` // YYYY-MM-DD
// 	Owner        string `json:"owner"`
// 	Location     string `json:"location"`
// 	Notes        string `json:"notes"`
// }

// // UpdateAssetRequest は資産情報更新リクエストの構造体 PUT /assets/edit/:id
// type UpdateAssetRequest struct {
// 	StatusID *int    `json:"status_id"`
// 	Owner    *string `json:"owner"`
// 	Location *string `json:"location"`
// 	Notes    *string `json:"notes"`
// }

// // CheckAssetRequest は資産点検リクエストの構造体 POST /assets/check/:id
// type CheckAssetRequest struct {
// 	LastChecker string `json:"last_checker"` // 点検者の学籍番号
// 	StatusID    int    `json:"status_id"`    // 点検後の状態
// 	Notes       string `json:"notes"`        // 点検に関する備考
// }

// // CreateLendRequest は貸出登録リクエストの構造体 POST /lend/:id
// type CreateLendRequest struct {
// 	Borrower           string `json:"borrower"` // 借主の学籍番号
// 	ExpectedReturnDate string `json:"expected_return_date"` // YYYY-MM-DD
// 	Notes              string `json:"notes"`
// 	Quantity           int    `json:"quantity"` // 通常は1
// }

// // UpdateLendRequest は貸出情報更新リクエストの構造体 PUT /lend/edit/:id
// type UpdateLendRequest struct {
// 	ExpectedReturnDate *string `json:"expected_return_date"`
// 	Notes              *string `json:"notes"`
// }

// // CreateReturnRequest は返却登録リクエストの構造体POST /lend/return/:id
// type CreateReturnRequest struct {
// 	Notes            string `json:"notes"`
// 	ReturnedQuantity int    `json:"returned_quantity"` // 通常は1
// }

// // CreateDisposalRequest は廃棄登録リクエストの構造体POST /disposal/:id
// type CreateDisposalRequest struct {
// 	Reason      string `json:"reason"`
// 	ProcessedBy string `json:"processed_by"` // 処理者の学籍番号
// 	Quantity    int    `json:"quantity"`     // 通常は1
// }