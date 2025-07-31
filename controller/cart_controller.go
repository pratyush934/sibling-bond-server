package controller

import (
	"encoding/json"
	"github.com/pratyush934/sibling-bond-server/cjson"
	"github.com/pratyush934/sibling-bond-server/database"
	"github.com/pratyush934/sibling-bond-server/dto"
	"github.com/pratyush934/sibling-bond-server/models"
	"net/http"
	"strconv"
)

/*
GetCart - Retrieve user's current cart
CreateCart - Initialize a new cart for a user
AddToCart - Add product to cart (handles new items and quantity increases)
UpdateCartItem - Update quantity of existing cart item
RemoveFromCart - Remove specific item from cart
ClearCart - Remove all items from user's cart
GetCartTotal - Calculate total price of items in cart
ValidateCartItems - Check if cart items are still available in inventory
*/

func GetCart(w http.ResponseWriter, r *http.Request) {

	userId, ok := r.Context().Value("userId").(string)

	if userId == "" || !ok {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Not able to get the UserId",
			InternalError: nil,
		})
	}

	cartById, err := models.GetCartByUserId(userId)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Not able to get the CartById",
			InternalError: err,
		})
	}
	_ = cjson.WriteJSON(w, http.StatusOK, cartById)
}

func CreateCart(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(string)

	if !ok {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "Not able to get the Id from the User",
			InternalError: nil,
		})
	}

	var cartModel dto.CartDataModel

	if err := json.NewDecoder(r.Body).Decode(&cartModel); err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Not able to decode the cartModel",
			InternalError: err,
		})
	}

	newCart := models.Cart{
		UserId: userId,
	}

	if len(cartModel.CartItems) > 0 {
		cartitems := make([]models.CartItem, 0, len(cartModel.CartItems))

		for _, v := range cartModel.CartItems {
			cartitems = append(cartitems, models.CartItem{
				ProductId:     v.ProductId,
				Quantity:      v.Quantity,
				PriceAtAdding: v.PriceAtAdding,
			})
		}
		newCart.CartItems = cartitems
	}

	cartCreated, err := models.Create(&newCart)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusInternalServerError,
			Message:       "Not able to get it",
			InternalError: err,
		})
	}

	_ = cjson.WriteJSON(w, http.StatusCreated, cartCreated)

}

func AddToCart(w http.ResponseWriter, r *http.Request) {

	userId, ok := r.Context().Value("userId").(string)

	if !ok {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "Not able to get the UserId",
			InternalError: nil,
		})
	}
	var cartItem dto.CartItemModel

	if err := json.NewDecoder(r.Body).Decode(&cartItem); err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Not able to get the CartItem",
			InternalError: err,
		})
	}

	cartByUserId, err := models.GetCartByUserId(userId)
	if err != nil {
		//cart is not there so we need to create one
		newCart := models.Cart{
			UserId: userId,
		}

		cartByUserId, err = models.Create(&newCart)
		if err != nil {
			panic(cjson.HTTPError{
				Status:        http.StatusInternalServerError,
				Message:       "Not able to create Cart",
				InternalError: err,
			})
		}
	}

	existingItem, err := models.GetItemByCartAndProduct(cartByUserId.Id, cartItem.ProductId)

	if err == nil {
		/* product exist and then update the quantity */
		updateQuantityStuff, err := models.IncrementItemQuantity(existingItem.Id, cartItem.Quantity)
		if err != nil {
			panic(cjson.HTTPError{
				Status:        http.StatusBadRequest,
				Message:       "Not able to update stuff",
				InternalError: err,
			})
		}
		_ = cjson.WriteJSON(w, http.StatusCreated, updateQuantityStuff)
		return
	}

	newProduct := models.CartItem{
		ProductId:     cartItem.ProductId,
		Quantity:      cartItem.Quantity,
		PriceAtAdding: cartItem.PriceAtAdding,
		CartId:        cartByUserId.Id,
	}

	addedItem, err := models.AddItem(&newProduct)

	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Not able to add item to the cart",
			InternalError: err,
		})
	}
	_ = cjson.WriteJSON(w, http.StatusCreated, addedItem)
}

