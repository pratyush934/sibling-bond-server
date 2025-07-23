package main

import (
	"github.com/gorilla/mux"
	"github.com/pratyush934/sibling-bond-server/database"
	"github.com/pratyush934/sibling-bond-server/models"
	"github.com/pratyush934/sibling-bond-server/utils"
	"net/http"
)

var (
	httpAddr = ":5000"
)

func LoadDB() {
	if err := database.InitDB(); err != nil {
		panic(models.HTTPError{
			Status:        http.StatusInternalServerError,
			Message:       "Issue while connecting LOAD DB",
			InternalError: err,
		})
	}
	if err := database.DB.AutoMigrate(&models.User{}, &models.Address{}, &models.Role{}, &models.OrderItem{}, &models.Order{}, &models.Product{}, &models.Product{}, &models.Category{}); err != nil {
		panic(models.HTTPError{
			Status:        http.StatusInternalServerError,
			Message:       "Issue while migrating models to DB",
			InternalError: err,
		})
	}
}

func SeedData() {

}

func Server() {

	router := mux.NewRouter()
	router.Use(utils.ErrorHandler)

	server := &http.Server{
		Addr:    httpAddr,
		Handler: router,
	}
	if err := server.ListenAndServe(); err != nil {
		panic(models.NewError(http.StatusInternalServerError, "Issue while starting the server", err))
	}
}

func main() {

	LoadDB()
	SeedData()
	Server()
}
