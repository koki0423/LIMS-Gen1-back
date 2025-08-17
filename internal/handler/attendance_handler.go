package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"
	"context"
	"strings"
	
	"github.com/gin-gonic/gin"

	"equipmentManager/service"
	"equipmentManager/domain/attendance"
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

func (ath *AttendanceHandler) GetAttendanceRankingHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	ranking, err := ath.Service.GetRanking(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取得失敗"})
		return
	}
	c.JSON(http.StatusOK, ranking)
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

// GET /api/v1/attendance/byIdWithDate?student_id=MA99999&date=2025-08-17
//  または ?student_id=MA99999&month=2025-08
func (ath *AttendanceHandler) GetByIdWithDate(c *gin.Context) {
	studentID := strings.TrimSpace(c.Query("student_id"))
	dateStr   := strings.TrimSpace(c.Query("date"))  // YYYY-MM-DD (JST)
	monthStr  := strings.TrimSpace(c.Query("month")) // YYYY-MM (JST)
	sort      := strings.TrimSpace(c.DefaultQuery("sort", "date:desc"))
	page, _   := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _   := strconv.Atoi(c.DefaultQuery("size", "100"))

	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student_id is required"})
		return
	}
	if dateStr == "" && monthStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date or month is required"})
		return
	}
	if page <= 0 { page = 1 }
	if size <= 0 || size > 1000 { size = 100 }

	// JSTで開始・終了（[start, end)）を作り、UTCに変換して検索
	jst, _ := time.LoadLocation("Asia/Tokyo")
	var startJST, endJST time.Time
	var err error

	if dateStr != "" {
		// 単日
		t, e := time.ParseInLocation("2006-01-02", dateStr, jst)
		if e != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date (YYYY-MM-DD)"})
			return
		}
		startJST = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, jst)
		endJST   = startJST.AddDate(0, 0, 1)
	} else {
		// 月指定
		t, e := time.ParseInLocation("2006-01", monthStr, jst)
		if e != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month (YYYY-MM)"})
			return
		}
		startJST = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, jst)
		endJST   = startJST.AddDate(0, 1, 0)
	}

	startUTC := startJST.UTC()
	endUTC   := endJST.UTC()

	items, total, err := ath.Service.FindByStudentAndRange(
		c.Request.Context(), studentID, startUTC, endUTC, sort, page, size,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch attendance"})
		return
	}

	c.JSON(http.StatusOK, &attendance.AttendanceListResponse{
		Items:   items,
		Total:   total,
		Message: "Attendance fetched successfully",
	})
}