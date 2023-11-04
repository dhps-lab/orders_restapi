package main

import (
	"log"
	"net/http"

	"github.com/dhps-lab/orders_restapi/db"
	"github.com/dhps-lab/orders_restapi/models"
	"github.com/dhps-lab/orders_restapi/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

//   Product Api:
//    version: 0.1
//    title: Work Orders Api
//   Schemes: http
//   Host:
//   BasePath: /api/v1
//      Consumes:
//      - application/json
//   Produces:
//   - application/json
//   SecurityDefinitions:
//   swagger:meta
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db.DBConnection()
	db.RedisClient()

	db.Database.AutoMigrate(&models.Customer{})
	db.Database.AutoMigrate(&models.WorkOrder{})

	r := mux.NewRouter()

	routes.Swagger(r)
	routes.SetRoutes(r)

	http.ListenAndServe(":3000", r)
}
