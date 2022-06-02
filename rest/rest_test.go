package rest_test

import (
	"os"
	"testing"

	bcex "github.com/hmedkouri/go-bcex"
	"github.com/hmedkouri/go-bcex/rest"

	"github.com/stretchr/testify/require"
)

var (
	apiKey    = os.Getenv("BCEX_API_KEY")
	apiSecret = os.Getenv("BCEX_API_SECRET")
	bc                  *bcex.Bcex = bcex.New(apiKey, apiSecret)
	defaultErrorMessage string     = "There should be no error"
)

	func TestGetSymbols(t *testing.T) {
	symbols, err := bc.Api.GetSymbols()
	t.Logf("GetSymbols : %#v\n", symbols)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetSymbol(t *testing.T) {
	symbol, err := bc.Api.GetSymbol("BTC-USD")
	t.Logf("GetSymbol : %#v\n", symbol)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetTicker(t *testing.T) {
	ticker, err := bc.Api.GetTicker("BTC-USD")
	t.Logf("GetTicker : %#v\n", ticker)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetAllTicker(t *testing.T) {
	tickers, err := bc.Api.GetAllTicker()
	t.Logf("GetAllTicker : %v\n", tickers)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetL2Orderbook(t *testing.T) {
	orderbook, err := bc.Api.GetL2Orderbook("BTC-USD")
	t.Logf("GetL2Orderbook : %#v\n", orderbook)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetL3Orderbook(t *testing.T) {
	orderbook, err := bc.Api.GetL3Orderbook("BTC-USD")
	t.Logf("GetL3Orderbook : %#v\n", orderbook)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetFees(t *testing.T) {
	fees, err := bc.Api.GetFees()
	t.Logf("GetFees : %#v\n", fees)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetBalances(t *testing.T) {
	balances, err := bc.Api.GetBalances()
	t.Logf("GetBalances : %#v\n", balances)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetTrades(t *testing.T) {
	trades, err := bc.Api.GetTrades("BTC-USD")
	t.Logf("GetTrades : %#v\n", trades)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetOrders(t *testing.T) {
	options := rest.GetOrdersOpts {
		Symbol: "BTC-USD",
	}
	orders, err := bc.Api.GetOrders(&options)
	t.Logf("GetOrders : %#v\n", orders)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetFills(t *testing.T) {
	options := rest.GetFillsOpts {
		Symbol: "BTC-USD",
	}
	fills, err := bc.Api.GetFills(&options)
	t.Logf("GetFills : %#v\n", fills)
	require.NoError(t, err, defaultErrorMessage)
}

func TestDeleteAllOrders(t *testing.T) {
	options := rest.DeleteAllOrdersOpts {
		Symbol: "BTC-USD",
	}
	err := bc.Api.DeleteAllOrders(&options)
	require.NoError(t, err, defaultErrorMessage)
}