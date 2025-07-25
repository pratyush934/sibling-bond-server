package models

import "time"

type CartItem struct {
	Id            string    `gorm:"primaryKey" json:"id"`
	CartId        string    `gorm:"not null" json:"cartId"`
	ProductId     string    `gorm:"not null" json:"productId"`
	Product       Product   `gorm:"foreignKey:ProductId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"product"`
	Quantity      int       `gorm:"default:0"  json:"quantity"`
	PriceAtAdding int       `gorm:"default:0" json:"priceAtAdding"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
