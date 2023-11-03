package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dhps-lab/orders_restapi/db"
	"github.com/dhps-lab/orders_restapi/models"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CustomerInput struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Address   string `json:"address" validate:"required"`
}

func GetCustomersHandler(w http.ResponseWriter, r *http.Request) {
	customers := []models.Customer{}
	db.Database.Find(&customers)
	respondJSON(w, http.StatusOK, customers)
}

func GetCustomerHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]
	customer := getCustomerOrErrorById(id, w, r)
	if customer == nil {
		return
	}
	respondJSON(w, http.StatusOK, customer)
}
func CreateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	var input CustomerInput
	json.NewDecoder(r.Body).Decode(&input)
	validate = validator.New()
	err := validate.Struct(input)
	if err != nil {
		log.Println(err)
		respondError(w, http.StatusBadRequest, "Validation Error")
		return
	}

	var customer models.Customer
	customer.Id = uuid.New().String()
	customer.FirstName = input.FirstName
	customer.LastName = input.LastName
	customer.Address = input.Address
	customer.IsActive = false

	createdCustomer := db.Database.Create(&customer)
	err = createdCustomer.Error
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, customer)
}

func UpdateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]
	customer := getCustomerOrErrorById(id, w, r)
	if customer == nil {
		return
	}

	var input CustomerInput
	json.NewDecoder(r.Body).Decode(&input)
	validate = validator.New()
	err := validate.Struct(input)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Validation Error")
		return
	}

	customer.FirstName = input.FirstName
	customer.LastName = input.LastName
	customer.Address = input.Address

	updatedCustomer := db.Database.Save(&customer)
	err = updatedCustomer.Error
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, customer)
}

func DeleteCustomerHandler(w http.ResponseWriter, r *http.Request) {}

func getCustomerOrErrorById(id string, w http.ResponseWriter, r *http.Request) *models.Customer {
	customer := &models.Customer{}
	err := db.Database.First(&customer, models.Customer{Id: id}).Error
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	return customer
}
