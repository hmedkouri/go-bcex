package bcex

type OrderBook struct {
	// Blockchain symbol identifier
	Symbol string `json:"symbol,omitempty"`
	Bids []OrderBookEntry `json:"bids,omitempty"`
	Asks []OrderBookEntry `json:"asks,omitempty"`
}

type OrderBookEntry struct {
	Px float64 `json:"px,omitempty"`
	Qty float64 `json:"qty,omitempty"`
	// Either the quantity of orders on this price level for L2, or the individual order id for L3
	Num int64 `json:"num,omitempty"`
}