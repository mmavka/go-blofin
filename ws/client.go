// Package ws provides WebSocket client functionality.
//
// This file implements the base WebSocket client, message routing, and logging.
//
// NOTE: Reconnect logic is NOT implemented in the library. Connection loss and errors are returned to the caller.
// The application is responsible for reconnecting and resubscribing if needed.
package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mmavka/go-blofin/models"
)

type subscription struct {
	channel string
	instID  string
	stype   string // "callback" или "channel"
	handler any    // func(models.WSCandlestickMsg) или nil
	ch      any    // chan models.WSCandlestickMsg или nil
}

// Client is a base WebSocket client with reconnect and logging support.
type Client struct {
	url  string
	conn *websocket.Conn
	subs map[string]struct{}

	mu sync.Mutex

	handlersCandles     map[string][]func(models.WSCandlestickMsg)
	handlersTrades      map[string][]func(models.WSTradeMsg)
	handlersTickers     map[string][]func(models.WSTickerMsg)
	handlersOrderBook   map[string][]func(models.WSOrderBookMsg)
	handlersFundingRate map[string][]func(models.WSFundingRateMsg)

	channelsCandles     map[string]chan models.WSCandlestickMsg
	channelsTrades      map[string]chan models.WSTradeMsg
	channelsTickers     map[string]chan models.WSTickerMsg
	channelsOrderBook   map[string]chan models.WSOrderBookMsg
	channelsFundingRate map[string]chan models.WSFundingRateMsg

	subscriptions []subscription
	pingInterval  time.Duration
	pongTimeout   time.Duration
	lastPong      time.Time
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	onError       func(error) // Error callback
}

// NewClient creates a new WebSocket client.
func NewClient(url string) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		url:                 url,
		subs:                make(map[string]struct{}),
		handlersCandles:     make(map[string][]func(models.WSCandlestickMsg)),
		handlersTrades:      make(map[string][]func(models.WSTradeMsg)),
		handlersTickers:     make(map[string][]func(models.WSTickerMsg)),
		handlersOrderBook:   make(map[string][]func(models.WSOrderBookMsg)),
		handlersFundingRate: make(map[string][]func(models.WSFundingRateMsg)),
		channelsCandles:     make(map[string]chan models.WSCandlestickMsg),
		channelsTrades:      make(map[string]chan models.WSTradeMsg),
		channelsTickers:     make(map[string]chan models.WSTickerMsg),
		channelsOrderBook:   make(map[string]chan models.WSOrderBookMsg),
		channelsFundingRate: make(map[string]chan models.WSFundingRateMsg),
		subscriptions:       []subscription{},
		pingInterval:        25 * time.Second,
		pongTimeout:         30 * time.Second,
		lastPong:            time.Now(),
		ctx:                 ctx,
		cancel:              cancel,
	}
}

// Connect establishes the WebSocket connection.
func (c *Client) Connect(ctx context.Context) error {
	d := websocket.Dialer{
		EnableCompression: false,
	}
	conn, _, err := d.DialContext(ctx, c.url, nil)
	if err != nil {
		return err
	}
	c.conn = conn

	// Control frame handlers
	c.conn.SetPingHandler(func(appData string) error {
		return c.conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
	})
	c.conn.SetPongHandler(func(appData string) error {
		c.lastPong = time.Now()
		return nil
	})

	c.wg.Add(1)
	go c.pingLoop()

	c.wg.Add(1)
	go c.readLoop()
	return nil
}

