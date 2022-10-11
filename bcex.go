// Package bcex is a golang Api Blockchain.com V4 Rest API and V1 Websockets API
//
// https://blockchain.info/api
package bcex

import (
	"time"

	"github.com/hmedkouri/go-bcex/rest"
	"github.com/hmedkouri/go-bcex/ws"
)

// Client is the main go-bcex handle to the rest and ws clients
type Client struct {
	Rest        *rest.Client
	Ws          *ws.WebSocketClient
	websocketOn bool
}

// New returns an instantiated go-bcex Client struct
func New(apiKey, apiSecret string) *Client {
	api := rest.NewClient(apiKey, apiSecret)
	ws := ws.NewWebSocketClient(ws.Configuration{
		Host:      ws.WsEndpoint,
		ApiKey:    apiSecret,
		Timeout:   60 * time.Second,
		Keepalive: true,
		Env:       ws.PROD,
	})
	return &Client{api, ws, true}
}
