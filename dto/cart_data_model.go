package dto

type CartDataModel struct {
	UserId    string          `json:"userId"`
	CartItems []CartItemModel `json:"cartItems"`
}

type CartItemModel struct {
	Id            string       `json:"id,omitempty"`
	ProductId     string       `json:"productId"`
	Quantity      int          `json:"quantity"`
	PriceAtAdding int          `json:"priceAdding,omitempty"`
	Product       ProductModel `json:"product,omitempty"`
}
