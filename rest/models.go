/**
 * @file: models.go
 * @description: Модели запросов и ответов для публичных REST методов BloFin
 * @dependencies: client.go
 * @created: 2025-05-19
 */

package rest

import (
	"fmt"

	"github.com/goccy/go-json"
)

// Instrument публичная информация об инструменте
// Пример: https://docs.blofin.com/index.html#get-instruments

type Instrument struct {
	InstID        string `json:"instId"`
	BaseCurrency  string `json:"baseCurrency"`
	QuoteCurrency string `json:"quoteCurrency"`
	ContractValue string `json:"contractValue"`
	ListTime      string `json:"listTime"`
	ExpireTime    string `json:"expireTime"`
	MaxLeverage   string `json:"maxLeverage"`
	MinSize       string `json:"minSize"`
	LotSize       string `json:"lotSize"`
	TickSize      string `json:"tickSize"`
	InstType      string `json:"instType"`
	ContractType  string `json:"contractType"`
	MaxLimitSize  string `json:"maxLimitSize"`
	MaxMarketSize string `json:"maxMarketSize"`
	State         string `json:"state"`
}

// BaseResponse базовая структура ответа API
type BaseResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

// InstrumentsResponse представляет ответ на запрос инструментов
type InstrumentsResponse struct {
	BaseResponse
	Data []Instrument `json:"data"`
}

// Ticker публичная информация о тикере
// Пример: https://docs.blofin.com/index.html#get-tickers

type Ticker struct {
	InstID         string `json:"instId"`
	Last           string `json:"last"`
	LastSize       string `json:"lastSize"`
	AskPrice       string `json:"askPrice"`
	AskSize        string `json:"askSize"`
	BidPrice       string `json:"bidPrice"`
	BidSize        string `json:"bidSize"`
	High24h        string `json:"high24h"`
	Open24h        string `json:"open24h"`
	Low24h         string `json:"low24h"`
	VolCurrency24h string `json:"volCurrency24h"`
	Vol24h         string `json:"vol24h"`
	Ts             string `json:"ts"`
}

// TickersResponse представляет ответ на запрос тикеров
type TickersResponse struct {
	BaseResponse
	Data []Ticker `json:"data"`
}

// OrderBook публичная информация о стакане
// Пример: https://docs.blofin.com/index.html#get-order-book

type OrderBook struct {
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
	Ts   string     `json:"ts"`
}

// OrderBookResponse представляет ответ на запрос стакана
type OrderBookResponse struct {
	BaseResponse
	Data []OrderBook `json:"data"`
}

// Balance информация о балансе
// Пример: https://docs.blofin.com/index.html#get-balance

type Balance struct {
	Currency  string `json:"currency"`
	Balance   string `json:"balance"`
	Available string `json:"available"`
	Frozen    string `json:"frozen"`
	Bonus     string `json:"bonus"`
}

// GetBalanceResponse представляет ответ на запрос баланса
type GetBalanceResponse struct {
	BaseResponse
	Data []Balance `json:"data"`
}

// Trade публичная информация о сделке
// Пример: https://docs.blofin.com/index.html#get-trades

type Trade struct {
	TradeID string `json:"tradeId"`
	InstID  string `json:"instId"`
	Price   string `json:"price"`
	Size    string `json:"size"`
	Side    string `json:"side"`
	Ts      string `json:"ts"`
}

// TradesResponse представляет ответ на запрос сделок
type TradesResponse struct {
	BaseResponse
	Data []Trade `json:"data"`
}

// MarkPrice информация о mark/index price
// Пример: https://docs.blofin.com/index.html#get-mark-price

type MarkPrice struct {
	InstID     string `json:"instId"`
	IndexPrice string `json:"indexPrice"`
	MarkPrice  string `json:"markPrice"`
	Ts         string `json:"ts"`
}

// MarkPriceResponse представляет ответ на запрос mark price
type MarkPriceResponse struct {
	BaseResponse
	Data []MarkPrice `json:"data"`
}

// FundingRate информация о ставке финансирования
// Пример: https://docs.blofin.com/index.html#get-funding-rate

type FundingRate struct {
	InstID      string `json:"instId"`
	FundingRate string `json:"fundingRate"`
	FundingTime string `json:"fundingTime"`
}

