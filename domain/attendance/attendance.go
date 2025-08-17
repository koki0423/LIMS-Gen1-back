package attendance

import (
	"time"
)

type Attendance struct {
	ID            int       `gorm:"primaryKey" json:"id"`
	Timestamp     time.Time `gorm:"type:timestamp" json:"timestamp"` // ISO8601形式を前提
	StudentNumber string    `json:"student_number"`                  // 学籍番号
}

// レスポンスDTO
type AttendanceItem struct {
	ID        int64      `json:"id"`
	StudentID string     `json:"student_id"`
	Name      *string    `json:"name,omitempty"`
	Timestamp time.Time  `json:"timestamp"`          // ISO8601（UTC）
	Status    string     `json:"status"`
	Subject   *string    `json:"subject,omitempty"`
	Notes     *string    `json:"notes,omitempty"`
}

type AttendanceListResponse struct {
	Items   []AttendanceItem `json:"items"`
	Total   int64            `json:"total"`
	Message string           `json:"message,omitempty"`
}

type AttendanceRanking struct {
	StudentID    string  `json:"student_id"`
	AttendanceCount int     `json:"attendance_count"`
}