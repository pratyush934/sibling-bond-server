package db

import (
	"github.com/pratyush934/sibling-bond-server/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)

var DB *gorm.DB

func InitDB() error {
	db, err := connectToDB()
	if err != nil {
		panic(models.HTTPError{Status: http.StatusBadRequest, Message: "Unable to connect to the DB", InternalError: err})
		return err
	}
	DB = db
	return nil
}

func connectToDB() (*gorm.DB, error) {
	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(models.HTTPError{
			Status:        http.StatusInternalServerError,
			Message:       "Unable to connect the DB, error in connectToDB",
			InternalError: err,
		})
		return nil, err
	}
	return db, nil
}
