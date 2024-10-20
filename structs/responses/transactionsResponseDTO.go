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

type TransactionsCustomResponseDTO struct {
	StatusCode   int
	Transactions *[]TransactionsCustom
	StatusDesc   string
}

type OrdersCustom struct {
	OrderId  int64 `orm:"auto"`
	Quantity int
	Cost     float32
	// Currency     *Currencies `orm:"rel(fk)"`
	Currency     int64
	OrderDate    time.Time `orm:"type(datetime)"`
	DateCreated  time.Time `orm:"type(datetime)"`
	DateModified time.Time `orm:"type(datetime)"`
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
