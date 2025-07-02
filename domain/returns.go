package domain

type AssetReturnRequest struct {
    LendID   int64   `json:"lend_id" binding:"required"`
    Quantity int     `json:"quantity" binding:"required"`
    ReturnDate string `json:"return_date" binding:"required"`
    Notes    *string `json:"notes"`
}
