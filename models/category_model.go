package models

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/pratyush934/sibling-bond-server/database"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type Category struct {
	Id          string    `gorm:"primaryKey;type:varchar(150)" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Products    []Product `gorm:"foreignKey:category_id" json:"products,omitempty"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

/*
Create(category *Category) (*Category, error)

GetByID(id string) (*Category, error)

GetByName(name string) (*Category, error)

Update(category *Category) (*Category, error)

Delete(id string) error

GetAll() ([]*Category, error)

*/

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	c.Id = uuid.New().String()
	return nil
}

func (c *Category) CreateCategory() (*Category, error) {
	if err := database.DB.Create(c).Error; err != nil {
		log.Err(err).Msg("Issue exist in CreateCategory")
		return nil, err
	}
	return c, nil
}

func GetCategoryById(id string) (*Category, error) {
	var category Category
	if err := database.DB.Where(&Category{Id: id}).First(&category).Error; err != nil {
		log.Err(err).Msg("Issue exist in GetCategoryById")
		return nil, err
	}
	return &category, nil
}

func GetCategoryByName(name string) (*Category, error) {
	var category Category
	if err := database.DB.Where(&Category{Name: name}).First(&category).Error; err != nil {
		// Only log if it's not the expected "record not found" error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Err(err).Msg("Unexpected error in GetCategoryByName")
		}
		return nil, err
	}
	return &category, nil
}

func UpdateCategory(category *Category) (*Category, error) {
	if err := database.DB.Updates(category).Error; err != nil {
		log.Err(err).Msg("Issue exist in UpdateCategory")
		return nil, err
	}
	return category, nil
}

func DeleteCategory(id string) error {
	return database.DB.Where(&Category{Id: id}).Delete(&Category{}).Error
}

func GetAll(limit, offSet int) (*[]Category, error) {
	var categories []Category
	if err := database.DB.Find(&categories).Limit(limit).Offset(offSet).Error; err != nil {
		return nil, err
	}
	return &categories, nil
}

func CategoryHasProducts(categoryId string) (bool, error) {
	var count int64
	if err := database.DB.Model(&Product{}).Where("category_id = ?", categoryId).Count(&count).Error; err != nil {
		return false, err
	}
	fmt.Println("The count is there and it is ", count)
	return count > 0, nil
}
