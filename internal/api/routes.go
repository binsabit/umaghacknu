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

	return router
}
