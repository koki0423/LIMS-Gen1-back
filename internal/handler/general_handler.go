package handler

import (
	"equipmentManager/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	// "log"

	"equipmentManager/domain"
	"database/sql"
)

// --- Swaggerのためのレスポンス構造体定義 ---

// 汎用的なエラーレスポンス
type ErrorResponse struct {
	Error string `json:"error" example:"エラーメッセージ"`
}

// 汎用的な成功メッセージ
type SuccessResponse struct {
	Message string `json:"message" example:"処理が正常に完了しました"`
}

// 備品情報を含む成功レスポンス
type AssetResponse struct {
	Message string       `json:"message" example:"Equipment retrieved successfully"`
	Asset   domain.Asset `json:"asset"`
}

type AssetMasterResponse struct {
	Message string          `json:"message" example:"Asset master retrieved successfully"`
	Asset   domain.AssetMaster `json:"asset"`
}

// 備品リストを含む成功レスポンス
type AssetListResponse struct {
	Message string         `json:"message" example:"Assets list retrieved successfully"`
	Assets  []domain.Asset `json:"assets"`
}

// 貸出情報リストを含む成功レスポンス
type LendListResponse struct {
	Message string              `json:"message" example:"Lends retrieved successfully"`
	Lends   []domain.AssetsLend `json:"lends"`
}

// IDを含む成功レスポンス
type IDResponse struct {
	Message string `json:"message" example:"Asset created successfully"`
	AssetID int64  `json:"asset_id" example:"123"`
}

// -----------------------------------------

type Handler struct {
	DB      *sql.DB
	Service *service.AssetService
}

func NewHandler(db *sql.DB) *Handler {
	service := service.NewAssetService(db)
	return &Handler{
		DB:      db,
		Service: service,
	}
}

// GetAssetsByAssetIdHandler godoc
// @Summary      備品情報をIDで取得
// @Description  指定されたIDを持つ単一の備品情報を取得します。
// @Tags         Assets
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "備品ID"
// @Success      200  {object}  AssetResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /assets/{id} [get]
func (h *Handler) GetAssetsByAssetIdHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid asset ID"})
		return
	}

	asset, err := h.Service.GetAssetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	// レスポンス構造体を使ってJSONを返す
	c.JSON(http.StatusOK, AssetResponse{
		Message: "Equipment retrieved successfully",
		Asset:   asset,
	})
}

// GetAssetsAllHandler godoc
// @Summary      備品情報（すべて）を取得
// @Description  登録されているすべての備品情報を取得します。
// @Tags         Assets
// @Accept       json
// @Produce      json
// @Success      200  {object}  AssetListResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /assetsAll [get]
func (h *Handler) GetAssetsAllHandler(c *gin.Context) {
	assets, err := h.Service.GetAssetAll()
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, AssetListResponse{
		Message: "Assets list retrieved successfully",
		Assets:  assets,
	})
}

// GetAssetsMasterAllHandler godoc
// @Summary      備品マスタ（すべて）を取得
// @Description  登録されているすべての備品マスタ情報を取得します。
// @Tags         Assets
// @Accept       json
// @Produce      json
// @Success      200  {object}  AssetListResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /assets/master [get]
func (h *Handler) GetAssetsMasterAllHandler(c *gin.Context) {
	assetsMaster, err := h.Service.GetAssetMasterAll()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Assets master list retrieved successfully",
		"assets":  assetsMaster,
	})
}

// GetAssetsMasterByIdHandler godoc
// @Summary      備品マスタをIDで取得
// @Description  指定されたIDを持つ単一の備品マスタ情報を取得します。
// @Tags         Assets
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "備品マスタID"
// @Success      200  {object}  AssetResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /assets/master/{id} [get]
func (h *Handler) GetAssetsMasterByIdHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid asset master ID"})
		return
	}

	assetMaster, err := h.Service.GetAssetMasterById(id)

	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, AssetMasterResponse{
		Message: "Asset master retrieved successfully",
		Asset:   assetMaster,
	})
}

//　開発者限定機能のためdocsには残さない
func (h *Handler) DeleteAssetsHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid asset ID"})
		return
	}
	_, err = h.Service.DeleteAssetMasterById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete asset: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Asset deleted successfully",
	})
}

// PostAssetsHandler godoc
// @Summary      新規備品を登録
// @Description  新しい備品情報を登録します。
// @Tags         Assets
// @Accept       json
// @Produce      json
// @Param        asset body      domain.CreateAssetRequest true "登録する備品情報"
// @Success      201   {object}  IDResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /assets [post]
func (h *Handler) PostAssetsHandler(c *gin.Context) {
	var req domain.CreateAssetRequest
	// リクエストボディをバインド（バリデーション付き）
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// サービス層に処理を委譲
	id, err := h.Service.CreateAssetWithMaster(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Asset creation failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Asset created successfully",
		"asset_id": id,
	})
}

// PutAssetsEditHandler godoc
// @Summary      備品情報を更新
// @Description  指定されたIDの備品情報を更新します。
// @Tags         Assets
// @Accept       json
// @Produce      json
// @Param        id    path      int                     true "備品ID"
// @Param        asset body      domain.EditAssetRequest true "更新する備品情報"
// @Success      200   {object}  SuccessResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /assets/{id} [put]
func (h *Handler) PutAssetsEditHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid asset ID"})
		return
	}

	var reqBody domain.EditAssetRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request: " + err.Error()})
		return
	}

	reqBody.AssetID = int64(id)

	_, err = h.Service.PutAssetsEdit(reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Asset update failed"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Asset updated successfully",
	})
}

// PostLendHandler godoc
// @Summary      備品を貸出
// @Description  指定されたIDの備品を貸し出します。
// @Tags         Lends
// @Accept       json
// @Produce      json
// @Param        id   path      int                   true "備品ID"
// @Param        lend body      domain.LendAssetRequest true "貸出情報（貸出先ユーザーなど）"
// @Success      200  {object}  SuccessResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /assets/{id}/lend [post]
func (h *Handler) PostLendHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid asset ID"})
		return
	}

	var req domain.LendAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request: " + err.Error()})
		return
	}
	req.AssetID = int64(id)

	_, err = h.Service.PostLend(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Borrowing asset failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Asset borrowed successfully",
	})
}

// PostReturnHandler godoc
// @Summary      備品を返却
// @Description  指定されたIDの貸出情報を「返却済み」にします。
// @Tags         Lends
// @Accept       json
// @Produce      json
// @Param        id     path      int                       true "貸出ID"
// @Param        return body      domain.AssetReturnRequest true "返却情報（返却日など）"
// @Success      200    {object}  SuccessResponse
// @Failure      400    {object}  ErrorResponse
// @Failure      500    {object}  ErrorResponse
// @Router       /lends/{id}/return [post]
func (h *Handler) PostReturnHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid lend ID"})
		return
	}

	var req domain.AssetReturnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request: " + err.Error()})
		return
	}
	req.LendID = int64(id)

	_, err = h.Service.PostAssetReturnHistory(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Returning asset failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Asset returned successfully",
	})
}

// GetLendsHandler godoc
// @Summary      貸出情報（すべて）を取得
// @Description  すべての貸出情報を取得します。
// @Tags         Lends
// @Accept       json
// @Produce      json
// @Success      200  {object}  LendListResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /lends [get]
func (h *Handler) GetLendsHandler(c *gin.Context) {
	lends, err := h.Service.GetLends()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve lends: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lends retrieved successfully",
		"lends":   lends,
	})
}
