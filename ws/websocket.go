package ws

import (
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/json-iterator/go"
)

type Env string

var (
	PROD    Env = "PROD"
	STAGING Env = "STAGING"
)

type Configuration struct {
	Host      string
	ApiKey    string
	Env       Env
	Timeout   time.Duration
	Keepalive bool
	IsSecure  bool
}

const (
	BTCUSD Symbol = "BTC-USD"
	BTCEUR Symbol = "BTC-EUR"
	LTCUSD Symbol = "LTC-USD"
	LTCBTC Symbol = "LTC-BTC"
	ETHUSD Symbol = "ETH-USD"
	ETHBTC Symbol = "ETH-BTC"
	ETHEUR Symbol = "ETH-EUR"
	ETHGBP Symbol = "ETH-GBP"
)

const (
	eventSubscribed   eventType = "subscribed"
	eventRejected     eventType = "rejected"
	eventSnapshot     eventType = "snapshot"
	eventUnsubscribed eventType = "unsubscribed"
	eventUpdate       eventType = "updated"

	actionSubscribe actionType = "subscribe"
	newOrderSingle  actionType = "NewOrderSingle"
	cancelOrder     actionType = "CancelOrderRequest"
	bulkCancel      actionType = "BulkCancelOrderRequest"

	heartbeatChannel channel = "heartbeat"
	symbolsChannel   channel = "symbols"
	l3Channel        channel = "l3"
	l2Channel        channel = "l2"
	pricesChannel    channel = "prices"
	tickerChannel    channel = "ticker"
	tradesChannel    channel = "trades"
	balancesChannel  channel = "balances"
	authChannel      channel = "auth"
	tradingChannel   channel = "trading"

	Granularity60    Granularity = 60
	Granularity300   Granularity = 300
	Granularity900   Granularity = 900
	Granularity3600  Granularity = 3600
	Granularity21600 Granularity = 21600
	Granularity86400 Granularity = 86400

	GTC TimeInForce = "GTC" // Good Till Cancel. The order will rest on the order book until it is cancelled or filled
	GTD TimeInForce = "GTD" // Good Till Date. The order will reset on the order book until it is cancelled, filled, or expired
	FOK TimeInForce = "FOK" // Fill or Kill. The order is either completely filled or cancelled. No Partial Fills are permitted.
	IOC TimeInForce = "IOC" // Immediate or Cancel. The order is either a) completely filled, b) partially filled and the remaining quantity canceled, or c) the order is canceled.

	LIMIT      OrderType = "limit"     // order which has a price limit
	MARKET     OrderType = "market"    // order that will match at any price available in the market, starting from the best prices and filling up to the available balance
	STOP       OrderType = "stop"      // order which has a stop/trigger price, and when that price is reached, it triggers a market order
	STOP_LIMIT OrderType = "stopLimit" // order which has a stop price and limit price, and when the stop price is reached, it triggers a limit order at the limit price

	ORDER_STATUS_PENDING   OrderStatus = "pending"   // is pending acceptance. Only applicable to stop and stop-limit orders
	ORDER_STATUS_OPEN      OrderStatus = "open"      // has been accepted
	ORDER_STATUS_REJECTED  OrderStatus = "rejected"  // has been rejected	Limit and market orders can get rejected if you have no balance to fill the order even partially.
	ORDER_STATUS_CANCELLED OrderStatus = "cancelled" //  has been cancelled	A market order might get in state cancelled if you don’t have enough balance to fill it at market price. Both market orders and limit orders with IOC can have ordStatus ‘cancelled’ if there is no market for them, even without the user requesting the cancellation.
	ORDER_STATUS_FILLED    OrderStatus = "filled"    // has been filled	A limit order get in state cancelled after the user requested a cancellation.
	ORDER_STATUS_PARTIAL   OrderStatus = "partial"   // has been partially filled
	ORDER_STATUS_EXPIRED   OrderStatus = "expired"   // has been expired

	BUY  OrderSide = "buy"
	SELL OrderSide = "sell"

	ALO ExecInst = "ALO"

	ExecutionReport     MsgType = "8"
	OrderCancelRejected MsgType = "9"

	EXEC_TYPE_NEW          ExecType = "0"
	EXEC_TYPE_CANCELLED    ExecType = "4"
	EXEC_TYPE_EXPIRED      ExecType = "C"
	EXEC_TYPE_REJECTED     ExecType = "8"
	EXEC_TYPE_PARTIAL_FILL ExecType = "F"
	EXEC_TYPE_PENDING      ExecType = "A"
	EXEC_TYPE_TRADE_BREAK  ExecType = "H"
	EXEC_TYPE_ORDER_STATUS ExecType = "I"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary

	WsHeaders = http.Header{
		"Origin": {"https://exchange.blockchain.com"},
	}

	cookie = map[Env]string{
		PROD:    "auth_token=",
		STAGING: "staging_auth_token=",
	}

	PingFrequency = time.Second * 5
)

