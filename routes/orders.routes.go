package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dhps-lab/orders_restapi/db"
	"github.com/dhps-lab/orders_restapi/models"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var validate *validator.Validate

type WorkOrderCreateInput struct {
	CustomerId       string    `json:"customer_id"`
	Title            string    `json:"title" validate:"required"`
	PlannedDateBegin time.Time `json:"planned_date_begin" validate:"required"`
	PlannedDateEnd   time.Time `json:"planned_date_end" validate:"required"`
}

func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	var orders []models.WorkOrder
	db.Database.Find(&orders)
	respondJSON(w, http.StatusOK, orders)
}

func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {

	var input WorkOrderCreateInput
	json.NewDecoder(r.Body).Decode(&input)
	validate = validator.New()
	err := validate.Struct(input)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Validation Error")
		return
	}

	var order models.WorkOrder
	order.Id = uuid.New().String()
	order.Title = input.Title
	order.PlannedDateBegin = input.PlannedDateBegin
	order.PlannedDateEnd = input.PlannedDateEnd
	order.Status = models.New

	createdOrder := db.Database.Create(&order)
	log.Println(createdOrder.RowsAffected)
	err = createdOrder.Error
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	json.NewEncoder(w).Encode(order)
}

func GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]
	workOrder := getWorkOrderOrErrorById(id, w, r)
	if workOrder == nil {
		return
	}
	respondJSON(w, http.StatusOK, workOrder)
}

func UpdateOrderHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]
	workOrder := getWorkOrderOrErrorById(id, w, r)
	if workOrder == nil {
		return
	}

	input := &WorkOrderCreateInput{}
	json.NewDecoder(r.Body).Decode(&input)
	validate = validator.New()
	err := validate.Struct(input)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Validation Error")
		return
	}

	workOrder.Title = input.Title
	workOrder.CustomerId = input.CustomerId
	workOrder.PlannedDateBegin = input.PlannedDateBegin
	workOrder.PlannedDateEnd = input.PlannedDateEnd
	err = db.Database.Save(&workOrder).Error
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Internal server Error")
		return
	}
	respondJSON(w, http.StatusOK, workOrder)
}

func getWorkOrderOrErrorById(id string, w http.ResponseWriter, r *http.Request) *models.WorkOrder {
	workOrder := &models.WorkOrder{}
	err := db.Database.First(&workOrder, models.WorkOrder{Id: id}).Error
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	return workOrder
}
