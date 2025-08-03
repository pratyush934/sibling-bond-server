package models

import (
	"github.com/google/uuid"
	"github.com/pratyush934/sibling-bond-server/database"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type Cart struct {
	Id        string     `gorm:"primaryKey;type:varchar(191)" json:"id"`
	UserId    string     `gorm:"not null" json:"userId"`
	User      User       `gorm:"foreignKey:UserId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"user"`
	CartItems []CartItem `gorm:"foreignKey:CartId" json:"cartItems"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

/*

GetByUserID(userID string) (*Cart, error)

Create(cart *Cart) (*Cart, error)

Delete(cartID string) error

Save(cart *Cart) (*Cart, error)
*/

func (c *Cart) BeforeCreate(t *gorm.DB) error {
	c.Id = uuid.New().String()
	return nil
}

func Create(cart *Cart) (*Cart, error) {
	if err := database.DB.Create(cart).Error; err != nil {
		log.Err(err).Msg("issue exist in Create in Cart")
		return nil, err
	}
	return cart, nil
}

func GetCartByUserId(userId string) (*Cart, error) {
	var cart Cart
	if err := database.DB.Preload("CartItems.Product").Where(&Cart{UserId: userId}).First(&cart).Error; err != nil {
		log.Err(err).Msg("Issue persist in GetCartByUserId")
		return nil, err
	}
	return &cart, nil
}

func GetCartById(cartId string) (*Cart, error) {
	var cart Cart
	if err := database.DB.Preload("CartItems.Product").Where(&Cart{Id: cartId}).First(&cart).Error; err != nil {
		log.Err(err).Msg("Issue persist in GetCartById")
		return nil, err
	}
	return &cart, nil
}

func DeleteCart(cartId string) error {
	return database.DB.Where(&Cart{Id: cartId}).Delete(&Cart{}).Error
}

func UpdateCart(cart *Cart) (*Cart, error) {
	if err := database.DB.Save(cart).Error; err != nil {
		log.Err(err).Msg("issue exist in Create in Cart")
		return nil, err
	}
	return cart, nil
}
