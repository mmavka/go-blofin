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

// PositionsRequest represents the request for positions channel
type PositionsRequest struct {
	Op   string `json:"op"`
	Args []struct {
		Channel string `json:"channel"`
		InstId  string `json:"instId,omitempty"`
	} `json:"args"`
}

// PositionsResponse represents the response for positions channel
type PositionsResponse struct {
	Event string `json:"event"`
	Arg   struct {
		Channel string `json:"channel"`
		InstId  string `json:"instId,omitempty"`
	} `json:"arg"`
	Code string `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

// PositionsData represents the push data for positions channel
type PositionsData struct {
	Arg struct {
		Channel string `json:"channel"`
	} `json:"arg"`
	Data []Position `json:"data"`
}

// Position represents a single position
type Position struct {
	InstType           string `json:"instType"`
	InstId             string `json:"instId"`
	MarginMode         string `json:"marginMode"`
	PositionId         string `json:"positionId"`
	PositionSide       string `json:"positionSide"`
	Positions          string `json:"positions"`
	AvailablePositions string `json:"availablePositions"`
	AveragePrice       string `json:"averagePrice"`
	UnrealizedPnl      string `json:"unrealizedPnl"`
	UnrealizedPnlRatio string `json:"unrealizedPnlRatio"`
	Leverage           string `json:"leverage"`
	LiquidationPrice   string `json:"liquidationPrice"`
	MarkPrice          string `json:"markPrice"`
	InitialMargin      string `json:"initialMargin"`
	Margin             string `json:"margin"`
	MarginRatio        string `json:"marginRatio"`
	MaintenanceMargin  string `json:"maintenanceMargin"`
	Adl                string `json:"adl"`
	CreateTime         string `json:"createTime"`
	UpdateTime         string `json:"updateTime"`
}

// OrdersRequest represents the request for orders channel
type OrdersRequest struct {
	Op   string `json:"op"`
	Args []struct {
		Channel string `json:"channel"`
		InstId  string `json:"instId,omitempty"`
	} `json:"args"`
}

// OrdersResponse represents the response for orders channel
type OrdersResponse struct {
	Event string `json:"event"`
	Arg   struct {
		Channel string `json:"channel"`
		InstId  string `json:"instId,omitempty"`
	} `json:"arg"`
	Code string `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

// OrdersData represents the push data for orders channel
type OrdersData struct {
	Action string `json:"action"`
	Arg    struct {
		Channel string `json:"channel"`
	} `json:"arg"`
	Data []Order `json:"data"`
}

// Order represents a single order
type Order struct {
	InstType           string `json:"instType"`
	InstId             string `json:"instId"`
	OrderId            string `json:"orderId"`
	ClientOrderId      string `json:"clientOrderId"`
	Price              string `json:"price"`
	Size               string `json:"size"`
	OrderType          string `json:"orderType"`
	Side               string `json:"side"`
	PositionSide       string `json:"positionSide"`
	MarginMode         string `json:"marginMode"`
	FilledSize         string `json:"filledSize"`
	FilledAmount       string `json:"filledAmount"`
	AveragePrice       string `json:"averagePrice"`
	State              string `json:"state"`
	Leverage           string `json:"leverage"`
	TpTriggerPrice     string `json:"tpTriggerPrice"`
	TpTriggerPriceType string `json:"tpTriggerPriceType"`
	TpOrderPrice       string `json:"tpOrderPrice"`
	SlTriggerPrice     string `json:"slTriggerPrice"`
	SlTriggerPriceType string `json:"slTriggerPriceType"`
	SlOrderPrice       string `json:"slOrderPrice"`
	Fee                string `json:"fee"`
	Pnl                string `json:"pnl"`
	CancelSource       string `json:"cancelSource"`
	OrderCategory      string `json:"orderCategory"`
	CreateTime         string `json:"createTime"`
	UpdateTime         string `json:"updateTime"`
	ReduceOnly         string `json:"reduceOnly"`
	BrokerId           string `json:"brokerId"`
}

// AlgoOrdersRequest represents the request for algo orders channel
type AlgoOrdersRequest struct {
	Op   string `json:"op"`
	Args []struct {
		Channel string `json:"channel"`
		InstId  string `json:"instId,omitempty"`
	} `json:"args"`
}

// AlgoOrdersResponse represents the response for algo orders channel
type AlgoOrdersResponse struct {
	Event string `json:"event"`
	Arg   struct {
		Channel string `json:"channel"`
		InstId  string `json:"instId,omitempty"`
	} `json:"arg"`
	Code string `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

// AlgoOrdersData represents the push data for algo orders channel
type AlgoOrdersData struct {
	Action string `json:"action"`
	Arg    struct {
		Channel string `json:"channel"`
	} `json:"arg"`
	Data []AlgoOrder `json:"data"`
}

// AlgoOrder represents a single algo order
type AlgoOrder struct {
	InstType         string            `json:"instType"`
	InstId           string            `json:"instId"`
	TpslId           string            `json:"tpslId"`
	AlgoId           string            `json:"algoId"`
	ClientOrderId    string            `json:"clientOrderId"`
	Size             string            `json:"size"`
	OrderType        string            `json:"orderType"`
	Side             string            `json:"side"`
	PositionSide     string            `json:"positionSide"`
	MarginMode       string            `json:"marginMode"`
	Leverage         string            `json:"leverage"`
	State            string            `json:"state"`
	TpTriggerPrice   string            `json:"tpTriggerPrice"`
	TpOrderPrice     string            `json:"tpOrderPrice"`
	SlTriggerPrice   string            `json:"slTriggerPrice"`
	SlOrderPrice     string            `json:"slOrderPrice"`
	TriggerPrice     string            `json:"triggerPrice"`
	TriggerPriceType string            `json:"triggerPriceType"`
	OrderPrice       string            `json:"orderPrice"`
	ActualSize       string            `json:"actualSize"`
	ActualSide       string            `json:"actualSide"`
	ReduceOnly       string            `json:"reduceOnly"`
	CancelType       string            `json:"cancelType"`
	CreateTime       string            `json:"createTime"`
	UpdateTime       string            `json:"updateTime"`
	Tag              string            `json:"tag"`
	BrokerId         string            `json:"brokerId"`
	AttachAlgoOrders []AttachAlgoOrder `json:"attachAlgoOrders,omitempty"`
}

// AttachAlgoOrder represents an attached TP/SL order
type AttachAlgoOrder struct {
	TpTriggerPrice     string `json:"tpTriggerPrice"`
	TpTriggerPriceType string `json:"tpTriggerPriceType"`
	TpOrderPrice       string `json:"tpOrderPrice"`
	SlTriggerPriceType string `json:"slTriggerPriceType"`
	SlTriggerPrice     string `json:"slTriggerPrice"`
	SlOrderPrice       string `json:"slOrderPrice"`
}

// AccountRequest represents the request for account channel
type AccountRequest struct {
	Op   string `json:"op"`
	Args []struct {
		Channel string `json:"channel"`
	} `json:"args"`
}

// AccountResponse represents the response for account channel
type AccountResponse struct {
	Event string `json:"event"`
	Arg   struct {
		Channel string `json:"channel"`
	} `json:"arg"`
	Code string `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

// AccountData represents the push data for account channel
type AccountData struct {
	Arg struct {
		Channel string `json:"channel"`
	} `json:"arg"`
	Data Account `json:"data"`
}

// Account represents account information
type Account struct {
	Ts             string          `json:"ts"`
	TotalEquity    string          `json:"totalEquity"`
	IsolatedEquity string          `json:"isolatedEquity"`
	Details        []AccountDetail `json:"details"`
}

// AccountDetail represents detailed asset information for a currency
type AccountDetail struct {
	Currency              string `json:"currency"`
	Equity                string `json:"equity"`
	Balance               string `json:"balance"`
	Ts                    string `json:"ts"`
	IsolatedEquity        string `json:"isolatedEquity"`
	EquityUsd             string `json:"equityUsd"`
	AvailableEquity       string `json:"availableEquity"`
	Available             string `json:"available"`
	Frozen                string `json:"frozen"`
	OrderFrozen           string `json:"orderFrozen"`
	UnrealizedPnl         string `json:"unrealizedPnl"`
	IsolatedUnrealizedPnl string `json:"isolatedUnrealizedPnl"`
	CoinUsdPrice          string `json:"coinUsdPrice"`
	SpotAvailable         string `json:"spotAvailable"`
	Liability             string `json:"liability"`
	BorrowFrozen          string `json:"borrowFrozen"`
	MarginRatio           string `json:"marginRatio"`
}
