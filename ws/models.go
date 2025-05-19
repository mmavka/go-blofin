/**
 * @file: models.go
 * @description: Модели для WebSocket BloFin
 * @dependencies: -
 * @created: 2025-05-19
 */

package ws

import (
	rest "github.com/mmavka/go-blofin/rest"
)

// LoginRequest для аутентификации

type LoginRequest struct {
	Op   string      `json:"op"`
	Args []LoginArgs `json:"args"`
}

type LoginArgs struct {
	ApiKey     string `json:"apiKey"`
	Passphrase string `json:"passphrase"`
	Timestamp  string `json:"timestamp"`
	Sign       string `json:"sign"`
	Nonce      string `json:"nonce"`
}

// SubscribeRequest для подписки на каналы

type SubscribeRequest struct {
	Op   string        `json:"op"`
	Args []ChannelArgs `json:"args"`
}

type ChannelArgs struct {
	Channel string `json:"channel"`
	InstId  string `json:"instId"`
}

// UnsubscribeRequest для отписки от каналов

type UnsubscribeRequest struct {
	Op   string        `json:"op"`
	Args []ChannelArgs `json:"args"`
}

// EventResponse для обработки событий login/subscribe/unsubscribe

type EventResponse struct {
	Event string      `json:"event"`
	Arg   ChannelArgs `json:"arg"`
	Code  string      `json:"code"`
	Msg   string      `json:"msg"`
}

// TradeWS — данные о сделке из WS trades channel

type TradeWS struct {
	InstID  string `json:"instId"`
	TradeID string `json:"tradeId"`
	Price   string `json:"price"`
	Size    string `json:"size"`
	Side    string `json:"side"`
	Ts      string `json:"ts"`
}

type TradeWSMessage struct {
	Arg  ChannelArgs `json:"arg"`
	Data []TradeWS   `json:"data"`
}

// CandleWSMessage — push-сообщение для канала свечей

type CandleWSMessage struct {
	Arg  ChannelArgs `json:"arg"`
	Data []Candle    `json:"data"`
}

type Candle = rest.Candle

// OrderBookWSData — данные стакана из WS order book channel

type OrderBookWSData struct {
	Asks      [][]string `json:"asks"`
	Bids      [][]string `json:"bids"`
	Ts        string     `json:"ts"`
	PrevSeqId string     `json:"prevSeqId"`
	SeqId     string     `json:"seqId"`
}

// OrderBookWSMessage — push-сообщение для канала стакана

type OrderBookWSMessage struct {
	Arg    ChannelArgs     `json:"arg"`
	Action string          `json:"action"`
	Data   OrderBookWSData `json:"data"`
}

// TickerWS — данные тикера из WS tickers channel

type TickerWS struct {
	InstID         string `json:"instId"`
	Last           string `json:"last"`
	LastSize       string `json:"lastSize"`
	AskPrice       string `json:"askPrice"`
	AskSize        string `json:"askSize"`
	BidPrice       string `json:"bidPrice"`
	BidSize        string `json:"bidSize"`
	Open24h        string `json:"open24h"`
	High24h        string `json:"high24h"`
	Low24h         string `json:"low24h"`
	VolCurrency24h string `json:"volCurrency24h"`
	Vol24h         string `json:"vol24h"`
	Ts             string `json:"ts"`
}

type TickerWSMessage struct {
	Arg  ChannelArgs `json:"arg"`
	Data []TickerWS  `json:"data"`
}

// FundingRateWS — данные funding rate из WS funding-rate channel

type FundingRateWS struct {
	InstID      string `json:"instId"`
	FundingRate string `json:"fundingRate"`
	FundingTime string `json:"fundingTime"`
}

type FundingRateWSMessage struct {
	Arg  ChannelArgs     `json:"arg"`
	Data []FundingRateWS `json:"data"`
}
