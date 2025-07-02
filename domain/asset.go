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
