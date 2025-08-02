package routes

import (
	"github.com/gorilla/mux"
	"github.com/pratyush934/sibling-bond-server/controller"
	"github.com/pratyush934/sibling-bond-server/utils"
)

// SetupCategoryRoutes configures all category-related routes
func SetupCategoryRoutes(router *mux.Router) {
	// Public routes (no authentication required)
	router.HandleFunc("/api/categories", controller.GetAllCategories).Methods("GET")
	router.HandleFunc("/api/categories/category", controller.GetCategoryById).Methods("GET")

	// Admin routes (admin authentication required)
	adminRoutes := router.PathPrefix("/api/admin/categories").Subrouter()
	adminRoutes.Use(utils.ValidateAdmin)
	adminRoutes.HandleFunc("", controller.CreateCategory).Methods("POST")
	adminRoutes.HandleFunc("", controller.UpdateCategory).Methods("PUT")
	adminRoutes.HandleFunc("", controller.DeleteCategory).Methods("DELETE")
}
