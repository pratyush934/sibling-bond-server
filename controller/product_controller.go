package controller

import (
	"github.com/pratyush934/sibling-bond-server/cjson"
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
GetFeaturedProducts - Get highlighted products
*/

func GetAllProducts(w http.ResponseWriter, r *http.Request) {

	limit := 10
	offSet := 5

	limitStr := r.URL.Query().Get("limit")
	offSetStr := r.URL.Query().Get("offset")

	if limitStr == "" || offSetStr == "" {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Please provide limitStr or OffSetStr",
			InternalError: nil,
		})
	}

	limit, _ = strconv.Atoi(limitStr)
	offSet, _ = strconv.Atoi(offSetStr)

	products, err := models.GetAllProducts(limit, offSet)
	if err != nil {
		panic(cjson.HTTPError{
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
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Please provide productId",
			InternalError: nil,
		})
	}
	productById, err := models.GetProductById(productId)
	if err != nil {
		panic(cjson.HTTPError{
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
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Please provide limitStr or OffSetStr",
			InternalError: nil,
		})
	}

	limit, _ = strconv.Atoi(limitStr)
	offSet, _ = strconv.Atoi(offSetStr)

	allProducts, err := models.GetAllProductsWithQueries(limit, offSet, categoryId, search)
	if err != nil {
		panic(cjson.HTTPError{
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
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Please provide limitStr or OffSetStr",
			InternalError: nil,
		})
	}

	limit, _ = strconv.Atoi(limitStr)
	offSet, _ = strconv.Atoi(offSetStr)
	categoryId := r.URL.Query().Get("categoryId")

	if categoryId == "" {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Please provide limitStr or OffSetStr",
			InternalError: nil,
		})
	}

	productById, err := models.GetProductsByCategoryId(categoryId, limit, offSet)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Not able to get the products",
			InternalError: err,
		})
	}

	_ = cjson.WriteJSON(w, http.StatusOK, productById)
}
