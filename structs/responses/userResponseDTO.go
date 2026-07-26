package responses

import (
	"time"
)

type Currencies struct {
	CurrencyId   int64
	Symbol       string
	Currency     string
	Active       int
	DateCreated  time.Time
	DateModified time.Time
	CreatedBy    int
	ModifiedBy   int
}

type Countries struct {
	CountryId       int64
	Country         string
	Description     string
	CountryCode     string
	DefaultCurrency int64
	DateCreated     time.Time
	DateModified    time.Time
	CreatedBy       int
	ModifiedBy      int
}

type Branches struct {
	BranchId     int64
	Branch       string
	Country      int64
	Location     string
	PhoneNumber  string
	Active       int
	DateCreated  time.Time
	DateModified time.Time
	CreatedBy    int
	ModifiedBy   int
}

type Shops struct {
	ShopId              int64 `orm:"auto"`
	ShopName            string
	ShopDescription     string `orm:"size(255)"`
	ShopAssistantName   string `orm:"size(100)"`
	ShopAssistantNumber string `orm:"size(100)"`
	PhoneNumber         string
	Email               string
	Image               string    `orm:"size(100);omitempty"`
	DateCreated         time.Time `orm:"type(datetime)"`
	DateModified        time.Time `orm:"type(datetime)"`
	CreatedBy           int
	ModifiedBy          int
	Active              int
}

type UserExtraDetails struct {
	UserDetailsId int64
	Branch        *Branches
	Shop          *Shops
	Nickname      string
	DateCreated   time.Time
	DateModified  time.Time
	CreatedBy     int
	ModifiedBy    int
	Active        int
}

type Roles struct {
	RoleId       int64
	Role         string
	Description  string
	DateCreated  time.Time
	DateModified time.Time
	CreatedBy    int
	ModifiedBy   int
	Active       int
}

type Users struct {
	UserId        int64 `orm:"auto"`
	UserType      int
	UserDetails   *UserExtraDetails
	ImagePath     string
	FullName      string
	Username      string
	Password      string
	Email         string
	PhoneNumber   string
	Gender        string
	Dob           time.Time
	Address       string
	IdType        string
	IdNumber      string
	MaritalStatus string
	Role          *Roles
	Active        int
	IsVerified    bool
	DateCreated   time.Time
	DateModified  time.Time
	CreatedBy     int
	ModifiedBy    int
}

type UserResp struct {
	UserId        int64
	ImagePath     string
	UserType      int
	FullName      string
	Username      string
	Password      string
	Email         string
	PhoneNumber   string
	Gender        string
	Dob           time.Time
	Address       string
	IdType        string
	IdNumber      string
	MaritalStatus string
	Active        int
	Role          *Roles
	IsVerified    bool
	DateCreated   time.Time
	DateModified  time.Time
	CreatedBy     int
	ModifiedBy    int
	Branch        *Branches
}

type UserResponseDTO struct {
	StatusCode int
	User       *Users
	StatusDesc string
}

type UserTokenResponseDTO struct {
	IsValid bool
	User    *Users
}

type Identification_types struct {
	IdentificationTypeId int64
	Name                 string
	Code                 string
	DateCreated          time.Time
	DateModified         time.Time
	CreatedBy            int
	ModifiedBy           int
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

type Customer_emergency_contacts struct {
	CustomerEmergencyContactId int64
	Name                       string
	Contact                    string
	Customer                   *Customers
	DateCreated                time.Time
	DateModified               time.Time
	CreatedBy                  int
	ModifiedBy                 int
}

type Customer_guarantors struct {
	CustomerGuarantorId int64
	Name                string
	Contact             string
	Customer            *Customers
	DateCreated         time.Time
	DateModified        time.Time
	CreatedBy           int
	ModifiedBy          int
}

type Customers struct {
	CustomerId           int64
	CustomerNumber       string
	FullName             string
	ImagePath            string
	Email                string
	PhoneNumber          string
	Gender               string
	Location             string
	IdentificationType   *Identification_types
	IdentificationNumber string
	Branch               *Branches
	Shop                 *Shops
	CustomerCategory     *Customer_categories
	Nickname             string
	Dob                  time.Time
	DateCreated          time.Time
	DateModified         time.Time
	CreatedBy            int
	ModifiedBy           int
	Active               int
	LastTxnDate          time.Time
	EmergencyContacts    []*Customer_emergency_contacts
	Guarantors           []*Customer_guarantors
}

type CustomersAlt struct {
	CustomerId    string
	CustomerName  string
	CustomerEmail string
	CustomerPhone string
}

type CustomerTokenResponseDTO struct {
	IsValid  bool
	Customer *Customers
}

type CustomerResponseAltDTO struct {
	StatusCode int
	Result     *Customers
	StatusDesc string
}

type CustomerResponseDTO struct {
	StatusCode int
	Customer   *Customers
	StatusDesc string
}

type BranchResponseDTO struct {
	StatusCode int
	Branch     *Branches
	StatusDesc string
}

type LoginResponseDTO struct {
	StatusCode   int
	AccessToken  string
	RefreshToken string
	User         *Users
	StatusDesc   string
}
