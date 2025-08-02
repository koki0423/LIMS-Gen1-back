package handler

import (
	"database/sql"
	"equipmentManager/domain"
	"equipmentManager/service"
	"strconv"
	"log"

	"github.com/gin-gonic/gin"
)

type LendHandler struct {
	DB      *sql.DB
	Service *service.LendService
}

func NewLendHandler(db *sql.DB) *LendHandler {
	service := service.NewLendService(db)
	return &LendHandler{
		DB:      db,
		Service: service,
	}
}

// @Summary      貸出情報一覧取得
// @Tags         Lend
// @Produce      json
// @Success      200  {object}  handler.LendListResponse
// @Failure      500  {object}  handler.ErrorResponse
// @Router       /lend/all [get]
func (lh *LendHandler) GetLendsHandler(c *gin.Context) {
	lends, err := lh.Service.GetLends()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch lends: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{"lends": lends})
}

// @Summary      備品貸出処理
// @Tags         Lend
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "備品ID (Asset ID)"
// @Param        body body      domain.LendAssetRequest true  "貸出情報"
// @Success      201  {object}  handler.SuccessResponse
// @Failure      400  {object}  handler.ErrorResponse
// @Failure      500  {object}  handler.ErrorResponse
// @Router       /lend/{id} [post]
// POST /lend/:id idはassetsテーブルの主キー
func (lh *LendHandler) PostLendHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Asset ID is required"})
	}

	int_id, err := strconv.Atoi(id)
	var req domain.LendAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	req.AssetID = int64(int_id)

	success, err := lh.Service.PostLend(req)
	if err != nil {
		c.JSON(500, gin.H{"error": "Lend creation failed: " + err.Error()})
		return
	}

	if success {
		c.JSON(201, gin.H{"message": "Lend created successfully"})
	} else {
		c.JSON(500, gin.H{"error": "Failed to create lend"})
	}
}

// @Summary      備品返却処理
// @Tags         Lend
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "貸出ID (Lend ID)"
// @Param        body body      domain.ReturnAssetRequest true  "返却情報"
// @Success      200  {object}  handler.SuccessResponse
// @Failure      400  {object}  handler.ErrorResponse
// @Failure      500  {object}  handler.ErrorResponse
// @Router       /lend/return/{id} [post]
// POST /lend/return/:id idはrendsテーブルの主キー
func (lh *LendHandler) PostReturnHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Lend ID is required"})
		return
	}
	int_id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid lend ID format: " + err.Error()})
		return
	}
	var req domain.ReturnAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	req.LendID = int64(int_id)
	success, err := lh.Service.PostReturn(req,int64(int_id))
	if err != nil {
		c.JSON(500, gin.H{"error": "Return processing failed: " + err.Error()})
		return
	}
	if success {
		c.JSON(200, gin.H{"message": "Return processed successfully"})
	} else {
		c.JSON(500, gin.H{"error": "Failed to process return"})
	}
}

// @Summary      貸出情報取得 (ID指定)
// @Tags         Lend
// @Produce      json
// @Param        id   path      int  true  "貸出ID (Lend ID)"
// @Success      200  {object}  handler.LendResponse
// @Failure      400  {object}  handler.ErrorResponse
// @Failure      500  {object}  handler.ErrorResponse
// @Router       /lend/{id} [get]
// GET /lend/:id
func (lh *LendHandler) GetLendByIdHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Lend ID is required"})
		return
	}

	int_id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid lend ID format: " + err.Error()})
		return
	}

	lend, err := lh.Service.GetLendById(int64(int_id))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch lend: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Fetch completed", "lend": lend})
}

// @Summary      貸出情報一覧取得 (名称付き)
// @Tags         Lend
// @Produce      json
// @Success      200  {object}  handler.LendingDetailListResponse
// @Failure      500  {object}  handler.ErrorResponse
// @Router       /lend/all/with-name [get]
// GET /lend/all/with-name
func (lh *LendHandler) GetLendsWithNameHandler(c *gin.Context) {
	lends, err := lh.Service.GetLendsWithName()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch lends with names: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"lends": lends})
}

// @Summary      貸出情報編集
// @Tags         Lend
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "貸出ID (Lend ID)"
// @Param        body body      domain.EditLendRequest true  "更新情報"
// @Success      200  {object}  handler.SuccessResponse
// @Failure      400  {object}  handler.ErrorResponse
// @Failure      500  {object}  handler.ErrorResponse
// @Router       /lend/edit/{id} [put]
// PUT /lend/edit/:id
func (lh *LendHandler) PutLendEditHandler(c *gin.Context)               {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Lend ID is required"})
		return
	}

	int_id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid lend ID format: " + err.Error()})
		return
	}

	var req domain.EditLendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	log.Printf("Lend edit request: %+v", req)

	success, err := lh.Service.PutLendEdit(req, int64(int_id))
	if err != nil || !success {
		c.JSON(500, gin.H{"error": "Failed to update lend: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Lend updated successfully"})
}

// func (lh *LendHandler) GetLendByStudentIdHandler(c *gin.Context)        {}
// func (lh *LendHandler) GetLendHistoryByStudentIdHandler(c *gin.Context) {}
