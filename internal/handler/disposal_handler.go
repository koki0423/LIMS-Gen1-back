package handler

import (
	"database/sql"
	"equipmentManager/service"

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

func (dh *DisposalHandler) PostDisposalHandler(c *gin.Context) {}
func (dh *DisposalHandler) GetDisposalAllHandler(c *gin.Context) {}
func (dh *DisposalHandler) GetDisposalByIdHandler(c *gin.Context) {
}