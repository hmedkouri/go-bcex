package bcex

import (
	"time"

	"github.com/hmedkouri/go-bcex/rest"
	"github.com/hmedkouri/go-bcex/ws"
)

type Bcex struct {
	Api *rest.Client
	Ws  *ws.WebSocketClient
	websocketOn      bool
}

// New returns an instantiated HitBTC struct
func New(apiKey, apiSecret string) *Bcex {
	api := rest.NewClient(apiKey, apiSecret)
	ws := ws.NewWebSocketClient(ws.Configuration{
		Host: ws.WsEndpoint,
		ApiKey: apiSecret,
		Timeout: 5 * time.Second,
		Env:     ws.PROD,
	})
	return &Bcex{api, ws, true}
}