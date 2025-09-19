package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
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

// フロントのビルド出力を埋め込む（backend/public 配下）
//go:embed public
var embedded embed.FS

func main() {
	cfg, err := db.LoadConfig("config/config.yaml")
	if err != nil {
		panic(err)
	}

	// NOTE: 変数名を分けてシャドーイング回避（お好みで）
	conn, err := db.Connect(cfg.DB)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	log.Printf("[INFO] connected to DB: %s", cfg.DB.DBName)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	_ = r.SetTrustedProxies(nil) // 直配信なら nil でOK（逆プロキシ配下ならCIDRを設定）

	// CORS（開発中のみ必要。埋め込み配信に切り替えたら基本不要）
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Idempotency-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}))

	// ヘルス
	r.GET("/healthz", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	// /api/v2
	api := r.Group("/api/v2")
	assets.RegisterRoutes(api, assets.NewService(conn))
	lends.RegisterRoutes(api, lends.NewService(conn))
	disposals.RegisterRoutes(api, disposals.NewService(conn))
	attendance.RegisterRoutes(api, attendance.NewService(conn))
	printLabels.RegisterRoutes(api, printLabels.NewService())

	// ===== ここから静的配信（埋め込みSPA） =====
	sub, err := fs.Sub(embedded, "public")
	if err != nil {
		log.Fatal(err)
	}
	fileFS := http.FS(sub)

	// すべての非APIリクエストを静的 or SPA で処理
	r.NoRoute(func(c *gin.Context) {
		// API は対象外
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Status(http.StatusNotFound)
			return
		}

		reqPath := strings.TrimPrefix(c.Request.URL.Path, "/")
		if reqPath == "" {
			reqPath = "index.html"
		}

		// 実ファイルがあるならそれを返す（Content-Type を推測、キャッシュ付与）
		if f, err := fileFS.Open(reqPath); err == nil {
			defer f.Close()
			if ct := mime.TypeByExtension(path.Ext(reqPath)); ct != "" {
				c.Header("Content-Type", ct)
			}
			// index.html 以外はキャッシュ（SPAの基本運用）
			if !strings.HasSuffix(reqPath, "index.html") {
				c.Header("Cache-Control", "public, max-age=86400, immutable")
			}
			if fileInfo, err := f.Stat(); err == nil {
				http.ServeContent(c.Writer, c.Request, reqPath, fileInfo.ModTime(), f)
			} else {
				c.Status(http.StatusInternalServerError)
			}
			return
		}

		// なければ index.html にフォールバック（ヒストリーAPI対応）
		if idx, err := fileFS.Open("index.html"); err == nil {
			defer idx.Close()
			c.Header("Content-Type", "text/html; charset=utf-8")
			if fileInfo, err := idx.Stat(); err == nil {
				http.ServeContent(c.Writer, c.Request, "index.html", fileInfo.ModTime(), idx)
			} else {
				c.Status(http.StatusInternalServerError)
			}
			return
		}

		c.Status(http.StatusNotFound)
	})
	// ===== ここまで =====

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

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
}