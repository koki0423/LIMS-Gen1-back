package utils

import (
	"log"
	"time"
)

func MustParseDate(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		log.Fatalf("日付のパースに失敗: %v", err)
	}
	return t
}
