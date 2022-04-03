package bcex

//Ticker represents a Ticker from hitbtc API.
type Ticker struct {
	Symbol string `json:"symbol,omitempty"`
	Price24h float64 `json:"price_24h,omitempty"`
	Volume24h float64 `json:"volume_24h,omitempty"`
	LastTradePrice float64 `json:"last_trade_price,omitempty"`
}

// Tickers rapresents a set of a valid Tickers struct
type Tickers []Ticker