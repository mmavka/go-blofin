package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/mmavka/go-blofin"
	"github.com/mmavka/go-blofin/ws"
)

func main() {
	// Run public API example
	runPublicExample()

	// Run private API example
	runPrivateExample()
}

func runPublicExample() {

	// Create new client
	client := ws.NewDefaultClient()

	// Set error handler
	client.SetErrorHandler(func(err error) {
		fmt.Printf("WebSocket error: %v\n", err)
	})

	// Connect to WebSocket
	if err := client.Connect(); err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}

	// Subscribe to book ticker
	err := client.Subscribe([]ws.ChannelArgs{
		{
			Channel: blofin.ChannelCandle + "1m",
			InstId:  "BTC-USDT",
		},
		{
			Channel: blofin.ChannelCandle + "1m",
			InstId:  "ETH-USDT",
		},
	})
	if err != nil {
		fmt.Printf("Failed to subscribe: %v\n", err)
		return
	}

	// Handle messages
	go func() {
		for msg := range client.Candles() {
			fmt.Printf("Received message: %s\n", msg)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Close connection
	client.Close()
	fmt.Println("Connection closed")
}

func runPrivateExample() {
	// Use testnet
	ws.UseTestnet = true

	// Create new client
	client := ws.NewClient(blofin.TestnetWSPrivate)

	// Set error handler
	client.SetErrorHandler(func(err error) {
		fmt.Printf("WebSocket error: %v\n", err)
	})

	// Connect to WebSocket
	if err := client.Connect(); err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}

	// Login
	if err := client.Login(
		"your-api-key",
		"your-api-secret",
		"your-passphrase",
	); err != nil {
		fmt.Printf("Failed to login: %v\n", err)
		return
	}

	// Subscribe to private channels
	err := client.Subscribe([]ws.ChannelArgs{
		{
			Channel: blofin.ChannelOrders,
			InstId:  "BTC-USDT",
		},
		{
			Channel: blofin.ChannelPositions,
			InstId:  "BTC-USDT",
		},
		{
			Channel: blofin.ChannelAccount,
		},
	})
	if err != nil {
		fmt.Printf("Failed to subscribe: %v\n", err)
		return
	}

	// Handle messages
	go func() {
		for msg := range client.Messages() {
			fmt.Printf("Received message: %s\n", string(msg))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Close connection
	client.Close()
	fmt.Println("Connection closed")
}
