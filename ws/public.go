// Package ws provides public WebSocket API for Blofin.
//
// This file contains public subscription methods for candlestick channels (callback and channel interface).
package ws

import (
	"context"
	"encoding/json"

	"github.com/mmavka/go-blofin/models"
)

// MessageHandler is a callback for incoming messages.
type MessageHandler func(msg any)

// SubscribeCandlesticks subscribes to candlestick channel with callback.
// channel - candlestick channel name (e.g. "candle1m"), instID - instrument ID (e.g. "BTC-USDT").
// handler is called for each push message.
func (c *Client) SubscribeCandlesticks(ctx context.Context, channel, instID string, handler func(models.WSCandlestickMsg)) error {
	key := channel + ":" + instID
	c.mu.Lock()
	// Не добавлять дубликаты подписок
	found := false
	for _, s := range c.subscriptions {
		if s.channel == channel && s.instID == instID && s.stype == "callback" {
			found = true
			break
		}
	}
	if !found {
		c.subscriptions = append(c.subscriptions, subscription{channel: channel, instID: instID, stype: "callback", handler: handler, ch: nil})
	}
	if handler != nil {
		c.handlersCandles[key] = append(c.handlersCandles[key], handler)
	}
	c.mu.Unlock()
	// Формируем запрос подписки
	req := map[string]interface{}{
		"op":   "subscribe",
		"args": []map[string]string{{"channel": channel, "instId": instID}},
	}
	msg, _ := json.Marshal(req)
	if err := c.conn.WriteMessage(1, msg); err != nil {
		return err
	}
	return nil
}

// SubscribeCandlesticksChan subscribes to candlestick channel and returns channel for messages.
// channel - candlestick channel name (e.g. "candle1m"), instID - instrument ID (e.g. "BTC-USDT").
// Returns channel for push messages.
func (c *Client) SubscribeCandlesticksChan(ctx context.Context, channel, instID string) (<-chan models.WSCandlestickMsg, error) {
	key := channel + ":" + instID
	c.mu.Lock()
	// Не добавлять дубликаты подписок
	for _, s := range c.subscriptions {
		if s.channel == channel && s.instID == instID && s.stype == "channel" {
			if s.ch != nil {
				c.mu.Unlock()
				return s.ch.(chan models.WSCandlestickMsg), nil
			}
		}
	}
	if ch, ok := c.channelsCandles[key]; ok {
		// Если канал уже есть в channelsCandles, возвращаем его
		c.mu.Unlock()
		return ch, nil
	}
	ch := make(chan models.WSCandlestickMsg, 100)
	c.channelsCandles[key] = ch
	c.subscriptions = append(c.subscriptions, subscription{channel: channel, instID: instID, stype: "channel", handler: nil, ch: ch})
	c.mu.Unlock()
	// Формируем запрос подписки
	req := map[string]interface{}{
		"op":   "subscribe",
		"args": []map[string]string{{"channel": channel, "instId": instID}},
	}
	msg, _ := json.Marshal(req)
	if err := c.conn.WriteMessage(1, msg); err != nil {
		return nil, err
	}
	return ch, nil
}

// UnsubscribeCandlesticks unsubscribes from candlestick channel and removes handlers/channels.
// channel - candlestick channel name (e.g. "candle1m"), instID - instrument ID (e.g. "BTC-USDT").
func (c *Client) UnsubscribeCandlesticks(ctx context.Context, channel, instID string) error {
	key := channel + ":" + instID
	c.mu.Lock()
	// Удаляем из subscriptions
	newSubs := c.subscriptions[:0]
	for _, s := range c.subscriptions {
		if !(s.channel == channel && s.instID == instID) {
			newSubs = append(newSubs, s)
		}
	}
	c.subscriptions = newSubs
	delete(c.handlersCandles, key)
	if ch, ok := c.channelsCandles[key]; ok {
		close(ch)
		delete(c.channelsCandles, key)
	}
	c.mu.Unlock()
	if c.conn == nil {
		return nil
	}
	req := map[string]interface{}{
		"op":   "unsubscribe",
		"args": []map[string]string{{"channel": channel, "instId": instID}},
	}
	msg, _ := json.Marshal(req)
	if err := c.conn.WriteMessage(1, msg); err != nil {
		return err
	}
	return nil
}

