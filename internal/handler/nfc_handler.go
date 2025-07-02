package handler

import (
	"equipmentManager/service"
	"equipmentManager/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ---swagger用---
type GetNFCResponse struct {
	Message string                  `json:"message" example:"NFC data retrieved successfully"`
	Data    domain.StudentCardInfo `json:"data"`
}
//----------------


type NfcHandler struct {
	Service *service.NfcService
}

func NewNfcHandler(svc *service.NfcService) *NfcHandler {
	return &NfcHandler{Service: svc}
}


// GetNFC godoc
// @Summary NFCカード情報の取得
// @Description NFCリーダからカード情報を取得します。ダミーモード時はサンプル値を返します。
// @Tags nfc
// @Produce json
// @Success 200 {object} GetNFCResponse
// @Failure 500 {object} ErrorResponse
// @Router /nfc/read [get]
func (nh *NfcHandler) GetNFC(c *gin.Context) {
	if nh.Service.Dummy {
		// ダミーデータを返す
		c.JSON(http.StatusOK, gin.H{
			"message": "Dummy NFC data retrieved successfully",
			"data":    domain.StudentCardInfo{StudentNumber: "MA12345"},
		})
		return
	}
	cardInfo, err := nh.Service.GetNFCData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "NFC data retrieved successfully",
		"data":    cardInfo,
	})
}
