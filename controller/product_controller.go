package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pratyush934/sibling-bond-server/cjson"
	"github.com/pratyush934/sibling-bond-server/database"
	"github.com/pratyush934/sibling-bond-server/dto"
	"github.com/pratyush934/sibling-bond-server/models"
	"net/http"
	"strconv"
)

/*
GetAllProducts - List all products with pagination
GetProductById - Get details for specific product
SearchProducts - Search products by keywords/filters
GetProductsByCategory - List products in a category
CreateProduct - Add new product (admin only)
UpdateProduct - Update product details (admin only)
DeleteProduct - Remove a product (admin only)
*/

func CheckAdmin(w http.ResponseWriter, r *http.Request) error {
	role, ok := r.Context().Value("role").(float64)

	if !ok {
		return fmt.Errorf("not able to get the roleId from context")
	}

	if role != 2 {
		return fmt.Errorf("this guy is not admin ")
	}
	return nil
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {

	limit := 10
	offSet := 5

	limitStr := r.URL.Query().Get("limit")
	offSetStr := r.URL.Query().Get("offset")

	if limitStr == "" || offSetStr == "" {
		panic(&cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Please provide limitStr or OffSetStr",
			InternalError: nil,
		})
	}

	limit, _ = strconv.Atoi(limitStr)
	offSet, _ = strconv.Atoi(offSetStr)

	products, err := models.GetAllProducts(limit, offSet)
	if err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Not able to get the Products",
			InternalError: err,
		})
	}
	_ = cjson.WriteJSON(w, http.StatusOK, products)

}

func GetProductById(w http.ResponseWriter, r *http.Request) {
	productId := r.URL.Query().Get("id")

	if productId == "" {
		panic(&cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Please provide productId",
			InternalError: nil,
		})
	}
	productById, err := models.GetProductById(productId)
	if err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Not found the product, the Id may be wrong",
			InternalError: err,
		})
	}
	_ = cjson.WriteJSON(w, http.StatusOK, productById)
}

func SearchProduct(w http.ResponseWriter, r *http.Request) {
	limit, offSet := 10, 5

	limitStr := r.URL.Query().Get("limit")
	offSetStr := r.URL.Query().Get("offSet")
	categoryId := r.URL.Query().Get("categoryId")
	search := r.URL.Query().Get("search")

	if limitStr == "" || offSetStr == "" {
		panic(&cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Please provide limitStr or OffSetStr",
			InternalError: nil,
		})
	}

	limit, _ = strconv.Atoi(limitStr)
	offSet, _ = strconv.Atoi(offSetStr)

	allProducts, err := models.GetAllProductsWithQueries(limit, offSet, categoryId, search)
	if err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Not able to get the products",
			InternalError: err,
		})
	}
	_ = cjson.WriteJSON(w, http.StatusOK, allProducts)
}

