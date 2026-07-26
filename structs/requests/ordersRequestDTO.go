package requests

type Cart struct {
	ItemId   string
	Quantity int
	// Price    float64
	// Currency int64
}

type OrdersRequestDTO struct {
	Currency      string
	Items         []Cart
	RequestType   string
	Comment       string
	CreatedBy     string
	OrderDate     string
	OrderEndDate  string
	OrderLocation string
	Customer      string
	Branch        string
}

type GetUserOrdersRequest struct {
	Id       int64
	FromDate string
	ToDate   string
}

type ConfirmOrderDTO struct {
	TransactionId string
	Status        string
	Confirmedby   string
}

type UpdateOrderItemDTO struct {
	OrderItemId int64
	Status      string
	Comment     string
	ModifiedBy  int64
}
