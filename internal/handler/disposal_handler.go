package handler

import (
	"database/sql"
	"equipmentManager/domain"
	"equipmentManager/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DisposalHandler struct {
	DB      *sql.DB
	Service *service.DisposalService
}

func NewDisposalHandler(db *sql.DB) *DisposalHandler {
	service := service.NewDisposalService(db)
	return &DisposalHandler{
		DB:      db,
		Service: service,
	}
}

// @Summary      資産の廃棄処理
// @Description  指定された資産IDの資産を廃棄登録します。
// @Tags         Disposal
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "資産ID (Asset ID)"
// @Param        body body      domain.CreateDisposalRequest true  "廃棄登録リクエスト"
// @Success      201  {object}  handler.SuccessResponse
// @Failure      400  {object}  handler.ErrorResponse
// @Failure      500  {object}  handler.ErrorResponse
// @Router       /disposal/{id} [post]
// POST /disposal/:id
func (dh *DisposalHandler) PostDisposalHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID is required"})
		return
	}

	int_id, err := strconv.Atoi(id)

	var req domain.CreateDisposalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	err = dh.Service.CreateDisposal(req, int64(int_id))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create disposal: " + err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Disposal created successfully"})
}

// @Summary      全廃棄履歴の取得
// @Description  全ての資産廃棄履歴を取得します。
// @Tags         Disposal
// @Produce      json
// @Success      200  {object}  handler.DisposalListResponse
// @Failure      500  {object}  handler.ErrorResponse
// @Router       /disposal/all [get]
// GET /disposal/all
func (dh *DisposalHandler) GetDisposalAllHandler(c *gin.Context) {
	allDisposals, err := dh.Service.GetDisposalAll()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch disposals: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"message":   "Disposals fetched successfully",
		"disposals": allDisposals,
	})
}

// @Summary      廃棄履歴の取得 (ID指定)
// @Description  指定されたIDの廃棄履歴を取得します。
// @Tags         Disposal
// @Produce      json
// @Param        id   path      int  true  "廃棄ID (Disposal ID)"
// @Success      200  {object}  handler.DisposalResponseWrapper
// @Failure      400  {object}  handler.ErrorResponse
// @Failure      500  {object}  handler.ErrorResponse
// @Router       /disposal/{id} [get]
func (dh *DisposalHandler) GetDisposalByIdHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID is required"})
		return
	}
	int_id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format: " + err.Error()})
		return
	}

	disposal, err := dh.Service.GetDisposalByID(int64(int_id))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch disposal: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"message":  "Disposal fetched successfully",
		"disposal": disposal,
	})
}
