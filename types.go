package bcex

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Side \"buy\" for Buy, \"sell\" for Sell
type Side string

// List of side
const (
	BUY Side = "BUY"
	SELL Side = "SELL"
)

// OrderStatus the model 'OrderStatus'
type OrderStatus string

// List of OrderStatus
const (
	OPEN        OrderStatus = "OPEN"
	REJECTED    OrderStatus = "REJECTED"
	CANCELED    OrderStatus = "CANCELED"
	FILLED      OrderStatus = "FILLED"
	PART_FILLED OrderStatus = "PART_FILLED"
	EXPIRED     OrderStatus = "EXPIRED"
)

// OrdType the model 'OrdType'
type OrdType string

// List of ordType
const (
	MARKET    OrdType = "MARKET"
	LIMIT     OrdType = "LIMIT"
	STOP      OrdType = "STOP"
	STOPLIMIT OrdType = "STOPLIMIT"
)

type Balance struct {
	Currency       string  `json:"currency"`
	Balance        float64 `json:"balance"`
	Available      float64 `json:"available"`
	BalanceLocal   float64 `json:"balance_local"`
	AvailableLocal float64 `json:"available_local"`
	Rate           float64 `json:"rate"`
}

type BalanceMap struct {
	Primary []Balance `json:"primary"`
}

type Fees struct {
	MakerRate   float64 `json:"makerRate"`
	TakerRate   float64 `json:"takerRate"`
	VolumeInUSD float64 `json:"volumeInUSD"`
}

type OrderBook struct {
	// Blockchain symbol identifier
	Symbol string           `json:"symbol,omitempty"`
	Bids   []OrderBookEntry `json:"bids,omitempty"`
	Asks   []OrderBookEntry `json:"asks,omitempty"`
}

type OrderBookEntry struct {
	Px  float64 `json:"px,omitempty"`
	Qty float64 `json:"qty,omitempty"`
	// Either the quantity of orders on this price level for L2, or the individual order id for L3
	Num int64 `json:"num,omitempty"`
}

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
	MinPriceIncrement      int64 `json:"min_price_increment,omitempty"`
	MinPriceIncrementScale int32 `json:"min_price_increment_scale,omitempty"`
	// The minimum quantity for an order for this instrument must be min_order_size*(10^-min_order_size_scale)
	MinOrderSize      int64 `json:"min_order_size,omitempty"`
	MinOrderSizeScale int32 `json:"min_order_size_scale,omitempty"`
	// The maximum quantity for an order for this instrument is max_order_size*(10^-max_order_size_scale). If this equal to zero, there is no limit
	MaxOrderSize      int64 `json:"max_order_size,omitempty"`
	MaxOrderSizeScale int32 `json:"max_order_size_scale,omitempty"`
	LotSize           int64 `json:"lot_size,omitempty"`
	LotSizeScale      int32 `json:"lot_size_scale,omitempty"`
	// Symbol status; open, close, suspend, halt, halt-freeze.
	Status string `json:"status,omitempty"`
	Id     int64  `json:"id,omitempty"`
	// If the symbol is halted and will open on an auction, this will be the opening price.
	AuctionPrice float64 `json:"auction_price,omitempty"`
	// Opening size
	AuctionSize float64 `json:"auction_size,omitempty"`
	// Opening time in HHMM format
	AuctionTime string `json:"auction_time,omitempty"`
	// Auction imbalance. If > 0 then there will be buy orders left over at the auction price. If < 0 then there will be sell orders left over at the auction price.
	Imbalance float64 `json:"imbalance,omitempty"`
}

//Ticker represents a Ticker from hitbtc API.
type Ticker struct {
	Symbol         string  `json:"symbol,omitempty"`
	Price24h       float64 `json:"price_24h,omitempty"`
	Volume24h      float64 `json:"volume_24h,omitempty"`
	LastTradePrice float64 `json:"last_trade_price,omitempty"`
}

// Tickers rapresents a set of a valid Tickers struct
type Tickers []Ticker

// Trade represents a single trade made by a user.
type Trade struct {
	Id            uint64    `json:"id"`
	OrderId       uint64    `json:"orderId"`
	ClientOrderId string    `json:"clientOrderId"`
	Symbol        string    `json:"symbol"`
	Type          string    `json:"side"`
	Price         float64   `json:"price,string"`
	Quantity      float64   `json:"quantity,string"`
	Fee           float64   `json:"fee,string"`
	Timestamp     time.Time `json:"timestamp"`
}

// GetOrdersOpts Optional parameters for the method 'GetOrders'
type GetOrdersOpts struct {
	Symbol string
	From   int64
	To     int64
	Status interface{}
	Limit  int32
}

