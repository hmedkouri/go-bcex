package ws_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/hmedkouri/go-bcex"
	"github.com/hmedkouri/go-bcex/ws"
)

var (
	apiKey    = os.Getenv("BCEX_API_KEY")
	apiSecret = os.Getenv("BCEX_API_SECRET")
	bc                  *bcex.Bcex = bcex.New(apiKey, apiSecret)
)

func TestWs(t *testing.T) {
	err := bc.Ws.Start(true)
	
	if err != nil {
		log.Fatal(err)
	}
	defer bc.Ws.Stop()

	go listenForUpdates(bc.Ws)

	subscribeToHeartbeats(bc.Ws)
	subscribeToSymbols(bc.Ws)
	subscribeToL3(bc.Ws, ws.BTCUSD)
	subscribeToL2(bc.Ws, ws.ETHBTC)

	time.Sleep(5 * time.Second)
}

func listenForUpdates(wsClient *ws.WebSocketClient) {
	for {
		select {
		case balancesMsg := <-wsClient.Balances():
			log.Printf("received balances %+v", balancesMsg)
		case l3Msg := <-wsClient.L3Quotes():
			log.Printf("received l3 quote %+v", l3Msg)
		case l2Msg := <-wsClient.L2Quotes():
			log.Printf("received l2 quote %+v", l2Msg)
		case heartbeatMsg := <-wsClient.Heartbeats():
			log.Printf("received heartbeat %+v", heartbeatMsg)
		case pricesMsg := <-wsClient.Prices():
			log.Printf("received price %+v", pricesMsg)
		case tradesMsg := <-wsClient.Trades():
			log.Printf("received trades %+v", tradesMsg)
		case tickerMsg := <-wsClient.Ticker():
			log.Printf("received ticker %+v", tickerMsg)
		case tradingMsg := <-wsClient.Trading():
			if tradingMsg.IsSnapshot() {
				log.Printf("received trading snapshot %+v", tradingMsg)
			} else if tradingMsg.IsUpdate() {
				log.Printf("received trading update %+v", tradingMsg)
			} else if tradingMsg.IsReject() {
				log.Printf("received trading reject %+v", tradingMsg)
			} else {
				log.Printf("received unknown trading msg %+v", tradingMsg)
			}
		case symbolsMsg := <-wsClient.Symbols():
			log.Printf("received symbols %+v", symbolsMsg)
		}
	}
}

func subscribeToHeartbeats(wsClient *ws.WebSocketClient) {
	err := wsClient.SubscribeHeartbeat()

	if err != nil {
		log.Fatal(err)
	}
}

func subscribeToSymbols(wsClient *ws.WebSocketClient) {
	err := wsClient.SubscribeToSymbols()
	if err != nil {
		log.Fatal(err)
	}
}

func subscribeToL3(wsClient *ws.WebSocketClient, symbol ws.Symbol) {
	err := wsClient.SubscribeToL3(symbol)

	if err != nil {
		log.Fatal(err)
	}
}

func subscribeToL2(wsClient *ws.WebSocketClient, symbol ws.Symbol) {
	err := wsClient.SubscribeToL2(symbol)

	if err != nil {
		log.Fatal(err)
	}
}