// FundingRateResponse представляет ответ на запрос funding rate
type FundingRateResponse struct {
	BaseResponse
	Data []FundingRate `json:"data"`
}

// Candle информация о свече (каждая свеча — массив строк)
// Пример: https://docs.blofin.com/index.html#get-candlesticks

type Candle struct {
	Ts               string `json:"ts"`
	Open             string `json:"open"`
	High             string `json:"high"`
	Low              string `json:"low"`
	Close            string `json:"close"`
	Vol              string `json:"vol"`
	VolCurrency      string `json:"volCurrency"`
	VolCurrencyQuote string `json:"volCurrencyQuote"`
	Confirm          string `json:"confirm"`
}

// UnmarshalJSON реализует парсинг массива строк в структуру Candle
func (c *Candle) UnmarshalJSON(data []byte) error {
	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	if len(arr) < 9 {
		return fmt.Errorf("invalid candle data length: %d", len(arr))
	}
	c.Ts = arr[0]
	c.Open = arr[1]
	c.High = arr[2]
	c.Low = arr[3]
	c.Close = arr[4]
	c.Vol = arr[5]
	c.VolCurrency = arr[6]
	c.VolCurrencyQuote = arr[7]
	c.Confirm = arr[8]
	return nil
}

// CandlesResponse представляет ответ на запрос свечей
type CandlesResponse struct {
	BaseResponse
	Data []Candle `json:"data"`
}

// TransferRequest — запрос на трансфер средств

type TransferRequest struct {
	Currency    string `json:"currency"`
	Amount      string `json:"amount"`
	FromAccount string `json:"fromAccount"`
	ToAccount   string `json:"toAccount"`
	ClientId    string `json:"clientId,omitempty"`
}

// TransferResponse представляет ответ на запрос трансфера
type TransferResponse struct {
	BaseResponse
	Data struct {
		TransferId       string `json:"transferId"`
		ClientTransferId string `json:"clientTransferId"`
	} `json:"data"`
}

// TransferHistoryItem — элемент истории трансферов

type TransferHistoryItem struct {
	TransferId  string `json:"transferId"`
	Currency    string `json:"currency"`
	FromAccount string `json:"fromAccount"`
	ToAccount   string `json:"toAccount"`
	Amount      string `json:"amount"`
	Ts          string `json:"ts"`
	ClientId    string `json:"clientId"`
}

// TransferHistoryResponse представляет ответ на запрос истории трансферов
type TransferHistoryResponse struct {
	BaseResponse
	Data []TransferHistoryItem `json:"data"`
}

// WithdrawHistoryItem — элемент истории выводов средств

type WithdrawHistoryItem struct {
	Currency    string  `json:"currency"`
	Chain       string  `json:"chain"`
	Address     string  `json:"address"`
	Type        string  `json:"type"`
	TxId        string  `json:"txId"`
	Amount      string  `json:"amount"`
	Fee         string  `json:"fee"`
	FeeCurrency string  `json:"feeCurrency"`
	State       string  `json:"state"`
	ClientId    *string `json:"clientId"`
	Ts          string  `json:"ts"`
	Tag         *string `json:"tag"`
	Memo        *string `json:"memo"`
	WithdrawId  string  `json:"withdrawId"`
}

// WithdrawHistoryResponse представляет ответ на запрос истории выводов
type WithdrawHistoryResponse struct {
	BaseResponse
	Data []WithdrawHistoryItem `json:"data"`
}

// DepositHistoryItem — элемент истории депозитов

type DepositHistoryItem struct {
	Currency  string `json:"currency"`
	Chain     string `json:"chain"`
	Address   string `json:"address"`
	Type      string `json:"type"`
	TxId      string `json:"txId"`
	Amount    string `json:"amount"`
	State     string `json:"state"`
	Confirm   string `json:"confirm"`
	Ts        string `json:"ts"`
	DepositId string `json:"depositId"`
}

// DepositHistoryResponse представляет ответ на запрос истории депозитов
type DepositHistoryResponse struct {
	BaseResponse
	Data []DepositHistoryItem `json:"data"`
}

// AccountBalanceResponse представляет ответ на запрос баланса фьючерсного аккаунта
type AccountBalanceResponse struct {
	BaseResponse
	Data struct {
		Ts             string                 `json:"ts"`
		TotalEquity    string                 `json:"totalEquity"`
		IsolatedEquity string                 `json:"isolatedEquity"`
		Details        []AccountBalanceDetail `json:"details"`
	} `json:"data"`
}

