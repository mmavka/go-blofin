// Package ws provides WebSocket client for Blofin public channels.
//
// This file implements the base WebSocket client, message routing, reconnect logic, and logging.
package ws

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"runtime"

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
	url       string
	conn      *websocket.Conn
	log       Logger
	subs      map[string]struct{}
	reconnect bool

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
	reconnecting  bool

	pingLoopStarted bool // Флаг: запущен ли pingLoop
}

// Logger is an interface for event logging.
type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Tracef(format string, args ...interface{})
}

// NewClient creates a new WebSocket client.
func NewClient(url string, logger ...Logger) *Client {
	var l Logger
	if len(logger) == 0 || logger[0] == nil {
		l = NewDefaultLogger(LogLevelWarn)
	} else {
		l = logger[0]
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		url:                 url,
		log:                 l,
		subs:                make(map[string]struct{}),
		reconnect:           true,
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
	c.log.Debugf("Connect: dialing %s (goroutines: %d)", c.url, runtime.NumGoroutine())
	d := websocket.Dialer{
		EnableCompression: false,
	}
	conn, _, err := d.DialContext(ctx, c.url, nil)
	if err != nil {
		c.log.Errorf("Connect: error: %v", err)
		return err
	}
	c.conn = conn
	c.log.Debugf("Connect: dialed %s (goroutines: %d)", c.url, runtime.NumGoroutine())

	// Control frame handlers
	c.conn.SetPingHandler(func(appData string) error {
		return c.conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
	})
	c.conn.SetPongHandler(func(appData string) error {
		c.lastPong = time.Now()
		c.log.Debugf("Pong received (SetPongHandler): lastPong=%s", c.lastPong.Format("2006/01/02 15:04:05.000000000"))
		return nil
	})
	c.log.Debugf("SetPongHandler installed after connect (goroutines: %d)", runtime.NumGoroutine())

	c.wg.Add(1)
	c.log.Debugf("Connect: wg.Add(1) (goroutines: %d)", runtime.NumGoroutine())
	go c.readLoop()
	return nil
}

// pingLoop periodically sends ping and checks pong.
func (c *Client) pingLoop() {
	c.log.Debugf("pingLoop: entering loop (goroutines: %d, lastPong: %s)", runtime.NumGoroutine(), c.lastPong.Format("2006/01/02 15:04:05.000000000"))
	defer func() {
		c.log.Debugf("pingLoop: exiting loop (goroutines: %d)", runtime.NumGoroutine())
		c.wg.Done()
		c.log.Debugf("pingLoop: wg.Done (goroutines: %d)", runtime.NumGoroutine())
	}()
	for {
		select {
		case <-time.After(c.pingInterval):
			if err := c.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second)); err != nil {
				c.log.Errorf("Ping error: %v", err)
				go c.reconnectWS()
				c.log.Debugf("pingLoop: return after ping error (goroutines: %d)", runtime.NumGoroutine())
				return
			}
			c.log.Debugf("Ping sent (lastPong: %s)", c.lastPong.Format("2006/01/02 15:04:05.000000000"))
			if time.Since(c.lastPong) > c.pongTimeout {
				c.log.Errorf("Pong timeout, reconnecting...")
				go c.reconnectWS()
				c.log.Debugf("pingLoop: return after pong timeout (goroutines: %d)", runtime.NumGoroutine())
				return
			}
		case <-c.ctx.Done():
			c.log.Debugf("pingLoop: select ctx.Done (goroutines: %d)", runtime.NumGoroutine())
			return
		}
	}
}

