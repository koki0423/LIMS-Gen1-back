package disposals

import "time"

// ---- Requests ----

type CreateDisposalRequest struct {
	Quantity      uint    `json:"quantity" binding:"required"` // >0
	Reason        *string `json:"reason,omitempty"`
	ProcessedByID *string `json:"processed_by_id,omitempty"`
}

// ---- Responses ----

type DisposalResponse struct {
	DisposalULID     string    `json:"disposal_ulid"`
	ManagementNumber string    `json:"management_number"`
	Quantity         uint      `json:"quantity"`
	Reason           *string   `json:"reason,omitempty"`
	ProcessedByID    *string   `json:"processed_by_id,omitempty"`
	DisposedAt       time.Time `json:"disposed_at"`
}

// ---- List payload ----

type Page struct {
	Limit  int
	Offset int
	Order  string // "asc" or "desc"
}

type DisposalFilter struct {
	ManagementNumber *string
	ProcessedByID    *string
	From             *time.Time
	To               *time.Time
}
