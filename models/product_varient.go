package models

import "time"

type ProductVariant struct {
	Id           string    `gorm:"primaryKey" json:"id"`
	ProductId    string    `gorm:"not null" json:"productId"`
	Product      Product   `gorm:"foreignKey:ProductId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"product"`
	VariantName  string    `gorm:"not null" json:"variantName"`
	VariantValue string    `gorm:"not null" json:"variantValue"`
	Price        int       `json:"price"`
	Stock        int       `gorm:"not null;default:0" json:"stock"`
	SKU          string    `gorm:"unique" json:"sku"`
	IsActive     bool      `gorm:"default:true" json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
