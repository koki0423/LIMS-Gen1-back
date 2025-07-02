// handler/auth_handler.go
package handler

import (
	"database/sql"
	"equipmentManager/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// LoginRequest はログインリクエストのJSONボディをバインドするための構造体
type LoginRequest struct {
	StudentNumber string `json:"student_number" binding:"required"`
	Password      string `json:"password" binding:"required"`
}

// AuthHandler は依存するサービスを持つ構造体
type AuthHandler struct {
	DB      *sql.DB
	Service *service.AuthService
}

// NewAuthHandler はハンドラを初期化するファクトリ関数
func NewAuthHandler(db *sql.DB) *AuthHandler {
	service := service.NewAuthService(db)
	return &AuthHandler{
		DB:      db,
		Service: service,
	}
}

// Login はログイン処理のハンドラー
func (ah *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	// JSON→構造体バインド（バリデーション付き）
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// サービス層にログイン処理を委任
	token, err := ah.Service.Login(req.StudentNumber, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 成功レスポンス
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}
