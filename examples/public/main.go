package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mmavka/go-blofin/rest"
)

func main() {
	// Create client with default URL
	client := rest.NewDefaultRestClient()

	// Get all instruments
	instruments, err := client.NewGetInstrumentsService().Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to get instruments: %v", err)
	}
	fmt.Printf("Found %d instruments\n", len(instruments.Data))

	// Get order book for BTC-USDT
	orderBook, err := client.NewGetOrderBookService().
		InstId("BTC-USDT").
		Size(10).
		Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to get order book: %v", err)
	}
	fmt.Printf("Order book for BTC-USDT:\n")
	if len(orderBook.Data) > 0 {
		fmt.Printf("Bids: %v\n", orderBook.Data[0].Bids)
		fmt.Printf("Asks: %v\n", orderBook.Data[0].Asks)
	}

	// Get recent trades
	trades, err := client.NewGetTradesService().
		InstId("BTC-USDT").
		Limit(5).
		Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to get trades: %v", err)
	}
	fmt.Printf("Recent trades:\n")
	for _, trade := range trades.Data {
		fmt.Printf("Price: %s, Size: %s, Side: %s\n", trade.Price, trade.Size, trade.Side)
	}

	// Get candles
	candles, err := client.NewGetCandlesService().
		InstId("BTC-USDT").
		Bar("1m").
		After(fmt.Sprintf("%d", time.Now().Add(-1*time.Hour).UnixMilli())).
		Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to get candles: %v", err)
	}
	fmt.Printf("Candles for last hour:\n")
	for _, candle := range candles.Data {
		fmt.Printf("Time: %s, Open: %s, High: %s, Low: %s, Close: %s\n",
			candle.Ts, candle.Open, candle.High, candle.Low, candle.Close)
	}
}
