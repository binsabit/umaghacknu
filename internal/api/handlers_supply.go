package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/binsabit/umaghacknu/internal/repository/data"
	"github.com/binsabit/umaghacknu/pkg/helpers"
)

func (app *application) GetSupplies(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()

	fromTime, err := time.Parse("2006-01-02 15:04:05", queryValues.Get("fromTime"))
	toTime, _ := time.Parse("2006-01-02 15:04:05", queryValues.Get("toTime"))
	barcode, _ := strconv.Atoi(queryValues.Get("barcode"))
	if err != nil {
		fmt.Println(err)
		return
	}

	supply, err := app.models.Supplies.GetByQuery(int64(barcode), fromTime, toTime)

	if err != nil {
		fmt.Println(err)
		return
	}
	err = helpers.WriteJSON(w, http.StatusOK, supply, nil)
}

func (app *application) AddSupplies(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Barcode    int64  `json:"barcode"`
		Price      int64  `json:"price"`
		Quantity   int    `json:"quantity"`
		SupplyTime string `json:"supplyTime"`
	}

	err := helpers.ReadJSON(r, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Printf("%v", input)
	t, _ := time.Parse("2006-01-02 15:04:05", input.SupplyTime)
	supply := data.Supply{
		Barcode:    input.Barcode,
		Price:      input.Price,
		Quantity:   input.Quantity,
		SupplyTime: t,
	}
	id, err := app.models.Supplies.Insert(&supply)

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
func (app *application) GetSuppliesByID(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadIDParam(r)
	if err != nil {
		fmt.Println("id param error")
		return
	}
	result, err := app.models.Supplies.GetByID(id)
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

func (app *application) UpdateSuppliesByID(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadIDParam(r)
	if err != nil {
		fmt.Println("id param error")
		return
	}

	fmt.Println(id)
	var input struct {
		Barcode    int64  `json:"barcode"`
		Price      int64  `json:"price"`
		Quantity   int    `json:"quantity"`
		SupplyTime string `json:"supplyTime"`
	}
	t, _ := time.Parse("2006-01-02 15:04:05", input.SupplyTime)

	err = helpers.ReadJSON(r, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	supply := data.Supply{
		Barcode:    input.Barcode,
		Price:      input.Price,
		Quantity:   input.Quantity,
		SupplyTime: t,
	}
	fmt.Printf("%v", input)
	err = app.models.Supplies.Update(&supply,id)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = helpers.WriteJSON(w, http.StatusOK, nil, nil)
}

//TODO DELETION
func (app *application) DeleteSuppliesByID(w http.ResponseWriter, r *http.Request) {

	id, err := helpers.ReadIDParam(r)
	if err != nil {
		fmt.Println("id param error")
		return
	}

	err = app.models.Supplies.DeleteByID(id)
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
func (app *application) SupplyTest(w http.ResponseWriter, r *http.Request){
	queryValues := r.URL.Query()
	fromTime, _ := time.Parse("2006-01-02 15:04:05", queryValues.Get("fromTime"))
	toTime, _ := time.Parse("2006-01-02 15:04:05", queryValues.Get("toTime"))
	barcode, err := strconv.Atoi(queryValues.Get("barcode"))
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := app.models.Supplies.GetSupplyAmount1(int64(barcode), fromTime, toTime)

	if err != nil {
		fmt.Println(err)
		return
	}


	totalQuantity,totalRev ,err := app.models.Sale.GetSalesAmount(int64(barcode), fromTime, toTime)
	
	if err != nil {
		fmt.Println(err)
		return
	}
	totalCost := NewSolution(res,totalQuantity)
	totalProfit := totalRev - totalCost 
	err = helpers.WriteJSON(w, http.StatusOK, helpers.Enveleope{"barcode":barcode,"quantity":totalQuantity,"revenue":totalRev,"netProfit":totalProfit}, nil)
	if err != nil{
		fmt.Println(err)
		return
	}

}

func NewSolution(supplies []data.Output,totalQuantity int64)int64{
	var totalCost int64 = 0
	for i := 0;i < len(supplies);i++{
		if supplies[i].Running_Qty >= totalQuantity{
			totalCost = supplies[i].Running_Cost - (supplies[i].Running_Qty - totalQuantity)* supplies[i].Price
		}
	}
	if totalCost == 0{
		totalCost = supplies[len(supplies)-1].Running_Cost
	}
	return totalCost
}

