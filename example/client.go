package main

import (
	"log"
	"os"

	"github.com/hmedkouri/go-bcex"
)

func main() {
	apiKey := os.Getenv("BCEX_API_KEY")
	secretKey := os.Getenv("BCEX_API_SECRET")
	client := bcex.New(apiKey, secretKey)

	symbol, err := client.Rest.GetSymbol("BTC-USD")
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("GetSymbol : %#v\n", symbol)

	ticker, err := client.Rest.GetTicker("BTC-USD")
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("GetTicker : %#v\n", ticker)

	tickers, err := client.Rest.GetAllTicker()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("GetAllTicker : %v\n", tickers)

	err = client.Ws.Start(true)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Ws.Stop()

	// subscribe to heartbeats
	err = client.Ws.SubscribeHeartbeat()
	if err != nil {
		log.Fatalln(err)
	}

	// subscribe to symbols
	err = client.Ws.SubscribeToSymbols()
	if err != nil {
		log.Fatalln(err)
	}

	for {
		select {
		case balancesMsg := <-client.Ws.Balances():
			log.Printf("received balances %+v", balancesMsg)
		case l3Msg := <-client.Ws.L3Quotes():
			log.Printf("received l3 quote %+v", l3Msg)
		case l2Msg := <-client.Ws.L2Quotes():
			log.Printf("received l2 quote %+v", l2Msg)
		case heartbeatMsg := <-client.Ws.Heartbeats():
			log.Printf("received heartbeat %+v", heartbeatMsg)
		case pricesMsg := <-client.Ws.Prices():
			log.Printf("received price %+v", pricesMsg)
		case tradesMsg := <-client.Ws.Trades():
			log.Printf("received trades %+v", tradesMsg)
		case tickerMsg := <-client.Ws.Ticker():
			log.Printf("received ticker %+v", tickerMsg)
		case tradingMsg := <-client.Ws.Trading():
			if tradingMsg.IsSnapshot() {
				log.Printf("received trading snapshot %+v", tradingMsg)
			} else if tradingMsg.IsUpdate() {
				log.Printf("received trading update %+v", tradingMsg)
			} else if tradingMsg.IsReject() {
				log.Printf("received trading reject %+v", tradingMsg)
			} else {
				log.Printf("received unknown trading msg %+v", tradingMsg)
			}
		case symbolsMsg := <-client.Ws.Symbols():
			log.Printf("received symbols %+v", symbolsMsg)
		}
	}

}