type WebSocketClient struct {
	connMu         *sync.Mutex
	conn           *websocket.Conn
	config         Configuration
	heartbeatTimer *time.Timer
	quit           chan struct{}

	// anonymous channels
	chTrades  chan TradesMsg
	chSymbols chan SymbolMsg
	chTicker  chan TickerMsg
	chPrices  chan PricesMsg
	chL2      chan L2Msg
	chL3      chan L3Msg

	// authenticated channels
	chBalances                  chan BalancesSnapshot
	chTrading                   chan TradingMsg
	chHeartbeat                 chan HeartbeatMsg
	subscriptionResponseChannel chan SubscriptionError

	errorsChan chan error

	mutex *sync.RWMutex
}

var (
	ErrAlreadySubscribed = errors.New("already subscribed")
	ErrInvalidRequest    = errors.New("invalid request")
)

type SubscriptionError struct {
	SubscriptionName string
	ErrorString      string
}

func NewWebSocketClient(configuration Configuration) *WebSocketClient {
	return &WebSocketClient{
		config:                      configuration,
		errorsChan:                  make(chan error, 10),
		mutex:                       &sync.RWMutex{},
		connMu:                      &sync.Mutex{},
		heartbeatTimer:              time.NewTimer(PingFrequency),
		subscriptionResponseChannel: make(chan SubscriptionError),
		chHeartbeat:                 make(chan HeartbeatMsg),
		chSymbols:                   make(chan SymbolMsg),
		chL3:                        make(chan L3Msg),
		chL2:                        make(chan L2Msg),
		chPrices:                    make(chan PricesMsg),
		chTicker:                    make(chan TickerMsg),
		chTrades:                    make(chan TradesMsg),
		chBalances:                  make(chan BalancesSnapshot),
		chTrading:                   make(chan TradingMsg),
	}
}

func (ws *WebSocketClient) Heartbeats() chan HeartbeatMsg {
	return ws.chHeartbeat
}

func (ws *WebSocketClient) Symbols() chan SymbolMsg {
	return ws.chSymbols
}

func (ws *WebSocketClient) L3Quotes() chan L3Msg {
	return ws.chL3
}

func (ws *WebSocketClient) L2Quotes() chan L2Msg {
	return ws.chL2
}

func (ws *WebSocketClient) Prices() chan PricesMsg {
	return ws.chPrices
}

func (ws *WebSocketClient) Ticker() chan TickerMsg {
	return ws.chTicker
}

func (ws *WebSocketClient) Trades() chan TradesMsg {
	return ws.chTrades
}

func (ws *WebSocketClient) Balances() chan BalancesSnapshot {
	return ws.chBalances
}

func (ws *WebSocketClient) Trading() chan TradingMsg {
	return ws.chTrading
}

type actionType string
type eventType string
type channel string

func (a actionType) String() string {
	return string(a)
}

func (e eventType) String() string {
	return string(e)
}

func (c channel) String() string {
	return string(c)
}

func (ws *WebSocketClient) Errors() chan error {
	return ws.errorsChan
}

type privateConnect struct {
	Token   string `json:"token"`
	Action  string `json:"action"`
	Channel string `json:"channel"`
}

