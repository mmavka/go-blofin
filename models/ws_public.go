// Package models contains WebSocket message structures for public channels.
//
// This file defines structures for public WebSocket push messages (candlesticks, trades, tickers, etc.).
package models

// WSTradeMsg represents a push message from the trades channel.
type WSTradeMsg struct {
	Arg struct {
		Channel string `json:"channel"`
		InstID  string `json:"instId"`
	} `json:"arg"`
	Data [][]string `json:"data"` // [tradeId, price, size, side, ts]
}

// WSTickerMsg represents a push message from the tickers channel.
type WSTickerMsg struct {
	Arg struct {
		Channel string `json:"channel"`
		InstID  string `json:"instId"`
	} `json:"arg"`
	Data [][]string `json:"data"` // [last, lastSize, askPx, askSz, bidPx, bidSz, high24h, open24h, low24h, volCcy24h, vol24h, ts]
}

// WSOrderBookMsg represents a push message from the order book channel.
type WSOrderBookMsg struct {
	Arg struct {
		Channel string `json:"channel"`
		InstID  string `json:"instId"`
	} `json:"arg"`
	Action string `json:"action,omitempty"` // "snapshot" or "update"
	Data   struct {
		Asks      [][]string `json:"asks"`
		Bids      [][]string `json:"bids"`
		TS        string     `json:"ts"`
		PrevSeqID string     `json:"prevSeqId"`
		SeqID     string     `json:"seqId"`
	} `json:"data"`
}

// WSCandlestickMsg represents a push message from the candlesticks channel.
type WSCandlestickMsg struct {
	Arg struct {
		Channel string `json:"channel"`
		InstID  string `json:"instId"`
	} `json:"arg"`
	Data [][]string `json:"data"` // Each element is a candlestick as string array
}

// WSCandle represents a single candlestick (OHLCV) from WS push.
type WSCandle struct {
	Ts               string // Opening time, ms
	Open             string // Open price
	High             string // Highest price
	Low              string // Lowest price
	Close            string // Close price
	Vol              string // Trading volume (contracts)
	VolCurrency      string // Trading volume (base currency)
	VolCurrencyQuote string // Trading volume (quote currency)
	Confirm          string // 0 = uncompleted, 1 = completed
}

// WSFundingRateMsg represents a push message from the funding rate channel.
type WSFundingRateMsg struct {
	Arg struct {
		Channel string `json:"channel"`
		InstID  string `json:"instId"`
	} `json:"arg"`
	Data [][]string `json:"data"` // [fundingRate, fundingTime, instId]
}

// ParseWSCandle parses a string array from WS push into WSCandle struct.
func ParseWSCandle(arr []string) (WSCandle, bool) {
	if len(arr) < 9 {
		return WSCandle{}, false
	}
	return WSCandle{
		Ts:               arr[0],
		Open:             arr[1],
		High:             arr[2],
		Low:              arr[3],
		Close:            arr[4],
		Vol:              arr[5],
		VolCurrency:      arr[6],
		VolCurrencyQuote: arr[7],
		Confirm:          arr[8],
	}, true
}

// ParseWSCandlestickMsg parses WSCandlestickMsg and returns a slice of WSCandle.
func ParseWSCandlestickMsg(msg WSCandlestickMsg) []WSCandle {
	candles := make([]WSCandle, 0, len(msg.Data))
	for _, arr := range msg.Data {
		if c, ok := ParseWSCandle(arr); ok {
			candles = append(candles, c)
		}
	}
	return candles
}
