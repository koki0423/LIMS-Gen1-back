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

func (r *Repository) GetRanking(ctx context.Context) ([]attendance.AttendanceRanking, error) {
	// JSTの今月 [start, end) を計算
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(jst)
	startJST := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, jst)
	endJST := startJST.AddDate(0, 1, 0)

	// UTC保存なら UTC に変換して渡す（JST保存なら startJST, endJST をそのまま使う）
	start := startJST.UTC()
	end := endJST.UTC()

	const query = "SELECT student_number, COUNT(*) AS cnt " +
		"FROM attendances " +
		"WHERE `timestamp` >= ? AND `timestamp` < ? " +
		"GROUP BY student_number " +
		"ORDER BY cnt DESC, student_number " +
		"LIMIT 5"

	rows, err := r.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rankings := make([]attendance.AttendanceRanking, 0, 5)
	for rows.Next() {
		var rec attendance.AttendanceRanking
		// AttendanceRanking は ↓ みたいに int64 推奨
		// type AttendanceRanking struct {
		//   StudentID       string `json:"student_number"`
		//   AttendanceCount int64  `json:"count"`
		// }
		if err := rows.Scan(&rec.StudentID, &rec.AttendanceCount); err != nil {
			return nil, err
		}
		rankings = append(rankings, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return rankings, nil
}

func (r *Repository) QueryByStudentAndRange(
	ctx context.Context,
	studentID string,
	startUTC, endUTC time.Time,
	orderBy string,
	limit, offset int,
) ([]attendance.AttendanceItem, int64, error) {

	// 件数
	var total int64
	countSQL := `
	  SELECT COUNT(*)
	    FROM attendance
	   WHERE student_id = ?
	     AND timestamp >= ?
	     AND timestamp <  ?
	`
	if err := r.DB.QueryRowContext(ctx, countSQL, studentID, startUTC, endUTC).Scan(&total); err != nil {
		return nil, 0, err
	}

	// 本体
	q := `
	  SELECT id, student_id, name, timestamp, status, subject, notes
	    FROM attendance
	   WHERE student_id = ?
	     AND timestamp >= ?
	     AND timestamp <  ?
	   ORDER BY ` + orderBy + `
	   LIMIT ? OFFSET ?
	`

	rows, err := r.DB.QueryContext(ctx, q, studentID, startUTC, endUTC, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]attendance.AttendanceItem, 0, limit)
	for rows.Next() {
		var (
			id        int64
			sid       string
			nameNS    sql.NullString
			ts        time.Time
			status    string
			subjectNS sql.NullString
			notesNS   sql.NullString
		)
		if err := rows.Scan(&id, &sid, &nameNS, &ts, &status, &subjectNS, &notesNS); err != nil {
			return nil, 0, err
		}
		items = append(items, attendance.AttendanceItem{
			ID:        id,
			StudentID: sid,
			Name:      nsPtr(nameNS),
			Timestamp: ts.UTC(), // 返却はUTC（FEでローカル表示推奨）
			Status:    status,
			Subject:   nsPtr(subjectNS),
			Notes:     nsPtr(notesNS),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func nsPtr(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	s := v.String
	return &s
}
