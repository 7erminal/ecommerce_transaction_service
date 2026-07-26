package responses

import (
	"time"
)

type OrdersCustom struct {
	OrderId      int64 `orm:"auto"`
	OrderNumber  string
	Quantity     int
	Cost         float32
	Customer     *CustomersAlt
	Currency     string
	OrderDate    time.Time
	OrderEndDate time.Time
	ReturnedDate time.Time
	DateCreated  time.Time
	DateModified time.Time
	OrderDetails []OrderItemsCustom
}

type OrderResponseDTO struct {
	StatusCode int
	Order      *OrdersCustom
	StatusDesc string
}

type OrdersResponseDTO struct {
	StatusCode int
	Orders     *[]OrdersCustom
	StatusDesc string
}

type ItemAlt struct {
	ItemId       string
	ItemName     string
	Description  string
	Price        float32
	Category     string
	Currency     string
	DateCreated  time.Time
	DateModified time.Time
}

type OrderItemsCustom struct {
	OrderItemId int64
	OrderId     string
	Item        *ItemAlt
	Quantity    int
	Status      string
	OrderDate   time.Time
	Comment     string
}

type OrderItemsResponseDTO struct {
	StatusCode int
	OrderItems *[]OrderItemsCustom
	StatusDesc string
}

type OrderItemResponseDTO struct {
	StatusCode int
	OrderItem  *OrderItemsCustom
	StatusDesc string
}
