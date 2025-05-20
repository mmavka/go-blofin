/**
 * @file: client.go
 * @description: WebSocket client for BloFin (public/private, ping/pong, subscription)
 * @dependencies: models.go, signature.go
 * @created: 2025-05-19
 */

package ws

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mmavka/go-blofin"
)

const (
	// Ping/pong interval (10 seconds < 30 seconds limit)
	pingInterval = 20 * time.Second
	pongTimeout  = 10 * time.Second

	// New connection limit
	reconnectDelay = 500 * time.Millisecond
	maxRetries     = 20
)

var privateChannels = map[string]struct{}{
	blofin.ChannelOrders:     {},
	blofin.ChannelPositions:  {},
	blofin.ChannelOrdersAlgo: {},
	// can add other private channels as needed
}

// UseTestnet use testnet
var UseTestnet = false

// WsHandler handle raw websocket message
type WsHandler func(message []byte)

// ErrHandler handles errors
type ErrHandler func(err error)

// WsConfig webservice configuration
type WsConfig struct {
	Endpoint string
}

// WsClient is websocket client
type WsClient struct {
	config     *WsConfig
	conn       *websocket.Conn
	mu         sync.RWMutex
	stopC      chan struct{}
	doneC      chan struct{}
	handler    WsHandler
	errHandler ErrHandler
}

type Client struct {
	conn         *websocket.Conn
	url          string
	messages     chan []byte
	errors       chan error
	trades       chan TradeWSMessage
	candles      chan CandleWSMessage
	orderBooks   chan OrderBookWSMessage
	tickers      chan TickerWSMessage
	fundingRates chan FundingRateWSMessage
	closeChan    chan struct{}
	once         sync.Once
	errorHandler func(error)
	isLoggedIn   bool

	// For ping/pong
	pingTimer    *time.Timer
	pongTimer    *time.Timer
	lastPongTime time.Time

	// For reconnection
	reconnecting bool
	mu           sync.Mutex

	// Save subscriptions for reconnection
	subscriptions []ChannelArgs
	credentials   *LoginCredentials

	// Keepalive
	resetKeepAlive func()
	handlePong     func()
}

type LoginCredentials struct {
	ApiKey     string
	Secret     string
	Passphrase string
}

// NewWsClient creates a new websocket client
func NewWsClient(config *WsConfig) *WsClient {
	return &WsClient{
		config: config,
		stopC:  make(chan struct{}),
		doneC:  make(chan struct{}),
	}
}

// SetHandler sets the message handler
func (c *WsClient) SetHandler(handler WsHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handler = handler
}

// SetErrHandler sets the error handler
func (c *WsClient) SetErrHandler(errHandler ErrHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.errHandler = errHandler
}

// Connect connects to the websocket server
func (c *WsClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return fmt.Errorf("already connected")
	}

	conn, _, err := websocket.DefaultDialer.Dial(c.config.Endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to dial websocket: %w", err)
	}

	c.conn = conn
	go c.readMessages()
	return nil
}

// Close closes the websocket connection
func (c *WsClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return fmt.Errorf("not connected")
	}

	close(c.stopC)
	<-c.doneC

	err := c.conn.Close()
	c.conn = nil
	return err
}

// readMessages reads messages from the websocket connection
func (c *WsClient) readMessages() {
	defer close(c.doneC)

	for {
		select {
		case <-c.stopC:
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				c.mu.RLock()
				if c.errHandler != nil {
					c.errHandler(err)
				}
				c.mu.RUnlock()
				return
			}

			c.mu.RLock()
			if c.handler != nil {
				c.handler(message)
			}
			c.mu.RUnlock()
		}
	}
}

// WriteMessage writes a message to the websocket connection
func (c *WsClient) WriteMessage(messageType int, data []byte) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.conn == nil {
		return fmt.Errorf("not connected")
	}

	return c.conn.WriteMessage(messageType, data)
}

