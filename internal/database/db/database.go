package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
)

const (
	dbDriverName   = "mysql"
	configFilePath = "config/config.yaml"
)

type Config struct {
	Version  string         `yaml:"version"`
	Database DatabaseConfig `yaml:"database"`
}

// データベース設定部分の構造体
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

func loadConfig(path string) (*DatabaseConfig, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("設定ファイル '%s' の読み込みに失敗しました: %w", path, err)
	}

	var cfg Config
	
	if err := yaml.Unmarshal(buf, &cfg); err != nil {
		return nil, fmt.Errorf("設定ファイルのパースに失敗しました: %w", err)
	}

	return &cfg.Database, nil
}

func ConnectDB() (*sql.DB, error) {
	config, err := loadConfig(configFilePath)
	if err != nil {
		return nil, err
	}
	// DSN (Data Source Name) を構築
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&tls=false",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	db, err := sql.Open(dbDriverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("データベース接続の準備に失敗しました: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("データベースへの接続に失敗しました: %w", err)
	}

	log.Println("データベース接続：成功")
	return db, nil
}