func (ws *WebSocketClient) Start(authenticate bool) error {
	ws.quit = make(chan struct{})
	var d = websocket.Dialer{
		Subprotocols:    []string{"p1", "p2"},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		Proxy:           http.ProxyFromEnvironment,
	}

	d.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	conn, _, err := d.Dial(ws.config.Host, WsHeaders)
	if err != nil {
		return err
	}
	log.Println("Connected")
	ws.conn = conn
	if authenticate {
		//WsHeaders.Add("Cookie", cookie[ws.config.Env]+ws.config.ApiKey)
		connectMsg, _ := json.Marshal(&privateConnect{
			Channel: "auth",
			Token:   ws.config.ApiKey,
			Action:  "subscribe",
		})

		// Send auth message
		err = ws.conn.WriteMessage(websocket.TextMessage, connectMsg)
		if err != nil {
			return err
		}
	}
	go ws.listenForUpdates(ws.quit)

	if authenticate {
		return ws.expectSubscriptionResponse(authChannel)
	}
	return nil
}

func (ws *WebSocketClient) Stop() error {
	if ws.quit != nil {
		close(ws.quit)
		ws.quit = nil

		ws.connMu.Lock()
		defer ws.connMu.Unlock()

		return ws.conn.Close()
	}
	return nil
}

func (ws *WebSocketClient) SubscribeToSymbols() error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	subscribeRequest := symbolsSubscriptionRequest{
		Action:  actionSubscribe,
		Channel: symbolsChannel,
	}

	subscribeRequestBytes, err := json.Marshal(subscribeRequest)
	if err != nil {
		return err
	}

	ws.connMu.Lock()
	err = ws.conn.WriteMessage(websocket.TextMessage, subscribeRequestBytes)
	if err != nil {
		ws.connMu.Unlock()
		return err
	}
	ws.connMu.Unlock()

	err = ws.expectSubscriptionResponse(symbolsChannel)
	return err
}

func (ws *WebSocketClient) SubscribeToL3(symbol Symbol) error {
	return ws.quoteSubscription(symbol, l3Channel)
}

func (ws *WebSocketClient) quoteSubscription(symbol Symbol, level channel) error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	subscribeRequest := quoteSubscriptionRequest{
		Action:  actionSubscribe,
		Channel: level,
		Symbol:  symbol,
	}

	subscribeRequestBytes, err := json.Marshal(subscribeRequest)
	if err != nil {
		return err
	}

	ws.connMu.Lock()
	err = ws.conn.WriteMessage(websocket.TextMessage, subscribeRequestBytes)
	if err != nil {
		ws.connMu.Unlock()
		return err
	}
	ws.connMu.Unlock()

	err = ws.expectSubscriptionResponse(level)
	return err
}

func (ws *WebSocketClient) SubscribeToPrices(symbol Symbol, granularity Granularity) error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	subscribeRequest := pricesSubscriptionRequest{
		Action:      actionSubscribe,
		Channel:     pricesChannel,
		Symbol:      symbol,
		Granularity: granularity,
	}

	subscribeRequestBytes, err := json.Marshal(subscribeRequest)
	if err != nil {
		return err
	}

	ws.connMu.Lock()
	err = ws.conn.WriteMessage(websocket.TextMessage, subscribeRequestBytes)
	if err != nil {
		ws.connMu.Unlock()
		return err
	}
	ws.connMu.Unlock()

	err = ws.expectSubscriptionResponse(pricesChannel)
	return err
}

func (ws *WebSocketClient) SubscribeToBalances() error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	subscribeRequest := balancesSubscriptionRequest{
		Action:  actionSubscribe,
		Channel: balancesChannel,
	}

	subscribeRequestBytes, err := json.Marshal(subscribeRequest)
	if err != nil {
		return err
	}

	ws.connMu.Lock()
	err = ws.conn.WriteMessage(websocket.TextMessage, subscribeRequestBytes)
	if err != nil {
		ws.connMu.Unlock()
		return err
	}
	ws.connMu.Unlock()

	err = ws.expectSubscriptionResponse(balancesChannel)
	return err
}

