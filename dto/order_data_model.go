package dto

/*

	Id                string      `gorm:"primaryKey" json:"id"`
	UserId            string      `gorm:"not null" json:"userId"`
	User              User        `gorm:"foreignKey:UserId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"user"`
	OrderItems        []OrderItem `gorm:"foreignKey:OrderId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"orderItems"`
	OrderedAt         time.Time   `json:"orderedAt"`
	TotalAmount       int         `json:"totalAmount"`
	ShippingAddressId string      `gorm:"not null" json:"shippingAddressId"`
	ShippingAddress   Address     `gorm:"foreignKey:ShippingAddressId;constraint:onUpdate:CASCADE,onDelete:CASCADE"  json:"address"`
	Status            string      `json:"status"`
	PaymentStatus     string      `json:"paymentStatus"`
	PaymentMode       string      `json:"paymentMode"`
	TrackingNumber    int         `json:"trackingNumber"`
	CreatedAt         time.Time   `json:"createdAt"`
	UpdatedAt         time.Time   `json:"updatedAt"`
*/

/*
Id              string    `gorm:"primaryKey" json:"id"`
	OrderId         string    `gorm:"not null" json:"orderId"`
	ProductId       string    `gorm:"not null" json:"productId"`
	VariantId       *string   `json:"variantId"`
	Order           Order     `gorm:"foreignKey:OrderId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"order"`
	Product         Product   `gorm:"foreignKey:ProductId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"product"`
	Quantity        int       `gorm:"not null" json:"quantity"`
	PriceAtPurchase int       `gorm:"not null" json:"priceAtPurchase"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
*/

type OrderModel struct {
	UserId            string           `json:"userId"`
	TotalAmount       int              `json:"totalAmount"`
	ShippingAddressId string           `json:"shippingAddressId"`
	OrderItems        []OrderItemModel `json:"orderItems"`
	Status            string           `json:"status"`
	PaymentStatus     string           `json:"paymentStatus"`
	PaymentMode       string           `json:"paymentMode"`
	TrackingNumber    int              `json:"trackingNumber"`
}

type OrderItemModel struct {
	ProductId       string `json:"productId"`
	Quantity        int    `json:"quantity"`
	PriceAtPurchase int    `json:"priceAtPurchase"`
}
