package models

import "time"

type Order struct {
	Id                string    `gorm:"primaryKey" json:"id"`
	UserId            string    `gorm:"foreignKey not null" json:"userId"`
	OrderedAt         time.Time `json:"orderedAt"`
	TotalAmount       int       `json:"totalAmount"`
	ShippingAddressId string    `gorm:"not null" json:"shippingAddressId"`
	ShippingAddress   Address   `gorm:"constraint:onUpdate:CASCADE onDelete:CASCADE"  json:"address"`
	Status            string    `json:"status"`
	PaymentStatus     string    `json:"paymentStatus"`
	PaymentMode       string    `json:"paymentMode"`
	TrackingNumber    int       `json:"trackingNumber"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
