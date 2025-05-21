// Package models contains data structures for API payloads.
package models

// Instrument represents a trading instrument from Blofin public API.
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

// Instrument types
const (
	InstTypeSwap = "SWAP"
	// Add other types if needed
)

// Contract types
const (
	ContractTypeLinear  = "linear"
	ContractTypeInverse = "inverse"
)

// Instrument states
const (
	InstrumentStateLive    = "live"
	InstrumentStateSuspend = "suspend"
)

// InstrumentsResponse is the response for GET /api/v1/public/instruments
// and similar endpoints.
type InstrumentsResponse struct {
	Code string       `json:"code"`
	Msg  string       `json:"msg"`
	Data []Instrument `json:"data"`
}

// ApiError represents an error returned by the Blofin API.
type ApiError struct {
	Code    string
	Message string
}

func (e *ApiError) Error() string {
	return e.Message
}

// Candlestick represents a single candlestick (OHLCV) from Blofin API.
type Candlestick struct {
	Ts             string `json:"ts"` // ms
	Open           string `json:"open"`
	High           string `json:"high"`
	Low            string `json:"low"`
	Close          string `json:"close"`
	Volume         string `json:"volume"`
	VolumeCurrency string `json:"volumeCurrency"`
	VolumeQuote    string `json:"volumeQuote"`
	Confirm        string `json:"confirm"` // 0 = uncompleted, 1 = completed
}

// Bar sizes
const (
	Bar1m  = "1m"
	Bar3m  = "3m"
	Bar5m  = "5m"
	Bar15m = "15m"
	Bar30m = "30m"
	Bar1H  = "1H"
	Bar2H  = "2H"
	Bar4H  = "4H"
	Bar6H  = "6H"
	Bar8H  = "8H"
	Bar12H = "12H"
	Bar1D  = "1D"
	Bar3D  = "3D"
	Bar1W  = "1W"
	Bar1M  = "1M"
)

// Candlestick confirm status
const (
	CandlestickUncompleted = "0"
	CandlestickCompleted   = "1"
)

// Ticker represents a ticker snapshot from Blofin public API.
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

// OrderBookLevel represents a single price level in the order book.
type OrderBookLevel struct {
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
}

// OrderBook represents the order book snapshot from Blofin public API.
type OrderBook struct {
	Asks []OrderBookLevel `json:"asks"`
	Bids []OrderBookLevel `json:"bids"`
	Ts   string           `json:"ts"`
}

// Trade represents a recent transaction from Blofin public API.
type Trade struct {
	TradeID string `json:"tradeId"`
	InstID  string `json:"instId"`
	Price   string `json:"price"`
	Size    string `json:"size"`
	Side    string `json:"side"`
	Ts      string `json:"ts"`
}

// Trade sides
const (
	TradeSideBuy  = "buy"
	TradeSideSell = "sell"
)

// MarkPrice represents index and mark price from Blofin public API.
type MarkPrice struct {
	InstID     string `json:"instId"`
	IndexPrice string `json:"indexPrice"`
	MarkPrice  string `json:"markPrice"`
	Ts         string `json:"ts"`
}

// FundingRate represents funding rate info from Blofin public API.
type FundingRate struct {
	InstID      string `json:"instId"`
	FundingRate string `json:"fundingRate"`
	FundingTime string `json:"fundingTime"`
}
