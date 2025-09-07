package lends

import "time"

// ---- Requests ----

type CreateLendRequest struct {
	Quantity   uint    `json:"quantity" binding:"required"`      // >0 をサービス層で検証
	BorrowerID string  `json:"borrower_id" binding:"required"`   // 借受者
	DueOn      *string `json:"due_on,omitempty"`                 // "YYYY-MM-DD"
	LentByID   *string `json:"lent_by_id,omitempty"`
	Note       *string `json:"note,omitempty"`
}

type CreateReturnRequest struct {
	Quantity       uint    `json:"quantity" binding:"required"` // >0
	ProcessedByID  *string `json:"processed_by_id,omitempty"`
	Note           *string `json:"note,omitempty"`
}

// ---- Responses ----

type LendResponse struct {
	LendULID            string     `json:"lend_ulid"`
	AssetMasterID       uint64     `json:"asset_master_id"`
	ManagementNumber    string     `json:"management_number"`
	Quantity            uint       `json:"quantity"`
	BorrowerID          string     `json:"borrower_id"`
	DueOn               *string    `json:"due_on,omitempty"`
	LentByID            *string    `json:"lent_by_id,omitempty"`
	LentAt              time.Time  `json:"lent_at"`
	ReturnedQuantity    uint       `json:"returned_quantity"`
	OutstandingQuantity uint       `json:"outstanding_quantity"`
	Note                *string    `json:"note,omitempty"`
}

type ReturnResponse struct {
	ReturnULID     string    `json:"return_ulid"`
	LendULID       string    `json:"lend_ulid"`
	Quantity       uint      `json:"quantity"`
	ProcessedByID  *string   `json:"processed_by_id,omitempty"`
	ReturnedAt     time.Time `json:"returned_at"`
	Note           *string   `json:"note,omitempty"`
}

// ---- List payload ----

type Page struct {
	Limit  int
	Offset int
	Order  string // "asc" or "desc"
}

type LendFilter struct {
	ManagementNumber *string
	BorrowerID       *string
	From             *time.Time
	To               *time.Time
	OnlyOutstanding  bool
}
