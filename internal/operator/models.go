package operator

type WithdrawRequest struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	RefID    string `json:"refId"`
}

type WithdrawResponse struct {
	Status       string `json:"status"`
	Balance      int64  `json:"balance"`
	ErrorMessage string `json:"error,omitempty"`
}

type DepositRequest struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	RefID    string `json:"refId"`
}

type DepositResponse struct {
	Status       string `json:"status"`
	Balance      int64  `json:"balance"`
	ErrorMessage string `json:"error,omitempty"`
}
