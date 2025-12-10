package requests

type CallbackMeta struct {
	Commission string
}

type CallbackData struct {
	AmountCharged         float64
	TransactionId         *string
	ClientReference       *string
	Description           *string
	ExternalTransactionId *string
	Amount                float64
	Charges               float64
	AmountAfterCharges    float64
	PaymentDate           *string
	OrderId               *string
}

type CallbackRequest struct {
	AmountCharged         float64
	TransactionId         *string
	ClientReference       *string
	Description           *string
	ExternalTransactionId *string
	Amount                float64
	Charges               float64
	AmountAfterCharges    float64
	PaymentDate           *string
	OrderId               *string
	Status                string
}
