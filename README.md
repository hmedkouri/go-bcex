bcex
====
[![Go](https://github.com/hmedkouri/go-bcex/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/hmedkouri/go-bcex/actions/workflows/go.yml)
[![CodeQL](https://github.com/hmedkouri/go-bcex/actions/workflows/codeql.yml/badge.svg)](https://github.com/hmedkouri/go-bcex/actions/workflows/codeql.yml)

*NOTICE:*
> WORK IN PROGRESS USE AT YOUR OWN RISK


Blockchain Exchange Golang API

A complete golang wrapper for the [Blockchain.com](https://exchange.blockchain.com) Exchange Websockets V1 and Rest V3 API.

For more info about Blockchain.com API [read here](https://blockchain.info/api).

Installation
-----------------

```bash
go get github.com/hmedkouri/go-bcex
```

Usage
-----------

```go
package main

import (
	"log"

	"github.com/hmedkouri/go-bcex"
)

func main() {
    apiKey := "YOUR-API-KEY"
	secretKey := "YOUR-SECRET-KEY"
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
```

Supporting APIs
---------------

* [Rest](https://api.blockchain.com/v3/)

* [Ws](https://exchange.blockchain.com/api/#websocket-api)