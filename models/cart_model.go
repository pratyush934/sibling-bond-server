package models

import "time"

type Cart struct {
	Id        string     `gorm:"primaryKey" json:"id"`
	UserId    string     `gorm:"not null" json:"userId"`
	User      User       `gorm:"foreignKey:UserId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"user"`
	CartItems []CartItem `gorm:"foreignKey:CartId" json:"cartItems"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}