// readLoop reads and dispatches incoming messages.
func (c *Client) readLoop() {
	c.log.Tracef("readLoop: entering loop (goroutines: %d)", runtime.NumGoroutine())
	defer func() {
		c.log.Debugf("readLoop: exiting loop (goroutines: %d)", runtime.NumGoroutine())
		c.wg.Done()
		c.log.Debugf("readLoop: wg.Done (goroutines: %d)", runtime.NumGoroutine())
	}()
	for {
		select {
		case <-c.ctx.Done():
			c.log.Debugf("readLoop: select ctx.Done (goroutines: %d)", runtime.NumGoroutine())
			return
		default:
			c.log.Tracef("readLoop: ReadMessage start (goroutines: %d)", runtime.NumGoroutine())
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				c.log.Errorf("Read error: %v", err)
				go c.reconnectWS()
				c.log.Debugf("readLoop: return after read error (goroutines: %d)", runtime.NumGoroutine())
				return
			}
			c.log.Tracef("RAW: %s", string(msg)) // Логируем raw сообщение
			c.log.Tracef("readLoop: ReadMessage done (goroutines: %d)", runtime.NumGoroutine())
			// Обработка pong
			if string(msg) == `{"event":"pong"}` {
				c.lastPong = time.Now()
				c.log.Debugf("readLoop: continue pong (goroutines: %d)", runtime.NumGoroutine())
				continue
			}

			// Логируем event:subscribe и event:error
			var eventMsg struct {
				Event string `json:"event"`
				Arg   struct {
					Channel string `json:"channel"`
					InstID  string `json:"instId"`
				} `json:"arg"`
				Code string `json:"code"`
				Msg  string `json:"msg"`
			}
			if err := json.Unmarshal(msg, &eventMsg); err == nil && eventMsg.Event != "" {
				c.log.Infof("WS event: %s channel=%s instID=%s code=%s msg=%s", eventMsg.Event, eventMsg.Arg.Channel, eventMsg.Arg.InstID, eventMsg.Code, eventMsg.Msg)
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
				c.log.Debugf("readLoop: continue unmatched message (goroutines: %d)", runtime.NumGoroutine())
				continue
			}
			// Если есть поле data (push-сообщение), стартуем pingLoop
			// if len(base.Data) > 0 {
			// 	c.startPingLoopOnce(0)
			// }
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
				c.log.Debugf("readLoop: continue unmatched message (goroutines: %d)", runtime.NumGoroutine())
				continue
			}
		}
	}
}

// dispatchCandlestick routes candlestick messages to handlers and channels.
func (c *Client) dispatchCandlestick(msg models.WSCandlestickMsg) {
	c.log.Tracef("dispatchCandlestick: start")
	c.mu.Lock()
	handlers := c.handlersCandles[msg.Arg.Channel+":"+msg.Arg.InstID]
	ch := c.channelsCandles[msg.Arg.Channel+":"+msg.Arg.InstID]
	c.mu.Unlock()
	for _, h := range handlers {
		go h(msg)
	}
	if ch != nil {
		select {
		case ch <- msg:
		default:
		}
	}
	c.log.Tracef("dispatchCandlestick: end")
}

// dispatchTrade routes trade messages (заглушка)
func (c *Client) dispatchTrade(msg models.WSTradeMsg) {
	c.log.Tracef("dispatchTrade: start")
	key := "trades:" + msg.Arg.InstID
	c.mu.Lock()
	handlers := c.handlersTrades[key]
	ch := c.channelsTrades[key]
	c.mu.Unlock()
	for _, h := range handlers {
		go h(msg)
	}
	if ch != nil {
		select {
		case ch <- msg:
		default:
		}
	}
	c.log.Tracef("dispatchTrade: end")
}

// dispatchTicker routes ticker messages (заглушка)
func (c *Client) dispatchTicker(msg models.WSTickerMsg) {
	c.log.Tracef("dispatchTicker: start")
	key := "tickers:" + msg.Arg.InstID
	c.mu.Lock()
	handlers := c.handlersTickers[key]
	ch := c.channelsTickers[key]
	c.mu.Unlock()
	for _, h := range handlers {
		go h(msg)
	}
	if ch != nil {
		select {
		case ch <- msg:
		default:
		}
	}
	c.log.Tracef("dispatchTicker: end")
}

