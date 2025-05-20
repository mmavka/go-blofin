# WebSocket Example

This example demonstrates how to use the WebSocket API to subscribe to market data channels.

## Features

- Subscribe to trades for multiple instruments
- Subscribe to candles for multiple instruments
- Subscribe to orderbook for multiple instruments
- Subscribe to tickers for multiple instruments
- Subscribe to funding rate for multiple instruments

## Usage

```bash
go run main.go
```

## Output

The example will print market data to the console:

```
Trade: BTC-USDT buy 0.001 @ 50000.00
Trade: ETH-USDT sell 0.1 @ 3000.00
Trade: SOL-USDT buy 1.0 @ 100.00

Candle: 2024-03-20T10:00:00Z O:50000.00 H:50100.00 L:49900.00 C:50050.00 V:100.00
Candle: 2024-03-20T10:00:00Z O:3000.00 H:3010.00 L:2990.00 C:3005.00 V:1000.00
Candle: 2024-03-20T10:00:00Z O:100.00 H:101.00 L:99.00 C:100.50 V:10000.00

Orderbook snapshot: 100 asks, 100 bids
Ask 1: Price: 50050.00, Size: 0.10000000
Bid 1: Price: 50049.00, Size: 0.20000000

Ticker: BTC-USDT Last:50050.00 Bid:50049.00 Ask:50051.00 Vol24h:1000.00
Ticker: ETH-USDT Last:3005.00 Bid:3004.00 Ask:3006.00 Vol24h:5000.00
Ticker: SOL-USDT Last:100.50 Bid:100.40 Ask:100.60 Vol24h:10000.00

Funding Rate: BTC-USDT Rate:0.0001 Time:2024-03-20T10:00:00Z
Funding Rate: ETH-USDT Rate:0.0002 Time:2024-03-20T10:00:00Z
Funding Rate: SOL-USDT Rate:0.0003 Time:2024-03-20T10:00:00Z
```

Press Ctrl+C to stop the example. 