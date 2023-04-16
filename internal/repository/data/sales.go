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
		VALUES ($1,$2,$3, $4)
		RETURNING id`

	args := []interface{}{sale.Barcode, sale.Quantity, sale.Price ,sale.SaleTime}

	rows, err := s.DB.Query(query, args...)
	if err != nil {
		return -1, err
	}
	var id int64
	defer rows.Close()
	rows.Next()
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

	rows, err := s.DB.Query(query, barcode, fromTime, toTime)
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

func (s SaleModel) Update(sale *Sale,id int64) error {
	query := `
		UPDATE sale
		SET barcode = $1, price = $2, quantity = $3, sale_time = $4
		WHERE id = $5`
	res, err := s.DB.Exec(query, sale.Barcode, sale.Price, sale.Quantity, sale.SaleTime, id)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}


func (s SaleModel) GetSalesAmount(barcode int64, fromTime, toTime time.Time)(int64,int64, error){
	query := ` with data as (select barcode,price, sum(quantity) as total_count, sum(quantity)*price as rev from sale
	 where sale_time between $1 and $2 and barcode = $3
	  group by barcode, price) select sum(total_count) as total_quantity, sum(rev) as total_revenue from data 
		group by barcode;`
		

		rows, err := s.DB.Query(query, fromTime, toTime, barcode)
		if err != nil {
			return -1,-1, err
		}
		defer rows.Close()
		rows.Next()
		var totalRev,totalQuantity int64
		err = rows.Scan(&totalRev, &totalQuantity)
		if err != nil {
			return -1,-1,err
		}

		return totalRev,totalQuantity,nil
}