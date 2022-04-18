package rest

import (
	"encoding/json"
	"fmt"
	"strconv"
)

/*
CreateOrder Add an order
 * @param baseOrder Trade
@return OrderSummary
*/
func (client *Client) CreateOrder(requestOrder BaseOrder) (order OrderSummary, err error) {
	payload := make(map[string]string, 6)

	payload["clOrdId"] = requestOrder.ClOrdId
	payload["ordType"] = string(requestOrder.OrdType)
	payload["symbol"] = requestOrder.Symbol
	payload["side"] = parameterToString(requestOrder.Side, "")
	payload["orderQty"] = fmt.Sprintf("%.8f", requestOrder.OrderQty)
	payload["price"] = fmt.Sprintf("%.8f", requestOrder.Price)

	method := "POST"
	resource := "orders"

	r, err := client.do(method, resource, payload, true)
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

	err = json.Unmarshal(r, &order)
	return
}

/*
DeleteAllOrders Delete all open orders (of a symbol, if specified)
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param optional nil or *DeleteAllOrdersOpts - Optional Parameters:
 * @param "Symbol" (optional.String) - 
*/
func (client *Client) DeleteAllOrders(options *DeleteAllOrdersOpts) (err error) {
	var payload map[string]string
	if (options != nil) {
		payload = options.parse()
	}
	r, err := client.do("DELETE", "orders", payload, true)
	_ = r
	if err != nil {
		return
	}
	return
}

// GetFees is used to retrieve the fees from your account
func (client *Client) GetFees() (fees Fees, err error) {
	r, err := client.do("GET", "fees", nil, true)
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
func (client *Client) GetBalances() (balances BalanceMap, err error) {
	r, err := client.do("GET", "accounts", nil, true)
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
func (client *Client) GetTrades(currencyPair string) (trades []Trade, err error) {
	payload := make(map[string]string)
	if currencyPair != "all" {
		payload["symbol"] = currencyPair
	}
	r, err := client.do("GET", "trades", payload, true)
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
 * @param nil or *GetOrdersOpts - Optional Parameters:
 * @param "Symbol" (string) -  Only return results for this symbol
 * @param "From" (int64) -  Epoch timestamp in ms
 * @param "To" (int64) -  Epoch timestamp in ms
 * @param "Status" (interface of OrderStatus) -  Order Status
 * @param "Limit" (int32) -  Maximum amount of results to return in a single call. If omitted, 100 results are returned by default. 
@return []OrderSummary
*/
func (client *Client) GetOrders(options *GetOrdersOpts) (orders []OrderSummary, err error) {
	var payload map[string]string
	if (options != nil) {
		payload = options.parse()
	}
	r, err := client.do("GET", "orders", payload, true)
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

/*
GetFills Get a list of filled orders
Returns filled orders, including partial fills. Returns at most 100 results, use timestamp to paginate for further results
 * @param nil or *GetFillsOpts - Parameters:
 * @param "Symbol" (string) -  Only return results for this symbol
 * @param "From" (int64) -  Epoch timestamp in ms
 * @param "To" (optional.Int64) -  Epoch timestamp in ms
 * @param "Limit" (optional.Int32) -  Maximum amount of results to return in a single call. If omitted, 100 results are returned by default. 
@return []OrderSummary
*/
func (client *Client) GetFills(options *GetFillsOpts) (fills []OrderSummary, err error) {
	var payload map[string]string
	if (options != nil) {
		payload = options.parse()
	}
	r, err := client.do("GET", "fills", payload, true)
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
	err = json.Unmarshal(r, &fills)
	return
}

/*
GetOrderById Get a specific order
 * @param orderId Order ID
@return OrderSummary
*/
func (client *Client) GetOrderById(orderId int64) (order OrderSummary, err error) {
	r, err := client.do("GET", "orders/"+strconv.Itoa(int(orderId)), nil, true)
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
	err = json.Unmarshal(r, &order)
	return
}