// pingLoop periodically sends ping and checks pong.
func (c *Client) pingLoop() {
	defer c.wg.Done()
	for {
		select {
		case <-time.After(c.pingInterval):
			if err := c.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second)); err != nil {
				if c.onError != nil {
					c.onError(err)
				}
				return
			}
			if time.Since(c.lastPong) > c.pongTimeout {
				if c.onError != nil {
					c.onError(fmt.Errorf("pong timeout"))
				}
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// readLoop reads and dispatches incoming messages.
func (c *Client) readLoop() {
	defer c.wg.Done()
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				if c.onError != nil {
					c.onError(err)
				}
				return
			}
			// Обработка pong
			if string(msg) == `{"event":"pong"}` {
				c.lastPong = time.Now()
				continue
			}

			// Универсальный парсер для определения канала
			var base struct {
				Arg struct {
					Channel string `json:"channel"`
					InstID  string `json:"instId"`
				} `json:"arg"`
				Data json.RawMessage `json:"data"`
			}
			if err := json.Unmarshal(msg, &base); err != nil || base.Arg.Channel == "" {
				continue
			}
			ch := base.Arg.Channel
			matched := false
			switch {
			case ch == "trades":
				var tradeMsg models.WSTradeMsg
				if err := json.Unmarshal(msg, &tradeMsg); err == nil {
					c.dispatchTrade(tradeMsg)
					matched = true
				}
			case ch == "tickers":
				var tickerMsg models.WSTickerMsg
				if err := json.Unmarshal(msg, &tickerMsg); err == nil {
					c.dispatchTicker(tickerMsg)
					matched = true
				}
			case ch == "books" || ch == "books5":
				var obMsg models.WSOrderBookMsg
				if err := json.Unmarshal(msg, &obMsg); err == nil {
					c.dispatchOrderBook(obMsg)
					matched = true
				}
			case ch == "fundingrate":
				var frMsg models.WSFundingRateMsg
				if err := json.Unmarshal(msg, &frMsg); err == nil {
					c.dispatchFundingRate(frMsg)
					matched = true
				}
			case len(ch) >= 6 && ch[:6] == "candle":
				var wsMsg models.WSCandlestickMsg
				if err := json.Unmarshal(msg, &wsMsg); err == nil {
					c.dispatchCandlestick(wsMsg)
					matched = true
				}
			}
			if !matched {
				continue
			}
		}
	}
}

// dispatchCandlestick routes candlestick messages to handlers and channels.
func (c *Client) dispatchCandlestick(msg models.WSCandlestickMsg) {
	key := msg.Arg.Channel + ":" + msg.Arg.InstID
	c.mu.Lock()
	handlers := c.handlersCandles[key]
	ch := c.channelsCandles[key]
	c.mu.Unlock()

	for _, handler := range handlers {
		go handler(msg)
	}
	if ch != nil {
		select {
		case ch <- msg:
		default:
		}
	}
}

// dispatchTrade routes trade messages (заглушка)
func (c *Client) dispatchTrade(msg models.WSTradeMsg) {
	key := "trades:" + msg.Arg.InstID
	c.mu.Lock()
	handlers := c.handlersTrades[key]
	ch := c.channelsTrades[key]
	c.mu.Unlock()

	for _, handler := range handlers {
		go handler(msg)
	}
	if ch != nil {
		select {
		case ch <- msg:
		default:
		}
	}
}

// dispatchTicker routes ticker messages (заглушка)
func (c *Client) dispatchTicker(msg models.WSTickerMsg) {
	key := "tickers:" + msg.Arg.InstID
	c.mu.Lock()
	handlers := c.handlersTickers[key]
	ch := c.channelsTickers[key]
	c.mu.Unlock()

	for _, handler := range handlers {
		go handler(msg)
	}
	if ch != nil {
		select {
		case ch <- msg:
		default:
		}
	}
}

// dispatchOrderBook routes order book messages (заглушка)
func (c *Client) dispatchOrderBook(msg models.WSOrderBookMsg) {
	key := msg.Arg.Channel + ":" + msg.Arg.InstID
	c.mu.Lock()
	handlers := c.handlersOrderBook[key]
	ch := c.channelsOrderBook[key]
	c.mu.Unlock()

	for _, handler := range handlers {
		go handler(msg)
	}
	if ch != nil {
		select {
		case ch <- msg:
		default:
		}
	}
}

// dispatchFundingRate routes funding rate messages (заглушка)
func (c *Client) dispatchFundingRate(msg models.WSFundingRateMsg) {
	key := "fundingrate:" + msg.Arg.InstID
	c.mu.Lock()
	handlers := c.handlersFundingRate[key]
	ch := c.channelsFundingRate[key]
	c.mu.Unlock()

	for _, handler := range handlers {
		go handler(msg)
	}
	if ch != nil {
		select {
		case ch <- msg:
		default:
		}
	}
}

// Close closes the WebSocket connection and завершает все goroutine.
func (c *Client) Close() error {
	c.cancel()
	if c.conn != nil {
		c.conn.Close()
	}
	c.wg.Wait()
	return nil
}

// SetErrorHandler sets the error callback for connection errors.
func (c *Client) SetErrorHandler(handler func(error)) {
	c.onError = handler
}
