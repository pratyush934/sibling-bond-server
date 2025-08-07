package main

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/pratyush934/sibling-bond-server/cjson"
	"github.com/pratyush934/sibling-bond-server/database"
	"github.com/pratyush934/sibling-bond-server/models"
	"github.com/pratyush934/sibling-bond-server/routes"
	"github.com/pratyush934/sibling-bond-server/utils"
	"gorm.io/gorm"
	"net/http"
)

var (
	httpAddr = ":5000"
)

func LoadDB() {
	if err := database.InitDB(); err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusInternalServerError,
			Message:       "Issue while connecting LOAD DB",
			InternalError: err,
		})
	}

	// Only migrate if the orders table does not exist
	if !database.DB.Migrator().HasTable(&models.Order{}) {
		if err := database.DB.AutoMigrate(
			&models.User{},
			&models.Address{},
			&models.Role{},
			&models.OrderItem{},
			&models.Order{},
			&models.Product{},
			&models.Category{},
			&models.ProductVariant{},
			&models.Cart{},
			&models.CartItem{},
		); err != nil {
			panic(&cjson.HTTPError{
				Status:        http.StatusInternalServerError,
				Message:       "Issue while migrating models to DB",
				InternalError: err,
			})
		}
	}
}

func SeedData() {
	db := database.DB

	// Seed Roles
	roles := []models.Role{
		{Id: 1, RoleName: "User", Description: "Normal User"},
		{Id: 2, RoleName: "Admin", Description: "Administrator"},
		{Id: 3, RoleName: "Tenant", Description: "Tenant"},
	}
	for _, role := range roles {
		var existing models.Role
		if err := db.First(&existing, role.Id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			_ = db.Create(&role)
		}
	}

	users := []models.User{
		{FirstName: "Pratyush", LastName: "Admin", Email: "pratyush@example.com", PassWord: "adminPassword", RoleId: 2},
		{FirstName: "Mridula", LastName: "Admin", Email: "mridula@example.com", PassWord: "adminPassword", RoleId: 2},
		{FirstName: "Akash", LastName: "Tenant", Email: "akash@example.com", PassWord: "tenantPassword", RoleId: 3},
		{FirstName: "Ayush", LastName: "User", Email: "ayush@example.com", PassWord: "userPassword", RoleId: 1},
	}
	for _, user := range users {
		var existing models.User
		if err := db.Where("email = ?", user.Email).First(&existing).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			_ = db.Create(&user)
		}
	}
}

func Server() {

	router := mux.NewRouter()
	router.Use(utils.ErrorHandler)
	router.Use(utils.CORSMiddleware)

	routes.SetupUserRoutes(router)
	routes.SetupCartRoutes(router)
	routes.SetupCategoryRoutes(router)
	routes.SetupProductRoutes(router)
	routes.SetupOrderRoutes(router)
	routes.SetupRoleRoutes(router)

	server := &http.Server{
		Addr:    httpAddr,
		Handler: router,
	}
	if err := server.ListenAndServe(); err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusInternalServerError,
			Message:       "Not able to start the Server",
			InternalError: err,
		})
	}
}

func main() {

	LoadDB()
	SeedData()
	Server()
}
