package bcex

// Symbol represents data of a Currency Pair on a market.
type Symbol struct {
	// Blockchain symbol identifier
	BaseCurrency string `json:"base_currency,omitempty"`
	// The number of decimals the currency can be split in
	BaseCurrencyScale int32 `json:"base_currency_scale,omitempty"`
	// Blockchain symbol identifier
	CounterCurrency string `json:"counter_currency,omitempty"`
	// The number of decimals the currency can be split in
	CounterCurrencyScale int32 `json:"counter_currency_scale,omitempty"`
	// The price of the instrument must be a multiple of min_price_increment * (10^-min_price_increment_scale)
	MinPriceIncrement int64 `json:"min_price_increment,omitempty"`
	MinPriceIncrementScale int32 `json:"min_price_increment_scale,omitempty"`
	// The minimum quantity for an order for this instrument must be min_order_size*(10^-min_order_size_scale)
	MinOrderSize int64 `json:"min_order_size,omitempty"`
	MinOrderSizeScale int32 `json:"min_order_size_scale,omitempty"`
	// The maximum quantity for an order for this instrument is max_order_size*(10^-max_order_size_scale). If this equal to zero, there is no limit
	MaxOrderSize int64 `json:"max_order_size,omitempty"`
	MaxOrderSizeScale int32 `json:"max_order_size_scale,omitempty"`
	LotSize int64 `json:"lot_size,omitempty"`
	LotSizeScale int32 `json:"lot_size_scale,omitempty"`
	// Symbol status; open, close, suspend, halt, halt-freeze.
	Status string `json:"status,omitempty"`
	Id int64 `json:"id,omitempty"`
	// If the symbol is halted and will open on an auction, this will be the opening price.
	AuctionPrice float64 `json:"auction_price,omitempty"`
	// Opening size
	AuctionSize float64 `json:"auction_size,omitempty"`
	// Opening time in HHMM format
	AuctionTime string `json:"auction_time,omitempty"`
	// Auction imbalance. If > 0 then there will be buy orders left over at the auction price. If < 0 then there will be sell orders left over at the auction price.
	Imbalance float64 `json:"imbalance,omitempty"`
}