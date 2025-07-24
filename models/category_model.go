package models

import (
	"github.com/google/uuid"
	"github.com/pratyush934/sibling-bond-server/database"
	"gorm.io/gorm"
	"time"
)

type Category struct {
	Id          string    `gorm:"primaryKey" json:"id"`
	ProductId   string    `gorm:"foreignKey" json:"productId"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Products    []Product `gorm:"foreignKey:CategoryId" json:"products,omitempty"`
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
		return &Category{}, err
	}
	return c, nil
}

func GetCategoryById(id string) (*Category, error) {
	var category Category
	if err := database.DB.Where(&Category{Id: id}).First(&category).Error; err != nil {
		return &category, err
	}
	return &category, nil
}
