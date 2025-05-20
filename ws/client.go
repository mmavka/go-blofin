/**
 * @file: client.go
 * @description: WebSocket client for BloFin (public/private, ping/pong, subscription)
 * @dependencies: models.go, signature.go
 * @created: 2025-05-19
 */

package ws

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	json "github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/mmavka/go-blofin/rest"
)

const (
	// Ping/pong interval (25 seconds < 30 seconds limit)
	pingInterval = 25 * time.Second
	pongTimeout  = 5 * time.Second

	// New connection limit
	reconnectDelay = time.Second
)

var privateChannels = map[string]struct{}{
	"orders":      {},
	"positions":   {},
	"orders-algo": {},
	// can add other private channels as needed
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
}

type LoginCredentials struct {
	ApiKey     string
	Secret     string
	Passphrase string
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

func NewDefaultClient() *Client {
	return NewClient(rest.DefaultWSPublic)
}

func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.reconnecting {
		time.Sleep(reconnectDelay) // Limit: 1 connection per second
	}

	// Close previous connection if exists
	if c.conn != nil {
		c.conn.Close()
	}

	u, err := url.Parse(c.url)
	if err != nil {
		return err
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	c.conn = conn

	// Setup ping/pong
	c.setupPingPong()

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

func (c *Client) setupPingPong() {
	c.pingTimer = time.NewTimer(pingInterval)
	c.pongTimer = time.NewTimer(pongTimeout)

	// Pong message handler
	c.conn.SetPongHandler(func(string) error {
		c.mu.Lock()
		c.lastPongTime = time.Now()
		c.pongTimer.Reset(pongTimeout)
		c.mu.Unlock()
		return nil
	})

	// Goroutine for sending ping
	go func() {
		for {
			select {
			case <-c.closeChan:
				return
			case <-c.pingTimer.C:
				if err := c.Ping(); err != nil {
					if c.errorHandler != nil {
						c.errorHandler(fmt.Errorf("ping error: %w", err))
					}
					c.reconnect()
					return
				}
				c.pingTimer.Reset(pingInterval)
			case <-c.pongTimer.C:
				if time.Since(c.lastPongTime) > pongTimeout {
					if c.errorHandler != nil {
						c.errorHandler(fmt.Errorf("pong timeout"))
					}
					c.reconnect()
					return
				}
				c.pongTimer.Reset(pongTimeout)
			}
		}
	}()
}

func (c *Client) reconnect() {
	c.mu.Lock()
	if c.reconnecting {
		c.mu.Unlock()
		return
	}
	c.reconnecting = true
	c.mu.Unlock()

	for {
		if err := c.Connect(); err == nil {
			c.mu.Lock()
			c.reconnecting = false
			c.mu.Unlock()
			return
		}
		time.Sleep(reconnectDelay)
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
			if c.errorHandler != nil {
				c.errorHandler(fmt.Errorf("panic in readLoop: %v", r))
			}
		}
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if c.errorHandler != nil {
				c.errorHandler(err)
			}
			select {
			case c.errors <- err:
			default:
			}
			return
		}

		select {
		case c.messages <- msg:
		default:
			// Channel overflow or closed, skip
			continue
		}

		// Base structure for event type determination
		var base struct {
			Arg struct {
				Channel string `json:"channel"`
			} `json:"arg"`
			Op    string `json:"op"`
			Event string `json:"event"`
		}
		if err := json.Unmarshal(msg, &base); err != nil {
			continue
		}

		switch base.Arg.Channel {
		case "trades":
			var tradeMsg TradeWSMessage
			if err := json.Unmarshal(msg, &tradeMsg); err == nil {
				select {
				case c.trades <- tradeMsg:
				default:
				}
			}
		case "funding-rate":
			var fundingMsg FundingRateWSMessage
			if err := json.Unmarshal(msg, &fundingMsg); err == nil {
				select {
				case c.fundingRates <- fundingMsg:
				default:
				}
			}
		case "tickers":
			var tickerMsg TickerWSMessage
			if err := json.Unmarshal(msg, &tickerMsg); err == nil {
				select {
				case c.tickers <- tickerMsg:
				default:
				}
			}
		case "books", "books5":
			var obMsg OrderBookWSMessage
			if err := json.Unmarshal(msg, &obMsg); err == nil {
				select {
				case c.orderBooks <- obMsg:
				default:
				}
			}
		default:
			// Check for candles (channel starts with candle)
			if len(base.Arg.Channel) >= 6 && base.Arg.Channel[:6] == "candle" {
				var candleMsg CandleWSMessage
				if err := json.Unmarshal(msg, &candleMsg); err == nil {
					select {
					case c.candles <- candleMsg:
					default:
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

// Ping send ping
func (c *Client) Ping() error {
	return c.conn.WriteMessage(websocket.TextMessage, []byte("ping"))
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
