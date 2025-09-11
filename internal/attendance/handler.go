package attendance

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r gin.IRoutes, svc *Service) {

	r.POST("/attendances", handleCreateAttendance(svc))
	r.GET("/attendances", handleListAttendances(svc))
	r.GET("/attendances/stats", handleStats(svc))

	//なぜかHEADがうまく動かないのでv2.0ではコメントアウト
	// g.HEAD("/attendances", handleHeadAttendance(svc))

}

func handleCreateAttendance(svc *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateAttendanceRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			writeErr(c, ErrInvalid("invalid json: "+err.Error()))
			return
		}
		res, created, err := svc.UpsertAttendance(c.Request.Context(), req)
		if err != nil {
			writeErr(c, err)
			return
		}
		if created {
			c.Header("Location", "/attendances/"+strconv.FormatUint(res.AttendanceID, 10))
			c.JSON(http.StatusCreated, res)
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func handleHeadAttendance(svc *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.Query("user_id")
		on := c.Query("on")
		if on == "" {
			on = "today"
		}
		ok, err := svc.Exists(c.Request.Context(), user, on)
		if err != nil {
			writeErr(c, err)
			return
		}
		if ok {
			c.Status(http.StatusOK)
			return
		}
		c.Status(http.StatusNotFound)
	}
}

func handleListAttendances(svc *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		q := ListQuery{
			Limit:  atoiDefault(c.Query("limit"), DefaultPageLimit),
			Offset: atoiDefault(c.Query("offset"), 0),
			Sort:   strDefault(c.Query("sort"), DefaultSort),
			TZ:     strDefault(c.Query("tz"), DefaultTZ),
		}
		if v := c.Query("user_id"); v != "" {
			q.StudentNumber = &[]string{v}[0]
		}
		if v := c.Query("on"); v != "" {
			q.On = &[]string{v}[0]
		}
		if v := c.Query("from"); v != "" {
			q.From = &[]string{v}[0]
		}
		if v := c.Query("to"); v != "" {
			q.To = &[]string{v}[0]
		}

		rows, total, err := svc.List(c.Request.Context(), q)
		if err != nil {
			writeErr(c, err)
			return
		}
		c.Header("X-Total-Count", strconv.FormatInt(total, 10))
		c.JSON(http.StatusOK, gin.H{
			"items": rows,
			"page": gin.H{
				"limit":  q.Limit,
				"offset": q.Offset,
				"total":  total,
			},
		})
	}
}

func handleStats(svc *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := StatsRequest{
			From:  c.Query("from"),
			To:    c.Query("to"),
			Limit: atoiDefault(c.Query("limit"), 10),
		}
		if req.From == "" || req.To == "" {
			writeErr(c, ErrInvalid("from/to are required (YYYY-MM-DD)"))
			return
		}
		rows, err := svc.Stats(c.Request.Context(), req)
		if err != nil {
			writeErr(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"period": gin.H{"from": req.From, "to": req.To},
			"result": rows,
		})
	}
}

func writeErr(c *gin.Context, err error) {
	status := toHTTPStatus(err)
	switch e := err.(type) {
	case *APIError:
		c.JSON(status, gin.H{"error": e})
	default:
		c.JSON(500, gin.H{"error": APIError{Code: CodeInternal, Message: e.Error()}})
	}
}

func atoiDefault(s string, d int) int {
	if s == "" {
		return d
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return d
	}
	return n
}

func strDefault(s, d string) string {
	if s == "" {
		return d
	}
	return s
}
