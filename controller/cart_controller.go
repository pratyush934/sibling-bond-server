package controller

import (
	"github.com/gorilla/mux"
	"github.com/pratyush934/sibling-bond-server/cjson"
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

}
