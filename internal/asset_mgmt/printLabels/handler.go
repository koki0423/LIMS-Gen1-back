package printLabels

import (
	"net/http"


	"github.com/gin-gonic/gin"
)

type Handler struct{ svc *Service }

func RegisterRoutes(r gin.IRoutes, svc *Service) {
	h := &Handler{svc: svc}

	r.POST("/assets/print", h.PrintLabels)
}

func (h *Handler) PrintLabels(c *gin.Context) {
	var req PrintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiErr(CodeInvalidArgument, "invalid json"))
		return
	}
	res, err := h.svc.PrintLabels(c.Request.Context(), req)
	if err != nil {
		c.JSON(toHTTPStatus(err), apiErrFrom(err))
		return
	}
	c.JSON(http.StatusCreated, res)
}

// ===== helpers =====

type errDTO struct {
	Error struct {
		Code    Code   `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func apiErr(code Code, msg string) errDTO {
	var e errDTO
	e.Error.Code = code
	e.Error.Message = msg
	return e
}
func apiErrFrom(err error) errDTO {
	if api, ok := err.(*APIError); ok {
		return apiErr(api.Code, api.Message)
	}
	return apiErr(CodeInternal, err.Error())
}
