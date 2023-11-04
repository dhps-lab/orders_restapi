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

	routes.SetRoutes(r)

	http.ListenAndServe(":3000", r)
}
