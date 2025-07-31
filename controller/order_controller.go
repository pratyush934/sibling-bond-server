package controller

import (
	"encoding/json"
	"github.com/pratyush934/sibling-bond-server/cjson"
	"github.com/pratyush934/sibling-bond-server/dto"
	"github.com/pratyush934/sibling-bond-server/models"
	"net/http"
)

/*
CreateOrder - Convert cart to order
CancelOrder - Cancel existing order
UpdateOrderStatus - Change order status (admin only)
ProcessPayment - Handle payment for order
GetAllOrders - List all orders (admin only)
*/

func CreateOrder(w http.ResponseWriter, r *http.Request) {

	/*
		IT IS WRONG

	*/

	userId, ok := r.Context().Value("userId").(string)

	if userId == "" || !ok {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "Not found the UserId",
			InternalError: nil,
		})
	}

	var orderModel dto.OrderModel
	if err := json.NewDecoder(r.Body).Decode(&orderModel); err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Model is not there",
			InternalError: err,
		})
	}

	newOrder := models.Order{
		UserId:            userId,
		ShippingAddressId: orderModel.ShippingAddressId,
		TotalAmount:       orderModel.TotalAmount,
		PaymentStatus:     orderModel.PaymentStatus,
		Status:            orderModel.Status,
		PaymentMode:       orderModel.PaymentMode,
		TrackingNumber:    orderModel.TrackingNumber,
	}

	if len(orderModel.OrderItems) > 0 {
		newOrderItems := make([]models.OrderItem, 0, len(orderModel.OrderItems))
		for _, v := range orderModel.OrderItems {
			newOrderItems = append(newOrderItems, models.OrderItem{
				ProductId:       v.ProductId,
				Quantity:        v.Quantity,
				PriceAtPurchase: v.PriceAtPurchase,
			})
		}
		newOrder.OrderItems = newOrderItems
	}

	createdOrder, err := newOrder.Create()

	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Not able to create the Order",
			InternalError: err,
		})
	}
	_ = cjson.WriteJSON(w, http.StatusCreated, createdOrder)

}
