package controller

import (
	"encoding/json"
	"github.com/pratyush934/sibling-bond-server/cjson"
	"github.com/pratyush934/sibling-bond-server/database"
	"github.com/pratyush934/sibling-bond-server/models"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
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
	// Get authenticated user ID from context
	userId, ok := r.Context().Value("userId").(string)
	if userId == "" || !ok {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "Authentication required",
			InternalError: nil,
		})
	}

	// Get order ID from URL params
	orderId := r.URL.Query().Get("orderId")
	if orderId == "" {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Order ID is required",
			InternalError: nil,
		})
	}

	// Get the order and verify ownership
	order, err := models.GetOrderByUserIdAndOrderId(userId, orderId)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Order not found or doesn't belong to you",
			InternalError: err,
		})
	}

	// Check if order is in a cancellable state
	// Only allow cancellation for pending or confirmed orders
	if order.Status != "pending" && order.Status != "confirmed" {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       fmt.Sprintf("Order in '%s' state cannot be cancelled", order.Status),
			InternalError: nil,
		})
	}

	// Delete the order
	err = models.DeleteOrderById(orderId)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusInternalServerError,
			Message:       "Failed to cancel order",
			InternalError: err,
		})
	}

	_ = cjson.WriteJSON(w, http.StatusOK, "Order cancelled successfully")
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

func GetAllOrders(w http.ResponseWriter, r *http.Request) {

	role, ok := r.Context().Value("role").(int)

	if role != 2 || !ok {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "You are not authorized to get the orders",
			InternalError: nil,
		})
	}
	limitN := 0
	offSetN := 0

	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limitN = 10
	} else {
		limitN, _ = strconv.Atoi(limit)
	}
	offSet := r.URL.Query().Get("offSet")
	if offSet == "" {
		offSetN = 5
	} else {
		offSetN, _ = strconv.Atoi(offSet)
	}

	//fmt.Println(limitN, offSetN)

	orderAll, err := models.GetAllOrder(limitN, offSetN, "")
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Not able to get all ordes",
			InternalError: err,
		})
	}

	_ = cjson.WriteJSON(w, http.StatusOK, orderAll)

}

func ProcessPayment(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID
	userId, ok := r.Context().Value("userId").(string)
	if userId == "" || !ok {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "Authentication required",
			InternalError: nil,
		})
	}

	// Get order ID from URL params
	orderId := r.URL.Query().Get("orderId")
	if orderId == "" {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Order ID is required",
			InternalError: nil,
		})
	}

	// Get payment details from request
	var paymentDetails struct {
		PaymentMethod string `json:"paymentMethod"`
		TransactionID string `json:"transactionId,omitempty"`
		PaymentAmount int    `json:"paymentAmount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&paymentDetails); err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Invalid payment details",
			InternalError: err,
		})
	}

	// Get the order and verify ownership
	order, err := models.GetOrderByUserIdAndOrderId(userId, orderId)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Order not found or doesn't belong to you",
			InternalError: err,
		})
	}

	// Verify payment amount matches order total
	if paymentDetails.PaymentAmount != order.TotalAmount {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Payment amount doesn't match order total",
			InternalError: nil,
		})
	}

	// Validate payment method
	validPaymentMethods := []string{"credit_card", "debit_card", "upi", "netbanking", "cod", "wallet"}
	isValidMethod := false
	for _, method := range validPaymentMethods {
		if method == paymentDetails.PaymentMethod {
			isValidMethod = true
			break
		}
	}

	if !isValidMethod {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Invalid payment method",
			InternalError: nil,
		})
	}

	// Update order with payment information
	tx := database.DB.Begin()

	// Update payment method if different from existing
	if order.PaymentMode != paymentDetails.PaymentMethod {
		order.PaymentMode = paymentDetails.PaymentMethod
	}

	// Update payment status to completed
	order.PaymentStatus = "completed"

	// If payment is completed, update order status to confirmed
	if order.Status == "pending" {
		order.Status = "confirmed"
	}

	if err := tx.Save(order).Error; err != nil {
		tx.Rollback()
		log.Err(err).Msg("Failed to update payment status")
		panic(cjson.HTTPError{
			Status:        http.StatusInternalServerError,
			Message:       "Failed to process payment",
			InternalError: err,
		})
	}

	tx.Commit()

	// Return the updated order
	_ = cjson.WriteJSON(w, http.StatusOK, order)
}

func generateTrackingNumber() int {
	return (int)(time.Now().Unix() % 1000000)
}