// AccountBalanceDetail представляет детальную информацию о балансе по валюте
type AccountBalanceDetail struct {
	Currency              string `json:"currency"`
	Equity                string `json:"equity"`
	Balance               string `json:"balance"`
	Ts                    string `json:"ts"`
	IsolatedEquity        string `json:"isolatedEquity"`
	Available             string `json:"available"`
	AvailableEquity       string `json:"availableEquity"`
	Frozen                string `json:"frozen"`
	OrderFrozen           string `json:"orderFrozen"`
	EquityUsd             string `json:"equityUsd"`
	IsolatedUnrealizedPnl string `json:"isolatedUnrealizedPnl"`
	Bonus                 string `json:"bonus"`
}

// Position представляет информацию о позиции
type Position struct {
	PositionId         string `json:"positionId"`
	InstId             string `json:"instId"`
	InstType           string `json:"instType"`
	MarginMode         string `json:"marginMode"`   // MarginModeCross или MarginModeIsolated
	PositionSide       string `json:"positionSide"` // PositionSideLong, PositionSideShort или PositionSideNet
	Adl                string `json:"adl"`
	Positions          string `json:"positions"`
	AvailablePositions string `json:"availablePositions"`
	AveragePrice       string `json:"averagePrice"`
	Margin             string `json:"margin,omitempty"`
	MarkPrice          string `json:"markPrice"`
	MarginRatio        string `json:"marginRatio"`
	LiquidationPrice   string `json:"liquidationPrice"`
	UnrealizedPnl      string `json:"unrealizedPnl"`
	UnrealizedPnlRatio string `json:"unrealizedPnlRatio"`
	InitialMargin      string `json:"initialMargin,omitempty"`
	MaintenanceMargin  string `json:"maintenanceMargin"`
	CreateTime         string `json:"createTime"`
	UpdateTime         string `json:"updateTime"`
	Leverage           string `json:"leverage"`
}

// GetPositionsResponse представляет ответ на запрос позиций
type GetPositionsResponse struct {
	BaseResponse
	Data []Position `json:"data"`
}

// GetMarginModeResponse представляет ответ на запрос режима маржи
type GetMarginModeResponse struct {
	BaseResponse
	Data struct {
		MarginMode string `json:"marginMode"`
	} `json:"data"`
}

// SetMarginModeRequest представляет запрос на установку режима маржи
type SetMarginModeRequest struct {
	MarginMode string `json:"marginMode"` // MarginModeCross или MarginModeIsolated
}

// GetPositionModeResponse представляет ответ на запрос режима позиций
type GetPositionModeResponse struct {
	BaseResponse
	Data struct {
		PositionMode string `json:"positionMode"`
	} `json:"data"`
}

// SetPositionModeRequest представляет запрос на установку режима позиций
type SetPositionModeRequest struct {
	PositionMode string `json:"positionMode"` // PositionModeNet или PositionModeLongShort
}

// GetLeverageInfoResponse представляет ответ на запрос информации о плече (deprecated)
type GetLeverageInfoResponse struct {
	BaseResponse
	Data struct {
		InstId     string `json:"instId"`
		Leverage   string `json:"leverage"`
		MarginMode string `json:"marginMode"`
	} `json:"data"`
}

// BatchLeverageInfo представляет информацию о плече для инструмента
type BatchLeverageInfo struct {
	InstId       string `json:"instId"`
	Leverage     string `json:"leverage"`
	MarginMode   string `json:"marginMode"`   // MarginModeCross или MarginModeIsolated
	PositionSide string `json:"positionSide"` // PositionSideLong, PositionSideShort или PositionSideNet
}

// GetBatchLeverageInfoResponse представляет ответ на запрос информации о плече для нескольких инструментов
type GetBatchLeverageInfoResponse struct {
	BaseResponse
	Data []BatchLeverageInfo `json:"data"`
}

// SetLeverageRequest представляет запрос на установку плеча
type SetLeverageRequest struct {
	InstId       string `json:"instId"`                 // Instrument ID
	Leverage     string `json:"leverage"`               // Leverage value
	MarginMode   string `json:"marginMode"`             // MarginModeCross или MarginModeIsolated
	PositionSide string `json:"positionSide,omitempty"` // PositionSideLong или PositionSideShort (только для isolated в long/short mode)
}

