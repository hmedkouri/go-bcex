package rest_test

import (
	"os"
	"testing"

	bcex "github.com/hmedkouri/go-bcex"
	"github.com/hmedkouri/go-bcex/rest"

	"github.com/stretchr/testify/require"
)

var (
	apiKey                           = os.Getenv("BCEX_API_KEY")
	apiSecret                        = os.Getenv("BCEX_API_SECRET")
	bc                  *bcex.Client = bcex.New(apiKey, apiSecret)
	defaultErrorMessage string       = "There should be no error"
)

func TestGetSymbols(t *testing.T) {
	symbols, err := bc.Rest.GetSymbols()
	t.Logf("GetSymbols : %#v\n", symbols)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetSymbol(t *testing.T) {
	symbol, err := bc.Rest.GetSymbol("BTC-USD")
	t.Logf("GetSymbol : %#v\n", symbol)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetTicker(t *testing.T) {
	ticker, err := bc.Rest.GetTicker("BTC-USD")
	t.Logf("GetTicker : %#v\n", ticker)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetAllTicker(t *testing.T) {
	tickers, err := bc.Rest.GetAllTicker()
	t.Logf("GetAllTicker : %v\n", tickers)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetL2Orderbook(t *testing.T) {
	orderbook, err := bc.Rest.GetL2Orderbook("BTC-USD")
	t.Logf("GetL2Orderbook : %#v\n", orderbook)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetL3Orderbook(t *testing.T) {
	orderbook, err := bc.Rest.GetL3Orderbook("BTC-USD")
	t.Logf("GetL3Orderbook : %#v\n", orderbook)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetFees(t *testing.T) {
	fees, err := bc.Rest.GetFees()
	t.Logf("GetFees : %#v\n", fees)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetBalances(t *testing.T) {
	balances, err := bc.Rest.GetBalances()
	t.Logf("GetBalances : %#v\n", balances)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetTrades(t *testing.T) {
	options := rest.GetTradesOpts{
		Symbol: "BTC-USD",
	}
	trades, err := bc.Rest.GetTrades(&options)
	t.Logf("GetTrades : %#v\n", trades)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetOrders(t *testing.T) {
	options := rest.GetOrdersOpts{
		Symbol: "BTC-USD",
	}
	orders, err := bc.Rest.GetOrders(&options)
	t.Logf("GetOrders : %#v\n", orders)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetFills(t *testing.T) {
	options := rest.GetFillsOpts{
		Symbol: "BTC-USD",
	}
	fills, err := bc.Rest.GetFills(&options)
	t.Logf("GetFills : %#v\n", fills)
	require.NoError(t, err, defaultErrorMessage)
}

func TestDeleteAllOrders(t *testing.T) {
	options := rest.DeleteAllOrdersOpts{
		Symbol: "BTC-USD",
	}
	err := bc.Rest.DeleteAllOrders(&options)
	require.NoError(t, err, defaultErrorMessage)
}