// dispatchOrderBook routes order book messages (заглушка)
func (c *Client) dispatchOrderBook(msg models.WSOrderBookMsg) {
	c.log.Tracef("dispatchOrderBook: start")
	key := msg.Arg.Channel + ":" + msg.Arg.InstID
	c.mu.Lock()
	handlers := c.handlersOrderBook[key]
	ch := c.channelsOrderBook[key]
	c.mu.Unlock()
	for _, h := range handlers {
		go h(msg)
	}
	if ch != nil {
		select {
		case ch <- msg:
		default:
		}
	}
	c.log.Tracef("dispatchOrderBook: end")
}

// dispatchFundingRate routes funding rate messages (заглушка)
func (c *Client) dispatchFundingRate(msg models.WSFundingRateMsg) {
	c.log.Tracef("dispatchFundingRate: start")
	key := "fundingrate:" + msg.Arg.InstID
	c.mu.Lock()
	handlers := c.handlersFundingRate[key]
	ch := c.channelsFundingRate[key]
	c.mu.Unlock()
	for _, h := range handlers {
		go h(msg)
	}
	if ch != nil {
		select {
		case ch <- msg:
		default:
		}
	}
	c.log.Tracef("dispatchFundingRate: end")
}

// Close closes the WebSocket connection and завершает все goroutine.
func (c *Client) Close() error {
	c.cancel()
	if c.conn != nil {
		_ = c.conn.Close()
	}
	c.wg.Wait()
	return nil
}

// reconnectWS performs reconnect and resubscribe.
func (c *Client) reconnectWS() {
	c.log.Debugf("reconnectWS: started (goroutines: %d)", runtime.NumGoroutine())
	c.mu.Lock()
	if c.reconnecting {
		c.mu.Unlock()
		c.log.Debugf("reconnectWS: already reconnecting, return (goroutines: %d)", runtime.NumGoroutine())
		return
	}
	c.reconnecting = true
	c.pingLoopStarted = false // сбросить флаг при reconnect
	c.mu.Unlock()

	c.log.Debugf("reconnectWS: calling cancel (goroutines: %d)", runtime.NumGoroutine())
	c.cancel()
	c.log.Debugf("reconnectWS: cancel called (goroutines: %d)", runtime.NumGoroutine())
	if c.conn != nil {
		_ = c.conn.Close()
		c.log.Debugf("reconnectWS: conn.Close called (goroutines: %d)", runtime.NumGoroutine())
	}
	c.log.Debugf("reconnectWS: waiting for goroutines to finish (goroutines: %d)", runtime.NumGoroutine())
	c.wg.Wait()
	c.log.Debugf("reconnectWS: goroutines finished (goroutines: %d)", runtime.NumGoroutine())
	ctx, cancel := context.WithCancel(context.Background())
	c.ctx = ctx
	c.cancel = cancel
	c.log.Debugf("reconnectWS: calling Connect (goroutines: %d)", runtime.NumGoroutine())
	if err := c.Connect(ctx); err != nil {
		c.log.Errorf("Reconnect failed: %v", err)
		c.mu.Lock()
		c.reconnecting = false
		c.mu.Unlock()
		c.log.Debugf("reconnectWS: Connect failed, return (goroutines: %d)", runtime.NumGoroutine())
		return
	}
	c.log.Debugf("reconnectWS: calling resubscribeAll (goroutines: %d)", runtime.NumGoroutine())
	c.resubscribeAll()
	// c.startPingLoopOnce(2 * time.Second)
	c.mu.Lock()
	c.reconnecting = false
	c.mu.Unlock()
	c.log.Debugf("reconnectWS: finished (goroutines: %d)", runtime.NumGoroutine())
}

// deduplicateSubscriptions removes duplicate subscriptions.
func (c *Client) deduplicateSubscriptions() {
	c.mu.Lock()
	unique := make(map[string]subscription)
	for _, s := range c.subscriptions {
		key := s.channel + ":" + s.instID + ":" + s.stype
		unique[key] = s
	}
	c.subscriptions = make([]subscription, 0, len(unique))
	for _, s := range unique {
		c.subscriptions = append(c.subscriptions, s)
	}
	c.mu.Unlock()
}

