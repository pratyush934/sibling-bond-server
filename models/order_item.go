package models

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pratyush934/sibling-bond-server/database"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type OrderItem struct {
	Id              string    `gorm:"primaryKey" json:"id"`
	OrderId         string    `gorm:"not null" json:"orderId"`
	ProductId       string    `gorm:"not null" json:"productId"`
	Order           Order     `gorm:"foreignKey:OrderId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"order"`
	Product         Product   `gorm:"foreignKey:ProductId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"product"`
	Quantity        int       `gorm:"not null" json:"quantity"`
	PriceAtPurchase int       `gorm:"not null" json:"priceAtPurchase"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

/*
Create(orderItem *OrderItem) (*OrderItem, error)

CreateMany(orderItems []*OrderItem) error

GetByOrderID(orderID string) ([]*OrderItem, error)

*/

func (o *OrderItem) BeforeCreate(t *gorm.DB) error {
	o.Id = uuid.New().String()
	return nil
}

func (o *OrderItem) Create(orderItem *OrderItem) (*OrderItem, error) {
	if err := database.DB.Create(orderItem).Error; err != nil {
		log.Err(err).Msg("Issue persist in the Create")
		return nil, err
	}
	return orderItem, nil
}

func (oi *OrderItem) ValidateOrderItem() error {

	var product Product
	if err := database.DB.Where("id = ?", oi.ProductId).First(&product).Error; err != nil {
		return fmt.Errorf("Product Not found %v ", err)
	}

	/* check stock availability */
	if product.StockQuantity < oi.Quantity {
		return fmt.Errorf("issue persist as the quanity is of product is less")
	}

	if oi.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}

	if oi.PriceAtPurchase != product.Price {
		return fmt.Errorf("the price must be same in both product and order model")
	}
}

func (oi *OrderItem) UpdateProductStock(tx *gorm.DB) error {
	return tx.Model(&Product{}).Where("id = ?", oi.ProductId).UpdateColumn("stock", gorm.Expr("stock - ?", oi.Quantity)).Error

}

func CreateMany(orderItem []*OrderItem) error {
	if err := database.DB.Create(orderItem).Error; err != nil {
		log.Err(err).Msg("issue persist in CreateMany")
		return err
	}
	return nil
}

func GetByOrderID(orderID string) ([]*OrderItem, error) {
	var orderItems []*OrderItem
	if err := database.DB.Where("order_id = ?", orderID).Find(&orderItems).Error; err != nil {
		log.Err(err).Msg("Issue exist in GetByOrderID")
		return nil, err
	}
	return orderItems, nil
}
