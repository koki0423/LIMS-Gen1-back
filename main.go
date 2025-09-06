package main

import (
	"equipmentManager/internal/database/db"
	"equipmentManager/internal/handler"
	"equipmentManager/router"

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
	cfg, err := db.LoadConfig("config/config.yaml")
	if err != nil {
		panic(err)
	}

	assetsDB, err := db.Connect(cfg.AssetsDB)
	if err != nil {
		panic("Failed to connect to the asset database: " + err.Error())
	}

	attendanceDB, err := db.Connect(cfg.AttendanceDB)
	if err != nil {
		panic("Failed to connect to the attendance database: " + err.Error())
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			// "*", //開発検証用
		},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Idempotency-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}))

	// ハンドラーインスタンスを生成f
	sh := handler.NewSystemHandler()
	auh := handler.NewAuthHandler(assetsDB)
	ath := handler.NewAttendanceHandler(attendanceDB)
	ah := handler.NewAssetHandler(assetsDB)
	lh := handler.NewLendHandler(assetsDB)
	dh := handler.NewDisposalHandler(assetsDB)
	ph := handler.NewPrintHandler(assetsDB)

	router.InitRouter(r, sh, ath, auh, ah, lh, dh, ph)

	r.Run("0.0.0.0:8080")
}
