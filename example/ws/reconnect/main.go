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
	"log/slog"
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
		client := ws.NewClient(ws.WSURLProd)

		// Канал для обработки ошибок соединения
		errCh := make(chan error, 1)
		client.SetErrorHandler(func(err error) {
			errCh <- err
		})

		err := client.Connect(ctx)
		if err != nil {
			slog.Error("connect error", "error", err)
			continue
		}

		// Subscribe to all channels
		for _, sub := range subscriptions {
			err := client.SubscribeCandlesticks(ctx, sub.Channel, sub.InstID, sub.Handler)
			if err != nil {
				slog.Error("subscribe error", "error", err)
			}
		}

		select {
		case <-sigCh:
			slog.Info("shutting down...")
			client.Close()
			return
		case err := <-errCh:
			slog.Error("connection error, reconnecting...", "error", err)
			client.Close()
			continue
		}
	}
}
