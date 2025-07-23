package database

import (
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	db, err := connectToDB()
	if err != nil {
		log.Err(err).Msg("Issue while InitDB")
		return err
	}
	DB = db
	return nil
}

func connectToDB() (*gorm.DB, error) {
	dsn := "root:Pratyush@123@tcp(127.0.0.1:3306)/siblingbond?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Err(err).Msg("Issue while connecting the DB connectToDB")
		return nil, err
	}
	return db, nil
}