// SetLeverageResponse представляет ответ на запрос установки плеча
type SetLeverageResponse struct {
	BaseResponse
	Data struct {
		InstId       string `json:"instId"`       // Instrument ID
		Leverage     string `json:"leverage"`     // Leverage value
		MarginMode   string `json:"marginMode"`   // MarginModeCross или MarginModeIsolated
		PositionSide string `json:"positionSide"` // PositionSideLong, PositionSideShort или PositionSideNet
	} `json:"data"`
}

// PlaceOrderRequest запрос на размещение ордера
type PlaceOrderRequest struct {
	InstID         string `json:"instId"`                   // Instrument ID
	MarginMode     string `json:"marginMode"`               // MarginModeCross или MarginModeIsolated
	PositionSide   string `json:"positionSide"`             // PositionSideLong, PositionSideShort или PositionSideNet
	Side           string `json:"side"`                     // OrderSideBuy или OrderSideSell
	OrderType      string `json:"orderType"`                // Тип ордера (market, limit, post_only, fok, ioc)
	Price          string `json:"price"`                    // Цена ордера
	Size           string `json:"size"`                     // Количество контрактов
	ReduceOnly     string `json:"reduceOnly,omitempty"`     // true или false
	ClientOrderId  string `json:"clientOrderId,omitempty"`  // ID ордера клиента
	TpTriggerPrice string `json:"tpTriggerPrice,omitempty"` // Take-profit trigger price
	TpOrderPrice   string `json:"tpOrderPrice,omitempty"`   // Take-profit order price
	SlTriggerPrice string `json:"slTriggerPrice,omitempty"` // Stop-loss trigger price
	SlOrderPrice   string `json:"slOrderPrice,omitempty"`   // Stop-loss order price
	BrokerId       string `json:"brokerId,omitempty"`       // Broker ID
}

// OrderResult результат выполнения ордера
type OrderResult struct {
	OrderId       string `json:"orderId"`       // Order ID
	ClientOrderId string `json:"clientOrderId"` // Client Order ID
	Code          string `json:"code"`          // Код результата
	Msg           string `json:"msg"`           // Сообщение
}

// PlaceOrderResponse представляет ответ на запрос размещения ордера
type PlaceOrderResponse struct {
	BaseResponse
	Data []OrderResult `json:"data"`
}

// BatchOrderResult результат выполнения ордера в пакетном запросе
type BatchOrderResult struct {
	OrderId       string `json:"orderId"`       // Order ID
	ClientOrderId string `json:"clientOrderId"` // Client Order ID
}

// BatchOrdersResponse представляет ответ на запрос размещения нескольких ордеров
type BatchOrdersResponse struct {
	BaseResponse
	Data []BatchOrderResult `json:"data"`
}

// PlaceTPSLOrderRequest запрос на размещение TPSL ордера
type PlaceTPSLOrderRequest struct {
	InstID         string `json:"instId"`                   // Instrument ID
	MarginMode     string `json:"marginMode"`               // MarginModeCross или MarginModeIsolated
	PositionSide   string `json:"positionSide"`             // PositionSideLong, PositionSideShort или PositionSideNet
	Side           string `json:"side"`                     // OrderSideBuy или OrderSideSell
	TpTriggerPrice string `json:"tpTriggerPrice"`           // Take-profit trigger price
	TpOrderPrice   string `json:"tpOrderPrice,omitempty"`   // Take-profit order price
	SlTriggerPrice string `json:"slTriggerPrice,omitempty"` // Stop-loss trigger price
	SlOrderPrice   string `json:"slOrderPrice,omitempty"`   // Stop-loss order price
	Size           string `json:"size"`                     // Количество контрактов
	ReduceOnly     string `json:"reduceOnly,omitempty"`     // true или false
	ClientOrderId  string `json:"clientOrderId,omitempty"`  // ID ордера клиента
	BrokerId       string `json:"brokerId,omitempty"`       // Broker ID
}

// TPSLOrderResult результат выполнения TPSL ордера
type TPSLOrderResult struct {
	TpslId        string `json:"tpslId"`        // TP/SL order ID
	ClientOrderId string `json:"clientOrderId"` // Client Order ID as assigned by the client
	Code          string `json:"code"`          // The code of the event execution result, 0 means success
	Msg           string `json:"msg"`           // Rejection or success message of event execution
}

