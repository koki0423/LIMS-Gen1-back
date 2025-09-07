package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"IRIS-backend/internal/asset_mgmt/assets"
	"IRIS-backend/internal/asset_mgmt/disposals"
	"IRIS-backend/internal/asset_mgmt/lends"
	"IRIS-backend/internal/attendance"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		// 例: root:password@tcp(127.0.0.1:3306)/assetdb?parseTime=true&loc=UTC
		dsn = "devadmin:X$Q9zB2Wb2x2@tcp(192.168.0.61:3306)/lims_v1?parseTime=true&loc=UTC"
		log.Printf("[INFO] DB_DSN not set, using default: %s", dsn)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	// ヘルス
	r.GET("/healthz", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	// ルート登録
	assets.RegisterRoutes(r, assets.NewService(db))
	lends.RegisterRoutes(r, lends.NewService(db))
	disposals.RegisterRoutes(r, disposals.NewService(db))
	attendance.RegisterRoutes(r, attendance.NewService(db))

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
