package requests

type Cart struct {
	ItemId   int64
	Quantity int64
	// Price    float64
	// Currency int64
}

type OrdersRequestDTO struct {
	Currency      int64
	Items         []Cart
	RequestType   string
	Comment       string
	CreatedBy     int64
	OrderDate     string
	OrderEndDate  string
	OrderLocation string
	Customer      int64
	Branch        int64
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
