package models

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pratyush934/sibling-bond-server/database"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type Order struct {
	Id                string      `gorm:"primaryKey;type:varchar(191)" json:"id"`
	UserId            string      `gorm:"not null" json:"userId"`
	User              User        `gorm:"foreignKey:UserId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"user"`
	OrderItems        []OrderItem `gorm:"foreignKey:OrderId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"orderItems"`
	OrderedAt         time.Time   `json:"orderedAt"`
	TotalAmount       int         `json:"totalAmount"`
	ShippingAddressId string      `gorm:"type:varchar(191);not null" json:"shippingAddressId"`
	ShippingAddress   Address     `gorm:"foreignKey:ShippingAddressId;constraint:onUpdate:CASCADE,onDelete:CASCADE"  json:"address"`
	Status            string      `json:"status"`
	PaymentStatus     string      `json:"paymentStatus"`
	PaymentMode       string      `json:"paymentMode"`
	TrackingNumber    int         `json:"trackingNumber"`
	CreatedAt         time.Time   `json:"createdAt"`
	UpdatedAt         time.Time   `json:"updatedAt"`
}

/*
Create(order *Order) (*Order, error)

GetByID(id string) (*Order, error)

GetByUserID(userID string, offset, limit int) ([]*Order, error)

UpdateStatus(orderID string, newStatus string) (*Order, error)

GetAll(offset, limit int, statusFilter string) ([]*Order, error)
*/

func (o *Order) BeforeCreate(t *gorm.DB) error {

	o.Id = uuid.New().String()
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	o.OrderedAt = time.Now()

	if o.Status == "" {
		o.Status = "pending"
	}
	if o.PaymentStatus == "" {
		o.PaymentStatus = "pending"
	}

	return nil
}

func (o *Order) BeforeUpdate(t *gorm.DB) error {
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) Create() (*Order, error) {

	if err := o.ValidateOrder(); err != nil {
		log.Err(err).Msg("Order Validation failed")
		return nil, err
	}

	tx := database.DB.Begin()

	if err := tx.Create(o).Error; err != nil {
		tx.Rollback()
		log.Err(err).Msg("issue exist in create order")
		return nil, err
	}

	for _, item := range o.OrderItems {
		if err := item.UpdateProductStock(tx); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()
	return o, nil
}

func (o *Order) ValidateOrder() error {
	var user User
	if err := database.DB.Where("id = ?", o.UserId).First(&user).Error; err != nil {
		log.Err(err).Msg("Issue in ValidateUser getting user")
		return err
	}

	var address Address
	if err := database.DB.Where(&Address{UserId: o.UserId, Id: o.ShippingAddressId}).First(&address).Error; err != nil {
		log.Err(err).Msg("Issue in ValidateUser Getting Address")
		return err
	}

	if len(o.OrderItems) == 0 {
		return fmt.Errorf("order must have at least one item")
	}

	for _, item := range o.OrderItems {
		if err := item.ValidateOrderItem(); err != nil {
			log.Err(err).Msg("Issue persist")
			return err
		}
	}

	if !o.ValidateTotal() {
		return fmt.Errorf("total amount mismatch")
	}

	validStatuses := []string{"pending", "confirmed", "processing", "shipping"}
	if !contains(validStatuses, o.Status) {
		return fmt.Errorf("invalid order status : %s", o.Status)
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (o *Order) CalculateTotal() int {
	total := 0
	for _, item := range o.OrderItems {
		total += item.PriceAtPurchase * item.Quantity
	}
	return total
}

func (o *Order) UpdateTotalAmount() error {

	if len(o.OrderItems) == 0 {
		if err := database.DB.Preload("OrderItems").Where("id = ?", o.Id).First(o).Error; err != nil {
			return err
		}
	}
	o.TotalAmount = o.CalculateTotal()
	return database.DB.Save(o).Error
}

func (o *Order) ValidateTotal() bool {
	calculateTotal := o.CalculateTotal()
	return o.TotalAmount == calculateTotal
}

func GetOrderById(id string) (*Order, error) {
	var order Order
	if err := database.DB.Where(&Order{Id: id}).First(&order).Error; err != nil {
		log.Err(err).Msg("Issue exist in GetOrderById")
		return nil, err
	}
	return &order, nil

}

func UpdateStatus(orderId string, newStatus string) (*Order, error) {
	var order Order

	if err := database.DB.Model(&Order{}).Where("id = ?", orderId).Update("status", newStatus).Error; err != nil {
		log.Err(err).Msg("Issue exist in UpdateStatus")
		return nil, err
	}

	if err := database.DB.Where("id = ?", orderId).First(&order).Error; err != nil {
		log.Err(err).Msg("issue exist in fetching updated User")
		return nil, err
	}

	return &order, nil
}

func GetAllOrder(limit, offset int, status string) ([]Order, error) {
	var orders []Order
	query := database.DB.Limit(limit).Offset(offset)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&orders).Error; err != nil {
		log.Err(err).Msg("Issue persist in GetAll Order")
		return nil, err
	}
	return orders, nil
}

func GetOrdersByUserId(userId string, page, pageSize int) ([]Order, error) {

	if pageSize <= 0 {
		pageSize = 10
	}

	if page <= 0 {
		page = 1
	}

	offSet := (page - 1) * pageSize

	var order []Order
	if err := database.DB.Where(&Order{UserId: userId}).Order("created_at DESC").Limit(pageSize).Offset(offSet).Find(&order).Error; err != nil {
		log.Err(err).Msg("Issue exist in GetOrderByUserId")
		return nil, err
	}
	return order, nil
}

func GetOrderByUserIdAndOrderId(userId, orderId string) (*Order, error) {
	if userId == "" || orderId == "" {
		return nil, fmt.Errorf("user Id or Order Id is required")
	}

	var order Order
	if err := database.DB.Where("id = ? AND user_id = ?", orderId, userId).Preload("OrderItems").Preload("ShippingAddress").First(&order).Error; err != nil {
		log.Err(err).Msg("Issue exist in GetOrderByUserIdAndOrderId")
		return nil, err
	}
	return &order, nil
}

func DeleteOrderById(orderId string) error {
	return database.DB.Where(&Order{Id: orderId}).Delete(&Order{}).Error
}
