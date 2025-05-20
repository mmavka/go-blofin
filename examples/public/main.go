package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mmavka/go-blofin/rest"
)

func main() {
	ctx := context.Background()

	// Create new client
	client := rest.NewDefaultRestClient()

	// Get instruments
	instruments, err := client.NewGetInstrumentsService().Do(ctx)
	if err != nil {
		log.Fatalf("Failed to get instruments: %v", err)
	}

	// Print instruments
	for _, instrument := range instruments.Data {
		fmt.Printf("Instrument: %+v\n", instrument)
	}

	// Get tickers
	tickers, err := client.NewGetTickersService().Do(ctx)
	if err != nil {
		log.Fatalf("Failed to get tickers: %v", err)
	}

	// Print tickers
	for _, ticker := range tickers.Data {
		fmt.Printf("Ticker: %+v\n", ticker)
	}

	// Get order book
	orderBook, err := client.NewGetOrderBookService().
		InstId("BTC-USDT").
		Size(5).
		Do(ctx)
	if err != nil {
		log.Fatalf("Failed to get order book: %v", err)
	}

	// Print order book
	fmt.Printf("Order Book: %+v\n", orderBook)

	// Get trades
	trades, err := client.NewGetTradesService().
		InstId("BTC-USDT").
		Limit(5).
		Do(ctx)
	if err != nil {
		log.Fatalf("Failed to get trades: %v", err)
	}

	// Print trades
	for _, trade := range trades.Data {
		fmt.Printf("Trade: %+v\n", trade)
	}

	// Get mark price
	markPrice, err := client.NewGetMarkPriceService().
		InstId("BTC-USDT").
		Do(ctx)
	if err != nil {
		log.Fatalf("Failed to get mark price: %v", err)
	}

	// Print mark price
	fmt.Printf("Mark Price: %+v\n", markPrice)

	// Get candles
	candles, err := client.NewGetCandlesService().
		InstId("BTC-USDT").
		Bar("1m").
		Limit(5).
		Do(ctx)
	if err != nil {
		log.Fatalf("Failed to get candles: %v", err)
	}

	// Print candles
	for _, candle := range candles.Data {
		fmt.Printf("Candle: %+v\n", candle)
	}
}
