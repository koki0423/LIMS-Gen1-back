package main

import (
	"equipmentManager/internal/database/db"
	"equipmentManager/internal/handler"
	"equipmentManager/router"
	"equipmentManager/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	_ "equipmentManager/docs"
)

// @title           備品管理システム
// @version         1.0
// @description     研究室内の備品を管理するシステムです。
// @host            localhost:8080
// @BasePath        /api/v1

func main() {
	db, err := db.ConnectDB()
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"*", //開発検証用
		},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "DELETE", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	// ハンドラーインスタンスを生成
	h := handler.NewHandler(db)
	ah := handler.NewAuthHandler(db)

	// シングルトンなNfcServiceインスタンスを生成
	nfcService := service.NewNfcService()

	nfcService.Dummy = true // テスト時だけtrue

	// NFC読み取りgoroutineで常駐
	go nfcService.RunNFCReader()
	defer nfcService.Close()

	// NFCハンドラーにはこのインスタンスを渡す
	nh := handler.NewNfcHandler(nfcService)

	router.InitRouter(r, h, ah, nh)

	r.Run("0.0.0.0:8080")
}