func GetProductsByCategory(w http.ResponseWriter, r *http.Request) {

	limit := 10
	offSet := 5

	limitStr := r.URL.Query().Get("limit")
	offSetStr := r.URL.Query().Get("offset")

	if limitStr == "" || offSetStr == "" {
		panic(&cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Please provide limitStr or OffSetStr",
			InternalError: nil,
		})
	}

	limit, _ = strconv.Atoi(limitStr)
	offSet, _ = strconv.Atoi(offSetStr)
	categoryId := r.URL.Query().Get("categoryId")

	if categoryId == "" {
		panic(&cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Please provide limitStr or OffSetStr",
			InternalError: nil,
		})
	}

	productById, err := models.GetProductsByCategoryId(categoryId, limit, offSet)
	if err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Not able to get the products",
			InternalError: err,
		})
	}

	_ = cjson.WriteJSON(w, http.StatusOK, productById)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {

	err := CheckAdmin(w, r)
	if err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "Not able to verify the user as admin",
			InternalError: err,
		})
	}

	var productModel dto.ProductModel

	if err := json.NewDecoder(r.Body).Decode(&productModel); err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Not able to decode the Product",
			InternalError: err,
		})
	}

	if productModel.Name == "" || productModel.CategoryId == "" || productModel.Price <= 0 {
		panic(&cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Need to feed Name, CategoryId, Price",
			InternalError: nil,
		})
	}

	newProduct := models.Product{
		Name:          productModel.Name,
		Description:   productModel.Description,
		Price:         productModel.Price,
		Stock:         productModel.Stock,
		CategoryId:    productModel.CategoryId,
		IsActive:      productModel.IsActive,
		MinStockLevel: productModel.MinStockLevel,
		MaxStockLevel: productModel.MaxStockLevel,
		ReorderPoint:  productModel.ReorderPoint,
		SKU:           productModel.SKU,
		Barcode:       productModel.Barcode,
		Weight:        productModel.Weight,
		Dimensions:    productModel.Dimensions,
	}

	if len(productModel.Variants) > 0 {
		variants := make([]models.ProductVariant, 0, len(productModel.Variants))
		for _, v := range productModel.Variants {
			variants = append(variants, models.ProductVariant{
				VariantName:  v.Name,
				VariantValue: v.Value,
				Price:        v.PriceAdjustment,
			})
		}
		newProduct.Variants = variants
	}

	product, err := newProduct.CreateProduct()
	if err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Not able to Create Product",
			InternalError: err,
		})
	}

	if len(productModel.Images) > 0 {
		for i, imageData := range productModel.Images {
			image := models.Image{
				URL:       imageData.URL,
				FileName:  imageData.Name,
				FieldId:   imageData.FieldId,
				ProductId: product.Id,
				SortOrder: i,
				IsPrimary: i == 0,
			}

			_, err := image.CreateImage()
			if err != nil {
				panic(&cjson.HTTPError{
					Status:        http.StatusBadRequest,
					Message:       fmt.Sprintf("Not able to store the %vth image ", i),
					InternalError: err,
				})
			}
		}

	}
	productById, err := models.GetProductById(product.Id)

	if err != nil {
		_ = cjson.WriteJSON(w, http.StatusCreated, product)
		return
	}

	_ = cjson.WriteJSON(w, http.StatusOK, productById)

}
func UpdateProductDetails(w http.ResponseWriter, r *http.Request) {
	err := CheckAdmin(w, r)
	if err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "Not an admin",
			InternalError: err,
		})
	}

	vars := mux.Vars(r)
	productId := vars["id"]

	if productId == "" {
		panic(&cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Product ID is required",
			InternalError: nil,
		})
	}

	_, err = models.GetProductById(productId)
	if err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Product not found",
			InternalError: err,
		})
	}

	var productModel dto.ProductModel
	if err := json.NewDecoder(r.Body).Decode(&productModel); err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Not able to Decode the Product",
			InternalError: err,
		})
	}

	// Update product WITHOUT images field
	updateProduct := models.Product{
		Id:            productId,
		Name:          productModel.Name,
		Description:   productModel.Description,
		Price:         productModel.Price,
		Stock:         productModel.Stock,
		CategoryId:    productModel.CategoryId,
		IsActive:      productModel.IsActive,
		MinStockLevel: productModel.MinStockLevel,
		MaxStockLevel: productModel.MaxStockLevel,
		ReorderPoint:  productModel.ReorderPoint,
		SKU:           productModel.SKU,
		Barcode:       productModel.Barcode,
		Weight:        productModel.Weight,
		Dimensions:    productModel.Dimensions,
	}

	// Handle variants
	if len(productModel.Variants) > 0 {
		if err := database.DB.Where("product_id = ?", productId).Delete(&models.ProductVariant{}).Error; err != nil {
			panic(&cjson.HTTPError{
				Status:        http.StatusInternalServerError,
				Message:       "Failed to delete existing product variants",
				InternalError: err,
			})
		}

		variants := make([]models.ProductVariant, 0, len(productModel.Variants))
		for _, v := range productModel.Variants {
			variants = append(variants, models.ProductVariant{
				ProductId:    productId,
				VariantName:  v.Name,
				VariantValue: v.Value,
				Price:        v.PriceAdjustment,
			})
		}
		updateProduct.Variants = variants
	}

	// Handle images separately if provided
	if len(productModel.Images) > 0 {
		// Delete existing images
		if err := database.DB.Where("product_id = ?", productId).Delete(&models.Image{}).Error; err != nil {
			panic(&cjson.HTTPError{
				Status:        http.StatusInternalServerError,
				Message:       "Failed to delete existing images",
				InternalError: err,
			})
		}

		// Create new images
		for i, imageData := range productModel.Images {
			image := models.Image{
				URL:       imageData.URL,
				FileName:  imageData.Name,
				FieldId:   imageData.FieldId,
				ProductId: productId,
				SortOrder: i,
				IsPrimary: i == 0,
			}

			_, err := image.CreateImage()
			if err != nil {
				panic(&cjson.HTTPError{
					Status:        http.StatusInternalServerError,
					Message:       fmt.Sprintf("Failed to create image %d", i),
					InternalError: err,
				})
			}
		}
	}

	product, err := models.UpdateProduct(&updateProduct)
	if err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusInternalServerError,
			Message:       "Not able to Update the Product",
			InternalError: err,
		})
	}

	// Fetch complete product with images
	completeProduct, err := models.GetProductById(productId)
	if err != nil {
		_ = cjson.WriteJSON(w, http.StatusOK, product)
		return
	}

	_ = cjson.WriteJSON(w, http.StatusOK, completeProduct)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {

	err2 := CheckAdmin(w, r)

	if err2 != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "Not an admin",
			InternalError: err2,
		})
	}

	vars := mux.Vars(r)
	productId := vars["id"]

	if productId == "" {
		panic(&cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "ProductId is empty",
			InternalError: nil,
		})
	}

	err := models.DeleteProduct(productId)

	if err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Not able to delete this product",
			InternalError: err,
		})
	}
	_ = cjson.WriteJSON(w, http.StatusOK, "Product deleted successfully")
}
