/**
 * @file: constants.go
 * @description: Constants for BloFin API
 * @dependencies: -
 * @created: 2025-05-19
 */

package blofin

// UseTestnet determines whether to use testnet or mainnet
var UseTestnet = false

// API Endpoints
const (
	// Production endpoints
	DefaultBaseURL     = "https://openapi.blofin.com"
	DefaultWSPublic    = "wss://openapi.blofin.com/ws/public"
	DefaultWSPrivate   = "wss://openapi.blofin.com/ws/private"
	DefaultWSCopyTrade = "wss://openapi.blofin.com/ws/copytrading/private"

	// Testnet (Demo) endpoints
	TestnetBaseURL     = "https://demo-trading-openapi.blofin.com"
	TestnetWSPublic    = "wss://demo-trading-openapi.blofin.com/ws/public"
	TestnetWSPrivate   = "wss://demo-trading-openapi.blofin.com/ws/private"
	TestnetWSCopyTrade = "wss://demo-trading-openapi.blofin.com/ws/copytrading/private"

	// Demo trading endpoints
	DemoBaseURL     = "https://demo-trading.blofin.com"
	DemoWSPublic    = "wss://demo-trading.blofin.com/ws"
	DemoWSPrivate   = "wss://demo-trading.blofin.com/ws/private"
	DemoWSCopyTrade = "wss://demo-trading.blofin.com/ws/copytrading/private"
)

// API Key Permissions
const (
	PermissionRead     = "READ"     // Can request and view account info
	PermissionTrade    = "TRADE"    // Can place and cancel orders
	PermissionTransfer = "TRANSFER" // Can make funding transfers
)

// Authentication Headers
const (
	HeaderAccessKey        = "ACCESS-KEY"
	HeaderAccessSign       = "ACCESS-SIGN"
	HeaderAccessTimestamp  = "ACCESS-TIMESTAMP"
	HeaderAccessNonce      = "ACCESS-NONCE"
	HeaderAccessPassphrase = "ACCESS-PASSPHRASE"
)

// Margin Modes
const (
	MarginModeCross    = "cross"
	MarginModeIsolated = "isolated"
)

// Position Modes
const (
	PositionModeNet       = "net_mode"
	PositionModeLongShort = "long_short_mode"
)

// Order Types
const (
	OrderTypeLimit    = "limit"     // Limit order
	OrderTypeMarket   = "market"    // Market order
	OrderTypePostOnly = "post_only" // Post-only order
	OrderTypeFOK      = "fok"       // Fill-or-kill order
	OrderTypeIOC      = "ioc"       // Immediate-or-cancel order
)

// Order Side
const (
	OrderSideBuy  = "buy"
	OrderSideSell = "sell"
)

// Position Side
const (
	PositionSideLong  = "long"
	PositionSideShort = "short"
	PositionSideNet   = "net"
)

// Account Types
const (
	AccountTypeFunding = "funding"
	AccountTypeFutures = "futures"
)

// WebSocket channels
const (
	// Public channels
	ChannelBooks       = "books"        // Order book channel
	ChannelBooks5      = "books5"       // Top 5 order book channel
	ChannelCandle      = "candle"       // Candlestick channel
	ChannelTrades      = "trades"       // Trades channel
	ChannelTickers     = "tickers"      // Tickers channel
	ChannelMarkPrice   = "mark-price"   // Mark price channel
	ChannelFundingRate = "funding-rate" // Funding rate channel

	// Private channels
	ChannelOrders      = "orders"      // Orders channel
	ChannelPositions   = "positions"   // Positions channel
	ChannelOrdersAlgo  = "orders-algo" // Algo orders channel
	ChannelAccount     = "account"     // Account channel
	ChannelLiquidation = "liquidation"
)

// WebSocket Operations
const (
	OpSubscribe   = "subscribe"
	OpUnsubscribe = "unsubscribe"
	OpLogin       = "login"
	OpPing        = "ping"
	OpPong        = "pong"
)

// Error Codes
const (
	ErrCodeSuccess          = "0"      // Success
	ErrCodeAllFailed        = "1"      // All operations failed
	ErrCodePartialSuccess   = "2"      // Batch operation partially succeeded
	ErrCodeEmptyParam       = "152001" // Parameter cannot be empty
	ErrCodeParamError       = "152002" // Parameter error
	ErrCodeEitherParam      = "152003" // Either parameter is required
	ErrCodeJSONSyntax       = "152004" // JSON syntax error
	ErrCodeWrongParam       = "152005" // Parameter error: wrong or empty
	ErrCodeBatchLimit       = "152006" // Batch orders limit (20)
	ErrCodeBatchSameInstId  = "152007" // Batch orders same instId required
	ErrCodeBatchSameField   = "152008" // Only same field for bulk cancellation
	ErrCodeInvalidFormat    = "152009" // Invalid format (alphanumeric + underscore, max 32 chars)
	ErrCodeNoTransactionAPI = "152011" // Transaction API Key does not support brokerId
	ErrCodeBrokerRequired   = "152012" // BrokerId is required
	ErrCodeUnmatchedBroker  = "152013" // Unmatched brokerId
	ErrCodeInvalidInstId    = "152014" // Instrument ID does not exist
	ErrCodeTooManyInstId    = "152015" // Too many instId values (max 20)
)

// Time Formats
const (
	TimeFormat = "2006-01-02T15:04:05.000Z"
)

// Default Values
const (
	DefaultTimeout = 10 // Default timeout in seconds
	MaxBatchSize   = 20 // Maximum number of orders in batch
)

// HTTP Methods
const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"
)

// Content Types
const (
	ContentTypeJSON = "application/json"
)

// API Paths
const (
	// Public API paths
	PathInstruments = "/api/v1/public/instruments"
	PathTickers     = "/api/v1/public/tickers"
	PathOrderBook   = "/api/v1/public/orderbook"
	PathTrades      = "/api/v1/public/trades"
	PathMarkPrice   = "/api/v1/public/mark-price"
	PathFundingRate = "/api/v1/public/funding-rate"
	PathCandles     = "/api/v1/public/candles"

	// Private API paths
	PathBalance       = "/api/v1/account/balance"
	PathPositions     = "/api/v1/account/positions"
	PathMarginMode    = "/api/v1/account/margin-mode"
	PathPositionMode  = "/api/v1/account/position-mode"
	PathLeverage      = "/api/v1/account/leverage"
	PathBatchLeverage = "/api/v1/account/batch-leverage"
	PathOrders        = "/api/v1/trade/orders"
	PathBatchOrders   = "/api/v1/trade/batch-orders"
	PathTPSLOrders    = "/api/v1/trade/tpsl-orders"
	PathAlgoOrders    = "/api/v1/trade/algo-orders"
	PathPendingOrders = "/api/v1/trade/pending-orders"
	PathOrderHistory  = "/api/v1/trade/order-history"
	PathTradeHistory  = "/api/v1/trade/trade-history"
)
