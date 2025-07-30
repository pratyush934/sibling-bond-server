package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pratyush934/sibling-bond-server/cjson"
	"github.com/pratyush934/sibling-bond-server/dto"
	"github.com/pratyush934/sibling-bond-server/models"
	"net/http"
)

/*
GetCart - Retrieve user's current cart
CreateCart - Initialize a new cart for a user
AddToCart - Add product to cart (handles new items and quantity increases)
UpdateCartItem - Update quantity of existing cart item
RemoveFromCart - Remove specific item from cart
ClearCart - Remove all items from user's cart
ApplyCoupon - Apply discount code to cart
GetCartTotal - Calculate total price of items in cart
MergeCart - Merge guest cart with user cart after login
ValidateCartItems - Check if cart items are still available in inventory
*/

func GetCart(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	cartId := vars["id"]

	if cartId == "" {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Not able to get the UserId",
			InternalError: nil,
		})
	}
	cartById, err := models.GetCartById(cartId)
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