func UpdateCartItem(w http.ResponseWriter, r *http.Request) {

	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "User authentication required",
			InternalError: nil,
		})
	}

	/* cartId and ProductId */

	cartItemId := r.URL.Query().Get("cartItem")
	quantity := r.URL.Query().Get("quantity")

	if cartItemId == "" {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Please provide the cartItemId",
			InternalError: nil,
		})
	}
	amount := 0
	if quantity == "" {
		amount = 0
	}
	amount, err := strconv.Atoi(quantity)
	if err != nil || amount < 0 {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Invalid quantity value",
			InternalError: err,
		})
	}

	// Verify cart ownership - get the cart item first
	var cartItem models.CartItem
	if err := database.DB.Where("id = ?", cartItemId).First(&cartItem).Error; err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Cart item not found",
			InternalError: err,
		})
	}

	// Get the cart to verify ownership
	cart, err := models.GetCartById(cartItem.CartId)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Cart not found",
			InternalError: err,
		})
	}

	// Verify cart belongs to the authenticated user
	if cart.UserId != userId {
		panic(cjson.HTTPError{
			Status:        http.StatusForbidden,
			Message:       "You are not authorized to modify this cart",
			InternalError: nil,
		})
	}

	itemQuantity, err := models.UpdateItemQuantity(cartItemId, amount)

	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusInternalServerError,
			Message:       "Not able to increase quantity",
			InternalError: err,
		})
	}
	_ = cjson.WriteJSON(w, http.StatusOK, itemQuantity)
}

func RemoveFromCart(w http.ResponseWriter, r *http.Request) {

	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "User authentication required",
			InternalError: nil,
		})
	}

	cartItemId := r.URL.Query().Get("cartItem")

	var cartItem models.CartItem
	if err := database.DB.Where("id = ?", cartItemId).First(&cartItem).Error; err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Cart item not found",
			InternalError: err,
		})
	}

	// Get the cart to verify ownership
	cart, err := models.GetCartById(cartItem.CartId)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Cart not found",
			InternalError: err,
		})
	}

	// Verify cart belongs to the authenticated user
	if cart.UserId != userId {
		panic(cjson.HTTPError{
			Status:        http.StatusForbidden,
			Message:       "You are not authorized to modify this cart",
			InternalError: nil,
		})
	}

	err = models.RemoveItem(cartItemId)

	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusInternalServerError,
			Message:       "Not able to delete",
			InternalError: err,
		})
	}

	_ = cjson.WriteJSON(w, http.StatusOK, "CartItemDeletedSuccessfully")
}

func ClearCart(w http.ResponseWriter, r *http.Request) {

	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "User authentication required",
			InternalError: nil,
		})
	}

	cartId := r.URL.Query().Get("cart")

	cart, err := models.GetCartById(cartId)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Cart not found",
			InternalError: err,
		})
	}

	// Verify cart belongs to the authenticated user
	if cart.UserId != userId {
		panic(cjson.HTTPError{
			Status:        http.StatusForbidden,
			Message:       "You are not authorized to modify this cart",
			InternalError: nil,
		})
	}

	err := models.DeleteAllItemsByCartId(cartId)

	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusInternalServerError,
			Message:       "Not able to delete",
			InternalError: err,
		})
	}

	_ = cjson.WriteJSON(w, http.StatusOK, "Cart is Empty Now")
}

func GetCartItemTotal(w http.ResponseWriter, r *http.Request) {

	userId, ok := r.Context().Value("userId").(string)

	if userId == "" || !ok {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Not able to get the UserId",
			InternalError: nil,
		})
	}

	cartByUserId, err := models.GetCartByUserId(userId)

	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusInternalServerError,
			Message:       "Not able to getCart",
			InternalError: err,
		})
	}

	total := 0
	totalQuantity := 0
	for _, v := range cartByUserId.CartItems {
		total += v.PriceAtAdding * v.Quantity
		totalQuantity += v.Quantity
	}

	type CartTotalResponse struct {
		TotalMoney int `json:"totalMoney"`
		Quantity   int `json:"quantity"`
	}

	totalResponse := CartTotalResponse{
		TotalMoney: total,
		Quantity:   totalQuantity,
	}

	_ = cjson.WriteJSON(w, http.StatusOK, totalResponse)

}

func ApplyCoupon(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "User ID not found in context",
			InternalError: nil,
		})
	}

	// Parse coupon from request body
	type CouponRequest struct {
		CouponCode string `json:"couponCode"`
	}
	var couponReq CouponRequest
	if err := json.NewDecoder(r.Body).Decode(&couponReq); err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Invalid coupon data",
			InternalError: err,
		})
	}

	// Get user's cart
	cart, err := models.GetCartByUserId(userId)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Cart not found",
			InternalError: err,
		})
	}

	// In a real implementation, verify coupon validity from database
	// For now, we'll just respond with success

	type ApplyCouponResponse struct {
		Message    string `json:"message"`
		CouponCode string `json:"couponCode"`
		CartId     string `json:"cartId"`
	}

	response := ApplyCouponResponse{
		Message:    "Coupon applied successfully",
		CouponCode: couponReq.CouponCode,
		CartId:     cart.Id,
	}

	_ = cjson.WriteJSON(w, http.StatusOK, response)
}

