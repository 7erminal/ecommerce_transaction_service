package models

type OrdersRequestDTO struct {
	Currency    int64
	Items       []Cart
	RequestType string
	Comment     string
	CreatedBy   int64
}

type Cart struct {
	ItemId   int64
	Quantity int64
	Price    float64
	Currency int64
}

type OrderResponseDTO struct {
	StatusCode int
	Order      *Orders
	StatusDesc string
}

type OrdersResponseDTO struct {
	StatusCode int
	Orders     *[]Orders
	StatusDesc string
}

type ConfirmOrderDTO struct {
	TransactionId string
	Status        string
	Confirmedby   string
}
