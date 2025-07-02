// service/auth_service.go
package service

import (
	"database/sql"
	"errors"
	"time"

	"equipmentManager/internal/database/developer"
	"equipmentManager/internal/database/utils"
	"github.com/golang-jwt/jwt/v5"
)

// AuthService は認証処理を担当するサービス
type AuthService struct {
	DB *sql.DB
}

// NewAuthService は依存注入のためのコンストラクタ
func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{DB: db}
}

var jwtSecret = []byte("my_super_secret_key")

// Login はログイン認証＋JWT発行
func (s *AuthService) Login(studentNumber string, password string) (string, error) {
	dev, err := developer.GetDeveloperByStudentNumber(s.DB, studentNumber)
	if err != nil {
		return "", err
	}
	if dev == nil {
		return "", errors.New("invalid student number or password")
	}

	if !utils.CheckPasswordHash(password, dev.PasswordHash) {
		return "", errors.New("invalid student number or password")
	}

	claims := jwt.MapClaims{
		"dev_id": dev.ID,
		"exp":    time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
