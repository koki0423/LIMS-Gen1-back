package router

import (
	"equipmentManager/internal/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter は全てのルーターを初期化します
func InitRouter(r *gin.Engine, sh *handler.SystemHandler, auh *handler.AuthHandler, ah *handler.AssetHandler, lh *handler.LendHandler, dh *handler.DisposalHandler) {
	// Swagger UI: http://localhost:8080/swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// APIのベースグループ
	api := r.Group("/api/v1")
	{
		// --- システム ---
		api.GET("/ping", sh.PingHandler)

		// --- 資産管理 (Assets) ---
		assets := api.Group("/assets")
		{
			// 既存のエンドポイント
			assets.POST("", ah.PostAssetsHandler)            // POST /assets
			assets.GET("/all", ah.GetAssetsAllHandler)       // GET /assets/all
			assets.GET("/:id", ah.GetAssetsByAssetIdHandler) // GET /assets/:id
			assets.PUT("/edit/:id", ah.PutAssetsEditHandler) // PUT /assets/edit/:id

			// 備品マスタ
			master := assets.Group("/master")
			{
				master.GET("/all", ah.GetAssetsMasterAllHandler)    // GET /assets/master/all
				master.GET("/:id", ah.GetAssetsMasterByIdHandler)   // GET /assets/master/:id 主キーはmasterの主キー
				master.DELETE("/:id", ah.DeleteAssetsMasterHandler) // DELETE /assets/master/:id 主キーはmasterの主キー
			}

			// 気が向いたら実装
			// idはassetsテーブルの主キー
			assets.POST("/check/:id", ah.PostAssetsCheckHandler) // POST /assets/check/:id
		}

		// --- 貸出・返却管理 (Lend/Return) ---
		lend := api.Group("/lend")
		{
			// 既存のエンドポイント
			lend.GET("/all", lh.GetLendsHandler) // GET /lend/all

			/*ここのエンドポイント違和感に思うが，特定の備品を指定して貸出記録をつけるから備品テーブルの主キーを指定する*/
			// idはassetsテーブルの主キー
			lend.POST("/:id", lh.PostLendHandler) // POST /lend/:id
			
			// idはasset_lendsテーブルの主キー
			lend.POST("/return/:id", lh.PostReturnHandler) // POST /lend/return/:id

			// idはassetテーブルの主キー
			lend.GET("/:id", lh.GetLendByIdHandler) // GET /lend/:id
			lend.GET("/all/with-name", lh.GetLendsWithNameHandler) // GET /lend/all/with-name

			//貸出情報の更新
			// idはasset_lendsテーブルの主キー
			lend.PUT("/edit/:id", lh.PutLendEditHandler) // PUT /lend/edit/:id

			/*必要になったら実装*/
			// lend.GET("/user/:student_id", lh.GetLendByStudentIdHandler)           // GET /lend/user/:student_id
			// lend.GET("/history/:student_id", lh.GetLendHistoryByStudentIdHandler) // GET /lend/history/:student_id
		}

		// --- 廃棄管理 (Disposal) ---
		disposal := api.Group("/disposal")
		{
			// idはassetsテーブルの主キー
			disposal.POST("/:id", dh.PostDisposalHandler)  // POST /disposal/:id
			disposal.GET("/all", dh.GetDisposalAllHandler) // GET /disposal/all
			// idはasset_disposalsテーブルの主キー
			disposal.GET("/:id", dh.GetDisposalByIdHandler) // GET /disposal/:id
		}
		
		// 認証関連ルーターの初期化
		initAuthRouter(api, auh)
	}
}
