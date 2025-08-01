package controller

import (
	"github.com/pratyush934/sibling-bond-server/cjson"
	"github.com/pratyush934/sibling-bond-server/models"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"time"
)

/*
CreateOrder - Convert cart to order
CancelOrder - Cancel existing order
UpdateOrderStatus - Change order status (admin only)
ProcessPayment - Handle payment for order
GetAllOrders - List all orders (admin only)
*/

func CreateOrder(w http.ResponseWriter, r *http.Request) {

	userId, ok := r.Context().Value("userId").(string)

	if userId == "" || !ok {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "There is no UserId, if there is we are not able to get it",
			InternalError: nil,
		})
	}

	shippingAddressId := r.URL.Query().Get("shippingAddressId")
	if shippingAddressId == "" {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "There is not shipping Id found",
			InternalError: nil,
		})
	}
	paymentMethod := r.URL.Query().Get("paymentMethod")
	if paymentMethod == "" {
		paymentMethod = "cod"
	}

	cartByUserId, err := models.GetCartByUserId(userId)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "There are is not a single cart found",
			InternalError: err,
		})
	}

	if len(cartByUserId.CartItems) == 0 {
		panic(cjson.HTTPError{
			Status:        http.StatusExpectationFailed,
			Message:       "There are not Items in the cart",
			InternalError: nil,
		})
	}

	orderItemsSlice := make([]models.OrderItem, 0, len(cartByUserId.CartItems))

	totalAmount := 0
	countItem := 0

	for _, v := range cartByUserId.CartItems {

		orderItem := models.OrderItem{
			ProductId:       v.ProductId,
			Quantity:        v.Quantity,
			PriceAtPurchase: v.PriceAtAdding,
		}

		totalAmount += v.PriceAtAdding * v.Quantity
		countItem += v.Quantity

		orderItemsSlice = append(orderItemsSlice, orderItem)
	}

	newOrderModel := models.Order{
		UserId:            userId,
		OrderItems:        orderItemsSlice,
		TotalAmount:       totalAmount,
		ShippingAddressId: shippingAddressId,
		PaymentMode:       paymentMethod,
		PaymentStatus:     "pending",
		Status:            "pending",
		TrackingNumber:    generateTrackingNumber(),
	}

	createdOrder, err := newOrderModel.Create()

	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusExpectationFailed,
			Message:       "Not able to create Order",
			InternalError: err,
		})
	}

	err = models.DeleteCart(cartByUserId.Id)
	if err != nil {
		log.Err(err).Msg("not able to delete the cart but order is now created")
	}
	_ = cjson.WriteJSON(w, http.StatusCreated, createdOrder)

}

func CancelOrder(w http.ResponseWriter, r *http.Request) {

}

func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {

	role, ok := r.Context().Value("role").(int)

	if role != 2 || !ok {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "Not able to get the role or if got the person is not admin",
			InternalError: nil,
		})
	}

	orderStatus := r.URL.Query().Get("orderStatus")
	if orderStatus == "" {
		orderStatus = "delivered"
	}
	orderId := r.URL.Query().Get("orderId")
	if orderId == "" {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "please provide orderId",
			InternalError: nil,
		})
	}

	status, err := models.UpdateStatus(orderId, strings.ToLower(orderStatus))
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusInternalServerError,
			Message:       "Not able to update the status",
			InternalError: err,
		})
	}

	_ = cjson.WriteJSON(w, http.StatusOK, status)
}

func generateTrackingNumber() int {
	return (int)(time.Now().Unix() % 1000000)
}
