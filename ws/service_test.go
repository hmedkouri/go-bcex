package ws_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hmedkouri/go-bcex/ws"
)

func TestService(t *testing.T) {
	bc.Ws.WsL2Serve("BTC-USD", func(event *ws.L2Msg) {
		fmt.Printf("L2MsgA : %s, %d\n", event.Event, event.Seqnum)
	}, func(err error) {
		log.Fatal(err)
	})
	bc.Ws.WsTickerServe("BTC-USD", func(event *ws.TickerMsg) {
		fmt.Printf("TickerMsg : %d\n", event.Seqnum)
	}, func(err error) {
		log.Fatal(err)
	})
	bc.Ws.WsPriceServe("BTC-USD", ws.Granularity21600, func(event *ws.PricesMsg) {
		fmt.Printf("PricesMsg : %d\n", event.Seqnum)
	}, func(err error) {
		log.Fatal(err)
	})

	time.Sleep(5 * time.Second)
}

func TestServiceCombined(t *testing.T) {
	symbols := []string{"BTC-USD", "ETH-USD"}
	bc.Ws.WsL2ServeCombined(symbols, func(event *ws.L2Msg) {
		fmt.Printf("L2MsgA : %s, %d, %s\n", event.Event, event.Seqnum, event.Symbol)
	}, func(err error) {
		log.Fatal(err)
	})

	time.Sleep(5 * time.Second)
}