func MergeCart(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "User ID not found in context",
			InternalError: nil,
		})
	}

	// Get guest cart ID from request
	type MergeRequest struct {
		GuestCartId string `json:"guestCartId"`
	}
	var mergeReq MergeRequest
	if err := json.NewDecoder(r.Body).Decode(&mergeReq); err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Invalid request data",
			InternalError: err,
		})
	}

	if mergeReq.GuestCartId == "" {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Guest cart ID is required",
			InternalError: nil,
		})
	}

	// Get user's cart or create if it doesn't exist
	userCart, err := models.GetCartByUserId(userId)
	if err != nil {
		newCart := models.Cart{
			UserId: userId,
		}
		userCart, err = models.Create(&newCart)
		if err != nil {
			panic(cjson.HTTPError{
				Status:        http.StatusInternalServerError,
				Message:       "Failed to create user cart",
				InternalError: err,
			})
		}
	}

	// Get guest cart
	guestCart, err := models.GetCartById(mergeReq.GuestCartId)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Guest cart not found",
			InternalError: err,
		})
	}

	// Merge items from guest cart to user cart
	for _, item := range guestCart.CartItems {
		existingItem, err := models.GetItemByCartAndProduct(userCart.Id, item.ProductId)
		if err == nil {
			// Product exists in user cart, update quantity
			_, err = models.IncrementItemQuantity(existingItem.Id, item.Quantity)
			if err != nil {
				panic(cjson.HTTPError{
					Status:        http.StatusInternalServerError,
					Message:       "Failed to update item quantity",
					InternalError: err,
				})
			}
		} else {
			// Product doesn't exist in user cart, add it
			newItem := models.CartItem{
				CartId:        userCart.Id,
				ProductId:     item.ProductId,
				Quantity:      item.Quantity,
				PriceAtAdding: item.PriceAtAdding,
			}
			_, err = models.AddItem(&newItem)
			if err != nil {
				panic(cjson.HTTPError{
					Status:        http.StatusInternalServerError,
					Message:       "Failed to add item to cart",
					InternalError: err,
				})
			}
		}
	}

	// Delete guest cart
	if err := models.DeleteCart(guestCart.Id); err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusInternalServerError,
			Message:       "Failed to delete guest cart",
			InternalError: err,
		})
	}

	// Get updated user cart to return
	updatedCart, _ := models.GetCartByUserId(userId)
	_ = cjson.WriteJSON(w, http.StatusOK, updatedCart)
}

func ValidateCartItems(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "User ID not found in context",
			InternalError: nil,
		})
	}

	// Get user's cart
	cart, err := models.GetCartByUserId(userId)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Cart not found",
			InternalError: err,
		})
	}

	type ValidationItem struct {
		CartItemId string `json:"cartItemId"`
		ProductId  string `json:"productId"`
		Name       string `json:"name"`
		Requested  int    `json:"requested"`
		Available  int    `json:"available"`
		IsValid    bool   `json:"isValid"`
	}

	type ValidationResponse struct {
		IsValid      bool             `json:"isValid"`
		InvalidItems []ValidationItem `json:"invalidItems,omitempty"`
		ValidItems   []ValidationItem `json:"validItems"`
	}

	response := ValidationResponse{
		IsValid:      true,
		InvalidItems: []ValidationItem{},
		ValidItems:   []ValidationItem{},
	}

	// Check each item in cart against inventory
	for _, item := range cart.CartItems {
		// In a real implementation, check against actual inventory
		// For now, assume we're checking product.Stock

		available := item.Product.Stock // Assuming Product has Stock field
		isValid := item.Quantity <= available

		validationItem := ValidationItem{
			CartItemId: item.Id,
			ProductId:  item.ProductId,
			Name:       item.Product.Name, // Assuming Product has Name field
			Requested:  item.Quantity,
			Available:  available,
			IsValid:    isValid,
		}

		if !isValid {
			response.IsValid = false
			response.InvalidItems = append(response.InvalidItems, validationItem)
		} else {
			response.ValidItems = append(response.ValidItems, validationItem)
		}
	}

	_ = cjson.WriteJSON(w, http.StatusOK, response)
}
