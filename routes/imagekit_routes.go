package routes

import (
	"github.com/gorilla/mux"
	"github.com/pratyush934/sibling-bond-server/ikprovider"
	"github.com/pratyush934/sibling-bond-server/utils"
)

func SetUpImageKitRoutes(router *mux.Router) {

	imagekitRoutes := router.PathPrefix("/api/admin").Subrouter()
	imagekitRoutes.Use(utils.ValidateAdmin)

	imagekitRoutes.HandleFunc("/images", ikprovider.GetImageKitAuthHandler).Methods("POST")
}
