# go-blofin

Go client for the BloFin API.

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

## Disclaimer
This library is not affiliated with, endorsed by, or sponsored by BloFin. Use at your own risk. The authors and contributors are not responsible for any financial losses or damages that may occur from using this library. Always verify the accuracy of the data and test thoroughly before using in production.

## License
MIT 

## Installation

```bash
go get github.com/mmavka/go-blofin
```

## Usage

### REST API

#### Public API

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mmavka/go-blofin/rest"
)

func main() {
	// Create a new client
	client := rest.NewDefaultRestClient()

	// Use testnet
	rest.UseTestnet = true

	// Get funding rate history
	fundingRateHistory, err := client.NewGetFundingRateHistoryService().
		InstId("BTC-USDT").
		Limit(10).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Funding rate history: %+v\n", fundingRateHistory)

	// Get candles
	candles, err := client.NewGetCandlesService().
		InstId("BTC-USDT").
		Bar("1m").
		Limit(10).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Candles: %+v\n", candles)

	// Get mark price
	markPrice, err := client.NewGetMarkPriceService().
		InstId("BTC-USDT").
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Mark price: %+v\n", markPrice)
}
```

#### Private API

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mmavka/go-blofin/rest"
)

func main() {
	// Create a new client
	client := rest.NewDefaultRestClient()

	// Use testnet
	rest.UseTestnet = true

	// Set API credentials
	client.SetAuth("your-api-key", "your-api-secret", "your-passphrase")

	// Get account balance
	balance, err := client.NewGetBalancesService().
		AccountType("cross").
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Account balance: %+v\n", balance)

	// Get leverage info
	leverageInfo, err := client.NewGetBatchLeverageInfoService().
		InstIds([]string{"BTC-USDT"}).
		MarginMode("cross").
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Leverage info: %+v\n", leverageInfo)

	// Place order
	order, err := client.NewPlaceOrderService().
		InstId("BTC-USDT").
		MarginMode("cross").
		PositionSide("long").
		Side("buy").
		OrderType("limit").
		Price("50000").
		Size("0.1").
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Order placed: %+v\n", order)

	// Get pending orders
	pendingOrders, err := client.NewGetPendingOrdersService().
		InstId("BTC-USDT").
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Pending orders: %+v\n", pendingOrders)

	// Cancel order
	cancelOrder, err := client.NewCancelOrderService().
		InstId("BTC-USDT").
		OrderId(order.Data[0].OrderId).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Order canceled: %+v\n", cancelOrder)
}
```

### WebSocket API

#### Public API

```go
package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/mmavka/go-blofin/ws"
)

func main() {
	// Use testnet
	ws.UseTestnet = true

	// Market data stream
	bookTickerHandler := func(event *ws.WsBookTickerEvent) {
		fmt.Printf("BookTicker event: %+v\n", event)
	}
	errHandler := func(err error) {
		fmt.Printf("Error: %v\n", err)
	}

	// Start book ticker stream
	doneC, stopC, err := ws.WsBookTickerServe("BTC-USDT", bookTickerHandler, errHandler)
	if err != nil {
		fmt.Printf("Error starting book ticker stream: %v\n", err)
		return
	}

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Wait for interrupt signal
	<-quit

	// Stop stream
	stopC <- struct{}{}

	// Wait for stream to close
	<-doneC

	fmt.Println("Stream stopped")
}
```

#### Private API

```go
package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/mmavka/go-blofin/ws"
)

func main() {
	// Use testnet
	ws.UseTestnet = true

	// Account data stream
	accountHandler := func(event *ws.WsAccountEvent) {
		fmt.Printf("Account event: %+v\n", event)
	}
	errHandler := func(err error) {
		fmt.Printf("Error: %v\n", err)
	}

	// Start account stream
	doneC, stopC, err := ws.WsAccountServe(
		"your-api-key",
		"your-api-secret",
		accountHandler,
		errHandler,
	)
	if err != nil {
		fmt.Printf("Error starting account stream: %v\n", err)
		return
	}

	// Position stream
	positionHandler := func(event *ws.WsPositionEvent) {
		fmt.Printf("Position event: %+v\n", event)
	}

	// Start position stream
	doneC2, stopC2, err := ws.WsPositionServe(
		"your-api-key",
		"your-api-secret",
		positionHandler,
		errHandler,
	)
	if err != nil {
		fmt.Printf("Error starting position stream: %v\n", err)
		return
	}

	// Order stream
	orderHandler := func(event *ws.WsOrderEvent) {
		fmt.Printf("Order event: %+v\n", event)
	}

	// Start order stream
	doneC3, stopC3, err := ws.WsOrderServe(
		"your-api-key",
		"your-api-secret",
		orderHandler,
		errHandler,
	)
	if err != nil {
		fmt.Printf("Error starting order stream: %v\n", err)
		return
	}

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Wait for interrupt signal
	<-quit

	// Stop all streams
	stopC <- struct{}{}
	stopC2 <- struct{}{}
	stopC3 <- struct{}{}

	// Wait for streams to close
	<-doneC
	<-doneC2
	<-doneC3

	fmt.Println("All streams stopped")
}
```

## Testnet

The library supports both mainnet and testnet. To use the testnet, set `rest.UseTestnet = true` for REST API or `ws.UseTestnet = true` for WebSocket API.

The testnet endpoints are:
- WebSocket: `wss://testnet.blofin.com/ws`
- REST API: `https://testnet.blofin.com`

The mainnet endpoints are:
- WebSocket: `wss://ws.blofin.com/ws`
- REST API: `https://api.blofin.com`

## Examples

See the [examples](examples) directory for more examples. 