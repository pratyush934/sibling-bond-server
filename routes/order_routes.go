package routes

import (
	"github.com/gorilla/mux"
	"github.com/pratyush934/sibling-bond-server/controller"
	"github.com/pratyush934/sibling-bond-server/utils"
)

// SetupOrderRoutes configures all order-related routes
func SetupOrderRoutes(router *mux.Router) {
	// User routes (require authentication)
	userOrderRoutes := router.PathPrefix("/api/orders").Subrouter()
	userOrderRoutes.Use(utils.ValidateUser)
	userOrderRoutes.HandleFunc("", controller.CreateOrder).Methods("POST")
	userOrderRoutes.HandleFunc("/cancel", controller.CancelOrder).Methods("POST")
	userOrderRoutes.HandleFunc("/payment", controller.ProcessPayment).Methods("POST")

	// Note: GetOrderHistory and GetOrderDetails are already defined in user_routes.go
	// under /api/users/orders and /api/users/orders/{id}

	// Admin routes (require admin authentication)
	adminOrderRoutes := router.PathPrefix("/api/admin/orders").Subrouter()
	adminOrderRoutes.Use(utils.ValidateAdmin)
	adminOrderRoutes.HandleFunc("", controller.GetAllOrders).Methods("GET")
	adminOrderRoutes.HandleFunc("/status", controller.UpdateOrderStatus).Methods("PUT")
}
