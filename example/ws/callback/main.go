package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mmavka/go-blofin/models"
	"github.com/mmavka/go-blofin/ws"
)

func main() {
	client := ws.NewClient(ws.WSURLProd)
	errCh := make(chan error, 1)
	client.SetErrorHandler(func(err error) {
		errCh <- err
	})
	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		fmt.Println("connect error:", err)
		return
	}
	defer client.Close()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	handler := func(msg models.WSCandlestickMsg) {
		for _, c := range models.ParseWSCandlestickMsg(msg) {
			fmt.Printf("%s %s %s %s %s %s %s %s\n", msg.Arg.InstID, c.Ts, c.Open, c.High, c.Low, c.Close, c.VolCurrency, c.Confirm)
		}
	}

	err := client.SubscribeCandlesticks(ctx, ws.ChannelCandle1m, "BTC-USDT", handler)
	if err != nil {
		fmt.Println("subscribe error:", err)
		return
	}

	err = client.SubscribeCandlesticks(ctx, ws.ChannelCandle1m, "ETH-USDT", handler)
	if err != nil {
		fmt.Println("subscribe error:", err)
		return
	}

	go func() {
		time.Sleep(30 * time.Second)
		err := client.UnsubscribeCandlesticks(ctx, ws.ChannelCandle1m, "ETH-USDT")
		if err != nil {
			fmt.Println("unsubscribe error:", err)
		}
	}()

	fmt.Println("Subscribed to BTC-USDT and ETH-USDT 1m candles (callback). Press Ctrl+C to exit.")
	select {
	case <-ch:
		client.UnsubscribeCandlesticks(ctx, ws.ChannelCandle1m, "BTC-USDT")
	case err := <-errCh:
		fmt.Println("connection error:", err)
	}
}
