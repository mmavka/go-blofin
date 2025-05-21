// Package ws contains constants for Blofin WebSocket public API.
//
// This file defines channel names, limits, error codes, and other constants for Blofin WebSocket API.
package ws

const (
	// WebSocket URLs
	WSURLProd = "wss://openapi.blofin.com/ws/public"
	WSURLDemo = "wss://demo-trading-openapi.blofin.com/ws/public"

	// Channel names (base)
	ChannelTrades      = "trades"
	ChannelTickers     = "tickers"
	ChannelOrderBook   = "books"
	ChannelFundingRate = "fundingrate"

	// Event types
	EventSubscribe   = "subscribe"
	EventUnsubscribe = "unsubscribe"
	EventError       = "error"
	EventInfo        = "info"
	EventPong        = "pong"
	EventPing        = "ping"

	// Limits
	// Max 30 subscriptions per WebSocket connection (see Blofin API docs)
	MaxSubscriptionsPerConn = 30
	// Max 10 connections per IP (see Blofin API docs)
	MaxConnectionsPerIP = 10
	// Max total size of args in one subscribe/unsubscribe request (bytes)
	MaxArgsSizeBytes = 4096

	// Ping/pong and connection management
	PingIntervalSec = 25 // Recommended ping interval (seconds)
	ConnTimeoutSec  = 30 // Connection closed if no data for 30s

	// Candlestick channel timeframes
	ChannelCandle1m  = "candle1m"
	ChannelCandle3m  = "candle3m"
	ChannelCandle5m  = "candle5m"
	ChannelCandle15m = "candle15m"
	ChannelCandle30m = "candle30m"
	ChannelCandle1H  = "candle1H"
	ChannelCandle2H  = "candle2H"
	ChannelCandle4H  = "candle4H"
	ChannelCandle6H  = "candle6H"
	ChannelCandle8H  = "candle8H"
	ChannelCandle12H = "candle12H"
	ChannelCandle1D  = "candle1D"
	ChannelCandle3D  = "candle3D"
	ChannelCandle1W  = "candle1W"
	ChannelCandle1M  = "candle1M"
)

// Error codes
const (
	ErrorCodeRateLimit   = "429"
	ErrorCodeLoginFailed = "60009"
)

// Subscription message format
const (
	OpSubscribe   = "subscribe"
	OpUnsubscribe = "unsubscribe"
	OpPing        = "ping"
	OpPong        = "pong"
)