// WsServe starts websocket connection and handles messages
func WsServe(config *WsConfig, handler WsHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	client := NewWsClient(config)
	client.SetHandler(handler)
	client.SetErrHandler(errHandler)

	if err := client.Connect(); err != nil {
		return nil, nil, err
	}

	return client.doneC, client.stopC, nil
}

// getWsEndpoint return the base endpoint of the WS according the UseTestnet flag
func getWsEndpoint() string {
	if UseTestnet {
		return blofin.TestnetWSPublic
	}
	return blofin.DefaultWSPublic
}

// WsUserDataServe starts user data stream websocket connection
func WsUserDataServe(listenKey string, handler WsHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s", blofin.DefaultWSPrivate, listenKey)
	config := &WsConfig{
		Endpoint: endpoint,
	}
	return WsServe(config, handler, errHandler)
}

// WsMarketDataServe starts market data stream websocket connection
func WsMarketDataServe(symbol string, handler WsHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/market/%s", blofin.DefaultWSPublic, symbol)
	config := &WsConfig{
		Endpoint: endpoint,
	}
	return WsServe(config, handler, errHandler)
}

// WsAccountDataServe starts account data stream websocket connection
func WsAccountDataServe(apiKey, secretKey string, handler WsHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/account/%s", blofin.DefaultWSPrivate, apiKey)
	config := &WsConfig{
		Endpoint: endpoint,
	}
	return WsServe(config, handler, errHandler)
}

func NewClient(url string) *Client {
	return &Client{
		url:          url,
		messages:     make(chan []byte, 100),
		errors:       make(chan error, 10),
		trades:       make(chan TradeWSMessage, 100),
		candles:      make(chan CandleWSMessage, 100),
		orderBooks:   make(chan OrderBookWSMessage, 100),
		tickers:      make(chan TickerWSMessage, 100),
		fundingRates: make(chan FundingRateWSMessage, 100),
		closeChan:    make(chan struct{}),
		lastPongTime: time.Now(),
	}
}

// NewDefaultClient creates a new WebSocket client with default settings
func NewDefaultClient() *Client {
	return NewClient(blofin.DefaultWSPublic)
}

func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.reconnecting {
		time.Sleep(reconnectDelay)
	}

	// Close previous connection if exists
	if c.conn != nil {
		c.conn.Close()
	}

	u, err := url.Parse(c.url)
	if err != nil {
		return err
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
		ReadBufferSize:   32768,
		WriteBufferSize:  32768,
		// Enable compression
		EnableCompression: true,
	}

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	// Set read deadline
	conn.SetReadDeadline(time.Now().Add(pongTimeout * 2))

	c.conn = conn

	// Setup ping/pong
	c.startKeepAlive()

	// Start message reading
	go c.readLoop()

	// Restore state after reconnection
	if c.credentials != nil {
		if err := c.Login(c.credentials.ApiKey, c.credentials.Secret, c.credentials.Passphrase); err != nil {
			return fmt.Errorf("failed to restore login: %w", err)
		}
	}

	if len(c.subscriptions) > 0 {
		if err := c.Subscribe(c.subscriptions); err != nil {
			return fmt.Errorf("failed to restore subscriptions: %w", err)
		}
	}

	return nil
}

func (c *Client) startKeepAlive() {
	pongCh := make(chan struct{}, 1)
	c.pingTimer = time.NewTimer(pingInterval)

	go func() {
		for {
			select {
			case <-c.closeChan:
				return
			case <-c.pingTimer.C:
				// Отправить текстовый ping
				if err := c.conn.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
					if c.errorHandler != nil {
						c.errorHandler(fmt.Errorf("ping error: %w", err))
					}
					c.reconnect()
					return
				}
				// Ждать pong
				pongTimer := time.NewTimer(pongTimeout)
				select {
				case <-pongCh:
					// pong получен, продолжаем
					pongTimer.Stop()
				case <-pongTimer.C:
					if c.errorHandler != nil {
						c.errorHandler(fmt.Errorf("pong timeout"))
					}
					c.reconnect()
					return
				}
				c.pingTimer.Reset(pingInterval)
			}
		}
	}()

	// Сброс таймера на любое входящее сообщение
	c.resetKeepAlive = func() {
		if !c.pingTimer.Stop() {
			select {
			case <-c.pingTimer.C:
			default:
			}
		}
		c.pingTimer.Reset(pingInterval)
	}

	// Обработка pong
	c.handlePong = func() {
		select {
		case pongCh <- struct{}{}:
		default:
		}
	}
}

