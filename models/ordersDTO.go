package models

type OrdersRequestDTO struct {
	TotalQuantity          string
	Currency               string
	Items                  *[]Cart
	RequestType            string
	Comment                string
	SenderAccountNumber    string
	RecipientAccountNumber string
	Cost                   string
	CreatedBy              string
}

type Cart struct {
	ItemId   string
	Quantity string
	Price    string
	Currency string
}

type OrdersResponseDTO struct {
	StatusCode int
	Order      *Orders
	StatusDesc string
}
