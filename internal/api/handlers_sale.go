package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/binsabit/umaghacknu/internal/repository/data"
	"github.com/binsabit/umaghacknu/pkg/helpers"
)

func (app *application) GetSales(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()

	fromTime, err := time.Parse("2006-01-02 15:04:05", queryValues.Get("fromTime"))
	toTime, _ := time.Parse("2006-01-02 15:04:05", queryValues.Get("toTime"))
	barcode, _ := strconv.Atoi(queryValues.Get("barcode"))
	if err != nil {
		fmt.Println(err)
		return
	}

	sale, err := app.models.Sale.GetByQuery(int64(barcode), fromTime, toTime)

	if err != nil {
		fmt.Println("error sale query")
		return
	}
	err = helpers.WriteJSON(w, http.StatusOK, sale, nil)
}

func (app *application) AddSales(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Barcode  int64  `json:"barcode"`
		Price    int64  `json:"price"`
		Quantity int    `json:"quantity"`
		SaleTime string `json:"saleTime"`
	}

	err := helpers.ReadJSON(r, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Printf("%v", input)
	t, _ := time.Parse("2006-01-02 15:04:05", input.SaleTime)
	sale := data.Sale{
		Barcode:  input.Barcode,
		Price:    input.Price,
		Quantity: input.Quantity,
		SaleTime: t,
	}
	fmt.Println(input)
	id, err := app.models.Sale.Insert(&sale)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, helpers.Enveleope{"id": id}, nil)
	if err != nil {
		fmt.Println("Create response error")
		return
	}
}

//TODO GET
func (app *application) GetSalesByID(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadIDParam(r)
	if err != nil {
		fmt.Println("id param error")
		return
	}
	result, err := app.models.Sale.GetByID(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, result, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (app *application) UpdateSalesByID(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadIDParam(r)
	if err != nil {
		fmt.Println("id param error")
		return
	}

	fmt.Println(id)
	var input struct {
		Barcode  int64  `json:"barcode"`
		Price    int64  `json:"price"`
		Quantity int    `json:"quantity"`
		SaleTime string `json:"saleTime"`
	}
	
	t, _ := time.Parse("2006-01-02 15:04:05", input.SaleTime)

	err = helpers.ReadJSON(r, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	sale := data.Sale{
		Barcode:  input.Barcode,
		Price:    input.Price,
		Quantity: input.Quantity,
		SaleTime: t,
	}
	fmt.Printf("%v", input)
	err = app.models.Sale.Update(&sale,id)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = helpers.WriteJSON(w, http.StatusOK, nil, nil)
}

//TODO DELETION
func (app *application) DeleteSalesByID(w http.ResponseWriter, r *http.Request) {

	id, err := helpers.ReadIDParam(r)
	if err != nil {
		fmt.Println("id param error")
		return
	}

	err = app.models.Sale.DeleteByID(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, nil, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (app *application) Tester(w http.ResponseWriter, r *http.Request){
	queryValues := r.URL.Query()
	fromTime, _ := time.Parse("2006-01-02 15:04:05", queryValues.Get("fromTime"))
	toTime, _ := time.Parse("2006-01-02 15:04:05", queryValues.Get("toTime"))
	barcode, err := strconv.Atoi(queryValues.Get("barcode"))
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := app.models.Supplies.GetSupplyAmount(int64(barcode), fromTime, toTime)

	if err != nil {
		fmt.Println(err)
		return
	}


	totalQuantity,totalRev ,err := app.models.Sale.GetSalesAmount(int64(barcode), fromTime, toTime)
	
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(totalRev)
	totalCost := Solution(res,totalQuantity)
	totalProfit := totalRev - totalCost 
	err = helpers.WriteJSON(w, http.StatusOK, helpers.Enveleope{"barcode":barcode,"quantity":totalQuantity,"revenue":totalRev,"netProfit":totalProfit}, nil)
	if err != nil{
		fmt.Println(err)
		return
	}

}

func Solution(supplies []data.Output, totalQuant int64)int64{
	fmt.Println(supplies, totalQuant)
	var totalCost int64 = 0
	for i := 0; i < len(supplies); i++ {
		totalQuant = totalQuant - supplies[i].Quantity
		if   totalQuant >= 0 {
			totalCost += supplies[i].Quantity * supplies[i].Price
		}else{
			totalQuant = totalQuant + supplies[i].Quantity
			fmt.Println(totalQuant, supplies[i].Price)
			totalCost += totalQuant * supplies[i].Price
			break
		}
	}
	return totalCost
}