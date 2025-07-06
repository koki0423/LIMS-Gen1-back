package handler
import (
	"database/sql"
	"equipmentManager/service"

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

func (lh *LendHandler) GetLendsHandler(c *gin.Context) {}
func (lh *LendHandler) PostLendHandler(c *gin.Context) {}
func (lh *LendHandler) PostReturnHandler(c *gin.Context) {}
func (lh *LendHandler) GetLendByIdHandler(c *gin.Context) {}
func (lh *LendHandler) PutLendEditHandler(c *gin.Context) {}
func (lh *LendHandler) GetLendOverdueHandler(c *gin.Context) {}
func (lh *LendHandler) GetLendByStudentIdHandler(c *gin.Context) {}
func (lh *LendHandler) GetLendHistoryByStudentIdHandler(c *gin.Context) {}