func (opts GetOrdersOpts) parse() map[string]string {
	payload := make(map[string]string)
	if opts.Symbol != "" {
		payload["symbol"] = opts.Symbol
	}
	if opts.From != 0 {
		payload["from"] = parameterToString(opts.From, "")
	}
	if opts.To != 0 {
		payload["to"] = parameterToString(opts.To, "")
	}
	if opts.Status != nil {
		payload["status"] = parameterToString(opts.Status, "")
	}
	if opts.Limit != 0 {
		payload["limit"] = parameterToString(opts.Limit, "")
	}
	return payload
}

// GetFillsOpts Optional parameters for the method 'GetFills'
type GetFillsOpts struct {
    Symbol string
    From int64
    To int64
    Limit int32
}

func (opts GetFillsOpts) parse() map[string]string {
	payload := make(map[string]string)
	if opts.Symbol != "" {
		payload["symbol"] = opts.Symbol
	}
	if opts.From != 0 {
		payload["from"] = parameterToString(opts.From, "")
	}
	if opts.To != 0 {
		payload["to"] = parameterToString(opts.To, "")
	}
	if opts.Limit != 0 {
		payload["limit"] = parameterToString(opts.Limit, "")
	}
	return payload
}

// DeleteAllOrdersOpts Optional parameters for the method 'DeleteAllOrders'
type DeleteAllOrdersOpts struct {
    Symbol string
}

func (opts DeleteAllOrdersOpts) parse() map[string]string {
	payload := make(map[string]string)
	if opts.Symbol != "" {
		payload["symbol"] = opts.Symbol
	}
	return payload
}

// OrderSummary struct for OrderSummary
type OrderSummary struct {
	// The unique order id assigned by the exchange
	ExOrdId int64 `json:"exOrdId,omitempty"`
	// Reference field provided by client. Cannot exceed 20 characters, only alphanumeric characters are allowed.
	ClOrdId string `json:"clOrdId"`
	OrdType OrdType `json:"ordType"`
	OrdStatus OrderStatus `json:"ordStatus"`
	Side Side `json:"side"`
	// The limit price for the order
	Price float64 `json:"price,omitempty"`
	// The reason for rejecting the order, if applicable
	Text string `json:"text,omitempty"`
	// Blockchain symbol identifier
	Symbol string `json:"symbol"`
	// The executed quantity for the order's last fill
	LastShares float64 `json:"lastShares,omitempty"`
	// The executed price for the last fill
	LastPx float64 `json:"lastPx,omitempty"`
	// For Open and Partially Filled orders this is the remaining quantity open for execution. For Canceled and Expired orders this is the quantity than was still open before cancellation/expiration. For Rejected order this is equal to orderQty. For other states this is always zero.
	LeavesQty float64 `json:"leavesQty,omitempty"`
	// The quantity of the order which has been filled
	CumQty float64 `json:"cumQty,omitempty"`
	// Calculated the Volume Weighted Average Price of all fills for this order
	AvgPx float64 `json:"avgPx,omitempty"`
	// Time in ms since 01/01/1970 (epoch)
	Timestamp int64 `json:"timestamp,omitempty"`
}

// TimeInForce \"GTC\" for Good Till Cancel, \"IOC\" for Immediate or Cancel, \"FOK\" for Fill or Kill, \"GTD\" Good Till Date
type TimeInForce string

// List of TimeInForce
const (
	GTC TimeInForce = "GTC"
	IOC TimeInForce = "IOC"
	FOK TimeInForce = "FOK"
	GTD TimeInForce = "GTD"
)

// BaseOrder struct for BaseOrder
type BaseOrder struct {
	// Reference field provided by client. Cannot exceed 20 characters, only alphanumeric characters are allowed.
	ClOrdId string `json:"clOrdId"`
	OrdType OrdType `json:"ordType"`
	// Blockchain symbol identifier
	Symbol string `json:"symbol"`
	Side Side `json:"side"`
	// The order size in the terms of the base currency
	OrderQty float64 `json:"orderQty"`
	TimeInForce TimeInForce `json:"timeInForce,omitempty"`
	// The limit price for the order
	Price float64 `json:"price,omitempty"`
	// expiry date in the format YYYYMMDD
	ExpireDate int32 `json:"expireDate,omitempty"`
	// The minimum quantity required for an IOC fill
	MinQty float64 `json:"minQty,omitempty"`
	// The limit price for the order
	StopPx float64 `json:"stopPx,omitempty"`
}

// parameterToString convert interface{} parameters to string, using a delimiter if format is provided.
func parameterToString(obj interface{}, collectionFormat string) string {
	var delimiter string

	switch collectionFormat {
	case "pipes":
		delimiter = "|"
	case "ssv":
		delimiter = " "
	case "tsv":
		delimiter = "\t"
	case "csv":
		delimiter = ","
	}

	if reflect.TypeOf(obj).Kind() == reflect.Slice {
		return strings.Trim(strings.Replace(fmt.Sprint(obj), " ", delimiter, -1), "[]")
	} else if t, ok := obj.(time.Time); ok {
		return t.Format(time.RFC3339)
	}

	return fmt.Sprintf("%v", obj)
}