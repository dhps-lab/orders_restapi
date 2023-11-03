package main

import (
	"net/http"

	"github.com/dhps-lab/orders_restapi/db"
	"github.com/dhps-lab/orders_restapi/models"
	"github.com/dhps-lab/orders_restapi/routes"
	"github.com/gorilla/mux"
)

func main() {
	db.DBConnection()

	db.Database.AutoMigrate(&models.Customer{})
	db.Database.AutoMigrate(&models.WorkOrder{})

	r := mux.NewRouter()

	routes.SetRoutes(r)

	http.ListenAndServe(":3000", r)
}
