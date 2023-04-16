package data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Supply struct {
	ID         int64     `json:"id"`
	Price      int64     `json:"price"`
	Barcode    int64     `json:"barcode"`
	Quantity   int       `json:"quantity"`
	SupplyTime time.Time `json:"supplyTime"`
}

type Output struct{
	Quantity int64
	Price int64
	Running_Qty int64
	Running_Cost int64
}

type SupplyModel struct {
	DB *sql.DB
}

//insert

func (s SupplyModel) Insert(supply *Supply) (int64, error) {
	query := `
		INSERT INTO supply (barcode, quantity, price, supply_time)
		VALUES ($1,$2,$3, $4)
		RETURNING id`

	args := []interface{}{supply.Barcode, supply.Quantity, supply.Price ,supply.SupplyTime}

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

	// args := []interface{}{barcode, fromTime, toTime}
	rows, err := s.DB.Query(query, barcode,fromTime,toTime)
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
		DELETE from supply
		WHERE id = $1`
	_, err := s.DB.Exec(query, id)
	return err
}

//update/id

func (s SupplyModel) Update(supply *Supply, id int64) (error) {
	query := `
		UPDATE supply
		SET barcode = $1, price = $2, quantity = $3, supply_time = $4
		WHERE id = $5`
	fmt.Println(supply, id)
	res, err := s.DB.Exec(query, supply.Barcode, supply.Price, supply.Quantity, supply.SupplyTime, id)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}


func (s SupplyModel) GetSupplyAmount(barcode int64, fromTime, toTime time.Time)([]Output,error){
	query := `with data as (select barcode, sum(quantity) as quantity, price, supply_time from supply 
	where supply_time between $1 and $2 and barcode = $3 
	group by quantity, price, barcode, supply_time order by supply_time asc) 
	select quantity, price 
	from data`

	
	result := []Output{}
	rows, err := s.DB.Query(query,fromTime,toTime,barcode)
	if err != nil{
		return nil,err
	}
	for rows.Next(){
		var temp Output
		err = rows.Scan(&temp.Quantity,&temp.Price)
		if err != nil {
			return nil, err
		}
		// fmt.Println(temp)
		result = append(result, temp)
	}
	// fmt.Println(result)

	return result,nil
}


func (s SupplyModel) GetSupplyAmount1(barcode int64, fromTime, toTime time.Time)([]Output,error){
	query := `select quantity, price, running_qty, running_cost from total_cost 
	where barcode = $1 and supply_time between $2 and $3;`

	
	result := []Output{}
	rows, err := s.DB.Query(query,barcode,fromTime,toTime)
	if err != nil{
		return nil,err
	}
	for rows.Next(){
		var temp Output
		err = rows.Scan(&temp.Quantity,&temp.Price,&temp.Running_Qty,&temp.Running_Cost)
		if err != nil {
			return nil, err
		}
		// fmt.Println(temp)
		result = append(result, temp)
	}
	// fmt.Println(result)

	return result,nil
}
