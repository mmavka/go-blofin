package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mmavka/go-blofin/rest"
)

func main() {
	// Create client with default URL
	client := rest.NewDefaultRestClient()

	// Set authentication credentials
	client.SetAuth(
		os.Getenv("BLOFIN_API_KEY"),
		os.Getenv("BLOFIN_API_SECRET"),
		os.Getenv("BLOFIN_PASSPHRASE"),
	)

	// Get account balance
	balance, err := client.GetAccountBalance().Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to get account balance: %v", err)
	}
	fmt.Printf("Account balance:\n")
	for _, detail := range balance.Data.Details {
		fmt.Printf("Currency: %s, Equity: %s, Available: %s\n",
			detail.Currency, detail.Equity, detail.Available)
	}

	// Get positions
	positions, err := client.GetPositions().Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to get positions: %v", err)
	}
	fmt.Printf("Open positions:\n")
	for _, position := range positions.Data {
		fmt.Printf("Instrument: %s, Size: %s, Unrealized PnL: %s\n",
			position.InstId, position.Positions, position.UnrealizedPnl)
	}

	// Place a limit order
	order, err := client.NewPlaceOrderService().
		InstId("BTC-USDT").
		Side("buy").
		OrderType("limit").
		Price("30000").
		Size("0.001").
		Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to place order: %v", err)
	}
	fmt.Printf("Order placed: %s\n", order.Data[0].OrderId)

	// Cancel the order
	_, err = client.NewCancelOrderService().
		InstId("BTC-USDT").
		OrderId(order.Data[0].OrderId).
		Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to cancel order: %v", err)
	}
	fmt.Printf("Order cancelled: %s\n", order.Data[0].OrderId)
}
