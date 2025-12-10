package requests

type ExtraData struct {
	ExtraData1 string
	ExtraData2 string
	ExtraData3 string
}

type BilTransactionRequestDTO struct {
	Source          string
	PhoneNumber     string
	Amount          float64
	Destination     string
	ClientReference string
	Package         string
	ServiceCode     string
	RequestId       string
	ExtraData       ExtraData
	BillerCode      string
}
