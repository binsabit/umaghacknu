package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Supplies SupplyModel
	Sale     SaleModel
}

func NewModels(db *sql.DB) Models {
	return Models{Supplies: SupplyModel{DB: db},Sale: SaleModel{DB:db}}
}
