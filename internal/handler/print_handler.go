package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"
	"context"

	"github.com/gin-gonic/gin"

	"equipmentManager/service"
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
func (ph *PrintHandler) PostPrintHandler(c *gin.Context) {
	var input struct {
		StudentNumber string `json:"student_number"`
	}

	if err := c.ShouldBindJSON(&input); err != nil || input.StudentNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なJSON（必須）"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	a, err := ph.Service.Create(ctx, input.StudentNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "登録失敗"})
		return
	}
	c.JSON(http.StatusOK, a)
}
