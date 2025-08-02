package routes

import (
	"github.com/gorilla/mux"
	"github.com/pratyush934/sibling-bond-server/controller"
	"github.com/pratyush934/sibling-bond-server/utils"
)

func SetupProductRoutes(router *mux.Router) {
	// Public routes for product access
	productsRouter := router.PathPrefix("/products").Subrouter()
	productsRouter.HandleFunc("", controller.GetAllProducts).Methods("GET")
	productsRouter.HandleFunc("/search", controller.SearchProduct).Methods("GET")
	productsRouter.HandleFunc("/category", controller.GetProductsByCategory).Methods("GET")
	productsRouter.HandleFunc("/{id}", controller.GetProductById).Methods("GET")

	// Admin-only routes that require authentication
	adminProductsRouter := router.PathPrefix("/admin/products").Subrouter()
	adminProductsRouter.Use(utils.ValidateAdmin)
	adminProductsRouter.HandleFunc("", controller.CreateProduct).Methods("POST")
	adminProductsRouter.HandleFunc("/{id}", controller.UpdateProductDetails).Methods("PUT")
	adminProductsRouter.HandleFunc("/{id}", controller.DeleteProduct).Methods("DELETE")
}
