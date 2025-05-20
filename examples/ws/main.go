package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mmavka/go-blofin/ws"
)

func main() {
	// Create WebSocket client
	config := &ws.WsConfig{
		Endpoint: ws.PublicWebSocketURL,
	}
	service := ws.NewMarketService(config)

	// Set error handler
	service.SetErrHandler(func(err error) {
		log.Printf("Error: %v", err)
	})

	// Connect to WebSocket server
	err := service.Connect()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer service.Close()

	// Subscribe to trades for multiple instruments
	err = service.SubscribeTradesMulti([]string{"BTC-USDT", "ETH-USDT", "SOL-USDT"}, func(msg []byte) {
		log.Printf("Received trade message: %s", string(msg))
		var tradeMsg ws.TradeMessage
		if err := json.Unmarshal(msg, &tradeMsg); err != nil {
			log.Printf("Failed to unmarshal trade message: %v", err)
			return
		}
		for _, trade := range tradeMsg.Data {
			log.Printf("Trade: %s %s %s @ %s", trade.InstId, trade.Side, trade.Size, trade.Price)
		}
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to trades: %v", err)
	}

	// Subscribe to candles for multiple instruments
	err = service.SubscribeCandlesMulti([]string{"BTC-USDT", "ETH-USDT", "SOL-USDT"}, ws.CandleInterval1m, func(msg []byte) {
		log.Printf("Received candle message: %s", string(msg))
		var candleMsg ws.CandleMessage
		if err := json.Unmarshal(msg, &candleMsg); err != nil {
			log.Printf("Failed to unmarshal candle message: %v", err)
			return
		}
		for _, candle := range candleMsg.Data {
			log.Printf("Candle: %s O:%s H:%s L:%s C:%s V:%s", candle[0], candle[1], candle[2], candle[3], candle[4], candle[5])
		}
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to candles: %v", err)
	}

	// Subscribe to orderbook for multiple instruments
	err = service.SubscribeOrderbookMulti([]string{"BTC-USDT", "ETH-USDT", "SOL-USDT"}, func(msg []byte) {
		log.Printf("Received orderbook message: %s", string(msg))
		var orderbookMsg ws.OrderBookMessage
		if err := json.Unmarshal(msg, &orderbookMsg); err != nil {
			log.Printf("Failed to unmarshal orderbook message: %v", err)
			return
		}
		log.Printf("Orderbook %s: %d asks, %d bids", orderbookMsg.Action, len(orderbookMsg.Data.Asks), len(orderbookMsg.Data.Bids))

		// Print top 5 asks and bids
		for i := 0; i < 5 && i < len(orderbookMsg.Data.Asks); i++ {
			ask := orderbookMsg.Data.Asks[i]
			log.Printf("Ask %d: Price: %.2f, Size: %.8f", i+1, ask[0], ask[1])
		}
		for i := 0; i < 5 && i < len(orderbookMsg.Data.Bids); i++ {
			bid := orderbookMsg.Data.Bids[i]
			log.Printf("Bid %d: Price: %.2f, Size: %.8f", i+1, bid[0], bid[1])
		}
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to orderbook: %v", err)
	}

	// Subscribe to tickers for multiple instruments
	err = service.SubscribeTickersMulti([]string{"BTC-USDT", "ETH-USDT", "SOL-USDT"}, func(msg []byte) {
		log.Printf("Received ticker message: %s", string(msg))
		var tickerMsg ws.TickerMessage
		if err := json.Unmarshal(msg, &tickerMsg); err != nil {
			log.Printf("Failed to unmarshal ticker message: %v", err)
			return
		}
		for _, ticker := range tickerMsg.Data {
			log.Printf("Ticker: %s Last:%s Bid:%s Ask:%s Vol24h:%s", ticker.InstId, ticker.Last, ticker.BidPrice, ticker.AskPrice, ticker.Vol24h)
		}
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to tickers: %v", err)
	}

	// Subscribe to funding rate for multiple instruments
	err = service.SubscribeFundingRateMulti([]string{"BTC-USDT", "ETH-USDT", "SOL-USDT"}, func(msg []byte) {
		log.Printf("Received funding rate message: %s", string(msg))
		var fundingRateMsg ws.FundingRateMessage
		if err := json.Unmarshal(msg, &fundingRateMsg); err != nil {
			log.Printf("Failed to unmarshal funding rate message: %v", err)
			return
		}
		for _, rate := range fundingRateMsg.Data {
			log.Printf("Funding Rate: %s Rate:%s Time:%s", rate.InstId, rate.FundingRate, rate.FundingTime)
		}
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to funding rate: %v", err)
	}

	log.Println("Subscribed to all channels. Waiting for messages...")

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Unsubscribing from all channels...")

	// Unsubscribe from all channels
	service.Unsubscribe(ws.ChannelTrades, "BTC-USDT")
	service.Unsubscribe(ws.ChannelTrades, "ETH-USDT")
	service.Unsubscribe(ws.ChannelTrades, "SOL-USDT")
	service.Unsubscribe(ws.CandleInterval1m, "BTC-USDT")
	service.Unsubscribe(ws.CandleInterval1m, "ETH-USDT")
	service.Unsubscribe(ws.CandleInterval1m, "SOL-USDT")
	service.Unsubscribe(ws.ChannelOrderbook, "BTC-USDT")
	service.Unsubscribe(ws.ChannelOrderbook, "ETH-USDT")
	service.Unsubscribe(ws.ChannelOrderbook, "SOL-USDT")
	service.Unsubscribe(ws.ChannelTickers, "BTC-USDT")
	service.Unsubscribe(ws.ChannelTickers, "ETH-USDT")
	service.Unsubscribe(ws.ChannelTickers, "SOL-USDT")
	service.Unsubscribe(ws.ChannelFundingRate, "BTC-USDT")
	service.Unsubscribe(ws.ChannelFundingRate, "ETH-USDT")
	service.Unsubscribe(ws.ChannelFundingRate, "SOL-USDT")

	// Wait for messages to be processed
	time.Sleep(time.Second)
	log.Println("Done")
}