// PlaceTPSLOrderResponse представляет ответ на запрос размещения TPSL ордера
type PlaceTPSLOrderResponse struct {
	BaseResponse
	Data TPSLOrderResult `json:"data"`
}

// AttachAlgoOrder представляет информацию о прикрепленных TP/SL ордерах
type AttachAlgoOrder struct {
	TpTriggerPrice     string `json:"tpTriggerPrice,omitempty"`
	TpOrderPrice       string `json:"tpOrderPrice,omitempty"`
	TpTriggerPriceType string `json:"tpTriggerPriceType,omitempty"`
	SlTriggerPrice     string `json:"slTriggerPrice,omitempty"`
	SlOrderPrice       string `json:"slOrderPrice,omitempty"`
	SlTriggerPriceType string `json:"slTriggerPriceType,omitempty"`
}

// PlaceAlgoOrderRequest запрос на размещение алго-ордера
type PlaceAlgoOrderRequest struct {
	InstID           string            `json:"instId"`
	MarginMode       string            `json:"marginMode"`
	PositionSide     string            `json:"positionSide"`
	Side             string            `json:"side"`
	Size             string            `json:"size"`
	ClientOrderId    string            `json:"clientOrderId,omitempty"`
	OrderPrice       string            `json:"orderPrice,omitempty"`
	OrderType        string            `json:"orderType"`
	TriggerPrice     string            `json:"triggerPrice"`
	TriggerPriceType string            `json:"triggerPriceType,omitempty"`
	ReduceOnly       string            `json:"reduceOnly,omitempty"`
	BrokerId         string            `json:"brokerId,omitempty"`
	AttachAlgoOrders []AttachAlgoOrder `json:"attachAlgoOrders,omitempty"`
}

// AlgoOrderResult результат выполнения алго-ордера
type AlgoOrderResult struct {
	AlgoId        string `json:"algoId"`
	ClientOrderId string `json:"clientOrderId"`
	Code          string `json:"code"`
	Msg           string `json:"msg"`
}

// PlaceAlgoOrderResponse представляет ответ на запрос размещения алго-ордера
type PlaceAlgoOrderResponse struct {
	BaseResponse
	Data AlgoOrderResult `json:"data"`
}

// CancelOrderRequest запрос на отмену ордера
type CancelOrderRequest struct {
	InstID        string `json:"instId,omitempty"`
	OrderId       string `json:"orderId"`
	ClientOrderId string `json:"clientOrderId,omitempty"`
}

// CancelOrderResult результат отмены ордера
type CancelOrderResult struct {
	OrderId       string `json:"orderId"`       // Order ID
	ClientOrderId string `json:"clientOrderId"` // Client Order ID
	Code          string `json:"code"`          // The code of the event execution result, 0 means success
	Msg           string `json:"msg"`           // Rejection or success message of event execution
}

// CancelOrderResponse представляет ответ на запрос отмены ордера
type CancelOrderResponse struct {
	BaseResponse
	Data CancelOrderResult `json:"data"`
}

// CancelBatchOrdersResponse представляет ответ на запрос отмены нескольких ордеров
type CancelBatchOrdersResponse struct {
	BaseResponse
	Data []CancelOrderResult `json:"data"`
}

// CancelTPSLOrderRequest запрос на отмену TPSL ордера
type CancelTPSLOrderRequest struct {
	InstID        string `json:"instId,omitempty"`
	TpslId        string `json:"tpslId,omitempty"`
	ClientOrderId string `json:"clientOrderId,omitempty"`
}

// CancelTPSLOrderResult результат отмены TPSL ордера
type CancelTPSLOrderResult struct {
	TpslId        string `json:"tpslId"`        // TP/SL order ID
	ClientOrderId string `json:"clientOrderId"` // Client Order ID
	Code          string `json:"code"`          // The code of the event execution result, 0 means success
	Msg           string `json:"msg"`           // Rejection or success message of event execution
}

// CancelTPSLOrderResponse представляет ответ на запрос отмены TPSL ордера
type CancelTPSLOrderResponse struct {
	BaseResponse
	Data []CancelTPSLOrderResult `json:"data"`
}

