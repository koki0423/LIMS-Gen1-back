package handler

import (
	"database/sql"
	"equipmentManager/domain"
	"equipmentManager/service"
	// "log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AssetHandler struct {
	DB      *sql.DB
	Service *service.AssetService
}

func NewAssetHandler(db *sql.DB) *AssetHandler {
	service := service.NewAssetService(db)
	return &AssetHandler{
		DB:      db,
		Service: service,
	}
}

// POST /assets
func (ah *AssetHandler) PostAssetsHandler(c *gin.Context) {
	var req domain.CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// log.Printf("Received request to create asset: %+v", req)
	// log.Printf("purchace_date: %s", *req.PurchaseDate)

	id, err := ah.Service.CreateAssetWithMaster(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Asset creation failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":         "Asset created successfully",
		"asset_master_id": id,
	})
}

// GET /assets/all
func (ah *AssetHandler) GetAssetsAllHandler(c *gin.Context) {
	assets, err := ah.Service.GetAssetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assets: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Assets fetched successfully",
		"assets":  assets,
	})
}

// GET /assets/:id
func (ah *AssetHandler) GetAssetsByAssetIdHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset ID"})
		return
	}

	int_id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset ID format: " + err.Error()})
		return
	}
	asset, err := ah.Service.GetAssetById(int_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch asset: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Asset fetched successfully",
		"asset":   asset,
	})
}

// PUT /assets/edit/:id
func (ah *AssetHandler) PutAssetsEditHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset ID"})
		return
	}

	int_id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset ID format: " + err.Error()})
		return
	}

	var domain_req domain.EditAssetRequest
	if err := c.ShouldBindJSON(&domain_req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	status, err := ah.Service.PutAssetsEdit(domain_req, int64(int_id))
	if err != nil || !status {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Asset update failed: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Asset updated successfully"})
}

// GET /assets/master/all
func (ah *AssetHandler) GetAssetsMasterAllHandler(c *gin.Context) {
	masters, err := ah.Service.GetAssetMasterAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch asset masters: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Asset masters fetched successfully",
		"masters": masters,
	})
}

// GET /assets/master/:id
func (ah *AssetHandler) GetAssetsMasterByIdHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset master ID"})
		return
	}

	int_id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset master ID format: " + err.Error()})
		return
	}

	master, err := ah.Service.GetAssetMasterById(int_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch asset master: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Asset master fetched successfully",
		"master":  master,
	})
}

// DELETE /assets/master/:id
func (ah *AssetHandler) DeleteAssetsMasterHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset master ID"})
		return
	}

	int_id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset master ID format: " + err.Error()})
		return
	}

	success, err := ah.Service.DeleteAssetMasterById(int_id)
	if err != nil || !success {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Asset master deletion failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset master deleted successfully"})
}

func (ah *AssetHandler) PostAssetsCheckHandler(c *gin.Context) {}
