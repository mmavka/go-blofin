package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mmavka/go-blofin/ws"
)

func main() {
	// Create WebSocket client
	wsClient := ws.NewDefaultClient()

	// Set error handler
	wsClient.SetErrorHandler(func(err error) {
		log.Printf("WebSocket error: %v", err)
	})

	// Connect to WebSocket server
	if err := wsClient.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer wsClient.Close()

	// Subscribe to trades channel
	err := wsClient.Subscribe([]ws.ChannelArgs{
		{Channel: "trades", InstId: "BTC-USDT"},
	})
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	// Subscribe to order book channel
	err = wsClient.Subscribe([]ws.ChannelArgs{
		{Channel: "books", InstId: "BTC-USDT"},
	})
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	// Handle trades
	go func() {
		for trade := range wsClient.Trades() {
			for _, t := range trade.Data {
				fmt.Printf("Trade: Price=%s, Size=%s, Side=%s\n",
					t.Price, t.Size, t.Side)
			}
		}
	}()

	// Handle order book updates
	go func() {
		for book := range wsClient.OrderBooks() {
			fmt.Printf("Order book update: Bids=%v, Asks=%v\n",
				book.Data.Bids, book.Data.Asks)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