// CancelAlgoOrderRequest запрос на отмену алго-ордера
type CancelAlgoOrderRequest struct {
	InstID        string `json:"instId,omitempty"`
	AlgoId        string `json:"algoId,omitempty"`
	ClientOrderId string `json:"clientOrderId,omitempty"`
}

// CancelAlgoOrderResult результат отмены алго-ордера
type CancelAlgoOrderResult struct {
	AlgoId        string `json:"algoId"`        // Algo order ID
	ClientOrderId string `json:"clientOrderId"` // Client Order ID
	Code          string `json:"code"`          // The code of the event execution result, 0 means success
	Msg           string `json:"msg"`           // Rejection or success message of event execution
}

// CancelAlgoOrderResponse представляет ответ на запрос отмены алго-ордера
type CancelAlgoOrderResponse struct {
	BaseResponse
	Data CancelAlgoOrderResult `json:"data"`
}

// PendingOrder представляет информацию об активном ордере
type PendingOrder struct {
	OrderId           string `json:"orderId"`           // Order ID
	ClientOrderId     string `json:"clientOrderId"`     // Client Order ID
	InstId            string `json:"instId"`            // Instrument ID
	MarginMode        string `json:"marginMode"`        // Margin mode
	PositionSide      string `json:"positionSide"`      // Position side
	Side              string `json:"side"`              // Order side
	OrderType         string `json:"orderType"`         // Order type
	Price             string `json:"price"`             // Price
	Size              string `json:"size"`              // Number of contracts
	ReduceOnly        string `json:"reduceOnly"`        // Whether orders can only reduce in position size
	Leverage          string `json:"leverage"`          // Leverage
	State             string `json:"state"`             // State
	FilledSize        string `json:"filledSize"`        // Accumulated fill quantity
	FilledAmount      string `json:"filled_amount"`     // Filled amount
	AveragePrice      string `json:"averagePrice"`      // Average filled price
	Fee               string `json:"fee"`               // Fee and rebate
	Pnl               string `json:"pnl"`               // Profit and loss
	CreateTime        string `json:"createTime"`        // Creation time
	UpdateTime        string `json:"updateTime"`        // Update time
	OrderCategory     string `json:"orderCategory"`     // Order category
	TpTriggerPrice    string `json:"tpTriggerPrice"`    // Take-profit trigger price
	TpOrderPrice      string `json:"tpOrderPrice"`      // Take-profit order price
	SlTriggerPrice    string `json:"slTriggerPrice"`    // Stop-loss trigger price
	SlOrderPrice      string `json:"slOrderPrice"`      // Stop-loss order price
	AlgoClientOrderId string `json:"algoClientOrderId"` // Algo client order ID
	AlgoId            string `json:"algoId"`            // Algo ID
	BrokerId          string `json:"brokerId"`          // Broker ID
}

// GetPendingOrdersResponse представляет ответ на запрос активных ордеров
type GetPendingOrdersResponse struct {
	BaseResponse
	Data []PendingOrder `json:"data"`
}

// PendingTPSLOrder представляет информацию об активном TPSL ордере
type PendingTPSLOrder struct {
	TpslId         string `json:"tpslId"`         // TP/SL order ID
	InstId         string `json:"instId"`         // Instrument ID
	MarginMode     string `json:"marginMode"`     // Margin mode
	PositionSide   string `json:"positionSide"`   // Position side
	Side           string `json:"side"`           // Order side
	TpTriggerPrice string `json:"tpTriggerPrice"` // Take-profit trigger price
	TpOrderPrice   string `json:"tpOrderPrice"`   // Take-profit order price
	SlTriggerPrice string `json:"slTriggerPrice"` // Stop-loss trigger price
	SlOrderPrice   string `json:"slOrderPrice"`   // Stop-loss order price
	Size           string `json:"size"`           // Number of contracts
	State          string `json:"state"`          // State
	Leverage       string `json:"leverage"`       // Leverage
	ReduceOnly     string `json:"reduceOnly"`     // Whether orders can only reduce in position size
	ActualSize     string `json:"actualSize"`     // Actual order quantity
	ClientOrderId  string `json:"clientOrderId"`  // Client Order ID
	CreateTime     string `json:"createTime"`     // Creation time
	BrokerId       string `json:"brokerId"`       // Broker ID
}

