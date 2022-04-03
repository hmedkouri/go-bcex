package bcex_test

import (
	bcex "hmedkouri/go-bcex"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	apiKey    = os.Getenv("BCEX_API_KEY")
	apiSecret = os.Getenv("BCEX_API_SECRET")
)

var (
	bc                  *bcex.Bcex = bcex.New(apiKey, apiSecret)
	defaultErrorMessage string     = "There should be no error"
)

func TestGetSymbols(t *testing.T) {
	symbols, err := bc.GetSymbols()
	t.Logf("GetSymbols : %#v\n", symbols)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetSymbol(t *testing.T) {
	symbol, err := bc.GetSymbol("BTC-USD")
	t.Logf("GetSymbol : %#v\n", symbol)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetTicker(t *testing.T) {
	ticker, err := bc.GetTicker("BTC-USD")
	t.Logf("GetTicker : %#v\n", ticker)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetAllTicker(t *testing.T) {
	tickers, err := bc.GetAllTicker()
	t.Logf("GetAllTicker : %v\n", tickers)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetL2Orderbook(t *testing.T) {
	orderbook, err := bc.GetL2Orderbook("BTC-USD")
	t.Logf("GetL2Orderbook : %#v\n", orderbook)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetL3Orderbook(t *testing.T) {
	orderbook, err := bc.GetL3Orderbook("BTC-USD")
	t.Logf("GetL3Orderbook : %#v\n", orderbook)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetFees(t *testing.T) {
	fees, err := bc.GetFees()
	t.Logf("GetFees : %#v\n", fees)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetBalances(t *testing.T) {
	balances, err := bc.GetBalances()
	t.Logf("GetBalances : %#v\n", balances)
	require.NoError(t, err, defaultErrorMessage)
}