package responses

type CallbackResponse struct {
	StatusCode    int
	StatusMessage string
	Result        *Bil_transactionCustom
}