// GetPendingTPSLOrdersResponse представляет ответ на запрос активных TPSL ордеров
type GetPendingTPSLOrdersResponse struct {
	BaseResponse
	Data []PendingTPSLOrder `json:"data"`
}

// PendingAlgoOrder представляет информацию об активном алго-ордере
type PendingAlgoOrder struct {
	AlgoId           string            `json:"algoId"`           // Algo order ID
	ClientOrderId    string            `json:"clientOrderId"`    // Client Order ID
	InstId           string            `json:"instId"`           // Instrument ID
	MarginMode       string            `json:"marginMode"`       // Margin mode
	PositionSide     string            `json:"positionSide"`     // Position side
	Side             string            `json:"side"`             // Order side
	OrderType        string            `json:"orderType"`        // Algo type
	Size             string            `json:"size"`             // Number of contracts
	Leverage         string            `json:"leverage"`         // Leverage
	State            string            `json:"state"`            // State
	TriggerPrice     string            `json:"triggerPrice"`     // Trigger price
	TriggerPriceType string            `json:"triggerPriceType"` // Trigger price type
	BrokerId         string            `json:"brokerId"`         // Broker ID
	CreateTime       string            `json:"createTime"`       // Creation time
	AttachAlgoOrders []AttachAlgoOrder `json:"attachAlgoOrders"` // Attached SL/TP orders
}

// GetPendingAlgoOrdersResponse представляет ответ на запрос активных алго-ордеров
type GetPendingAlgoOrdersResponse struct {
	BaseResponse
	Data []PendingAlgoOrder `json:"data"`
}

// ClosePositionRequest запрос на закрытие позиции
type ClosePositionRequest struct {
	InstID        string `json:"instId"`        // Instrument ID
	MarginMode    string `json:"marginMode"`    // Margin mode (cross/isolated)
	PositionSide  string `json:"positionSide"`  // Position side (net/long/short)
	ClientOrderId string `json:"clientOrderId"` // Client Order ID
	BrokerId      string `json:"brokerId"`      // Broker ID
}

// ClosePositionResponse представляет ответ на запрос закрытия позиции
type ClosePositionResponse struct {
	BaseResponse
	Data struct {
		InstID        string `json:"instId"`        // Instrument ID
		PositionSide  string `json:"positionSide"`  // Position side
		ClientOrderId string `json:"clientOrderId"` // Client Order ID
	} `json:"data"`
}

// OrderHistory представляет информацию об ордере в истории
type OrderHistory struct {
	OrderId            string `json:"orderId"`            // Order ID
	ClientOrderId      string `json:"clientOrderId"`      // Client Order ID
	InstId             string `json:"instId"`             // Instrument ID
	MarginMode         string `json:"marginMode"`         // Margin mode
	PositionSide       string `json:"positionSide"`       // Position side
	Side               string `json:"side"`               // Order side
	OrderType          string `json:"orderType"`          // Order type
	Price              string `json:"price"`              // Price
	Size               string `json:"size"`               // Number of contracts
	ReduceOnly         string `json:"reduceOnly"`         // Whether orders can only reduce in position size
	Leverage           string `json:"leverage"`           // Leverage
	State              string `json:"state"`              // State
	FilledSize         string `json:"filledSize"`         // Accumulated fill quantity
	Pnl                string `json:"pnl"`                // Profit and loss
	AveragePrice       string `json:"averagePrice"`       // Average filled price
	Fee                string `json:"fee"`                // Fee and rebate
	CreateTime         string `json:"createTime"`         // Creation time
	UpdateTime         string `json:"updateTime"`         // Update time
	OrderCategory      string `json:"orderCategory"`      // Order category
	TpTriggerPrice     string `json:"tpTriggerPrice"`     // Take-profit trigger price
	TpOrderPrice       string `json:"tpOrderPrice"`       // Take-profit order price
	SlTriggerPrice     string `json:"slTriggerPrice"`     // Stop-loss trigger price
	SlOrderPrice       string `json:"slOrderPrice"`       // Stop-loss order price
	CancelSource       string `json:"cancelSource"`       // Type of the cancellation source
	CancelSourceReason string `json:"cancelSourceReason"` // Reason for the cancellation
	AlgoClientOrderId  string `json:"algoClientOrderId"`  // Algo client order ID
	AlgoId             string `json:"algoId"`             // Algo ID
	BrokerId           string `json:"brokerId"`           // Broker ID
}

