package ws

// ChannelArgs represents channel subscription arguments
type ChannelArgs struct {
	Channel string `json:"channel"`
	InstId  string `json:"instId"`
}

// SubscribeRequest represents subscription request
type SubscribeRequest struct {
	Op   string        `json:"op"`
	Args []ChannelArgs `json:"args"`
}

// SubscribeResponse represents subscription response
type SubscribeResponse struct {
	Event string      `json:"event"`
	Arg   ChannelArgs `json:"arg"`
	Code  string      `json:"code"`
	Msg   string      `json:"msg"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Op   string `json:"op"`
	Args []struct {
		APIKey     string `json:"apiKey"`
		Passphrase string `json:"passphrase"`
		Timestamp  string `json:"timestamp"`
		Sign       string `json:"sign"`
		Nonce      string `json:"nonce"`
	} `json:"args"`
}

// LoginResponse represents login response
type LoginResponse struct {
	Event string `json:"event"`
	Code  string `json:"code"`
	Msg   string `json:"msg"`
}

// Trade represents a trade message
type Trade struct {
	InstId  string `json:"instId"`
	TradeId string `json:"tradeId"`
	Price   string `json:"price"`
	Size    string `json:"size"`
	Side    string `json:"side"`
	Ts      string `json:"ts"`
}

// TradeMessage represents a trade message with channel info
type TradeMessage struct {
	Arg  ChannelArgs `json:"arg"`
	Data []Trade     `json:"data"`
}

// Candle represents a candlestick data
type Candle struct {
	Ts               string `json:"ts"`               // Opening time
	Open             string `json:"open"`             // Open price
	High             string `json:"high"`             // Highest price
	Low              string `json:"low"`              // Lowest price
	Close            string `json:"close"`            // Close price
	Vol              string `json:"vol"`              // Trading volume in contracts
	VolCurrency      string `json:"volCurrency"`      // Trading volume in base currency
	VolCurrencyQuote string `json:"volCurrencyQuote"` // Trading volume in quote currency
	Confirm          string `json:"confirm"`          // Candlestick state (0: uncompleted, 1: completed)
}

// CandleMessage represents a candle message with channel info
type CandleMessage struct {
	Arg  ChannelArgs `json:"arg"`
	Data [][]string  `json:"data"` // Array of candle data arrays
}

// OrderBookLevel represents a single level in the order book
type OrderBookLevel struct {
	Price  string `json:"price"`  // Price level
	Amount string `json:"amount"` // Amount at this price level
}

// OrderBookData represents orderbook data
type OrderBookData struct {
	Asks      [][]float64 `json:"asks"`      // [price, size]
	Bids      [][]float64 `json:"bids"`      // [price, size]
	Ts        string      `json:"ts"`        // Order book generation time
	PrevSeqId string      `json:"prevSeqId"` // Previous sequence ID
	SeqId     string      `json:"seqId"`     // Current sequence ID
}

// OrderBookMessage represents orderbook message
type OrderBookMessage struct {
	Action string        `json:"action"` // "snapshot" or "update"
	Arg    ChannelArgs   `json:"arg"`
	Data   OrderBookData `json:"data"`
}

// Ticker represents a ticker data
type Ticker struct {
	InstId         string `json:"instId"`         // Instrument ID
	Last           string `json:"last"`           // Last traded price
	LastSize       string `json:"lastSize"`       // Last traded size
	AskPrice       string `json:"askPrice"`       // Best ask price
	AskSize        string `json:"askSize"`        // Best ask size
	BidPrice       string `json:"bidPrice"`       // Best bid price
	BidSize        string `json:"bidSize"`        // Best bid size
	Open24h        string `json:"open24h"`        // Open price in the past 24 hours
	High24h        string `json:"high24h"`        // Highest price in the past 24 hours
	Low24h         string `json:"low24h"`         // Lowest price in the past 24 hours
	VolCurrency24h string `json:"volCurrency24h"` // 24h trading volume in base currency
	Vol24h         string `json:"vol24h"`         // 24h trading volume in contracts
	Ts             string `json:"ts"`             // Ticker data generation time
}

// TickerMessage represents a ticker message with channel info
type TickerMessage struct {
	Arg  ChannelArgs `json:"arg"`
	Data []Ticker    `json:"data"`
}

// FundingRate represents a funding rate data
type FundingRate struct {
	InstId      string `json:"instId"`      // Instrument ID
	FundingRate string `json:"fundingRate"` // Current funding rate
	FundingTime string `json:"fundingTime"` // Funding time of the upcoming settlement
}

// FundingRateMessage represents a funding rate message with channel info
type FundingRateMessage struct {
	Arg  ChannelArgs   `json:"arg"`
	Data []FundingRate `json:"data"`
}
