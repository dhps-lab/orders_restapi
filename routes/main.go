package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetRoutes(r *mux.Router) {

	// Orders routes
	r.HandleFunc("/orders", GetOrdersHandler).Methods(http.MethodGet)
	r.HandleFunc("/order", CreateOrderHandler).Methods(http.MethodPost)
	r.HandleFunc("/orders/filter", GetOrdersByFilterHandler).Methods(http.MethodGet)
	r.HandleFunc("/order/{id}", GetOrderHandler).Methods(http.MethodGet)
	r.HandleFunc("/order/{id}", UpdateOrderHandler).Methods(http.MethodPut)
	r.HandleFunc("/order/{id}/complete", completeOrderHandler).Methods(http.MethodPut)
	r.HandleFunc("/order/{id}/cancel", cancelOrderHandler).Methods(http.MethodPut)
	r.HandleFunc("/orders/customer/{id}", getOrdersByCustomerHandler).Methods(http.MethodGet)
	r.HandleFunc("/order/customer/{id}/deactivate", deactivateCustomerOrderHandler).Methods(http.MethodPost)

	// Customers routes
	r.HandleFunc("/customers", GetCustomersHandler).Methods(http.MethodGet)
	r.HandleFunc("/customer", CreateCustomerHandler).Methods(http.MethodPost)
	r.HandleFunc("/customers/active", getAllCustomersActive).Methods(http.MethodGet)
	r.HandleFunc("/customer/{id}", GetCustomerHandler).Methods(http.MethodGet)
	r.HandleFunc("/customer/{id}", UpdateCustomerHandler).Methods(http.MethodPut)
}
