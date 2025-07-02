package domain

type StudentCardInfo struct {
	StudentNumber string `json:"student_number"`
	// 今後ここにフィールド追加する
	// Name         string `json:"name,omitempty"`
	// Expiry       string `json:"expiry,omitempty"`
}
