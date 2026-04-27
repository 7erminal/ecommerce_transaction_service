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
	Transactions *[]interface{}
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
	Status              string
	DateCreated         time.Time `orm:"type(datetime)"`
	DateModified        time.Time `orm:"type(datetime)"`
	CreatedBy           int
	ModifiedBy          int
	Active              int
	Branch              *models.Branches `orm:"rel(fk)"`
}

type BilTransactionResponseDTO struct {
	StatusCode int
	Result     *models.Bil_transactions
	StatusDesc string
}

type UserTransactionResponseDTO struct {
	StatusCode int
	Result     *UserTransactions
	StatusDesc string
}

type UserTransactionsResponseDTO struct {
	StatusCode int
	Result     *[]UserTransactions
	StatusDesc string
}

type BilTransactionsResponseDTO struct {
	StatusCode int
	Result     *[]interface{}
	StatusDesc string
}

type Bil_ins_transactionCustom struct {
	BilInsTransactionId    int64
	Amount                 string
	Biller                 string
	SenderAccountNumber    string
	RecipientAccountNumber string
	Network                string
	Request                string
	Response               string
	Active                 int
}

type Bil_transactionCustom struct {
	TransactionId           string
	TransactionRefNumber    string
	Service                 string
	BillerCode              string
	Amount                  string
	TransactingCurrency     string
	SourceChannel           string
	Source                  string
	Destination             string
	Package                 string
	Charge                  string
	Commission              string
	ExternalReferenceNumber string
	Status                  string
	ExtraDetails1           string
	ExtraDetails2           string
	ExtraDetails3           string
	DateProcessed           time.Time
	Active                  int
	InsTxns                 []*Bil_ins_transactionCustom
}

type UserInsTransactions struct {
	UserInsTransactionId   string
	UserTransactionId      string
	Amount                 float64
	Data                   string
	SenderAccountNumber    string
	RecipientAccountNumber string
	Service                string
	Status                 string
	Request                string
	Response               string
}

type UserTransactions struct {
	TransactionId                string
	Service                      string
	TransactionCustomerReference string
	Amount                       float64
	TransactingCurrency          string
	SourceChannel                string
	Source                       string
	Destination                  string
	Package                      string
	Charge                       float64
	Commission                   float64
	Reference                    string
	ExternalReferenceNumber      string
	Status                       string
	ExtraDetails1                string
	ExtraDetails2                string
	ExtraDetails3                string
	ClientResponseCode           string
	DateCreated                  time.Time
	DateModified                 time.Time
	CreatedBy                    string
	ModifiedBy                   string
	Active                       int
	TransactionDetails           *[]UserInsTransactions
}
