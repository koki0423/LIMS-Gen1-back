package domain

// User はユーザーのデータ構造を表します
type User struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"` // パスワードはJSONレスポンスに含めないように'-'を指定
}