package service

import (
	"database/sql"

)

type DisposalService struct {
	DB *sql.DB
}

func NewDisposalService(db *sql.DB) *DisposalService {
	return &DisposalService{DB: db}
}

