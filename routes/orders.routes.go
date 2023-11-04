package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/dhps-lab/orders_restapi/db"
	"github.com/dhps-lab/orders_restapi/models"
	"github.com/dhps-lab/orders_restapi/utils"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var validate *validator.Validate

type WorkOrderCreateInput struct {
	CustomerId       string    `json:"customer_id" validate:"required"`
	Title            string    `json:"title" validate:"required"`
	PlannedDateBegin time.Time `json:"planned_date_begin" validate:"required"`
	PlannedDateEnd   time.Time `json:"planned_date_end" validate:"required"`
}

type filterInput struct {
	Since  string
	Until  string
	Status string `validate:"status"`
}

func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	var orders []models.WorkOrder
	db.Database.Preload("Customer").Find(&orders)
	respondJSON(w, http.StatusOK, orders)
}

func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {

	var input WorkOrderCreateInput
	json.NewDecoder(r.Body).Decode(&input)
	validate = validator.New()
	err := validate.Struct(input)
	differenceHours := input.PlannedDateEnd.Sub(input.PlannedDateBegin).Hours()
	if err != nil || !(differenceHours <= 2 && differenceHours >= 0) {
		respondError(w, http.StatusBadRequest, "Validation Error")
		return
	}

	customer := getCustomerOrErrorById(input.CustomerId, w, r)
	if customer == nil {
		return
	}

	var order models.WorkOrder
	order.Id = uuid.New().String()
	order.Title = input.Title
	order.CustomerId = input.CustomerId
	order.PlannedDateBegin = input.PlannedDateBegin
	order.PlannedDateEnd = input.PlannedDateEnd
	order.Status = models.New

	createdOrder := db.Database.Create(&order)
	err = createdOrder.Error
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, order)
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

func completeOrderHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]
	workOrder := getWorkOrderOrErrorById(id, w, r)
	if workOrder == nil {
		return
	}

	isFirst := workOrderIsFirstByCustomer(workOrder.Customer.Id, w, r)
	if isFirst {
		customer := activateCustomer(workOrder.Customer.Id, w, r)
		if customer == nil {
			respondError(w, http.StatusInternalServerError, "Internal server Error")
			return
		}
		workOrder = getWorkOrderOrErrorById(id, w, r)
		if workOrder == nil {
			return
		}
	}

	if workOrder.Status != models.New {
		respondError(w, http.StatusBadRequest, "Bad request")
		return
	}

	workOrder.Status = models.Done
	err := db.Database.Save(&workOrder).Error
	if err != nil {
		log.Println(err)
		respondError(w, http.StatusInternalServerError, "Internal server Error")
		return
	}
	log.Print("I'm working with redis")
	err = db.Publish_order(*workOrder)
	if err != nil {
		log.Println(err)
		respondError(w, http.StatusInternalServerError, "Internal server Error")
		return
	}

	respondJSON(w, http.StatusOK, workOrder)
}

func cancelOrderHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]
	workOrder := getWorkOrderOrErrorById(id, w, r)
	if workOrder == nil {
		return
	}

	if workOrder.Status == models.Done {
		respondError(w, http.StatusBadRequest, "BadRequest server Error")
		return
	}
	workOrder.Status = models.Cancelled
	err := db.Database.Save(&workOrder).Error
	if err != nil {
		log.Println(err)
		respondError(w, http.StatusInternalServerError, "Internal server Error")
		return
	}

	respondJSON(w, http.StatusOK, workOrder)
}

func getOrdersByCustomerHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]
	ordersByCustomer := getWorkOrdersOrErrorByCustomerId(id, w, r)
	if ordersByCustomer == nil {
		if len(ordersByCustomer) == 0 {
			respondError(w, http.StatusBadRequest, "BadRequest server Error")
			return
		}
		respondError(w, http.StatusInternalServerError, "Internal server Error")
		return
	}

	respondJSON(w, http.StatusOK, ordersByCustomer)
}

func deactivateCustomerOrderHandler(w http.ResponseWriter, r *http.Request) {
	order := &models.WorkOrder{}
	params := mux.Vars(r)

	id := params["id"]
	customer := getCustomerOrErrorById(id, w, r)
	if customer == nil {
		return
	}

	customer = deactivateCustomer(customer.Id, w, r)
	if customer == nil {
		respondError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	order.Id = uuid.New().String()
	order.Title = models.DeactivateCustomer
	order.CustomerId = customer.Id
	order.PlannedDateEnd = time.Now()
	order.PlannedDateBegin = time.Now().Add(time.Duration(-1) * time.Hour)
	order.Status = models.Done

	createdOrder := db.Database.Create(&order)
	log.Println(createdOrder.RowsAffected)
	err := createdOrder.Error
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, order)
}

func GetOrdersByFilterHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	since := queryParams.Get("since")
	until := r.URL.Query().Get("until")
	status := r.URL.Query().Get("status")

	regexpStatus := regexp.MustCompile("(new|done|cancelled)$")

	timeSince, errSince := utils.IsDateValue(since)
	timeUntil, errUntil := utils.IsDateValue(until)
	comparer := timeSince.Compare(timeUntil)

	if since == "" && until == "" && status == "" {
		respondError(w, http.StatusBadRequest, "Bad params")
		return
	}

	if !regexpStatus.MatchString(status) && status != "" {
		log.Println("Error in status filter")
		respondError(w, http.StatusBadRequest, "Bad params")
		return
	}

	if (since != "" || until != "") && (!errSince || !errUntil) {
		log.Println("Error in date filters")
		respondError(w, http.StatusBadRequest, "Bad params")
		return
	}

	var orders []models.WorkOrder
	var dates *gorm.DB
	if comparer == -1 {
		dates = db.Database.Where("planned_date_begin >= ?", timeSince).Where("planned_date_end <= ?", timeUntil)
	} else {
		dates = db.Database
	}
	var statusDB *gorm.DB
	if status != "" {
		statusDB = db.Database.Where("status = ?", status)
	} else {
		statusDB = db.Database
	}
	dates.Where(statusDB).Preload("Customer").Find(&orders)

	respondJSON(w, http.StatusOK, orders)
}

func getWorkOrderOrErrorById(id string, w http.ResponseWriter, r *http.Request) *models.WorkOrder {
	workOrder := &models.WorkOrder{}
	err := db.Database.Preload("Customer").First(&workOrder, models.WorkOrder{Id: id}).Error
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	return workOrder
}

func workOrderIsFirstByCustomer(customerId string, w http.ResponseWriter, r *http.Request) bool {
	ordersByCustomer := getWorkOrdersOrErrorByCustomerId(customerId, w, r)
	var isFirst bool
	isFirst = true
	for _, order := range ordersByCustomer {
		if order.Status == models.Done || order.Customer.IsActive {
			isFirst = false
			log.Println("It's not the first order", order.Status, " and ", order.Customer.IsActive)
			return isFirst
		}
	}
	log.Println("It's first order")
	return isFirst
}
func getWorkOrdersOrErrorByCustomerId(id string, w http.ResponseWriter, r *http.Request) []models.WorkOrder {
	workOrder := []models.WorkOrder{}
	err := db.Database.Preload("Customer").Find(&workOrder, models.WorkOrder{CustomerId: id}).Error
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	return workOrder
}
