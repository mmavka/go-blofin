package ws

// WebSocket endpoints
const (
	// Public WebSocket endpoint
	PublicWebSocketURL = "wss://openapi.blofin.com/ws/public"
	// Private WebSocket endpoint
	PrivateWebSocketURL = "wss://openapi.blofin.com/ws/private"
)

// WebSocket operations
const (
	// Operations
	OpSubscribe   = "subscribe"
	OpUnsubscribe = "unsubscribe"
	OpLogin       = "login"
	OpPing        = "ping"
	OpPong        = "pong"
)

// WebSocket channels
const (
	// Public channels
	ChannelTrades      = "trades"
	ChannelCandles     = "candles"
	ChannelOrderbook   = "books"
	ChannelTickers     = "tickers"
	ChannelFundingRate = "funding-rate"

	// Candle intervals
	CandleInterval1m  = "candle1m"
	CandleInterval3m  = "candle3m"
	CandleInterval5m  = "candle5m"
	CandleInterval15m = "candle15m"
	CandleInterval30m = "candle30m"
	CandleInterval1h  = "candle1H"
	CandleInterval2h  = "candle2H"
	CandleInterval4h  = "candle4H"
	CandleInterval6h  = "candle6H"
	CandleInterval8h  = "candle8H"
	CandleInterval12h = "candle12H"
	CandleInterval1d  = "candle1D"
	CandleInterval3d  = "candle3D"
	CandleInterval1w  = "candle1W"
	CandleInterval1M  = "candle1M"
)

// Rate limits
const (
	// WebSocket connection limits
	MaxNewConnectionsPerSecond = 1 // Maximum number of new connections per second per IP
	MaxConnectionsPerIP        = 5 // Maximum number of concurrent connections per IP

	// WebSocket limits
	MaxSubscriptionLength = 4096 // Maximum length of subscription in bytes

	// WebSocket ping/pong
	PingInterval = 20 // Send ping every 20 seconds
	PongTimeout  = 5  // Wait for pong response for 5 seconds

	// WebSocket reconnection
	ReconnectDelay = 1 // Delay between reconnection attempts in seconds
	MaxReconnects  = 5 // Maximum number of reconnection attempts
)

// Error codes
const (
	ErrCodeRateLimit = "429" // Rate limit reached
	ErrCodeForbidden = "403" // Network firewall restriction
)