func (ws *WebSocketClient) SubscribeToTrading() error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	subscribeRequest := tradingSubscriptionRequest{
		Action:  actionSubscribe,
		Channel: tradingChannel,
	}

	subscribeRequestBytes, err := json.Marshal(subscribeRequest)
	if err != nil {
		return err
	}

	ws.connMu.Lock()
	err = ws.conn.WriteMessage(websocket.TextMessage, subscribeRequestBytes)
	if err != nil {
		ws.connMu.Unlock()
		return err
	}
	ws.connMu.Unlock()

	err = ws.expectSubscriptionResponse(tradingChannel)
	return err
}

func (ws *WebSocketClient) Authenticate(token string) error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	subscribeRequest := authSubscriptionRequest{
		Action:  actionSubscribe,
		Channel: authChannel,
		Token:   token,
	}

	subscribeRequestBytes, err := json.Marshal(subscribeRequest)
	if err != nil {
		return err
	}

	ws.connMu.Lock()
	err = ws.conn.WriteMessage(websocket.TextMessage, subscribeRequestBytes)
	if err != nil {
		ws.connMu.Unlock()
		return err
	}
	ws.connMu.Unlock()

	err = ws.expectSubscriptionResponse(authChannel)
	return err
}

func (ws *WebSocketClient) SubscribeToTicker(symbol Symbol) error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	subscribeRequest := tickerSubscriptionRequest{
		Action:  actionSubscribe,
		Channel: tickerChannel,
		Symbol:  symbol,
	}

	subscribeRequestBytes, err := json.Marshal(subscribeRequest)
	if err != nil {
		return err
	}

	ws.connMu.Lock()
	err = ws.conn.WriteMessage(websocket.TextMessage, subscribeRequestBytes)
	if err != nil {
		ws.connMu.Unlock()
		return err
	}
	ws.connMu.Unlock()

	err = ws.expectSubscriptionResponse(tickerChannel)
	return err
}

func (ws *WebSocketClient) SubscribeToTrades(symbol Symbol) error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	subscribeRequest := tradesSubscriptionRequest{
		Action:  actionSubscribe,
		Channel: tradesChannel,
		Symbol:  symbol,
	}

	subscribeRequestBytes, err := json.Marshal(subscribeRequest)
	if err != nil {
		return err
	}

	ws.connMu.Lock()
	err = ws.conn.WriteMessage(websocket.TextMessage, subscribeRequestBytes)
	if err != nil {
		ws.connMu.Unlock()
		return err
	}
	ws.connMu.Unlock()

	err = ws.expectSubscriptionResponse(tradesChannel)
	return err
}

func (ws *WebSocketClient) SubscribeToL2(symbol Symbol) error {
	return ws.quoteSubscription(symbol, l2Channel)
}

func (ws *WebSocketClient) SubscribeHeartbeat() error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	subscribeRequest := heartbeatSubscriptionRequest{
		Action:  actionSubscribe,
		Channel: heartbeatChannel,
	}

	subscribeRequestBytes, err := json.Marshal(subscribeRequest)
	if err != nil {
		return err
	}

	ws.connMu.Lock()
	err = ws.conn.WriteMessage(websocket.TextMessage, subscribeRequestBytes)
	if err != nil {
		ws.connMu.Unlock()
		return err
	}
	ws.connMu.Unlock()

	err = ws.expectSubscriptionResponse(heartbeatChannel)
	return err
}

func (ws *WebSocketClient) resetHeartbeat() {
	ws.heartbeatTimer.Reset(PingFrequency)
}

func (ws *WebSocketClient) expectSubscriptionResponse(subChannel channel) error {
	for {
		select {
		case subMsg := <-ws.subscriptionResponseChannel:
			if subMsg.SubscriptionName == subChannel.String() && len(subMsg.ErrorString) == 0 {
				log.Printf("Successfully subscribed to %s", subMsg.SubscriptionName)
				return nil
			} else if subMsg.SubscriptionName == subChannel.String() && len(subMsg.ErrorString) > 0 {
				log.Printf("Failed to subscribe: %s %s", subChannel.String(), subMsg.ErrorString)
				return errors.New(subMsg.ErrorString)
			}
		case <-time.After(ws.config.Timeout):
			log.Printf("timed out waiting for subscription response (channel: %s)", subChannel.String())
			return errors.New("timed out waiting for subscription response")
		}
	}
}

