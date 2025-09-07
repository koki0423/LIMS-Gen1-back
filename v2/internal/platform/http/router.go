package http

import (
  "database/sql"
  "github.com/gin-gonic/gin"
  "IRIS-backend/internal/lends"
)

func NewRouter(db *sql.DB) *gin.Engine {
  r := gin.Default()
  lends.RegisterRoutes(r, lends.NewService(db))
  return r
}
