package models

type OrdersRequestDTO struct {
	TotalQuantity string
	Currency      string
	Items         []Cart
	RequestType   string
	Comment       string
	Cost          string
	CreatedBy     string
}

type Cart struct {
	ItemId   int64
	Quantity int64
	Price    float64
	Currency int64
}

type OrdersResponseDTO struct {
	StatusCode int
	Order      *Orders
	StatusDesc string
}

type ConfirmOrderDTO struct {
	TransactionId string
	Status        string
	Confirmedby   string
}
