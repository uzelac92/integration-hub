package operator

type WithdrawRequest struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	RefID    string `json:"refId"`
}

func (w WithdrawRequest) ToMap() map[string]any {
	return map[string]any{
		"amount":   w.Amount,
		"currency": w.Currency,
		"refId":    w.RefID,
	}
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

func (d DepositRequest) ToMap() map[string]any {
	return map[string]any{
		"amount":   d.Amount,
		"currency": d.Currency,
		"refId":    d.RefID,
	}
}

type DepositResponse struct {
	Status       string `json:"status"`
	Balance      int64  `json:"balance"`
	ErrorMessage string `json:"error,omitempty"`
}
