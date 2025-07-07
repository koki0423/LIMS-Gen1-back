package utils

import (
	"database/sql"
	"log"
	"time"
)




//
// --- 文字列 <-> Null型 変換ユーティリティ ---
//

// NullStringToString は sql.NullString から値を取り出し、無効なら空文字を返します。
//
//	ns: 変換元の sql.NullString
func NullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// StringToNullString は *string を sql.NullString に変換します。
// nilまたは空文字の場合は無効扱いで返します。
//
//	s: 変換元の *string
func StringToNullString(s *string) sql.NullString {
	if s != nil && *s != "" {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{}
}

// DerefString は *string を値に変換します。
// nilの場合は空文字列を返します。
//
//	s: 変換元の *string
func DerefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

//
// --- 日付・時間関連ユーティリティ ---
//

// NullTimeToPtr は sql.NullTime から *time.Time への変換を行います。
// 無効な場合は nil を返します。
//
//	nt: 変換元の sql.NullTime
func NullTimeToPtr(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}

// StringToNullTime は "YYYY-MM-DD" 形式の*stringをsql.NullTimeに変換します。
// 無効またはパースエラー時は Valid=false となります。
//
//	s: 変換元の *string
func StringToNullTime(s *string) sql.NullTime {
	if s == nil || *s == "" {
		return sql.NullTime{Valid: false}
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
}

// NullTimeToString は sql.NullTime を "YYYY-MM-DD" 形式の *string に変換します。
// 無効な場合は nil を返します。
//
//	nt: 変換元の sql.NullTime
func NullTimeToString(nt sql.NullTime) *string {
	if nt.Valid {
		s := nt.Time.Format("2006-01-02")
		return &s
	}
	return nil
}

// MustParseDate は "YYYY-MM-DD" 形式の文字列を time.Time に変換します。
// パース失敗時はpanicします（主にテスト・初期化用）。
//
//	dateStr: 変換元の文字列
func MustParseDate(dateStr string) time.Time {
	const layout = "2006-01-02"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		log.Panicf("invalid date format (expected YYYY-MM-DD): %s", dateStr)
	}
	return t
}

//
// --- 数値型 <-> Null型 変換ユーティリティ ---
//

// Int64ToNullInt64 は *int64 を sql.NullInt64 に変換します。
// nilの場合は Valid=false となります。
//
//	i: 変換元の *int64
func Int64ToNullInt64(i *int64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: *i, Valid: true}
}

//
// --- Null型 <-> 文字列型 変換ユーティリティ ---
//

// NullStringToPtr は sql.NullString を *string に変換します。
// Valid が true ならポインタ、false なら nil を返します。
func NullStringToPtr(ns sql.NullString) *string {
	if ns.Valid {
		s := ns.String // ローカル変数に値コピー
		return &s
	}
	return nil
}

// NullTimeToPtrString は sql.NullTime を "YYYY-MM-DD" 形式の *string に変換します。
// Valid が false の場合は nil を返します。
func NullTimeToPtrString(nt sql.NullTime) *string {
	// 1. Validフラグをチェックします。
	//    これがfalseの場合、DBの値はNULLです。
	if !nt.Valid {
		// NULLの場合は、文字列ポインタのゼロ値であるnilを返します。
		return nil
	}

	// 2. 有効な場合、nt.Time（time.Time型）を
	//    指定したフォーマットの文字列に変換します。
	formatted := nt.Time.Format("2006-01-02")

	// 3. 生成した文字列変数のアドレス（ポインタ）を返します。
	return &formatted
}

// PtrInt64ToInt64 は *int64 を int64 に変換します。
func PtrInt64ToInt64(ptr *int64) int64 {
	// 1. ポインタがnilかどうかをチェックする
	if ptr == nil {
		// nilの場合は、デフォルト値（ここでは0）を返す
		return 0
	}
	// 2. nilでなければ、`*`で値を取り出す
	return *ptr
}

func PtrIntToInt(ptr *int) int {
	if ptr == nil {
		return 0 // nilの場合はデフォルト値0を返す
	}
	return *ptr // nilでなければポインタの値を返す
}

func PtrStringToString(ptr *string) string {
	if ptr == nil {
		return "" // nilの場合は空文字を返す
	}
	return *ptr // nilでなければポインタの値を返す
}

