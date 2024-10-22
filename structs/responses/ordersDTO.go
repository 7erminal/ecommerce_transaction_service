package responses

import (
	"time"
	"transaction_service/models"
)

type OrdersCustom struct {
	OrderId     int64 `orm:"auto"`
	OrderNumber int64
	Quantity    int
	Cost        float32
	// Currency     *Currencies `orm:"rel(fk)"`
	Currency     int64
	OrderDate    time.Time `orm:"type(datetime)"`
	DateCreated  time.Time `orm:"type(datetime)"`
	DateModified time.Time `orm:"type(datetime)"`
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
