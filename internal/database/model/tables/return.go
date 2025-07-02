package model

import (
	"database/sql"
	"time"
)

type AssetReturn struct {
	ID               int64          `db:"id"`
	LendID           int64          `db:"lend_id"`
	ReturnedQuantity int            `db:"returned_quantity"`
	ReturnDate       time.Time      `db:"return_date"`
	Notes            sql.NullString `db:"notes"`
}
