package responses

type CallbackResponse struct {
	StatusCode    bool
	StatusMessage string
	Result        *Bil_transactionCustom
}