// resubscribeAll resends all active subscriptions.
func (c *Client) resubscribeAll() {
	c.log.Debugf("resubscribeAll: started (goroutines: %d)", runtime.NumGoroutine())
	c.deduplicateSubscriptions()
	c.mu.Lock()
	subs := make([]subscription, len(c.subscriptions))
	copy(subs, c.subscriptions)
	c.mu.Unlock()

	// Не отправляем unsubscribe, если соединение было восстановлено — сервер уже "забыл" старые подписки

	// Просто отправляем subscribe как обычно
	for _, s := range subs {
		c.log.Debugf("resubscribeAll: before subscribe %s:%s type=%s (goroutines: %d)", s.channel, s.instID, s.stype, runtime.NumGoroutine())
		switch s.channel {
		case "trades":
			if s.stype == "callback" {
				err := c.SubscribeTrades(context.Background(), s.instID, s.handler.(func(models.WSTradeMsg)))
				if err != nil {
					c.log.Errorf("Resubscribe trades error: %v", err)
				}
			} else {
				_, err := c.SubscribeTradesChan(context.Background(), s.instID)
				if err != nil {
					c.log.Errorf("Resubscribe trades error: %v", err)
				}
			}
		case "tickers":
			if s.stype == "callback" {
				err := c.SubscribeTickers(context.Background(), s.instID, s.handler.(func(models.WSTickerMsg)))
				if err != nil {
					c.log.Errorf("Resubscribe tickers error: %v", err)
				}
			} else {
				_, err := c.SubscribeTickersChan(context.Background(), s.instID)
				if err != nil {
					c.log.Errorf("Resubscribe tickers error: %v", err)
				}
			}
		case "fundingrate":
			if s.stype == "callback" {
				err := c.SubscribeFundingRate(context.Background(), s.instID, s.handler.(func(models.WSFundingRateMsg)))
				if err != nil {
					c.log.Errorf("Resubscribe fundingrate error: %v", err)
				}
			} else {
				_, err := c.SubscribeFundingRateChan(context.Background(), s.instID)
				if err != nil {
					c.log.Errorf("Resubscribe fundingrate error: %v", err)
				}
			}
		case "books", "books5":
			if s.stype == "callback" {
				err := c.SubscribeOrderBook(context.Background(), s.channel, s.instID, s.handler.(func(models.WSOrderBookMsg)))
				if err != nil {
					c.log.Errorf("Resubscribe orderbook error: %v", err)
				}
			} else {
				_, err := c.SubscribeOrderBookChan(context.Background(), s.channel, s.instID)
				if err != nil {
					c.log.Errorf("Resubscribe orderbook error: %v", err)
				}
			}
		default:
			if s.stype == "callback" {
				err := c.SubscribeCandlesticks(context.Background(), s.channel, s.instID, s.handler.(func(models.WSCandlestickMsg)))
				if err != nil {
					c.log.Errorf("Resubscribe candles error: %v", err)
				}
			} else {
				// Только отправляем subscribe-запрос, не создаём новый канал
				req := map[string]interface{}{
					"op":   "subscribe",
					"args": []map[string]string{{"channel": s.channel, "instId": s.instID}},
				}
				msg, _ := json.Marshal(req)
				err := c.conn.WriteMessage(1, msg)
				if err != nil {
					c.log.Errorf("Resubscribe candles error: %v", err)
				}
			}
		}
		c.log.Debugf("resubscribeAll: after subscribe %s:%s type=%s (goroutines: %d)", s.channel, s.instID, s.stype, runtime.NumGoroutine())
	}
	c.log.Debugf("resubscribeAll: finished (goroutines: %d)", runtime.NumGoroutine())
}

// Запуск pingLoop только если он ещё не запущен
func (c *Client) startPingLoopOnce(delay time.Duration) {
	c.mu.Lock()
	if c.pingLoopStarted {
		c.mu.Unlock()
		return
	}
	c.pingLoopStarted = true
	c.mu.Unlock()
	c.wg.Add(1)
	go func() {
		time.Sleep(delay)
		c.pingLoop()
	}()
}
