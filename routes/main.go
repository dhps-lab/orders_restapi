package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetRoutes(r *mux.Router) {

	// Orders routes
	r.HandleFunc("/orders", GetOrdersHandler).Methods(http.MethodGet)
	r.HandleFunc("/orders", CreateOrderHandler).Methods(http.MethodPost)
	r.HandleFunc("/orders/{id}", GetOrderHandler).Methods(http.MethodGet)
	r.HandleFunc("/orders/{id}", UpdateOrderHandler).Methods(http.MethodPut)

	// Customers routes
	r.HandleFunc("/customers", GetCustomersHandler).Methods(http.MethodGet)
	r.HandleFunc("/customers", CreateCustomerHandler).Methods(http.MethodPost)
	r.HandleFunc("/customers/{id}", GetCustomerHandler).Methods(http.MethodGet)
	r.HandleFunc("/customers/{id}", UpdateCustomerHandler).Methods(http.MethodPut)
}
