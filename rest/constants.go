// Package rest contains constants for Blofin API endpoints and settings.
package rest

const (
	// Base URLs
	BaseURLProd = "https://openapi.blofin.com"
	BaseURLDemo = "https://demo-trading-openapi.blofin.com"

	// Public REST endpoints
	EndpointInstruments     = "/api/v1/market/instruments"
	EndpointTickers         = "/api/v1/market/tickers"
	EndpointOrderBook       = "/api/v1/market/books"
	EndpointTrades          = "/api/v1/market/trades"
	EndpointMarkPrice       = "/api/v1/market/mark-price"
	EndpointFundingRate     = "/api/v1/market/funding-rate"
	EndpointFundingRateHist = "/api/v1/market/funding-rate-history"
	EndpointCandlesticks    = "/api/v1/market/candles"

	// WebSocket URLs
	WSURLProd = "wss://openapi.blofin.com/ws/public"
	WSURLDemo = "wss://demo-trading-openapi.blofin.com/ws/public"

	// API response codes
	CodeSuccess = "0"
)
