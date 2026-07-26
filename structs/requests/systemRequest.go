package requests

type GetCurrencyRequest struct {
	CurrencyId string
}

type GetCountryRequest struct {
	CountryId string
}

type GetStatusRequest struct {
	StatusId string
}

type GetItemRequest struct {
	ItemId string
}

type UpdateItemQuantityRequest struct {
	ItemId   string
	Quantity int
}

type GetBillerRequest struct {
	BillerId string
}

type GetOperatorRequest struct {
	OperatorId string
}

type GetServiceRequest struct {
	ServiceId string
}
