package controller

import (
	"encoding/json"
	"github.com/pratyush934/sibling-bond-server/cjson"
	"github.com/pratyush934/sibling-bond-server/models"
	"net/http"
	"strconv"
)

// GetAllCategories - List all product categories
func GetAllCategories(w http.ResponseWriter, r *http.Request) {
	// Handle pagination parameters
	limit := 10
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	categories, err := models.GetAll(limit, offset)
	if err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusInternalServerError,
			Message:       "Failed to fetch categories",
			InternalError: err,
		})
	}

	_ = cjson.WriteJSON(w, http.StatusOK, categories)
}

// GetCategoryById - Get specific category details
func GetCategoryById(w http.ResponseWriter, r *http.Request) {
	categoryId := r.URL.Query().Get("categoryId")
	if categoryId == "" {
		panic(&cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Category ID is required",
			InternalError: nil,
		})
	}

	category, err := models.GetCategoryById(categoryId)
	if err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Category not found",
			InternalError: err,
		})
	}

	_ = cjson.WriteJSON(w, http.StatusOK, category)
}

// CreateCategory - Add new category (admin only)
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	// Verify admin role
	role, ok := r.Context().Value("role").(int)
	if !ok || role != 2 {
		panic(&cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "Admin privileges required",
			InternalError: nil,
		})
	}

	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Invalid category data",
			InternalError: err,
		})
	}

	// Validate category name
	if category.Name == "" {
		panic(&cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Category name is required",
			InternalError: nil,
		})
	}

	// Check if category with same name already exists
	existing, err := models.GetCategoryByName(category.Name)
	if err == nil && existing != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusConflict,
			Message:       "Category with this name already exists",
			InternalError: nil,
		})
	}

	createdCategory, err := category.CreateCategory()
	if err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusInternalServerError,
			Message:       "Failed to create category",
			InternalError: err,
		})
	}

	_ = cjson.WriteJSON(w, http.StatusCreated, createdCategory)
}

// UpdateCategory - Update category details (admin only)
func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	// Verify admin role
	role, ok := r.Context().Value("role").(int)
	if !ok || role != 2 {
		panic(&cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "Admin privileges required",
			InternalError: nil,
		})
	}

	categoryId := r.URL.Query().Get("categoryId")
	if categoryId == "" {
		panic(&cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Category ID is required",
			InternalError: nil,
		})
	}

	// Get existing category
	existingCategory, err := models.GetCategoryById(categoryId)
	if err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Category not found",
			InternalError: err,
		})
	}

	// Parse update data
	var updatedData models.Category
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Invalid update data",
			InternalError: err,
		})
	}

	// Update fields if provided
	if updatedData.Name != "" {
		// Check if the new name would conflict with an existing category
		if existingCategory.Name != updatedData.Name {
			existing, err := models.GetCategoryByName(updatedData.Name)
			if err == nil && existing != nil && existing.Id != categoryId {
				panic(&cjson.HTTPError{
					Status:        http.StatusConflict,
					Message:       "Another category with this name already exists",
					InternalError: nil,
				})
			}
		}
		existingCategory.Name = updatedData.Name
	}

	if updatedData.Description != "" {
		existingCategory.Description = updatedData.Description
	}

	// Save updates
	updatedCategory, err := models.UpdateCategory(existingCategory)
	if err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusInternalServerError,
			Message:       "Failed to update category",
			InternalError: err,
		})
	}

	_ = cjson.WriteJSON(w, http.StatusOK, updatedCategory)
}

// DeleteCategory - Remove a category (admin only)
func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	// Verify admin role
	role, ok := r.Context().Value("role").(int)
	if !ok || role != 2 {
		panic(&cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "Admin privileges required",
			InternalError: nil,
		})
	}

	categoryId := r.URL.Query().Get("categoryId")
	if categoryId == "" {
		panic(&cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Category ID is required",
			InternalError: nil,
		})
	}

	// Check if category exists
	_, err := models.GetCategoryById(categoryId)
	if err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Category not found",
			InternalError: err,
		})
	}

	// Check if category has associated products
	hasProducts, err := models.CategoryHasProducts(categoryId)
	if err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusInternalServerError,
			Message:       "Failed to check category usage",
			InternalError: err,
		})
	}

	if hasProducts {
		panic(&cjson.HTTPError{
			Status:        http.StatusConflict,
			Message:       "Cannot delete category with associated products",
			InternalError: nil,
		})
	}

	// Delete the category
	if err := models.DeleteCategory(categoryId); err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusInternalServerError,
			Message:       "Failed to delete category",
			InternalError: err,
		})
	}

	_ = cjson.WriteJSON(w, http.StatusOK, "Category deleted successfully")
}
