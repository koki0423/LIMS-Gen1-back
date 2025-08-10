package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword は平文のパスワードをbcryptでハッシュ化します。
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPasswordHash は入力パスワードとハッシュを比較して一致するか判定します。
func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
