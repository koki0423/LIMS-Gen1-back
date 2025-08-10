package service

import (
	"context"
	"database/sql"
	"time"

	"equipmentManager/domain/attendance"
	attendance_repo "equipmentManager/internal/database/attendance"
)

type AttendanceService struct {
	repo *attendance_repo.Repository
}

func NewAttendanceService(db *sql.DB) *AttendanceService {
	return &AttendanceService{repo: attendance_repo.New(db)}
}

func (s *AttendanceService) Create(ctx context.Context, sn string) (*attendance.Attendance, error) {
	a := &attendance.Attendance{
		StudentNumber: sn,
		Timestamp:     time.Now(), // サーバ時刻
	}
	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *AttendanceService) ListAll(ctx context.Context) ([]attendance.Attendance, error) {
	return s.repo.FindAll(ctx)
}

func (s *AttendanceService) GetByID(ctx context.Context, id int64) (*attendance.Attendance, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *AttendanceService) ListTodayJST(ctx context.Context, studentNumber string) ([]attendance.Attendance, error) {
	return s.repo.FindTodayJST(ctx, studentNumber)
}
