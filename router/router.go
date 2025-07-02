package router

import (
	"equipmentManager/internal/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter は全てのルーターを初期化します
func InitRouter(r *gin.Engine, h *handler.Handler, ah *handler.AuthHandler, nh *handler.NfcHandler) {
	// Swagger UI: http://localhost:8080/swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// APIのベースグループ
	api := r.Group("/api/v1")
	{
		// --- システム ---
		api.GET("/ping", h.PingHandler)

		// --- 資産管理 (Assets) ---
		assets := api.Group("/assets")
		{
			// 既存のエンドポイント
			assets.POST("", h.PostAssetsHandler)            // POST /assets
			assets.GET("/all", h.GetAssetsAllHandler)       // GET /assets/all (ユーザーの元コード /assetsAll に相当)
			assets.GET("/:id", h.GetAssetsByAssetIdHandler) // GET /assets/:id
			assets.PUT("/edit/:id", h.PutAssetsEditHandler) // PUT /assets/edit/:id

			// 備品マスタ
			master := assets.Group("/master")
			{
				master.GET("/all", h.GetAssetsMasterAllHandler)  // GET /assets/master/all
				master.GET("/:id", h.GetAssetsMasterByIdHandler) // GET /assets/master/:id
				master.DELETE("/:id", h.DeleteAssetsHandler)     // DELETE /assets/master/:id
			}

			//【追加】棚卸・点検
			// idはassetsテーブルの主キー
			assets.POST("/check/:id", h.PostAssetsCheckHandler) // POST /assets/check/:id
		}

		// --- 貸出・返却管理 (Lend/Return) ---
		lend := api.Group("/lend")
		{
			// 既存のエンドポイント
			lend.GET("/all", h.GetLendsHandler) // GET /lend/all (ユーザーの元コード /lends に相当)
			// idはassetsテーブルの主キー
			lend.POST("/:id", h.PostLendHandler) // POST /lend/:id
			// idはasset_lendsテーブルの主キー
			lend.POST("/return/:id", h.PostReturnHandler) // POST /lend/return/:id
			// idはasset_lendsテーブルの主キー
			lend.GET("/:id", h.GetLendByIdHandler) // GET /lend/:id

			//【追加】貸出情報の更新
			// idはasset_lendsテーブルの主キー
			lend.PUT("/edit/:id", h.PutLendEditHandler) // PUT /lend/edit/:id

			//【追加】利便性向上
			lend.GET("/overdue", h.GetLendOverdueHandler)                        // GET /lend/overdue
			lend.GET("/user/:student_id", h.GetLendByStudentIdHandler)           // GET /lend/user/:student_id
			lend.GET("/history/:student_id", h.GetLendHistoryByStudentIdHandler) // GET /lend/history/:student_id
		}

		// --- 廃棄管理 (Disposal) ---
		disposal := api.Group("/disposal")
		{
			//【追加】
			// idはassetsテーブルの主キー
			disposal.POST("/:id", h.PostDisposalHandler)  // POST /disposal/:id
			disposal.GET("/all", h.GetDisposalAllHandler) // GET /disposal/all
			// idはasset_disposalsテーブルの主キー
			disposal.GET("/:id", h.GetDisposalByIdHandler) // GET /disposal/:id
		}

		// --- NFC ---
		nfc := api.Group("/nfc")
		{
			nfc.GET("/read", nh.GetNFC) // GET /nfc/read
		}
		// 認証関連ルーターの初期化
		initAuthRouter(api, ah)
	}
}
