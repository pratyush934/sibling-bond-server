package routes

import (
	"github.com/gorilla/mux"
	"github.com/pratyush934/sibling-bond-server/controller"
	"github.com/pratyush934/sibling-bond-server/utils"
)

// SetupCartRoutes configures all cart-related routes
func SetupCartRoutes(router *mux.Router) {
	// All cart routes require user authentication
	cartRoutes := router.PathPrefix("/api/cart").Subrouter()
	cartRoutes.Use(utils.ValidateUser)

	// Cart operations
	cartRoutes.HandleFunc("", controller.GetCart).Methods("GET")
	cartRoutes.HandleFunc("/items", controller.AddToCart).Methods("POST")
	cartRoutes.HandleFunc("/items/{id}", controller.UpdateCartItem).Methods("PUT")
	cartRoutes.HandleFunc("/items/{id}", controller.RemoveFromCart).Methods("DELETE")
	cartRoutes.HandleFunc("", controller.ClearCart).Methods("DELETE")

	// Cart summary/checkout preparation
	//cartRoutes.HandleFunc("/summary", controller.GetCartSummary).Methods("GET")
}
