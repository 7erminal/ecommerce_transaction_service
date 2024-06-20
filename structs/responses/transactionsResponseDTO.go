package responses

import (
	"transaction_service/models"
)

type TransactionResponseDTO struct {
	StatusCode  int
	Transaction *models.Transactions
	StatusDesc  string
}
