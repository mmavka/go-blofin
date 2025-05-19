/**
 * @file: constants.go
 * @description: Константы для работы с BloFin API
 * @dependencies: -
 * @created: 2025-05-19
 */

package rest

// API Endpoints
const (
	// Production endpoints
	DefaultBaseURL     = "https://openapi.blofin.com"
	DefaultWSPublic    = "wss://openapi.blofin.com/ws/public"
	DefaultWSPrivate   = "wss://openapi.blofin.com/ws/private"
	DefaultWSCopyTrade = "wss://openapi.blofin.com/ws/copytrading/private"

	// Demo trading endpoints
	DemoBaseURL   = "https://demo-trading-openapi.blofin.com"
	DemoWSPublic  = "wss://demo-trading-openapi.blofin.com/ws/public"
	DemoWSPrivate = "wss://demo-trading-openapi.blofin.com/ws/private"
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
	OrderTypeLimit    = "limit"     // Лимитный ордер
	OrderTypeMarket   = "market"    // Рыночный ордер
	OrderTypePostOnly = "post_only" // Post-only ордер
	OrderTypeFOK      = "fok"       // Fill-or-kill ордер
	OrderTypeIOC      = "ioc"       // Immediate-or-cancel ордер
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
