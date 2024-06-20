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
