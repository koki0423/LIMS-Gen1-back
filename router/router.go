package router

import (
	"equipmentManager/internal/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter は全てのルーターを初期化する
func InitRouter(r *gin.Engine, h *handler.Handler, ah *handler.AuthHandler, nh *handler.NfcHandler) {
	// http://localhost:8080/swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// APIのベースグループ
	api := r.Group("/api/v1")
	{
		// システム
		// GET /api/v1/ping
		api.GET("/ping", h.PingHandler) //OK

		// 資産管理
		//GET /api/v1/assets/:id
		// id はassetsテーブルのID（主キー）を指定
		// 資産IDを指定して資産情報を取得するエンドポイント
		api.GET("/assets/:id", h.GetAssetsByAssetIdHandler) //OK

		// GET /api/v1/assetsAll
		// 全資産一覧を取得するエンドポイント
		api.GET("/assetsAll", h.GetAssetsAllHandler) //OK

		// GET /api/v1/assets/master
		// 全資産マスター一覧を取得するエンドポイント
		api.GET("/assets/master", h.GetAssetsMasterAllHandler) //OK

		//DELETE /api/v1/assets/master/:id
		// id はassets_masterテーブルのID（主キー）を指定．管理者権限が必要
		api.DELETE("/assets/master/:id", h.DeleteAssetsHandler) //OK

		// GET /api/v1/assets/master/:id
		// id はassets_masterテーブルのID（主キー）を指定．
		api.GET("/assets/master/:id", h.GetAssetsMasterByIdHandler)

		// POST /api/v1/assets
		// 新規資産登録のエンドポイント
		api.POST("/assets", h.PostAssetsHandler) //OK

		//PUT /api/v1/assets/edit/:id
		// id はassetsテーブルのID（主キー）を指定
		// 資産情報を更新するエンドポイント
		api.PUT("/assets/edit/:id", h.PutAssetsEditHandler) //OK

		// 貸出管理
		// POST /api/v1/borrow/:id
		// id はassetsテーブルのID（主キー）を指定．
		api.POST("/lend/:id", h.PostLendHandler) //ok

		// POST /api/v1/return/:id
		// id はasset_lendsテーブルのID（主キー）を指定．
		api.POST("/lend/return/:id", h.PostReturnHandler) //OK

		// GET /api/v1/borrows
		// 貸出中一覧を取得するエンドポイント
		api.GET("/lends", h.GetLendsHandler) //OK

		// GET /api/v1/auth/nfc
		api.GET("/nfc/read", nh.GetNFC)

		// 認証関連ルーターの初期化
		initAuthRouter(api, ah)

		// 他のルーターがあればここに追加
		// initUserRouter(api)
		// initProductRouter(api)
	}
}
