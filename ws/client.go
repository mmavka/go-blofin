/**
 * @file: client.go
 * @description: WebSocket клиент для BloFin (public/private, ping/pong, подписки)
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
)

var privateChannels = map[string]struct{}{
	"orders":      {},
	"positions":   {},
	"orders-algo": {},
	// можно добавить другие приватные каналы по мере необходимости
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
	}
}

func NewDefaultClient() *Client {
	return NewClient("wss://ws.blofin.com/ws")
}

func (c *Client) Connect() error {
	u, err := url.Parse(c.url)
	if err != nil {
		return err
	}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	c.conn = conn
	go c.readLoop()
	return nil
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
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if c.errorHandler != nil {
				c.errorHandler(err)
			}
			c.errors <- err
			return
		}
		c.messages <- msg

		// Базовая структура для определения типа события
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
					// канал переполнен, пропускаем
				}
			}
		case "funding-rate":
			var fundingMsg FundingRateWSMessage
			if err := json.Unmarshal(msg, &fundingMsg); err == nil {
				select {
				case c.fundingRates <- fundingMsg:
				default:
					// канал переполнен, пропускаем
				}
			}
		case "tickers":
			var tickerMsg TickerWSMessage
			if err := json.Unmarshal(msg, &tickerMsg); err == nil {
				select {
				case c.tickers <- tickerMsg:
				default:
					// канал переполнен, пропускаем
				}
			}
		case "books", "books5":
			var obMsg OrderBookWSMessage
			if err := json.Unmarshal(msg, &obMsg); err == nil {
				select {
				case c.orderBooks <- obMsg:
				default:
					// канал переполнен, пропускаем
				}
			}
		default:
			// Проверка на свечи (канал начинается с candle)
			if len(base.Arg.Channel) >= 6 && base.Arg.Channel[:6] == "candle" {
				var candleMsg CandleWSMessage
				if err := json.Unmarshal(msg, &candleMsg); err == nil {
					select {
					case c.candles <- candleMsg:
					default:
						// канал переполнен, пропускаем
					}
				}
			}
			// Можно добавить обработку других каналов
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

// Login для приватных каналов
func (c *Client) Login(apiKey, secret, passphrase string) error {
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
		c.isLoggedIn = true // простая логика, можно доработать по событию login success
	}
	return err
}

// Subscribe к каналам
func (c *Client) Subscribe(channels []ChannelArgs) error {
	// Проверка приватных каналов
	for _, ch := range channels {
		if _, ok := privateChannels[ch.Channel]; ok && !c.isLoggedIn {
			return fmt.Errorf("subscription to private channel '%s' requires login", ch.Channel)
		}
	}
	// Проверка длины запроса
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

// Unsubscribe от каналов
func (c *Client) Unsubscribe(channels []ChannelArgs) error {
	req := UnsubscribeRequest{
		Op:   "unsubscribe",
		Args: channels,
	}
	return c.Send(req)
}

// Ping отправка ping
func (c *Client) Ping() error {
	return c.conn.WriteMessage(websocket.TextMessage, []byte("ping"))
}

// Messages возвращает канал для чтения сообщений
func (c *Client) Messages() <-chan []byte {
	return c.messages
}

// Errors возвращает канал ошибок
func (c *Client) Errors() <-chan error {
	return c.errors
}

// Trades возвращает канал для push trades
func (c *Client) Trades() <-chan TradeWSMessage {
	return c.trades
}

// Candles возвращает канал для push свечей
func (c *Client) Candles() <-chan CandleWSMessage {
	return c.candles
}

// OrderBooks возвращает канал для push order book
func (c *Client) OrderBooks() <-chan OrderBookWSMessage {
	return c.orderBooks
}

// Tickers возвращает канал для push тикеров
func (c *Client) Tickers() <-chan TickerWSMessage {
	return c.tickers
}

// FundingRates возвращает канал для push funding-rate
func (c *Client) FundingRates() <-chan FundingRateWSMessage {
	return c.fundingRates
}

// Пример использования:
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
