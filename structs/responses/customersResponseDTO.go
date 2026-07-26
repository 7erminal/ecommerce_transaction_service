package responses

import "time"

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

type CustomersResponseDTO struct {
	StatusCode int
	Customers  *[]Customer
	StatusDesc string
}