// GetOrderHistoryResponse представляет ответ на запрос истории ордеров
type GetOrderHistoryResponse struct {
	BaseResponse
	Data []OrderHistory `json:"data"`
}

// TPSLOrderHistory represents information about a TPSL order in history
type TPSLOrderHistory struct {
	TpslId         string `json:"tpslId"`
	ClientOrderId  string `json:"clientOrderId"`
	InstId         string `json:"instId"`
	MarginMode     string `json:"marginMode"`
	PositionSide   string `json:"positionSide"`
	Side           string `json:"side"`
	OrderType      string `json:"orderType"`
	Size           string `json:"size"`
	ReduceOnly     string `json:"reduceOnly"`
	Leverage       string `json:"leverage"`
	State          string `json:"state"`
	ActualSize     string `json:"actualSize"`
	TriggerType    string `json:"triggerType"`
	OrderCategory  string `json:"orderCategory"`
	TpTriggerPrice string `json:"tpTriggerPrice"`
	TpOrderPrice   string `json:"tpOrderPrice"`
	SlTriggerPrice string `json:"slTriggerPrice"`
	SlOrderPrice   string `json:"slOrderPrice"`
	CreateTime     string `json:"createTime"`
	BrokerId       string `json:"brokerId"`
}

// GetTPSLOrderHistoryResponse represents the response structure for TPSL order history request
type GetTPSLOrderHistoryResponse struct {
	BaseResponse
	Data []TPSLOrderHistory `json:"data"`
}

// AlgoOrderHistory represents information about an algo order in history
type AlgoOrderHistory struct {
	AlgoId           string            `json:"algoId"`           // Algo order ID
	ClientOrderId    string            `json:"clientOrderId"`    // Client Order ID
	InstId           string            `json:"instId"`           // Instrument ID
	MarginMode       string            `json:"marginMode"`       // Margin mode
	PositionSide     string            `json:"positionSide"`     // Position side
	Side             string            `json:"side"`             // Order side
	ReduceOnly       string            `json:"reduceOnly"`       // Whether orders can only reduce in position size
	OrderType        string            `json:"orderType"`        // Algo type
	Size             string            `json:"size"`             // Number of contracts
	Leverage         string            `json:"leverage"`         // Leverage
	State            string            `json:"state"`            // State
	ActualSize       string            `json:"actualSize"`       // Actual order quantity
	CreateTime       string            `json:"createTime"`       // Creation time
	TriggerPrice     string            `json:"triggerPrice"`     // Trigger price
	TriggerPriceType string            `json:"triggerPriceType"` // Trigger price type
	BrokerId         string            `json:"brokerId"`         // Broker ID
	AttachAlgoOrders []AttachAlgoOrder `json:"attachAlgoOrders"` // Attached SL/TP orders
}

// GetAlgoOrderHistoryResponse represents the response structure for algo order history request
type GetAlgoOrderHistoryResponse struct {
	BaseResponse
	Data []AlgoOrderHistory `json:"data"`
}

// TradeHistory represents information about a trade in history
type TradeHistory struct {
	InstId       string `json:"instId"`       // Instrument ID
	TradeId      string `json:"tradeId"`      // Trade ID
	OrderId      string `json:"orderId"`      // Order ID
	FillPrice    string `json:"fillPrice"`    // Filled price
	FillSize     string `json:"fillSize"`     // Filled quantity
	FillPnl      string `json:"fillPnl"`      // Last filled profit and loss
	PositionSide string `json:"positionSide"` // Position side
	Side         string `json:"side"`         // Order side
	Fee          string `json:"fee"`          // Fee
	Ts           string `json:"ts"`           // Data generation time
	BrokerId     string `json:"brokerId"`     // Broker ID
}

// GetTradeHistoryResponse represents the response structure for trade history request
type GetTradeHistoryResponse struct {
	BaseResponse
	Data []TradeHistory `json:"data"`
}

// GetOrderPriceRangeResponse represents the response structure for the order price range request
type GetOrderPriceRangeResponse struct {
	BaseResponse
	Data struct {
		MaxPrice string `json:"maxPrice"`
		MinPrice string `json:"minPrice"`
	} `json:"data"`
}
