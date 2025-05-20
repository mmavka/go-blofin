# Examples

This directory contains examples of using the go-blofin library.

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