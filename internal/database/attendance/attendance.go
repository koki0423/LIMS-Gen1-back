package attendance

import (
	"context"
	"database/sql"
	"time"

	"equipmentManager/domain/attendance"
)

type Repository struct {
	DB *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Create(ctx context.Context, a *attendance.Attendance) error {
	// Timestampは呼び出し側でtime.Now()を入れる前提
	const q = `INSERT INTO attendances (timestamp, student_number) VALUES (?, ?)`
	res, err := r.DB.ExecContext(ctx, q, a.Timestamp, a.StudentNumber)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	a.ID = int(id)
	return nil
}

func (r *Repository) FindAll(ctx context.Context) ([]attendance.Attendance, error) {
	const q = `SELECT id, timestamp, student_number FROM attendances ORDER BY id DESC`
	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]attendance.Attendance, 0, 128)
	for rows.Next() {
		var a attendance.Attendance
		if err := rows.Scan(&a.ID, &a.Timestamp, &a.StudentNumber); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) FindByID(ctx context.Context, id int64) (*attendance.Attendance, error) {
	const q = `SELECT id, timestamp, student_number FROM attendances WHERE id = ?`
	var a attendance.Attendance
	if err := r.DB.QueryRowContext(ctx, q, id).Scan(&a.ID, &a.Timestamp, &a.StudentNumber); err != nil {
		return nil, err
	}
	return &a, nil
}

// JSTの“今日”を返す。studentNumberフィルタが空なら全件。
func (r *Repository) FindTodayJST(ctx context.Context, studentNumber string) ([]attendance.Attendance, error) {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, err
	}
	now := time.Now().In(jst)
	y, m, d := now.Date()
	start := time.Date(y, m, d, 0, 0, 0, 0, jst)
	end := start.Add(24 * time.Hour)

	q := `SELECT id, timestamp, student_number FROM attendances WHERE timestamp >= ? AND timestamp < ?`
	args := []any{start, end}
	if studentNumber != "" {
		q += ` AND student_number = ?`
		args = append(args, studentNumber)
	}
	q += ` ORDER BY timestamp ASC, id ASC`

	rows, err := r.DB.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]attendance.Attendance, 0, 64)
	for rows.Next() {
		var a attendance.Attendance
		if err := rows.Scan(&a.ID, &a.Timestamp, &a.StudentNumber); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}
