package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mmavka/go-blofin"
	"github.com/mmavka/go-blofin/rest"
)

func main() {
	// Run public API example
	runPublicExample()

	// Run private API example
	runPrivateExample()
}

func runPublicExample() {
	// Create a new client
	client := rest.NewDefaultRestClient()

	// Use testnet
	blofin.UseTestnet = true

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

func runPrivateExample() {
	// Create a new client
	client := rest.NewDefaultRestClient()

	// Use testnet
	blofin.UseTestnet = true

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
