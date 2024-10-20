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
