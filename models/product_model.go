package models

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pratyush934/sibling-bond-server/database"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Product struct {
	Id          string `gorm:"primaryKey;type:varchar(191)" json:"id"`
	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`
	Price       int    `gorm:"not null" json:"price"`
	Stock       int    `gorm:"not null;default:0" json:"stock"`

	CategoryId string    `gorm:"not null;type:varchar(150);column:category_id" json:"categoryId"`
	Category   *Category `gorm:"foreignKey:category_id;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"category"`

	Images []Image `gorm:"foreignKey:ProductId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"images"`

	IsActive  bool           `gorm:"default:true" json:"isActive"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`

	Variants []ProductVariant `gorm:"foreignKey:ProductId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"variants"`

	MinStockLevel int     `gorm:"default:5" json:"minStockLevel"`
	MaxStockLevel int     `gorm:"default:100" json:"maxStockLevel"`
	ReorderPoint  int     `gorm:"default:10" json:"reorderPoint"`
	SKU           string  `gorm:"unique" json:"sku"`
	Barcode       string  `json:"barcode"`
	Weight        float64 `json:"weight"`
	Dimensions    string  `json:"dimensions"`
}

func (p *Product) BeforeCreate(t *gorm.DB) error {
	p.Id = uuid.New().String()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	if p.SKU == "" {
		p.SKU = generateSKU(p.Name, p.CategoryId)
	}

	if p.MinStockLevel == 0 {
		p.MinStockLevel = 5
	}
	if p.MaxStockLevel == 0 {
		p.MaxStockLevel = 100
	}

	if p.ReorderPoint == 0 {
		p.ReorderPoint = 10
	}

	return nil
}

func (p *Product) BeforeUpdate(t *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}

func generateSKU(productName, categoryId string) string {
	prefixLen := min(3, len(productName))
	categoryLen := min(8, len(categoryId))
	prefix := strings.ToUpper(productName[:prefixLen])
	suffix := uuid.New().String()[:8]
	return fmt.Sprintf("%s-%s-%s", prefix, categoryId[:categoryLen], suffix)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (p *Product) CreateProduct() (*Product, error) {
	if err := database.DB.Create(p).Error; err != nil {
		log.Err(err).Msg("Issue persist in the CreateProduct")
		return &Product{}, err
	}
	return p, nil
}

func (p *Product) IsInStock() bool {
	return p.Stock > 0 && p.IsActive
}

func (p *Product) IsLowStock() bool {
	return p.Stock <= p.ReorderPoint
}

func (p *Product) CanFulFillOrder(quantity int) bool {
	return p.Stock >= quantity && p.IsActive
}

func (p *Product) UpdateStock(quantity int, operation string) (int, error) {
	switch operation {
	case "add":
		p.Stock += quantity
	case "subtract":
		if p.Stock < quantity {
			return p.Stock, fmt.Errorf("stock is 0 can't add the stuff")
		}
		p.Stock -= quantity
	case "set":
		p.Stock = quantity
	default:
		return p.Stock, fmt.Errorf("please add valid operation")
	}
	// Only update the stock field, not the whole struct with associations
	return p.Stock, database.DB.Model(&Product{}).Where("id = ?", p.Id).Update("stock", p.Stock).Error
}

func (p *Product) GetStockStatus() string {
	if p.Stock == 0 {
		return "out_of_stocks"
	} else if p.IsLowStock() {
		return "low_stocks"
	} else {
		return "in_stock"
	}
}

//func (p *Product) ReserveStock(quantity int) error {
//	if !p.CanFulFillOrder(quantity) {
//		return fmt.Errorf("can not full fill the order")
//	}
//}

func (p *Product) RestoreStock(quantity int) (int, error) {
	return p.UpdateStock(quantity, "add")
}

func (p *Product) SoftDelete() error {
	return database.DB.Delete(p).Error
}

func (p *Product) Restore() error {
	return database.DB.Unscoped().Model(p).Update("deleted_at", nil).Error
}

func (p *Product) ToggleActive() error {
	p.IsActive = !p.IsActive
	return database.DB.Save(p).Error
}

func GetDeleteProducts(limit, offset int) ([]Product, error) {
	var products []Product
	if err := database.DB.Unscoped().Where("deleted_at IS NOT NULL").Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		log.Err(err).Msg("Issue getting deleted products")
		return nil, err
	}
	return products, nil
}

func GetProductById(id string) (*Product, error) {
	var product Product
	if err := database.DB.Preload("Category").Preload("Variants").Preload("Images").Where(&Product{Id: id}).First(&product).Error; err != nil {
		log.Err(err).Msg("Issue exist in GetProductById")
		return nil, err
	}
	return &product, nil
}

func GetProductsByCategoryId(categoryId string, limit, offset int) ([]Product, error) {
	var products []Product

	if err := database.DB.Where(&Product{CategoryId: categoryId}).
		Preload("Category").
		Preload("Variants").
		Preload("Images").
		Limit(limit).
		Offset(offset).
		Find(&products).Error; err != nil {
		log.Err(err).Msg("Issue getting products by category id")
		return nil, err
	}

	return products, nil
}

func GetLowStockProducts() ([]Product, error) {
	var products []Product
	if err := database.DB.Where("stock <= reorder_point AND is_active = ?", true).Find(&products).Error; err != nil {
		log.Err(err).Msg("Issue exist in getting LowStockProduct")
		return nil, err
	}
	return products, nil
}

func GetOutOfStockProducts() ([]Product, error) {
	var products []Product
	if err := database.DB.Where("stock = 0 AND is_active = ?", true).Find(&products).Error; err != nil {
		log.Err(err).Msg("Issue exist in getting OutOfStock")
		return nil, err
	}
	return products, nil
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

func GetAllProducts(limit, offSet int) ([]Product, error) {
	var products []Product
	query := database.DB.Limit(limit).Offset(offSet)

	if err := query.Preload("Images").Find(&products).Error; err != nil {
		log.Err(err).Msg("Issue exist in GetAllProducts")
		return nil, err
	}

	return products, nil
}

func GetAllProductsWithQueries(limit, offSet int, categoryId, searchQuery string) ([]Product, error) {
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
	if err := database.DB.Model(&Product{}).Where("id = ?", productId).Update("stock", gorm.Expr("stock + ?", quantityChange)).Error; err != nil {
		log.Err(err).Msg("Issue exist in UpdateStock")
		return err
	}
	return nil
}

func (p *Product) GetPrimaryImage() *Image {
	for _, img := range p.Images {
		if img.IsPrimary {
			return &img
		}
	}
	if len(p.Images) > 0 {
		return &p.Images[0] // fallback to first image
	}
	return nil
}

func (p *Product) SetPrimaryImage(imageId string) error {
	// Reset all images to non-primary
	err := database.DB.Model(&Image{}).Where("product_id = ?", p.Id).Update("is_primary", false).Error
	if err != nil {
		return err
	}

	// Set specified image as primary
	return database.DB.Model(&Image{}).Where("id = ? AND product_id = ?", imageId, p.Id).Update("is_primary", true).Error
}

func (p *Product) AddImage(url, fileName, altText string) error {
	image := Image{
		URL:       url,
		FileName:  fileName,
		AltText:   altText,
		ProductId: p.Id,
		SortOrder: len(p.Images),
	}
	return database.DB.Create(&image).Error
}
