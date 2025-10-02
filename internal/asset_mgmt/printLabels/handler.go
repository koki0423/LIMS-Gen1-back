package printLabels

import (
	"errors"
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
		c.JSON(http.StatusBadRequest, newErrDTO(ErrInvalid("invalid json")))
		return
	}

	res, err := h.svc.PrintLabels(c.Request.Context(), req)
	if err != nil {
		c.JSON(toHTTPStatus(err), newErrDTO(err))
		return
	}
	c.JSON(http.StatusCreated, res)
}

// ===== helpers =====
type errDTO struct {
	Error *APIError `json:"error"`
}

func newErrDTO(err error) errDTO {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return errDTO{Error: apiErr}
	}

	return errDTO{Error: ErrInternal(err.Error())}
}
