/*将来的に使うかもしれない認証必須ページへのルーター*/
package router

import (
	"github.com/gin-gonic/gin"
	"equipmentManager/internal/handler"
)

// curl -X POST http://localhost:8080/api/v1/auth/login -H "Content-Type: application/json" -d '{"email": "test@example.com", "password": "password123"}'

func initAuthRouter(apiRouter *gin.RouterGroup,ah *handler.AuthHandler) {
	// /api/v1/auth グループ
	auth := apiRouter.Group("/auth")
	{
		// POST /api/v1/auth/login
		auth.POST("/login", ah.Login)

		// POST /api/v1/auth/register
		//auth.POST("/register", handler.RegisterHandler)
	}
}