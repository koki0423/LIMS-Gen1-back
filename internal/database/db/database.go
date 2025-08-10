package db

import (
	"database/sql"
	"fmt"
	// "log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
)

const (
	driverName     = "mysql"
	configFilePath = "config/config.yaml"
)

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type Config struct {
	Version       string         `yaml:"version"`
	AssetsDB      DatabaseConfig `yaml:"assets_db"`      // 備品管理
	AttendanceDB  DatabaseConfig `yaml:"attendance_db"`  // 出席管理
}

func LoadConfig(path string) (*Config, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("設定ファイルの読み込み失敗: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(buf, &cfg); err != nil {
		return nil, fmt.Errorf("設定ファイルのパース失敗: %w", err)
	}
	return &cfg, nil
}

func Connect(c DatabaseConfig) (*sql.DB, error) {
	// タイムアウトも付けて握りっぱなし対策
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&tls=false&timeout=3s&readTimeout=5s&writeTimeout=5s",
		c.Username, c.Password, c.Host, c.Port, c.DBName)

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("接続準備に失敗: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("DB接続に失敗: %w", err)
	}

	// 接続プール（※合算がMySQLの max_connections を超えないよう配分する）
	db.SetMaxOpenConns(80)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	return db, nil
}