// SubscribeTrades subscribes to trades channel with callback.
func (c *Client) SubscribeTrades(ctx context.Context, instID string, handler func(models.WSTradeMsg)) error {
	key := "trades:" + instID
	c.mu.Lock()
	if handler != nil {
		c.handlersTrades[key] = append(c.handlersTrades[key], handler)
	}
	c.subscriptions = append(c.subscriptions, subscription{channel: "trades", instID: instID})
	c.mu.Unlock()
	req := map[string]interface{}{
		"op":   "subscribe",
		"args": []map[string]string{{"channel": "trades", "instId": instID}},
	}
	msg, _ := json.Marshal(req)
	if err := c.conn.WriteMessage(1, msg); err != nil {
		return err
	}
	return nil
}

// SubscribeTradesChan subscribes to trades channel and returns channel for messages.
func (c *Client) SubscribeTradesChan(ctx context.Context, instID string) (<-chan models.WSTradeMsg, error) {
	key := "trades:" + instID
	c.mu.Lock()
	ch := make(chan models.WSTradeMsg, 100)
	c.channelsTrades[key] = ch
	c.subscriptions = append(c.subscriptions, subscription{channel: "trades", instID: instID})
	c.mu.Unlock()
	req := map[string]interface{}{
		"op":   "subscribe",
		"args": []map[string]string{{"channel": "trades", "instId": instID}},
	}
	msg, _ := json.Marshal(req)
	if err := c.conn.WriteMessage(1, msg); err != nil {
		return nil, err
	}
	return ch, nil
}

// UnsubscribeTrades unsubscribes from trades channel and removes handlers/channels.
func (c *Client) UnsubscribeTrades(ctx context.Context, instID string) error {
	key := "trades:" + instID
	c.mu.Lock()
	delete(c.handlersTrades, key)
	if ch, ok := c.channelsTrades[key]; ok {
		close(ch)
		delete(c.channelsTrades, key)
	}
	c.mu.Unlock()
	if c.conn == nil {
		return nil
	}
	req := map[string]interface{}{
		"op":   "unsubscribe",
		"args": []map[string]string{{"channel": "trades", "instId": instID}},
	}
	msg, _ := json.Marshal(req)
	if err := c.conn.WriteMessage(1, msg); err != nil {
		return err
	}
	return nil
}

// SubscribeTickers subscribes to tickers channel with callback.
func (c *Client) SubscribeTickers(ctx context.Context, instID string, handler func(models.WSTickerMsg)) error {
	key := "tickers:" + instID
	c.mu.Lock()
	if handler != nil {
		c.handlersTickers[key] = append(c.handlersTickers[key], handler)
	}
	c.subscriptions = append(c.subscriptions, subscription{channel: "tickers", instID: instID})
	c.mu.Unlock()
	req := map[string]interface{}{
		"op":   "subscribe",
		"args": []map[string]string{{"channel": "tickers", "instId": instID}},
	}
	msg, _ := json.Marshal(req)
	if err := c.conn.WriteMessage(1, msg); err != nil {
		return err
	}
	return nil
}

func (c *Client) SubscribeTickersChan(ctx context.Context, instID string) (<-chan models.WSTickerMsg, error) {
	key := "tickers:" + instID
	c.mu.Lock()
	ch := make(chan models.WSTickerMsg, 100)
	c.channelsTickers[key] = ch
	c.subscriptions = append(c.subscriptions, subscription{channel: "tickers", instID: instID})
	c.mu.Unlock()
	req := map[string]interface{}{
		"op":   "subscribe",
		"args": []map[string]string{{"channel": "tickers", "instId": instID}},
	}
	msg, _ := json.Marshal(req)
	if err := c.conn.WriteMessage(1, msg); err != nil {
		return nil, err
	}
	return ch, nil
}

func (c *Client) UnsubscribeTickers(ctx context.Context, instID string) error {
	key := "tickers:" + instID
	c.mu.Lock()
	delete(c.handlersTickers, key)
	if ch, ok := c.channelsTickers[key]; ok {
		close(ch)
		delete(c.channelsTickers, key)
	}
	c.mu.Unlock()
	if c.conn == nil {
		return nil
	}
	req := map[string]interface{}{
		"op":   "unsubscribe",
		"args": []map[string]string{{"channel": "tickers", "instId": instID}},
	}
	msg, _ := json.Marshal(req)
	if err := c.conn.WriteMessage(1, msg); err != nil {
		return err
	}
	return nil
}

