package service
import (
	"database/sql"

)

type LendService struct {
	DB *sql.DB
}

func NewLendService(db *sql.DB) *LendService {
	return &LendService{DB: db}
}