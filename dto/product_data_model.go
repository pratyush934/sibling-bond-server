package dto

type ProductModel struct {
	Name          string              `json:"name"`
	Description   string              `json:"description"`
	Price         int                 `json:"price"`
	Stock         int                 `json:"stock"`
	CategoryId    string              `json:"categoryId"`
	Images        []ImageMode         `json:"images"`
	IsActive      bool                `json:"isActive"`
	MinStockLevel int                 `json:"minStockLevel"`
	MaxStockLevel int                 `json:"maxStockLevel"`
	ReorderPoint  int                 `json:"reorderPoint"`
	SKU           string              `json:"sku"` // Optional, will be auto-generated if empty
	Barcode       string              `json:"barcode"`
	Weight        float64             `json:"weight"`
	Dimensions    string              `json:"dimensions"`
	Variants      []ProductVariantDTO `json:"variants"`
}

type ProductVariantDTO struct {
	Name            string `json:"name"`
	Value           string `json:"value"`
	PriceAdjustment int    `json:"priceAdjustment"`
}