// SubscribeOrderBook subscribes to order book channel with callback.
func (c *Client) SubscribeOrderBook(ctx context.Context, channel, instID string, handler func(models.WSOrderBookMsg)) error {
	key := channel + ":" + instID
	c.mu.Lock()
	if handler != nil {
		c.handlersOrderBook[key] = append(c.handlersOrderBook[key], handler)
	}
	c.subscriptions = append(c.subscriptions, subscription{channel: channel, instID: instID})
	c.mu.Unlock()
	req := map[string]interface{}{
		"op":   "subscribe",
		"args": []map[string]string{{"channel": channel, "instId": instID}},
	}
	msg, _ := json.Marshal(req)
	if err := c.conn.WriteMessage(1, msg); err != nil {
		return err
	}
	return nil
}

func (c *Client) SubscribeOrderBookChan(ctx context.Context, channel, instID string) (<-chan models.WSOrderBookMsg, error) {
	key := channel + ":" + instID
	c.mu.Lock()
	ch := make(chan models.WSOrderBookMsg, 100)
	c.channelsOrderBook[key] = ch
	c.subscriptions = append(c.subscriptions, subscription{channel: channel, instID: instID})
	c.mu.Unlock()
	req := map[string]interface{}{
		"op":   "subscribe",
		"args": []map[string]string{{"channel": channel, "instId": instID}},
	}
	msg, _ := json.Marshal(req)
	if err := c.conn.WriteMessage(1, msg); err != nil {
		return nil, err
	}
	return ch, nil
}

func (c *Client) UnsubscribeOrderBook(ctx context.Context, channel, instID string) error {
	key := channel + ":" + instID
	c.mu.Lock()
	delete(c.handlersOrderBook, key)
	if ch, ok := c.channelsOrderBook[key]; ok {
		close(ch)
		delete(c.channelsOrderBook, key)
	}
	c.mu.Unlock()
	if c.conn == nil {
		return nil
	}
	req := map[string]interface{}{
		"op":   "unsubscribe",
		"args": []map[string]string{{"channel": channel, "instId": instID}},
	}
	msg, _ := json.Marshal(req)
	if err := c.conn.WriteMessage(1, msg); err != nil {
		return err
	}
	return nil
}

// SubscribeFundingRate subscribes to funding rate channel with callback.
func (c *Client) SubscribeFundingRate(ctx context.Context, instID string, handler func(models.WSFundingRateMsg)) error {
	key := "fundingrate:" + instID
	c.mu.Lock()
	if handler != nil {
		c.handlersFundingRate[key] = append(c.handlersFundingRate[key], handler)
	}
	c.subscriptions = append(c.subscriptions, subscription{channel: "fundingrate", instID: instID})
	c.mu.Unlock()
	req := map[string]interface{}{
		"op":   "subscribe",
		"args": []map[string]string{{"channel": "fundingrate", "instId": instID}},
	}
	msg, _ := json.Marshal(req)
	if err := c.conn.WriteMessage(1, msg); err != nil {
		return err
	}
	return nil
}

func (c *Client) SubscribeFundingRateChan(ctx context.Context, instID string) (<-chan models.WSFundingRateMsg, error) {
	key := "fundingrate:" + instID
	c.mu.Lock()
	ch := make(chan models.WSFundingRateMsg, 100)
	c.channelsFundingRate[key] = ch
	c.subscriptions = append(c.subscriptions, subscription{channel: "fundingrate", instID: instID})
	c.mu.Unlock()
	req := map[string]interface{}{
		"op":   "subscribe",
		"args": []map[string]string{{"channel": "fundingrate", "instId": instID}},
	}
	msg, _ := json.Marshal(req)
	if err := c.conn.WriteMessage(1, msg); err != nil {
		return nil, err
	}
	return ch, nil
}

func (c *Client) UnsubscribeFundingRate(ctx context.Context, instID string) error {
	key := "fundingrate:" + instID
	c.mu.Lock()
	delete(c.handlersFundingRate, key)
	if ch, ok := c.channelsFundingRate[key]; ok {
		close(ch)
		delete(c.channelsFundingRate, key)
	}
	c.mu.Unlock()
	if c.conn == nil {
		return nil
	}
	req := map[string]interface{}{
		"op":   "unsubscribe",
		"args": []map[string]string{{"channel": "fundingrate", "instId": instID}},
	}
	msg, _ := json.Marshal(req)
	if err := c.conn.WriteMessage(1, msg); err != nil {
		return err
	}
	return nil
}
