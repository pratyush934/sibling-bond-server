package models

import "time"

type OrderItem struct {
	Id              string    `gorm:"primaryKey" json:"id"`
	OrderId         string    `gorm:"not null" json:"orderId"`
	ProductId       string    `gorm:"not null" json:"productId"`
	Quantity        int       `gorm:"not null" json:"quantity"`
	PriceAtPurchase int       `gorm:"not null" json:"priceAtPurchase"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
