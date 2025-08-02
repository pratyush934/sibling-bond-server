package routes

import (
	"github.com/gorilla/mux"
	"github.com/pratyush934/sibling-bond-server/controller"
	"github.com/pratyush934/sibling-bond-server/utils"
)

// SetupUserRoutes configures all user-related routes
func SetupUserRoutes(router *mux.Router) {
	// Public routes (no authentication required)
	router.HandleFunc("/api/users/register", controller.Register).Methods("POST")
	router.HandleFunc("/api/users/login", controller.Login).Methods("POST")
	router.HandleFunc("/api/users/forgot-password", controller.ForgotPassWord).Methods("POST")
	router.HandleFunc("/api/users/reset-password", controller.ResetPasswordFromToken).Methods("POST")

	// Authenticated user routes
	userRoutes := router.PathPrefix("/api/users").Subrouter()
	userRoutes.Use(utils.ValidateUser) // Assuming you have this middleware to validate JWT
	userRoutes.HandleFunc("/logout", controller.LogOut).Methods("POST")
	userRoutes.HandleFunc("/profile", controller.GetProfile).Methods("GET")
	userRoutes.HandleFunc("/change-password", controller.ChangePassword).Methods("POST")

	// Address routes
	userRoutes.HandleFunc("/addresses", controller.GetAddresses).Methods("GET")
	userRoutes.HandleFunc("/addresses", controller.AddAddress).Methods("POST")
	userRoutes.HandleFunc("/addresses", controller.UpdateAddress).Methods("PUT")
	userRoutes.HandleFunc("/addresses/{id}", controller.DeleteAddress).Methods("DELETE")

	// Order routes
	userRoutes.HandleFunc("/orders", controller.GetOrderHistory).Methods("GET")
	userRoutes.HandleFunc("/orders/{id}", controller.GetOrderDetails).Methods("GET")

	// Admin routes
	adminRoutes := router.PathPrefix("/api/admin/users").Subrouter()
	adminRoutes.Use(utils.ValidateAdmin) // This should verify both auth and admin role
	adminRoutes.HandleFunc("", controller.GetAllUsersByAdmin).Methods("GET")
	adminRoutes.HandleFunc("/{id}", controller.GetUserById).Methods("GET")
	adminRoutes.HandleFunc("/{id}", controller.DeleteUserById).Methods("DELETE")
}
