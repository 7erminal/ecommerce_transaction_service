package requests

import (
	"transaction_service/models"
)

type TransactionRequestDTO struct {
	StatusCode  int
	Transaction *models.Transactions
	StatusDesc  string
}

type UpdateTransactionRequestDTO struct {
	SenderAccountNumber    string
	RecipientAccountNumber string
}

type GetUserTransactionsByDateRequest struct {
	Id       int64
	FromDate string
	ToDate   string
}

type GetUserTransactionsRequest struct {
	Id    int64
	Limit int
}

type UserTransactionRequestDTO struct {
	SourceChannel            string
	SourceAccountNumber      string
	PhoneNumber              string
	Amount                   float64
	DestinationAccountNumber string
	Reference                string
	ClientReference          string
	Package                  string
	ServiceCode              string
	RequestId                string
	ExtraData                ExtraData
	CreatedBy                string
	Status                   string
}

type UpdateUserTransactionRequest struct {
	ClientReference    string
	Status             string
	ClientResponseCode string
}
