package handler
import (
	"database/sql"
	"equipmentManager/service"
	"equipmentManager/domain"
	"net/http"

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


func (ah *AssetHandler) PostAssetsHandler(c *gin.Context) {
		var req domain.CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	id, err := ah.Service.CreateAssetWithMaster(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Asset creation failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Asset created successfully",
		"asset_id": id,
	})
}
func (ah *AssetHandler) GetAssetsAllHandler(c *gin.Context) {}
func (ah *AssetHandler) GetAssetsByAssetIdHandler(c *gin.Context) {}
func (ah *AssetHandler) PutAssetsEditHandler(c *gin.Context) {}
func (ah *AssetHandler) GetAssetsMasterAllHandler(c *gin.Context) {}
func (ah *AssetHandler) GetAssetsMasterByIdHandler(c *gin.Context) {}
func (ah *AssetHandler) DeleteAssetsMasterHandler(c *gin.Context) {}
func (ah *AssetHandler) PostAssetsCheckHandler(c *gin.Context) {}