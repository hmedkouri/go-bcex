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