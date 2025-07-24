package models

import (
	"github.com/google/uuid"
	"github.com/pratyush934/sibling-bond-server/database"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

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

/*
Create(order *Order) (*Order, error)

GetByID(id string) (*Order, error)

GetByUserID(userID string, offset, limit int) ([]*Order, error)

UpdateStatus(orderID string, newStatus string) (*Order, error)

GetAll(offset, limit int, statusFilter string) ([]*Order, error)
*/

func (o *Order) BeforeCreate(t *gorm.DB) error {
	o.Id = uuid.New().String()
	return nil
}

func (o *Order) Create() (*Order, error) {
	if err := database.DB.Create(o).Error; err != nil {
		log.Err(err).Msg("Issue exist in Create")
		return nil, err
	}
	return o, nil
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
