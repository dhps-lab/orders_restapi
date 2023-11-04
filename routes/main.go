package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetRoutes(r *mux.Router) {

	routes := r.PathPrefix("/api/v1").Subrouter()

	// Orders routes
	routes.HandleFunc("/orders", GetOrdersHandler).Methods(http.MethodGet)
	routes.HandleFunc("/order", CreateOrderHandler).Methods(http.MethodPost)
	routes.HandleFunc("/orders/filter", GetOrdersByFilterHandler).Methods(http.MethodGet)
	routes.HandleFunc("/order/{id}", GetOrderHandler).Methods(http.MethodGet)
	routes.HandleFunc("/order/{id}", UpdateOrderHandler).Methods(http.MethodPut)
	routes.HandleFunc("/order/{id}/complete", completeOrderHandler).Methods(http.MethodPut)
	routes.HandleFunc("/order/{id}/cancel", cancelOrderHandler).Methods(http.MethodPut)
	routes.HandleFunc("/orders/customer/{id}", getOrdersByCustomerHandler).Methods(http.MethodGet)
	routes.HandleFunc("/order/customer/{id}/deactivate", deactivateCustomerOrderHandler).Methods(http.MethodPost)

	// Customers routes
	routes.HandleFunc("/customers", GetCustomersHandler).Methods(http.MethodGet)
	routes.HandleFunc("/customer", CreateCustomerHandler).Methods(http.MethodPost)
	routes.HandleFunc("/customers/active", getAllCustomersActive).Methods(http.MethodGet)
	routes.HandleFunc("/customer/{id}", GetCustomerHandler).Methods(http.MethodGet)
	routes.HandleFunc("/customer/{id}", UpdateCustomerHandler).Methods(http.MethodPut)
}
