package models

import "time"

type Product struct {
	Id            string    `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"not null" json:"name"`
	Description   string    `gorm:"not null" json:"description"`
	Price         int       `gorm:"not null" json:"price"`
	StockQuantity int       `gorm:"not null" json:"stockQuantity"`
	CategoryId    string    `json:"categoryId" gorm:"not null"`
	Category      Category  `json:"category" gorm:"foreignKey:CategoryId"`
	ImageUrl      string    `json:"imageUrl"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
