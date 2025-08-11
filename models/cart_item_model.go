package models

import (
	"github.com/google/uuid"
	"github.com/pratyush934/sibling-bond-server/database"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type CartItem struct {
	Id            string    `gorm:"primaryKey;type:varchar(191)" json:"id"`
	CartId        string    `gorm:"not null" json:"cartId"`
	ProductId     string    `gorm:"not null" json:"productId"`
	Product       Product   `gorm:"foreignKey:ProductId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"product"`
	Quantity      int       `gorm:"default:0"  json:"quantity"`
	PriceAtAdding int       `gorm:"default:0" json:"priceAtAdding"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

/*
AddItem(cartItem *CartItem) (*CartItem, error) (Handles both new additions and quantity increments)

GetItemByCartAndProduct(cartID, productID string) (*CartItem, error)

UpdateItemQuantity(cartItemID string, newQuantity int) (*CartItem, error)

RemoveItem(cartItemID string) error

GetItemsByCartID(cartID string) ([]*CartItem, error)

DeleteAllItemsByCartID(cartID string) error

*/

func (ct *CartItem) BeforeCreate(t *gorm.DB) error {
	ct.Id = uuid.New().String()
	return nil
}

func AddItem(item *CartItem) (*CartItem, error) {
	if err := database.DB.Create(item).Error; err != nil {
		log.Err(err).Msg("Issue persist in AddItem")
		return nil, err
	}
	return item, nil
}

func GetItemByCartAndProduct(cartID, productId string) (*CartItem, error) {
	var cartItem CartItem
	if err := database.DB.Preload("Product").Where(&CartItem{CartId: cartID, ProductId: productId}).First(&cartItem).Error; err != nil {

		log.Err(err).Msg("Issue persist in GetItemByCartAndProduct")
		return nil, err
	}
	return &cartItem, nil
}

func UpdateItemQuantity(cartItemId string, newQuantity int) (*CartItem, error) {
	var cartItem CartItem
	if err := database.DB.Model(&CartItem{}).Where(&CartItem{Id: cartItemId}).Update("quantity", newQuantity).Error; err != nil {

		log.Err(err).Msg("Issue persist in UpdateItemQuantity")
		return nil, err
	}
	if err := database.DB.Where(&CartItem{Id: cartItemId}).First(&cartItem).Error; err != nil {
		log.Err(err).Msg("Issue persist in UpdateItemQuantity Part 2")
		return nil, err
	}
	return &cartItem, nil
}

func RemoveItem(cartItemId string) error {
	return database.DB.Where(&CartItem{Id: cartItemId}).Delete(&CartItem{}).Error
}

func IncrementItemQuantity(cartItemId string, quantityToAdd int) (*CartItem, error) {
	var cartItem CartItem
	if err := database.DB.Model(&cartItem).Where("id = ?", cartItemId).Update("quantity", gorm.Expr("quantity + ?", quantityToAdd)).Error; err != nil {
		log.Err(err).Msg("Issue persist in IncrementItemQuantity")
		return nil, err
	}
	if err := database.DB.Where(&CartItem{Id: cartItemId}).First(&cartItem).Error; err != nil {
		log.Err(err).Msg("Issue persist in IncrementItemQuantity Part 2")
		return nil, err
	}
	return &cartItem, nil
}

func GetItemsByCartId(cartId string) ([]*CartItem, error) {
	var cartItem []*CartItem
	if err := database.DB.Preload("Product").Where(&CartItem{CartId: cartId}).Find(&cartItem).Error; err != nil {
		log.Err(err).Msg("Issue persist in GetItemsByCartIt")
		return nil, err
	}
	return cartItem, nil
}

func DeleteAllItemsByCartId(cartId string) error {
	return database.DB.Where(&CartItem{CartId: cartId}).Delete(&CartItem{}).Error
}
