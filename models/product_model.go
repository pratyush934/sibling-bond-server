package models

import (
	"github.com/google/uuid"
	"github.com/pratyush934/sibling-bond-server/database"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type Product struct {
	Id            string    `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"not null" json:"name"`
	Description   string    `gorm:"not null" json:"description"`
	Price         int       `gorm:"not null" json:"price"`
	StockQuantity int       `gorm:"not null" json:"stockQuantity"`
	CategoryId    string    `json:"categoryId" gorm:"not null"`
	Category      Category  `json:"category" gorm:"foreignKey:CategoryId"`
	ImageUrl      string    `json:"imageUrl"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func (p *Product) BeforeCreate(t *gorm.DB) error {
	p.Id = uuid.New().String()
	return nil
}

func (p *Product) CreateProduct() (*Product, error) {
	if err := database.DB.Create(p).Error; err != nil {
		log.Err(err).Msg("Issue persist in the CreateProduct")
		return &Product{}, err
	}
	return p, nil
}

func GetProductById(id string) (*Product, error) {
	var product Product
	if err := database.DB.Where(&Product{Id: id}).First(&product).Error; err != nil {
		log.Err(err).Msg("Issue exist in GetProductById")
		return &product, err
	}
	return &product, nil
}

func UpdateProduct(p *Product) (*Product, error) {
	if err := database.DB.Updates(p).Error; err != nil {
		log.Err(err).Msg("Issue exist in UpdateProduct")
		return &Product{}, err
	}
	return p, nil
}

func DeleteProduct(id string) error {
	return database.DB.Where(&Product{Id: id}).Delete(&Product{}).Error
}

func GetAllProducts(limit, offSet int, categoryId, searchQuery string) ([]Product, error) {
	var products []Product
	query := database.DB.Limit(limit).Offset(offSet)

	if categoryId != "" {
		query = query.Where(&Product{CategoryId: categoryId})
	}

	if searchQuery != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	if err := query.Find(&products).Error; err != nil {
		log.Err(err).Msg("Issue exist in GetAllProducts")
		return nil, err
	}

	return products, nil

}

func UpdateStock(productId string, quantityChange int) error {
	if err := database.DB.Model(&Product{}).Where("id = ?", productId).Update("stock_quantity", gorm.Expr("stock_quantity + ?", quantityChange)).Error; err != nil {
		log.Err(err).Msg("Issue exist in UpdateStock")
		return err
	}
	return nil
}
