package requests

type GetUserOrdersRequest struct {
	Id       int64
	FromDate string
	ToDate   string
}
