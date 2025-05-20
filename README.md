# go-blofin

Go library for working with BloFin cryptocurrency exchange

## Description

`go-blofin` is a modular library for working with BloFin REST and WebSocket API, inspired by go-binance architecture. Supports public and private methods, event subscriptions, signature generation, error handling and covered with unit tests.

## Features
- Public REST methods: instruments, tickers, order book, trades, candles, mark price, funding rate
- Private REST methods: balance, positions, orders, transfer and withdrawal history
- WebSocket: subscribe to trades, candles, tickers, order book, funding-rate
- Signature generation for private requests
- Full unit test coverage
- Modern architecture, extensibility, clean code

## Project Structure
- `rest/` — REST client, services, models
- `ws/` — WebSocket client, models, signatures
- `auth/` — signature generation
- `utils/` — base errors
- `docs/` — documentation (architecture, changelog, tasks, Q&A)

## Quick Start
```go
import "github.com/mmavka/go-blofin/rest"

client := rest.NewDefaultRestClient()
resp, err := client.NewGetInstrumentsService().Do(context.Background())
if err != nil {
    // handle error
}
for _, inst := range resp.Data {
    fmt.Println(inst.InstID)
}
```

## WebSocket Example
```go
import "github.com/mmavka/go-blofin/ws"

wsClient := ws.NewDefaultClient()
wsClient.SetErrorHandler(func(err error) {
    log.Printf("WebSocket error: %v", err)
})
if err := wsClient.Connect(); err != nil {
    panic(err)
}
_ = wsClient.Subscribe([]ws.ChannelArgs{{Channel: "trades", InstId: "BTC-USDT"}})
for trade := range wsClient.Trades() {
    fmt.Println(trade)
}
```

## WebSocket Error Handling
```go
wsClient := ws.NewDefaultClient()
wsClient.SetErrorHandler(func(err error) {
    log.Printf("WebSocket error: %v", err)
})
```

## Custom Endpoint
```go
import "github.com/mmavka/go-blofin/rest"
import "github.com/mmavka/go-blofin/ws"

client := rest.NewDefaultRestClient()
client.SetBaseURL("https://sandbox.blofin.com")

wsClient := ws.NewDefaultClient()
wsClient.SetURL("wss://sandbox-ws.blofin.com/ws")
```

## WebSocket Channel Subscription Features
- Public channels (trades, candles, tickers, order book, etc.) do not require authentication.
- Private channels (orders, positions, orders-algo, etc.) require prior Login call.
- Total subscription request length should not exceed 4096 bytes.
- Violation of these conditions will result in an error.

## Requirements
- Go 1.24+
- [resty](https://github.com/go-resty/resty)
- [goccy/go-json](https://github.com/goccy/go-json)
- [gorilla/websocket](https://github.com/gorilla/websocket)

## License
MIT 