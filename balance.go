package bcex

type Balance struct {
	Currency string `json:"currency"`
	Balance float64 `json:"balance"`
	Available float64 `json:"available"`
	BalanceLocal float64 `json:"balance_local"`
	AvailableLocal float64 `json:"available_local"`
	Rate float64 `json:"rate"`
}

type BalanceMap struct {
	Primary []Balance `json:"primary"`
}