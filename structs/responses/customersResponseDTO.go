package responses

import "time"

type Identification_types struct {
	IdentificationTypeId int64  `orm:"auto"`
	Name                 string `orm:"size(100)"`
	Code                 string `orm:"size(100)"`
	Active               int
}

type Customer_categories struct {
	CustomerCategoryId int64
	Category           string
	Description        string
	DateCreated        time.Time
	DateModified       time.Time
	CreatedBy          int
	ModifiedBy         int
	Active             int
}

type Customer struct {
	CustomerId int64
	// FullName             string
	// Email                string
	// PhoneNumber          string
	// Location             string
	// Nickname             string
	// Dob                  time.Time
	// DateCreated          time.Time
	// DateModified         time.Time
	// CreatedBy            int
	// ModifiedBy           int
	Active      int
	LastTxnDate time.Time
}

type CustomerResponseDTO struct {
	StatusCode int
	Customer   *Customer
	StatusDesc string
}

type CustomersResponseDTO struct {
	StatusCode int
	Customers  *[]Customer
	StatusDesc string
}
