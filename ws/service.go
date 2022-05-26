package ws

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// UseTestnet switch all the API endpoints from production to the testnet
var UseTestnet = false

// Endpoints
const (
	WsEndpoint     = "wss://ws.blockchain.info/mercury-gateway/v1/ws"
	WsTestEndpoint = "wss://ws.staging.blockchain.info/mercury-gateway/v1/ws"
)

// WsHandler handle raw websocket message
type WsHandler func(message []byte)

// ErrHandler handles errors
type ErrHandler func(err error)

// WsL2MsgHandler handle websocket l2 event
type WsL2MsgHandler func(event *L2Msg)

// WsL2MsgHandler handle websocket l3 event
type WsL3MsgHandler func(event *L3Msg)

// WsTickerMsgHandler handle websocket ticker event
type WsTickerMsgHandler func(event *TickerMsg)

// WsPriceMsgHandler handle websocket price event
type WsPriceMsgHandler func(event *PricesMsg)

// WsL2Serve serve websocket l2 handler with a symbol
func (ws *WebSocketClient) WsL2Serve(symbol string, handler WsL2MsgHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	cfg := ws.config
	cfg.IsSecure = false
	wsHandler := func(message []byte) {
		var m L2Msg
		err := json.Unmarshal(message, &m)
		if err != nil {
			errHandler(err)
			return
		}		
		handler(&m)
	}

	r := quoteSubscriptionRequest{
		Action:  actionSubscribe,
		Channel: l2Channel,
		Symbol:  Symbol(symbol),
	}

	return wsServe(cfg, r, wsHandler, errHandler)
}

// WsL3Serve serve websocket l3 handler with a symbol
func (ws *WebSocketClient) WsL3Serve(symbol string, handler WsL3MsgHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	cfg := ws.config
	cfg.IsSecure = false
	wsHandler := func(message []byte) {
		var m L3Msg
		err := json.Unmarshal(message, &m)
		if err != nil {
			errHandler(err)
			return
		}		
		handler(&m)
	}

	r := quoteSubscriptionRequest{
		Action:  actionSubscribe,
		Channel: l3Channel,
		Symbol:  Symbol(symbol),
	}

	return wsServe(cfg, r, wsHandler, errHandler)
}

// WsTickerServe serve websocket ticker handler with a symbol
func (ws *WebSocketClient) WsTickerServe(symbol string, handler WsTickerMsgHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	cfg := ws.config
	cfg.IsSecure = false
	wsHandler := func(message []byte) {
		var m TickerMsg
		err := json.Unmarshal(message, &m)
		if err != nil {
			errHandler(err)
			return
		}		
		handler(&m)
	}

	r := tickerSubscriptionRequest{
		Action:  actionSubscribe,
		Channel: tickerChannel,
		Symbol:  Symbol(symbol),
	}

	return wsServe(cfg, r, wsHandler, errHandler)
}

// WsPriceServe serve websocket ticker handler with a symbol
func (ws *WebSocketClient) WsPriceServe(symbol string, granularity Granularity, handler WsPriceMsgHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	cfg := ws.config
	cfg.IsSecure = false
	wsHandler := func(message []byte) {
		var m PricesMsg
		err := json.Unmarshal(message, &m)
		if err != nil {
			errHandler(err)
			return
		}		
		handler(&m)
	}

	r := pricesSubscriptionRequest{
		Action:      actionSubscribe,
		Channel:     pricesChannel,
		Symbol:      Symbol(symbol),
		Granularity: granularity,
	}

	return wsServe(cfg, r, wsHandler, errHandler)
}

var wsServe = func(cfg Configuration, request interface{}, handler WsHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	var d = websocket.Dialer{
		Subprotocols:    []string{"p1", "p2"},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		Proxy:           http.ProxyFromEnvironment,
	}

	d.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c, _, err := d.Dial(cfg.Host, WsHeaders)
	if err != nil {
		return nil, nil, err
	}
	c.SetReadLimit(655350)
	if cfg.IsSecure {
		//WsHeaders.Add("Cookie", cookie[ws.config.Env]+ws.config.ApiKey)
		connectMsg, _ := json.Marshal(&privateConnect{
			Channel: "auth",
			Token:   cfg.ApiKey,
			Action:  "subscribe",
		})

		// Send auth message
		err = c.WriteMessage(websocket.TextMessage, connectMsg)
		if err != nil {
			return nil, nil, err
		}
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, nil, err
	}

	// Send subscription request
	err = c.WriteMessage(websocket.TextMessage, requestBytes)
	if err != nil {
		return nil, nil, err
	}

	doneC = make(chan struct{})
	stopC = make(chan struct{})
	go func() {
		// This function will exit either on error from
		// websocket.Conn.ReadMessage or when the stopC channel is
		// closed by the client.
		defer close(doneC)
		if cfg.Keepalive {
			keepAlive(c, cfg.Timeout)
		}
		// Wait for the stopC channel to be closed.  We do that in a
		// separate goroutine because ReadMessage is a blocking
		// operation.
		silent := false
		go func() {
			select {
			case <-stopC:
				silent = true
			case <-doneC:
			}
			c.Close()
		}()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				if !silent {
					errHandler(err)
				}
				return
			}
			handler(message)
		}
	}()
	return
}

func keepAlive(c *websocket.Conn, timeout time.Duration) {
	ticker := time.NewTicker(timeout)

	lastResponse := time.Now()
	c.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})

	go func() {
		defer ticker.Stop()
		for {
			deadline := time.Now().Add(10 * time.Second)
			err := c.WriteControl(websocket.PingMessage, []byte{}, deadline)
			if err != nil {
				return
			}
			<-ticker.C
			if time.Since(lastResponse) > timeout {
				c.Close()
				return
			}
		}
	}()
}