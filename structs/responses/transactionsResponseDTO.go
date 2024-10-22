package responses

import (
	"time"
	"transaction_service/models"
)

type TransactionResponseDTO struct {
	StatusCode  int
	Transaction *models.Transactions
	StatusDesc  string
}

type TransactionsResponseDTO struct {
	StatusCode   int
	Transactions *[]models.Transactions
	StatusDesc   string
}

type TransactionCustomResponseDTO struct {
	StatusCode  int
	Transaction *TransactionsCustom
	StatusDesc  string
}

type TransactionsCustomResponseDTO struct {
	StatusCode   int
	Transactions *[]TransactionsCustom
	StatusDesc   string
}

type TransactionsCustom struct {
	TransactionId       int64         `orm:"auto"`
	Order               *OrdersCustom `orm:"rel(fk)"`
	Amount              float32
	TransactingCurrency int64
	StatusId            int64
	DateCreated         time.Time `orm:"type(datetime)"`
	DateModified        time.Time `orm:"type(datetime)"`
	CreatedBy           int
	ModifiedBy          int
	Active              int
}
