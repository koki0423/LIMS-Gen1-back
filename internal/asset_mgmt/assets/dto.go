package assets

import "time"

// ===== Requests =====

type CreateAssetMasterRequest struct {
	Name                 string  `json:"name" binding:"required"`
	ManagementCategoryID uint    `json:"management_category_id" binding:"required"`
	GenreID              uint    `json:"genre_id" binding:"required"`
	Manufacturer         string  `json:"manufacturer" binding:"required"`
	Model                *string `json:"model,omitempty"`
}

type UpdateAssetMasterRequest struct {
	Name                 *string `json:"name,omitempty"`
	ManagementCategoryID *uint   `json:"management_category_id,omitempty"`
	GenreID              *uint   `json:"genre_id,omitempty"`
	Manufacturer         *string `json:"manufacturer,omitempty"`
	Model                *string `json:"model,omitempty"`
}

type CreateAssetRequest struct {
	AssetMasterID   *uint64    `json:"asset_master_id,omitempty"`
	Serial          *string    `json:"serial,omitempty"`
	Quantity        uint       `json:"quantity"` // >=0, default 1 はDB側デフォルトでも可
	PurchasedAt     time.Time  `json:"purchased_at" binding:"required"`
	StatusID        uint       `json:"status_id" binding:"required"`
	Owner           string     `json:"owner" binding:"required"`
	DefaultLocation string     `json:"default_location" binding:"required"`
	Location        *string    `json:"location,omitempty"`
	LastCheckedAt   *time.Time `json:"last_checked_at,omitempty"`
	LastCheckedBy   *string    `json:"last_checked_by,omitempty"`
	Notes           *string    `json:"notes,omitempty"`
}

type UpdateAssetRequest struct {
	Serial          *string    `json:"serial,omitempty"`
	Quantity        *uint      `json:"quantity,omitempty"` // >=0
	PurchasedAt     *time.Time `json:"purchased_at,omitempty"`
	StatusID        *uint      `json:"status_id,omitempty"`
	Owner           *string    `json:"owner,omitempty"`
	DefaultLocation *string    `json:"default_location,omitempty"`
	Location        *string    `json:"location,omitempty"`
	LastCheckedAt   *time.Time `json:"last_checked_at,omitempty"`
	LastCheckedBy   *string    `json:"last_checked_by,omitempty"`
	Notes           *string    `json:"notes,omitempty"`
}

// ===== Responses =====

type AssetMasterResponse struct {
	AssetMasterID        uint64    `json:"asset_master_id"`
	ManagementNumber     string    `json:"management_number"`
	Name                 string    `json:"name"`
	ManagementCategoryID uint      `json:"management_category_id"`
	GenreID              uint      `json:"genre_id"`
	Manufacturer         string    `json:"manufacturer"`
	Model                *string   `json:"model,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
}

type AssetResponse struct {
	AssetID          uint64     `json:"asset_id"`
	AssetMasterID    uint64     `json:"asset_master_id"`
	ManagementNumber string     `json:"management_number"`
	Serial           *string    `json:"serial,omitempty"`
	Quantity         uint       `json:"quantity"`
	PurchasedAt      time.Time  `json:"purchased_at"`
	StatusID         uint       `json:"status_id"`
	Owner            string     `json:"owner"`
	DefaultLocation  string     `json:"default_location"`
	Location         *string    `json:"location,omitempty"`
	LastCheckedAt    *time.Time `json:"last_checked_at,omitempty"`
	LastCheckedBy    *string    `json:"last_checked_by,omitempty"`
	Notes            *string    `json:"notes,omitempty"`
}

// ===== Listing helpers =====

type Page struct {
	Limit  int
	Offset int
	Order  string // "asc" or "desc"
}

type AssetSearchQuery struct {
	ManagementNumber *string
	StatusID         *uint
	Owner            *string
	Location         *string
	PurchasedFrom    *time.Time
	PurchasedTo      *time.Time
}
