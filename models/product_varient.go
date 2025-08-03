package models

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pratyush934/sibling-bond-server/database"
	"gorm.io/gorm"
	"strings"
	"time"
)

type ProductVariant struct {
	Id           string    `gorm:"primaryKey;type:varchar(191)" json:"id"`
	ProductId    string    `gorm:"not null" json:"productId"`
	Product      Product   `gorm:"foreignKey:ProductId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"product"`
	VariantName  string    `gorm:"not null" json:"variantName"`
	VariantValue string    `gorm:"not null" json:"variantValue"`
	Price        int       `json:"price"`
	Stock        int       `gorm:"not null;default:0" json:"stock"`
	SKU          string    `gorm:"unique" json:"sku"`
	IsActive     bool      `gorm:"default:true" json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (pv *ProductVariant) BeforeCreate(t *gorm.DB) error {
	pv.Id = uuid.New().String()
	pv.CreatedAt = time.Now()
	pv.UpdatedAt = time.Now()

	if pv.SKU == "" {
		pv.SKU = fmt.Sprintf("%s-%s-%s", pv.ProductId[:8], strings.ToUpper(pv.VariantName[:3]), uuid.New().String()[:6])
	}
	return nil
}

func (pv *ProductVariant) BeforeUpdate(t *gorm.DB) error {
	pv.UpdatedAt = time.Now()
	return nil
}

func GetProductVariants(productId string) ([]ProductVariant, error) {
	var variants []ProductVariant
	if err := database.DB.Where("product_id = ? AND is_active = ?", productId, true).Find(&variants).Error; err != nil {
		return nil, err
	}
	return variants, nil
}
