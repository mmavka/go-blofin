# Examples

This directory contains examples of using the go-blofin library.

## REST API Example

The REST API example shows how to use the REST API to get market data and manage orders:

### Public API

- Funding rate history
- Candles (OHLCV data)
- Mark price

### Private API

- Account balance
- Leverage information
- Place order
- Get pending orders
- Cancel order

To run the example:

```bash
cd examples/rest
go run main.go
```

Note: You need to set your API key and secret key in the example before running it.

## WebSocket API Example

The WebSocket API example shows how to use the WebSocket API to receive real-time data:

### Public API

- Book ticker updates

### Private API

- Account updates
- Position updates
- Order updates

To run the example:

```bash
cd examples/ws
go run main.go
```

Note: You need to set your API key and secret key in the example before running it.

## Testnet

All examples use the testnet by default. To use the mainnet, set `rest.UseTestnet = false` or `ws.UseTestnet = false` in the example.

The testnet endpoints are:
- WebSocket: `wss://testnet.blofin.com/ws`
- REST API: `https://testnet.blofin.com`

The mainnet endpoints are:
- WebSocket: `wss://ws.blofin.com/ws`
- REST API: `https://api.blofin.com`

## Running Examples

All examples can be run using the main.go file:

```bash
# Run public methods example
go run main.go -example=public

# Run private methods example
export BLOFIN_API_KEY="your-api-key"
export BLOFIN_API_SECRET="your-api-secret"
export BLOFIN_PASSPHRASE="your-passphrase"
go run main.go -example=private

# Run WebSocket example
go run main.go -example=ws
```

## Examples Description

### Public Methods Example
Demonstrates how to use public REST API methods:
- Getting list of instruments
- Getting order book
- Getting recent trades
- Getting candles

### Private Methods Example
Shows how to use private REST API methods:
- Getting account balance
- Getting positions
- Placing and canceling orders

### WebSocket Example
Demonstrates WebSocket usage:
- Connecting to WebSocket server
- Subscribing to trades and order book channels
- Handling real-time updates
- Graceful shutdown 