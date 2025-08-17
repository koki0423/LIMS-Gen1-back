package service

import (
	"context"
	"database/sql"
	"time"
	"strings"

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

func (s *AttendanceService) FindByStudentAndRange(
	ctx context.Context,
	studentID string,
	startUTC, endUTC time.Time,
	sort string, page, size int,
) ([]attendance.AttendanceItem, int64, error) {

	orderBy := "timestamp DESC"
	switch strings.ToLower(sort) {
	case "date:asc", "timestamp:asc":
		orderBy = "timestamp ASC"
	case "date:desc", "timestamp:desc":
		orderBy = "timestamp DESC"
	// 将来 student_id ソートなどを足したければここで許可リスト化
	}

	offset := (page - 1) * size
	return s.repo.QueryByStudentAndRange(ctx, studentID, startUTC, endUTC, orderBy, size, offset)
}