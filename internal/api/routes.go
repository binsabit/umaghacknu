package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//  create url
//  go to url
//  delete url

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	//The routes for CRUD operations on Supplies
	router.HandlerFunc(http.MethodGet, "/alash/api/supplies", app.GetSupplies)
	router.HandlerFunc(http.MethodGet, "/alash/api/supplies/:id", app.GetSuppliesByID)
	router.HandlerFunc(http.MethodPut, "/alash/api/supplies/:id", app.UpdateSuppliesByID)
	router.HandlerFunc(http.MethodDelete, "/alash/api/supplies/:id", app.DeleteSuppliesByID)
	router.HandlerFunc(http.MethodPost, "/alash/api/supplies", app.AddSupplies)

	

	router.HandlerFunc(http.MethodGet, "/alash/api/sales", app.GetSales)
	router.HandlerFunc(http.MethodGet, "/alash/api/sales/:id", app.GetSalesByID)
	router.HandlerFunc(http.MethodPut, "/alash/api/sales/:id", app.UpdateSalesByID)
	router.HandlerFunc(http.MethodDelete, "/alash/api/sales/:id", app.DeleteSalesByID)
	router.HandlerFunc(http.MethodPost, "/alash/api/sales", app.AddSales)

	router.HandlerFunc(http.MethodGet, "/alash/api/reports", app.Tester)
	router.HandlerFunc(http.MethodGet, "/alash/api/test1", app.SupplyTest)

	// router.HandlerFunc(http.MethodGet, "/alash/api/test", app.Tester)

	return router
}
