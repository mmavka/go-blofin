package ws

import (
	"encoding/json"
	"fmt"

	"github.com/mmavka/go-blofin"
)

// WsBookTickerHandler handle websocket that pushes updates to the best bid or ask price or quantity
type WsBookTickerHandler func(event *WsBookTickerEvent)

// WsBookTickerEvent define websocket best book ticker event
type WsBookTickerEvent struct {
	Symbol    string `json:"symbol"`
	BidPrice  string `json:"bidPrice"`
	BidQty    string `json:"bidQty"`
	AskPrice  string `json:"askPrice"`
	AskQty    string `json:"askQty"`
	Timestamp int64  `json:"timestamp"`
}

// WsBookTickerServe serve websocket that pushes updates to the best bid or ask price or quantity
func WsBookTickerServe(symbol string, handler WsBookTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/market/%s/book-ticker", blofin.DefaultWSPublic, symbol)
	config := &WsConfig{
		Endpoint: endpoint,
	}
	wsHandler := func(message []byte) {
		var event WsBookTickerEvent
		if err := json.Unmarshal(message, &event); err != nil {
			errHandler(err)
			return
		}
		handler(&event)
	}
	return WsServe(config, wsHandler, errHandler)
}

// WsKlineHandler handle websocket kline event
type WsKlineHandler func(event *WsKlineEvent)

// WsKlineEvent define websocket kline event
type WsKlineEvent struct {
	Symbol    string `json:"symbol"`
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
	Interval  string `json:"interval"`
	Open      string `json:"open"`
	High      string `json:"high"`
	Low       string `json:"low"`
	Close     string `json:"close"`
	Volume    string `json:"volume"`
}

// WsKlineServe serve websocket kline handler with a symbol and interval like 15m, 30s
func WsKlineServe(symbol string, interval string, handler WsKlineHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/market/%s/kline/%s", blofin.DefaultWSPublic, symbol, interval)
	config := &WsConfig{
		Endpoint: endpoint,
	}
	wsHandler := func(message []byte) {
		var event WsKlineEvent
		if err := json.Unmarshal(message, &event); err != nil {
			errHandler(err)
			return
		}
		handler(&event)
	}
	return WsServe(config, wsHandler, errHandler)
}

// WsTradeHandler handle websocket trade event
type WsTradeHandler func(event *WsTradeEvent)

// WsTradeEvent define websocket trade event
type WsTradeEvent struct {
	Symbol   string `json:"symbol"`
	TradeID  string `json:"tradeId"`
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
	Time     int64  `json:"time"`
	IsBuyer  bool   `json:"isBuyer"`
	IsMaker  bool   `json:"isMaker"`
}

// WsTradeServe serve websocket that pushes trade information
func WsTradeServe(symbol string, handler WsTradeHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/market/%s/trade", blofin.DefaultWSPublic, symbol)
	config := &WsConfig{
		Endpoint: endpoint,
	}
	wsHandler := func(message []byte) {
		var event WsTradeEvent
		if err := json.Unmarshal(message, &event); err != nil {
			errHandler(err)
			return
		}
		handler(&event)
	}
	return WsServe(config, wsHandler, errHandler)
}
