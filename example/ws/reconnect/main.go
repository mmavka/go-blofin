/**
 * @file: main.go
 * @description: Example of reconnect and resubscribe logic for ws client (handled by application)
 * @dependencies: github.com/mmavka/go-blofin/ws, github.com/mmavka/go-blofin/models
 * @created: 2024-06-14
 */

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/mmavka/go-blofin/models"
	"github.com/mmavka/go-blofin/ws"
)

type Subscription struct {
	Channel string
	InstID  string
	Handler func(models.WSCandlestickMsg)
}

func main() {
	// For graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	// List of active subscriptions
	subscriptions := []Subscription{
		{
			Channel: ws.ChannelCandle1m,
			InstID:  "BTC-USDT",
			Handler: func(msg models.WSCandlestickMsg) {
				fmt.Printf("BTC candle: %+v\n", msg)
			},
		},
		{
			Channel: ws.ChannelCandle1m,
			InstID:  "ETH-USDT",
			Handler: func(msg models.WSCandlestickMsg) {
				fmt.Printf("ETH candle: %+v\n", msg)
			},
		},
	}

	for {
		ctx := context.Background()
		logger := ws.NewDefaultLogger(ws.LogLevelDebug)
		client := ws.NewClient(ws.WSURLProd, logger)

		// Канал для обработки ошибок соединения
		errCh := make(chan error, 1)
		client.SetErrorHandler(func(err error) {
			errCh <- err
		})

		err := client.Connect(ctx)
		if err != nil {
			log.Printf("connect error: %v", err)
			continue
		}

		// Subscribe to all channels
		for _, sub := range subscriptions {
			err := client.SubscribeCandlesticks(ctx, sub.Channel, sub.InstID, sub.Handler)
			if err != nil {
				log.Printf("subscribe error: %v", err)
			}
		}

		select {
		case <-sigCh:
			log.Println("shutting down...")
			client.Close()
			return
		case err := <-errCh:
			log.Printf("connection error: %v, reconnecting...", err)
			client.Close()
			continue
		}
	}
}
