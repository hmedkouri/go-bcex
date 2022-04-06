package bcex

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	API_BASE = "https://api.blockchain.com/v3/exchange" // BCEX API endpoint
)

type Bcex struct {
	client *client
}

// New returns an instantiated HitBTC struct
func New(apiKey, apiSecret string) *Bcex {
	client := NewClient(apiKey, apiSecret)
	return &Bcex{client}
}

// GetSymbols is used to get the open and available trading markets along with other meta data.
func (b *Bcex) GetSymbols() (symbols []Symbol, err error) {
	r, err := b.client.do("GET", "symbols", nil, false)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var m = make(map[string]Symbol)
	err = json.Unmarshal(r, &m)
	for _, value := range m {
        symbols = append(symbols, value)
    }
	return
}

// GetSymbol is used to get the current symbol data for a market.
func (b *Bcex) GetSymbol(market string) (symbol Symbol, err error) {
	r, err := b.client.do("GET", "symbols/"+strings.ToUpper(market), nil, false)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	err = json.Unmarshal(r, &symbol)
	return
}

// GetAllTicker is used to get the current ticker values for all markets.
func (b *Bcex) GetAllTicker() (tickers Tickers, err error) {
	r, err := b.client.do("GET", "tickers", nil, false)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	err = json.Unmarshal(r, &tickers)
	return
}

// GetTicker is used to get the current ticker values for a market.
func (b *Bcex) GetTicker(market string) (ticker Ticker, err error) {
	r, err := b.client.do("GET", "tickers/"+strings.ToUpper(market), nil, false)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	err = json.Unmarshal(r, &ticker)
	return
}

// GetL2Orderbook is used to get the current level 2 order book for a market.
func (b *Bcex) GetL2Orderbook(market string) (orderbook OrderBook, err error) {
	r, err := b.client.do("GET", "l2/"+strings.ToUpper(market), nil, false)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	err = json.Unmarshal(r, &orderbook)
	return
}

// GetL3Orderbook is used to get the current level 2 order book for a market.
func (b *Bcex) GetL3Orderbook(market string) (orderbook OrderBook, err error) {
	r, err := b.client.do("GET", "l3/"+strings.ToUpper(market), nil, false)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	err = json.Unmarshal(r, &orderbook)
	return
}

// GetFees is used to retrieve the fees from your account
func (b *Bcex) GetFees() (fees Fees, err error) {
	r, err := b.client.do("GET", "fees", nil, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	err = json.Unmarshal(r, &fees)
	return
}

// GetBalances is used to retrieve all balances from your account
func (b *Bcex) GetBalances() (balances BalanceMap, err error) {
	r, err := b.client.do("GET", "accounts", nil, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	err = json.Unmarshal(r, &balances)
	return
}

// GetTrades used to retrieve your trade history.
// market string literal for the market (ie. BTC/LTC). If set to "all", will return for all market
func (b *Bcex) GetTrades(currencyPair string) (trades []Trade, err error) {
	payload := make(map[string]string)
	if currencyPair != "all" {
		payload["symbol"] = currencyPair
	}
	r, err := b.client.do("GET", "trades", payload, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	err = json.Unmarshal(r, &trades)
	return
}

/*
GetOrders Get a list orders
Returns live and historic orders, defaulting to live orders. Returns at most 100 results, use timestamp to paginate for further results
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param nil or *GetOrdersOpts - Optional Parameters:
 * @param "Symbol" (string) -  Only return results for this symbol
 * @param "From" (int64) -  Epoch timestamp in ms
 * @param "To" (int64) -  Epoch timestamp in ms
 * @param "Status" (interface of OrderStatus) -  Order Status
 * @param "Limit" (int32) -  Maximum amount of results to return in a single call. If omitted, 100 results are returned by default. 
@return []OrderSummary
*/
func (b *Bcex) GetOrders(options *GetOrdersOpts) (orders []OrderSummary, err error) {
	var payload map[string]string
	if (options != nil) {
		payload = options.parse()
	}
	r, err := b.client.do("GET", "orders", payload, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	err = json.Unmarshal(r, &orders)
	return
}

// handleErr gets JSON response from API and deal with error
func handleErr(r interface{}) error {
	switch v := r.(type) {
	case map[string]interface{}:
		error := r.(map[string]interface{})["error"]
		if error != nil {
			switch v := error.(type) {
			case map[string]interface{}:
				errorMessage := error.(map[string]interface{})["message"]
				return errors.New(errorMessage.(string))
			default:
				return fmt.Errorf("I don't know about type %T!\n", v)
			}
		}
	case []interface{}:
		return nil
	default:
		return fmt.Errorf("I don't know about type %T!\n", v)
	}

	return nil
}