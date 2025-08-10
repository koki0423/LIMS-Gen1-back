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

type AttendanceHandler struct {
	DB      *sql.DB
	Service *service.AttendanceService
}

func NewAttendanceHandler(db *sql.DB) *AttendanceHandler {
	return &AttendanceHandler{
		DB:      db,
		Service: service.NewAttendanceService(db),
	}
}

// POST /attendance
func (ath *AttendanceHandler) PostAttendanceHandler(c *gin.Context) {
	var input struct {
		StudentNumber string `json:"student_number"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.StudentNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なJSON（student_number必須）"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	a, err := ath.Service.Create(ctx, input.StudentNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "登録失敗"})
		return
	}
	c.JSON(http.StatusOK, a)
}

// GET /attendance/all
func (ath *AttendanceHandler) GetAttendanceAllHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	list, err := ath.Service.ListAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取得失敗"})
		return
	}
	c.JSON(http.StatusOK, list)
}

// GET /attendance/:id
func (ath *AttendanceHandler) GetAttendanceByIdHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不正なID"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	a, err := ath.Service.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "見つかりません"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取得失敗"})
		return
	}
	c.JSON(http.StatusOK, a)
}

// GET /attendance/today?student_number=AL21034（任意）
func (ath *AttendanceHandler) GetAttendanceTodayHandler(c *gin.Context) {
	sn := c.Query("student_number")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	list, err := ath.Service.ListTodayJST(ctx, sn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取得失敗"})
		return
	}
	c.JSON(http.StatusOK, list)
}
