package routes

import (
	"github.com/gorilla/mux"
	"github.com/pratyush934/sibling-bond-server/controller"
	"github.com/pratyush934/sibling-bond-server/database"
)

func SetupRoleRoutes(router *mux.Router) {
	roleController := controller.NewRoleController(database.DB)

	router.HandleFunc("/roles", roleController.CreateRole).Methods("POST")
	router.HandleFunc("/roles/{id}", roleController.GetRoleByID).Methods("GET")
	router.HandleFunc("/roles", roleController.GetAllRoles).Methods("GET")
	router.HandleFunc("/roles/{id}", roleController.UpdateRole).Methods("PUT")
	router.HandleFunc("/roles/{id}", roleController.DeleteRole).Methods("DELETE")
}
