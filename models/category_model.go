package models

import "time"

type Category struct {
	Id          string    `gorm:"primaryKey" json:"id"`
	ProductId   string    `gorm:"foreignKey" json:"productId"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Products    []Product `gorm:"foreignKey:CategoryId" json:"products,omitempty"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
