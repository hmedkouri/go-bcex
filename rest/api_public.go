package rest

import (
	"encoding/json"
	"strings"
)

// GetSymbols is used to get the open and available trading markets along with other meta data.
func (client *Client) GetSymbols() (symbols []Symbol, err error) {
	r, err := client.do("GET", "symbols", nil, nil, false)
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
func (client *Client) GetSymbol(market string) (symbol Symbol, err error) {
	r, err := client.do("GET", "symbols/"+strings.ToUpper(market), nil, nil, false)
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
func (client *Client) GetAllTicker() (tickers Tickers, err error) {
	r, err := client.do("GET", "tickers", nil, nil, false)
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
func (client *Client) GetTicker(market string) (ticker Ticker, err error) {
	r, err := client.do("GET", "tickers/"+strings.ToUpper(market), nil, nil, false)
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
func (client *Client) GetL2Orderbook(market string) (orderbook OrderBook, err error) {
	r, err := client.do("GET", "l2/"+strings.ToUpper(market), nil, nil, false)
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
func (client *Client) GetL3Orderbook(market string) (orderbook OrderBook, err error) {
	r, err := client.do("GET", "l3/"+strings.ToUpper(market), nil, nil, false)
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