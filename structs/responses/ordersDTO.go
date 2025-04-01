package responses

import (
	"time"
	"transaction_service/models"
)

type OrdersCustom struct {
	OrderId     int64 `orm:"auto"`
	OrderNumber string
	Quantity    int
	Cost        float32
	Customer    *models.Customers
	// Currency     *Currencies `orm:"rel(fk)"`
	CurrencyId   int64                 `orm:"column(currency)"`
	OrderDate    time.Time             `orm:"type(datetime)"`
	DateCreated  time.Time             `orm:"type(datetime)"`
	DateModified time.Time             `orm:"type(datetime)"`
	OrderDetails []*models.Order_items `orm:"reverse(many);null;"`
}

type OrderResponseDTO struct {
	StatusCode int
	Order      *OrdersCustom
	StatusDesc string
}

type OrdersResponseDTO struct {
	StatusCode int
	Orders     *[]models.Orders
	StatusDesc string
}

type OrderItemsCustom struct {
	OrderItemId int64
	Order       *OrdersCustom
	Item        *models.Items
	Quantity    int64
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