func (c *Client) reconnect() {
	c.mu.Lock()
	if c.reconnecting {
		c.mu.Unlock()
		return
	}
	c.reconnecting = true
	c.mu.Unlock()

	retries := 0
	for {
		if err := c.Connect(); err == nil {
			c.mu.Lock()
			c.reconnecting = false
			c.mu.Unlock()
			return
		}

		retries++
		if retries >= maxRetries {
			if c.errorHandler != nil {
				c.errorHandler(fmt.Errorf("max reconnection attempts reached"))
			}
			c.mu.Lock()
			c.reconnecting = false
			c.mu.Unlock()
			return
		}

		// Exponential backoff with jitter
		delay := reconnectDelay * time.Duration(1<<uint(retries))
		jitter := time.Duration(rand.Int63n(int64(delay / 4)))
		time.Sleep(delay + jitter)
	}
}

func (c *Client) Close() {
	c.once.Do(func() {
		close(c.closeChan)
		if c.conn != nil {
			c.conn.Close()
		}
		close(c.messages)
		close(c.errors)
		close(c.trades)
		close(c.candles)
		close(c.orderBooks)
		close(c.tickers)
		close(c.fundingRates)
	})
}

func (c *Client) readLoop() {
	defer func() {
		if r := recover(); r != nil {
			c.errorHandler(fmt.Errorf("panic in readLoop: %v", r))
		}
	}()

	for {
		select {
		case <-c.closeChan:
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				c.errorHandler(fmt.Errorf("WebSocket error: %v", err))
				c.reconnect()
				return
			}

			// Reset keepalive timer on any message
			if c.resetKeepAlive != nil {
				c.resetKeepAlive()
			}

			// Handle pong message
			if string(message) == "pong" {
				if c.handlePong != nil {
					c.handlePong()
				}
				continue
			}

			// Try to parse as JSON
			var msg map[string]interface{}
			if err := json.Unmarshal(message, &msg); err != nil {
				c.errorHandler(fmt.Errorf("failed to parse message: %v", err))
				continue
			}

			// Handle ping message
			if event, ok := msg["event"].(string); ok && event == "ping" {
				if err := c.Send(map[string]string{"event": "pong"}); err != nil {
					c.errorHandler(fmt.Errorf("failed to send pong: %v", err))
				}
				continue
			}

			// Send to messages channel
			select {
			case c.messages <- message:
			default:
				// Channel is full, skip message
			}

			// Try to parse as specific message type
			if arg, ok := msg["arg"].(map[string]interface{}); ok {
				channel, _ := arg["channel"].(string)
				switch channel {
				case "trades":
					var tradeMsg TradeWSMessage
					if err := json.Unmarshal(message, &tradeMsg); err == nil {
						select {
						case c.trades <- tradeMsg:
						default:
							// Channel is full, skip message
						}
					}
				case "candle1m", "candle5m", "candle15m", "candle30m", "candle1h", "candle2h", "candle4h", "candle6h", "candle12h", "candle1d", "candle1w", "candle1M":
					var candleMsg CandleWSMessage
					if err := json.Unmarshal(message, &candleMsg); err == nil {
						select {
						case c.candles <- candleMsg:
						default:
							// Channel is full, skip message
						}
					}
				case "books":
					var bookMsg OrderBookWSMessage
					if err := json.Unmarshal(message, &bookMsg); err == nil {
						select {
						case c.orderBooks <- bookMsg:
						default:
							// Channel is full, skip message
						}
					}
				case "tickers":
					var tickerMsg TickerWSMessage
					if err := json.Unmarshal(message, &tickerMsg); err == nil {
						select {
						case c.tickers <- tickerMsg:
						default:
							// Channel is full, skip message
						}
					}
				case "funding-rate":
					var frMsg FundingRateWSMessage
					if err := json.Unmarshal(message, &frMsg); err == nil {
						select {
						case c.fundingRates <- frMsg:
						default:
							// Channel is full, skip message
						}
					}
				}
			}
		}
	}
}