func (ws *WebSocketClient) NewOrderSingleMessage(order NewOrderSingleMsg) error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	newOrderSingleMsgRequest := newOrderSingleRequest{
		Action:      newOrderSingle,
		Channel:     tradingChannel,
		ClOrdID:     order.ClOrdID,
		Symbol:      order.Symbol,
		OrdType:     order.OrdType,
		TimeInForce: order.TimeInForce,
		Side:        order.Side,
		OrderQty:    order.OrderQty,
		Price:       order.Price,
		ExecInst:    order.ExecInst,
	}

	newOrderSingleRequestBytes, err := json.Marshal(newOrderSingleMsgRequest)
	if err != nil {
		return err
	}

	ws.connMu.Lock()
	err = ws.conn.WriteMessage(websocket.TextMessage, newOrderSingleRequestBytes)
	if err != nil {
		ws.connMu.Unlock()
		return err
	}
	ws.connMu.Unlock()
	return nil
}

func (ws *WebSocketClient) CancelOrder(orderID string) error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	cancelOrderRequest := cancelOrderRequest{
		Action:  cancelOrder,
		Channel: tradingChannel,
		OrderID: orderID,
	}

	cancelOrderRequestBytes, err := json.Marshal(cancelOrderRequest)
	if err != nil {
		return err
	}

	ws.connMu.Lock()
	err = ws.conn.WriteMessage(websocket.TextMessage, cancelOrderRequestBytes)
	if err != nil {
		ws.connMu.Unlock()
		return err
	}
	ws.connMu.Unlock()
	return nil
}

func (ws *WebSocketClient) BulkCancel(symbol *Symbol) error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	bulkCancelRequest := bulkCancelOrderRequest{
		Action:  bulkCancel,
		Channel: tradingChannel,
	}

	if symbol != nil {
		bulkCancelRequest.Symbol = *symbol
	}

	bulkCancelOrderRequestBytes, err := json.Marshal(bulkCancelRequest)
	if err != nil {
		return err
	}

	ws.connMu.Lock()
	err = ws.conn.WriteMessage(websocket.TextMessage, bulkCancelOrderRequestBytes)
	if err != nil {
		ws.connMu.Unlock()
		return err
	}
	ws.connMu.Unlock()
	return nil
}

