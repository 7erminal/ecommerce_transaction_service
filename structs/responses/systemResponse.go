package responses

import "time"

type CurrencyObject struct {
	CurrencyId   int64
	Symbol       string
	Currency     string
	Active       int
	DateCreated  time.Time
	DateModified time.Time
	CreatedBy    int
	ModifiedBy   int
}

type CountryObject struct {
	CountryId       int64
	Country         string
	Description     string
	CountryCode     string
	DefaultCurrency *CurrencyObject
	DateCreated     time.Time
	DateModified    time.Time
	CreatedBy       int
	ModifiedBy      int
}

type StatusObject struct {
	Active       int
	CreatedBy    int
	DateCreated  time.Time
	DateModified time.Time
	ModifiedBy   int
	Status       string
	StatusCode   string
	StatusId     int64
}

type CurrencyResponseDTO struct {
	StatusCode int
	Result     CurrencyObject
	StatusDesc string
}

type CountryResponseDTO struct {
	StatusCode int
	Result     CountryObject
	StatusDesc string
}

type StatusResponseDTO struct {
	StatusCode int
	Result     StatusObject
	StatusDesc string
}

type Categories struct {
	CategoryId   int64
	CategoryName string
	ImagePath    string
	Icon         string
	Description  string
	Active       int8
	DateCreated  time.Time
	DateModified time.Time
	CreatedBy    int
	ModifiedBy   int
}

type Item_prices struct {
	ItemPriceId   int64
	ItemPrice     float32
	AltItemPrice  float32
	ShowAltPrice  bool
	Discount      string
	Discount_type string
	ExtraCharges  float32
	Currency      int64
	Active        int
	DateCreated   time.Time
	DateModified  time.Time
	CreatedBy     int
	ModifiedBy    int
}

type Item_quantity struct {
	ItemQuantityId int64
	Item           *ItemObject
	Quantity       int
	QuantityAlert  int
	Active         int
	DateCreated    time.Time
	DateModified   time.Time
	CreatedBy      int
	ModifiedBy     int
}

type ItemStatus struct {
	StatusId     int64
	Status       string
	StatusCode   string
	DateCreated  time.Time
	DateModified time.Time
	CreatedBy    int
	ModifiedBy   int
	Active       int
}

type Features struct {
	FeatureId    int64
	FeatureName  string
	ImagePath    string
	Visible      bool
	Description  string
	Active       int
	DateCreated  time.Time
	DateModified time.Time
	CreatedBy    int
	ModifiedBy   int
}

type Purposes struct {
	PurposeId    int64
	Purpose      string
	ImagePath    string
	Visible      bool
	Description  string
	Active       int
	DateCreated  time.Time
	DateModified time.Time
	CreatedBy    int
	ModifiedBy   int
}

type ItemObject struct {
	ItemId          int64
	ItemName        string
	Description     string
	Weight          string
	Category        *Categories
	ItemPrice       *Item_prices
	AvailableSizes  string
	AvailableColors string
	Material        string
	ImagePath       string
	Quantity        int
	Active          int
	DateCreated     time.Time
	DateModified    time.Time
	CreatedBy       int
	ModifiedBy      int
	Country         int64
	Branch          int64
	Status          *ItemStatus
	LastOrderDate   time.Time
	ItemQuantity    *Item_quantity
	ItemFeatures    []*Features
	ItemPurposes    []*Purposes
}

type ItemResponseDTO struct {
	StatusCode int
	Item       ItemObject
	StatusDesc string
}

type Billers struct {
	BillerId          int64
	BillerName        string
	BillerCode        string
	BillerReferenceId string
	Description       string
	Operator          string
	DateCreated       string
	DateModified      string
	CreatedBy         int
	ModifiedBy        int
	Active            int
}

type Services struct {
	ServiceId          int64
	ServiceName        string
	ServiceCode        string
	ServiceDescription string
	DateCreated        string
	DateModified       string
	CreatedBy          int
	ModifiedBy         int
	Active             int
}

type Operator struct {
	OperatorId   int64
	OperatorName string
	Description  string
	DateCreated  string
	DateModified string
	CreatedBy    int
	ModifiedBy   int
	Active       int
}

type BillerResponseDTO struct {
	StatusCode int
	Biller     Billers
	StatusDesc string
}

type BillersResponseDTO struct {
	StatusCode int
	Result     []Billers
	StatusDesc string
}

type ServiceResponseDTO struct {
	StatusCode int
	Service    Services
	StatusDesc string
}

type ServicesResponseDTO struct {
	StatusCode int
	Result     []Services
	StatusDesc string
}

type OperatorResponseDTO struct {
	StatusCode int
	Operator   Operator
	StatusDesc string
}
