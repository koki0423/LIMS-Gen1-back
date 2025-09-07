package disposals

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct{ svc *Service }

func RegisterRoutes(r *gin.Engine, svc *Service) {
	h := &Handler{svc: svc}
	// 登録
	r.POST("/assets/:management_number/disposals", h.CreateDisposal) //OK
	// 参照
	r.GET("/disposals", h.ListDisposals)              //OK
	r.GET("/disposals/:disposal_ulid", h.GetDisposal) //OK
}

func (h *Handler) CreateDisposal(c *gin.Context) {
	mng := c.Param("management_number")
	var req CreateDisposalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorBody(CodeInvalidArgument, "invalid json"))
		return
	}
	res, err := h.svc.CreateDisposal(c.Request.Context(), mng, req)
	if err != nil {
		c.JSON(ToHTTPStatus(err), errorFromErr(err))
		return
	}
	c.Header("Location", "/disposals/"+res.DisposalULID)
	c.JSON(http.StatusCreated, res)
}

func (h *Handler) GetDisposal(c *gin.Context) {
	ul := c.Param("disposal_ulid")
	res, err := h.svc.GetDisposalByULID(c.Request.Context(), ul)
	if err != nil {
		c.JSON(ToHTTPStatus(err), errorFromErr(err))
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) ListDisposals(c *gin.Context) {
	f := DisposalFilter{}
	if v := c.Query("management_number"); v != "" {
		f.ManagementNumber = &v
	}
	if v := c.Query("processed_by_id"); v != "" {
		f.ProcessedByID = &v
	}
	if v := c.Query("from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			f.From = &t
		}
	}
	if v := c.Query("to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			f.To = &t
		}
	}
	p := Page{
		Limit:  parseIntDefault(c.Query("limit"), 50),
		Offset: parseIntDefault(c.Query("offset"), 0),
		Order:  c.DefaultQuery("order", "desc"),
	}
	res, err := h.svc.ListDisposals(c.Request.Context(), f, p)
	if err != nil {
		c.JSON(ToHTTPStatus(err), errorFromErr(err))
		return
	}
	c.JSON(http.StatusOK, res)
}

// ---- helpers ----

func parseIntDefault(s string, d int) int {
	if s == "" {
		return d
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return d
	}
	return v
}

type errorDTO struct {
	Error struct {
		Code    Code   `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func errorBody(code Code, msg string) errorDTO {
	var e errorDTO
	e.Error.Code = code
	e.Error.Message = msg
	return e
}
func errorFromErr(err error) errorDTO {
	msg := err.Error()
	if api, ok := err.(*APIError); ok {
		return errorBody(api.Code, api.Message)
	}
	return errorBody(CodeInternal, msg)
}