func (c *Client) Send(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

// Login for private channels
func (c *Client) Login(apiKey, secret, passphrase string) error {
	c.credentials = &LoginCredentials{
		ApiKey:     apiKey,
		Secret:     secret,
		Passphrase: passphrase,
	}

	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
	nonce := timestamp
	sign := SignWebSocketLogin(secret, timestamp, nonce)
	login := LoginRequest{
		Op: "login",
		Args: []LoginArgs{{
			ApiKey:     apiKey,
			Passphrase: passphrase,
			Timestamp:  timestamp,
			Sign:       sign,
			Nonce:      nonce,
		}},
	}
	err := c.Send(login)
	if err == nil {
		c.isLoggedIn = true
	}
	return err
}

// Subscribe to channels
func (c *Client) Subscribe(channels []ChannelArgs) error {
	// Save subscriptions for reconnection
	c.subscriptions = append(c.subscriptions, channels...)

	// Check private channels
	for _, ch := range channels {
		if _, ok := privateChannels[ch.Channel]; ok && !c.isLoggedIn {
			return fmt.Errorf("subscription to private channel '%s' requires login", ch.Channel)
		}
	}

	req := SubscribeRequest{
		Op:   "subscribe",
		Args: channels,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	if len(data) > 4096 {
		return fmt.Errorf("subscription request exceeds 4096 bytes (actual: %d)", len(data))
	}
	return c.Send(req)
}

// Unsubscribe from channels
func (c *Client) Unsubscribe(channels []ChannelArgs) error {
	req := UnsubscribeRequest{
		Op:   "unsubscribe",
		Args: channels,
	}
	return c.Send(req)
}

// Messages returns channel for reading messages
func (c *Client) Messages() <-chan []byte {
	return c.messages
}

// Errors returns error channel
func (c *Client) Errors() <-chan error {
	return c.errors
}

// Trades returns channel for push trades
func (c *Client) Trades() <-chan TradeWSMessage {
	return c.trades
}

// Candles returns channel for push candles
func (c *Client) Candles() <-chan CandleWSMessage {
	return c.candles
}

// OrderBooks returns channel for push order book
func (c *Client) OrderBooks() <-chan OrderBookWSMessage {
	return c.orderBooks
}

// Tickers returns channel for push tickers
func (c *Client) Tickers() <-chan TickerWSMessage {
	return c.tickers
}

// FundingRates returns channel for push funding-rate
func (c *Client) FundingRates() <-chan FundingRateWSMessage {
	return c.fundingRates
}

// Usage example:
// ws := ws.NewDefaultClient()
// ws.SetErrorHandler(func(err error) {
//     log.Printf("WebSocket error: %v", err)
// })
// err := ws.Connect()
// if err != nil {
//     log.Fatal(err)
// }
// _ = ws.Subscribe([]ws.ChannelArgs{{Channel: "trades", InstId: "ETH-USDT"}})
// for trade := range ws.Trades() {
//     fmt.Println(trade)
// }

func (c *Client) SetURL(url string) {
	c.url = url
}

func (c *Client) SetErrorHandler(handler func(error)) {
	c.errorHandler = handler
}
