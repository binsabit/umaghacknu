package data

import (
	"database/sql"
	"errors"
	"time"
)

type Sale struct {
	ID       int64     `json:"id"`
	Price    int64     `json:"price"`
	Barcode  int64     `json:"barcode"`
	Quantity int       `json:"quantity"`
	SaleTime time.Time `json:"saleTime"`
}

type SaleModel struct {
	DB *sql.DB
}

//insert
func (s SaleModel) Insert(sale *Sale) (int64, error) {
	query := `
		INSERT INTO sale (barcode, quantity, price, sale_time)
		VALUES ($1,$2,$3 $4)
		RETURNING id`

	args := []interface{}{sale.Barcode, sale.Quantity, sale.SaleTime}

	rows, err := s.DB.Query(query, args...)
	if err != nil {
		return -1, err
	}
	var id int64
	defer rows.Close()
	err = rows.Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

//getByID
func (s SaleModel) GetByID(id int64) (*Sale, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id,barcode,price,quantity,sale_time
		FROM sale
		WHERE id = $1`

	var sale Sale

	err := s.DB.QueryRow(query, id).Scan(
		&sale.ID,
		&sale.Barcode,
		&sale.Price,
		&sale.Quantity,
		&sale.SaleTime,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &sale, nil
}

//getByQuery
func (s SaleModel) GetByQuery(barcode int64, fromTime, toTime time.Time) ([]Sale, error) {

	query := `
		SELECT id,barcode,price,quantity,sale_time
		FROM sale
		WHERE barcode = $1 AND sale_time between $2 AND $3 `

	sale := []Sale{}

	args := []interface{}{barcode, fromTime, toTime}
	rows, err := s.DB.Query(query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var temp Sale
		err = rows.Scan(
			&temp.ID,
			&temp.Barcode,
			&temp.Price,
			&temp.Quantity,
			&temp.SaleTime,
		)
		if err != nil {
			return nil, err
		}
		sale = append(sale, temp)
	}

	return sale, nil
}

//DeleteByID
func (s SaleModel) DeleteByID(id int64) error {
	query := `
		DELETE sale
		WHERE id = $1`
	_, err := s.DB.Exec(query, id)

	return err
}

//update/id

func (s SaleModel) Update(sale *Sale) error {
	query := `
		UPDATE sale
		SET barcode = $1, price = $2, quantity = $3, saletime = $4
		WHERE id = $5,
		RETURNING id`
	res, err := s.DB.Exec(query, sale.Barcode, sale.Price, sale.Quantity, sale.SaleTime, sale.ID)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}