func (ws *WebSocketClient) listenForUpdates(quitCh chan struct{}) {
	defer log.Println("listenForUpdates closed.")
	for {
		select {
		case <-quitCh:
			return
		default:
			_, msg, err := ws.conn.ReadMessage()
			if err != nil {
				select {
				case <-quitCh:
					return
				default:
				}
				log.Printf("Websocket read error: %s", err.Error())
				ws.Stop()
				ws.errorsChan <- err
				return
			} else {
				msgString := string(msg)
				if msgString == "ping" {
					log.Println("Received ping.")
					ws.resetHeartbeat()
				} else {
					var commonMsg msgCommon
					if err := json.Unmarshal(msg, &commonMsg); err != nil {
						log.Printf("Error un-marshalling common message: %s", err.Error())
						continue
					}
					switch commonMsg.Event {
					case eventSubscribed:
						ws.subscriptionResponseChannel <- SubscriptionError{
							SubscriptionName: commonMsg.Channel.String(),
						}
					case eventRejected:
						switch commonMsg.Channel {
						case tradingChannel:
							var rejectMsg TradingReject
							if err := json.Unmarshal(msg, &rejectMsg); err != nil {
								log.Printf("Error un-marshalling trading reject message: %s", err.Error())
								continue
							}
							ws.chTrading <- &rejectMsg
						default:
							var rejectMsg RejectMsg
							if err := json.Unmarshal(msg, &rejectMsg); err != nil {
								log.Printf("Error un-marshalling reject message: %s", err.Error())
								continue
							}
							ws.subscriptionResponseChannel <- SubscriptionError{
								SubscriptionName: commonMsg.Channel.String(),
								ErrorString:      rejectMsg.Text,
							}
						}
					case eventUpdate:
						switch commonMsg.Channel {
						case heartbeatChannel:
							var heartbeatMsg HeartbeatMsg
							if err := json.Unmarshal(msg, &heartbeatMsg); err != nil {
								log.Printf("Error un-marshalling heartbeat message: %s", err.Error())
								continue
							}
							ws.chHeartbeat <- heartbeatMsg
						case l3Channel:
							var l3Msg L3Msg
							if err := json.Unmarshal(msg, &l3Msg); err != nil {
								log.Printf("Error un-marshalling l3 update message: %s", err.Error())
								continue
							}
							ws.chL3 <- l3Msg
						case l2Channel:
							var l2Msg L2Msg
							if err := json.Unmarshal(msg, &l2Msg); err != nil {
								log.Printf("Error un-marshalling l2 update message: %s", err.Error())
								continue
							}
							ws.chL2 <- l2Msg
						case pricesChannel:
							var priceMsg PricesMsg
							if err := json.Unmarshal(msg, &priceMsg); err != nil {
								log.Printf("Error un-marshalling prices update message: %s", err.Error())
								continue
							}
							ws.chPrices <- priceMsg
						case tradesChannel:
							var tradesMsg TradesMsg
							if err := json.Unmarshal(msg, &tradesMsg); err != nil {
								log.Printf("Error un-marshalling trades update message: %s", err.Error())
								continue
							}
							ws.chTrades <- tradesMsg
						case tradingChannel:
							var tradingUpdate TradingUpdated
							if err := json.Unmarshal(msg, &tradingUpdate); err != nil {
								log.Printf("Error un-marshalling trading update message: %s", err.Error())
								continue
							}
							ws.chTrading <- &tradingUpdate
						}
					case eventSnapshot:
						switch commonMsg.Channel {
						case symbolsChannel:
							var symbolMsg SymbolsSnapshot
							if err := json.Unmarshal(msg, &symbolMsg); err != nil {
								log.Printf("Error un-marshalling symbols snapshot message: %s", err.Error())
								continue
							}
							for name, symbolData := range symbolMsg.Symbols {
								symbolData.Name = name
								ws.chSymbols <- symbolData
							}
						case l3Channel:
							var l3Msg L3Msg
							if err := json.Unmarshal(msg, &l3Msg); err != nil {
								log.Printf("Error un-marshalling l3 snapshot message: %s", err.Error())
								continue
							}
							ws.chL3 <- l3Msg
						case l2Channel:
							var l2Msg L2Msg
							if err := json.Unmarshal(msg, &l2Msg); err != nil {
								log.Printf("Error un-marshalling l2 snapshot message: %s", err.Error())
								continue
							}
							ws.chL2 <- l2Msg
						case tickerChannel:
							var tickerMsg TickerMsg
							if err := json.Unmarshal(msg, &tickerMsg); err != nil {
								log.Printf("Error un-marshalling ticker snapshot message: %s", err.Error())
								continue
							}
							ws.chTicker <- tickerMsg
						case balancesChannel:
							var balanceMsg BalancesSnapshot
							if err := json.Unmarshal(msg, &balanceMsg); err != nil {
								log.Printf("Error un-marshalling balances snapshot message: %s", err.Error())
								continue
							}
							ws.chBalances <- balanceMsg
						case tradingChannel:
							var tradingSnapShot TradingSnapshot
							if err := json.Unmarshal(msg, &tradingSnapShot); err != nil {
								log.Printf("Error un-marshalling trading snapshot message: %s", err.Error())
								continue
							}
							ws.chTrading <- &tradingSnapShot
						}
					}
				}
				//log.Printf("Received message: %s", msg)
				ws.resetHeartbeat()
			}
		}
	}
}
