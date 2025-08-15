package handler

import (
	"database/sql"
	"net/http"
	"time"
	"context"

	"github.com/gin-gonic/gin"

	"equipmentManager/service"
	"equipmentManager/domain"
)

type PrintHandler struct {
	DB      *sql.DB
	Service *service.PrintService
}

func NewPrintHandler(db *sql.DB) *PrintHandler {
	return &PrintHandler{
		DB:      db,
		Service: service.NewPrintService(db),
	}
}

// POST /print
func (ph *PrintHandler) PostPrinthandler(c *gin.Context) {
	var inputDomain domain.PrintRequest
	if err := c.ShouldBindJSON(&inputDomain); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なJSON"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	_, err := ph.Service.Create(ctx, inputDomain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "印刷エラー"})
		return
	}
	c.JSON(http.StatusOK, "印刷リクエストが正常に処理されました")
}
