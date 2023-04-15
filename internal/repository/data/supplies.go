package data

import (
	"database/sql"
	"errors"
	"time"
)

type Supply struct {
	ID         int64     `json:"id"`
	Price      int64     `json:"price"`
	Barcode    int64     `json:"barcode"`
	Quantity   int       `json:"quantity"`
	SupplyTime time.Time `json:"supplyTime"`
}

type SupplyModel struct {
	DB *sql.DB
}

//insert

func (s SupplyModel) Insert(supply *Supply) (int64, error) {
	query := `
		INSERT INTO supply (barcode, quantity, price, supply_time)
		VALUES ($1,$2,$3 $4)
		RETURNING id`

	args := []interface{}{supply.Barcode, supply.Quantity, supply.SupplyTime}

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
func (s SupplyModel) GetByID(id int64) (*Supply, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id,barcode,price,quantity,supply_time
		FROM supply
		WHERE id = $1`

	var supply Supply

	err := s.DB.QueryRow(query, id).Scan(
		&supply.ID,
		&supply.Barcode,
		&supply.Price,
		&supply.Quantity,
		&supply.SupplyTime,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &supply, nil
}

//getByQuery
func (s SupplyModel) GetByQuery(barcode int64, fromTime, toTime time.Time) ([]Supply, error) {

	query := `
		SELECT id,barcode,price,quantity,supply_time
		FROM supply
		WHERE barcode = $1 AND supply_time between $2 AND $3 `

	supply := []Supply{}

	args := []interface{}{barcode, fromTime, toTime}
	rows, err := s.DB.Query(query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var temp Supply
		err = rows.Scan(
			&temp.ID,
			&temp.Barcode,
			&temp.Price,
			&temp.Quantity,
			&temp.SupplyTime,
		)
		if err != nil {
			return nil, err
		}
		supply = append(supply, temp)
	}

	return supply, nil
}

//DeleteByID
func (s SupplyModel) DeleteByID(id int64) error {
	query := `
		DELETE supply
		WHERE id = $1`
	_, err := s.DB.Exec(query, id)

	return err
}

//update/id

func (s SupplyModel) Update(supply *Supply) error {
	query := `
		UPDATE supply
		SET barcode = $1, price = $2, quantity = $3, saletime = $4
		WHERE id = $5,
		RETURNING id`
	res, err := s.DB.Exec(query, supply.Barcode, supply.Price, supply.Quantity, supply.SupplyTime, supply.ID)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}
