package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mmavka/go-blofin/internal/models"
	"github.com/mmavka/go-blofin/internal/ws"
)

func main() {
	// Можно явно указать уровень логирования:
	// logger := ws.NewDefaultLogger(ws.LogLevelDebug)
	// client := ws.NewClient(ws.WSURLProd, logger)
	// Или не передавать logger вовсе, тогда будет только error-логирование:
	logger := ws.NewDefaultLogger(ws.LogLevelDebug)
	client := ws.NewClient(ws.WSURLProd, logger)
	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		fmt.Println("connect error:", err)
		return
	}
	defer client.Close()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	msgChBtc, err := client.SubscribeCandlesticksChan(ctx, ws.ChannelCandle1m, "BTC-USDT")
	if err != nil {
		fmt.Println("subscribe error:", err)
		return
	}

	msgChEth, err := client.SubscribeCandlesticksChan(ctx, ws.ChannelCandle1m, "ETH-USDT")
	if err != nil {
		fmt.Println("subscribe error:", err)
		return
	}

	fmt.Println("Subscribed to BTC-USDT and ETH-USDT 1m candles (channel). Press Ctrl+C to exit.")

	go func() {
		for msg := range msgChBtc {
			for _, c := range models.ParseWSCandlestickMsg(msg) {
				fmt.Printf("%s %s %s %s %s %s %s %s\n", msg.Arg.InstID, c.Ts, c.Open, c.High, c.Low, c.Close, c.VolCurrency, c.Confirm)
			}
		}
	}()

	go func() {
		for msg := range msgChEth {
			for _, c := range models.ParseWSCandlestickMsg(msg) {
				fmt.Printf("%s %s %s %s %s %s %s %s\n", msg.Arg.InstID, c.Ts, c.Open, c.High, c.Low, c.Close, c.VolCurrency, c.Confirm)
			}
		}
	}()

	<-ch
	client.UnsubscribeCandlesticks(ctx, ws.ChannelCandle1m, "BTC-USDT")
}
