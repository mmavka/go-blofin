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

type InstrumentsResponse struct {
	Code string       `json:"code"`
	Msg  string       `json:"msg"`
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

type TickersResponse struct {
	Code string   `json:"code"`
	Msg  string   `json:"msg"`
	Data []Ticker `json:"data"`
}

// OrderBook публичная информация о стакане
// Пример: https://docs.blofin.com/index.html#get-order-book

type OrderBook struct {
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
	Ts   string     `json:"ts"`
}

type OrderBookResponse struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
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

type GetBalanceResponse struct {
	Code string    `json:"code"`
	Msg  string    `json:"msg"`
	Data []Balance `json:"data"`
}

// Position информация о позиции
// Пример: https://docs.blofin.com/index.html#get-positions

type Position struct {
	InstID string `json:"instId"`
	Pos    string `json:"pos"`
	Side   string `json:"side"`
}

type PositionsResponse struct {
	Code string     `json:"code"`
	Msg  string     `json:"msg"`
	Data []Position `json:"data"`
}

// OrderRequest запрос на размещение ордера

type OrderRequest struct {
	InstID     string `json:"instId"`
	Side       string `json:"side"`
	OrderType  string `json:"orderType"`
	Price      string `json:"price,omitempty"`
	Size       string `json:"size"`
	MarginMode string `json:"marginMode,omitempty"`
}

type OrderResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
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

type TradesResponse struct {
	Code string  `json:"code"`
	Msg  string  `json:"msg"`
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

type MarkPriceResponse struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data []MarkPrice `json:"data"`
}

// FundingRate информация о ставке финансирования
// Пример: https://docs.blofin.com/index.html#get-funding-rate

type FundingRate struct {
	InstID      string `json:"instId"`
	FundingRate string `json:"fundingRate"`
	FundingTime string `json:"fundingTime"`
}

type FundingRateResponse struct {
	Code string        `json:"code"`
	Msg  string        `json:"msg"`
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

type CandlesResponse struct {
	Code string   `json:"code"`
	Msg  string   `json:"msg"`
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

type TransferResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
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

type TransferHistoryResponse struct {
	Code string                `json:"code"`
	Msg  string                `json:"msg"`
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

type WithdrawHistoryResponse struct {
	Code string                `json:"code"`
	Msg  string                `json:"msg"`
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

type DepositHistoryResponse struct {
	Code string               `json:"code"`
	Msg  string               `json:"msg"`
	Data []DepositHistoryItem `json:"data"`
}
