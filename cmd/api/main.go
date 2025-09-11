package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"IRIS-backend/internal/asset_mgmt/assets"
	"IRIS-backend/internal/asset_mgmt/disposals"
	"IRIS-backend/internal/asset_mgmt/lends"
	"IRIS-backend/internal/asset_mgmt/printLabels"
	"IRIS-backend/internal/attendance"
	"IRIS-backend/internal/platform/db"
)

func main() {
	cfg, err := db.LoadConfig("config/config.yaml")
	if err != nil {
		panic(err)
	}

	db, err := db.Connect(cfg.DB)
	if err != nil {
		panic(err)
	}

	log.Printf("[INFO] connected to DB: %s", cfg.DB.DBName)

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

	// ヘルス
	r.GET("/healthz", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	// /api/v2 グループ
	api := r.Group("/api/v2")
	
	// ルート登録
	assets.RegisterRoutes(api, assets.NewService(db))
	lends.RegisterRoutes(api, lends.NewService(db))
	disposals.RegisterRoutes(api, disposals.NewService(db))
	attendance.RegisterRoutes(api, attendance.NewService(db))
	printLabels.RegisterRoutes(api, printLabels.NewService())

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// デバッグ用
	// for _, ri := range r.Routes() {
	// 	log.Printf("ROUTE: %s %s -> %s", ri.Method, ri.Path, ri.Handler)
	// }

	// 起動
	go func() {
		log.Println("[INFO] listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("[INFO] shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	_ = db.Close()
}
