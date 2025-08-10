package attendance

import (
	"time"
)

type Attendance struct {
	ID            int       `gorm:"primaryKey" json:"id"`
	Timestamp     time.Time `gorm:"type:timestamp" json:"timestamp"` // ISO8601形式を前提
	StudentNumber string    `json:"student_number"`                  // 学籍番号
}
