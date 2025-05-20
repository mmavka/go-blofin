package ws

import (
	"encoding/json"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestMarketService_Subscribe(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	config := &WsConfig{Endpoint: ts.URL()}
	service := NewMarketService(config)

	err := service.Connect()
	assert.NoError(t, err)

	// Subscribe to trades
	err = service.SubscribeTrades("BTC-USDT", func(msg []byte) {})
	assert.NoError(t, err)

	// Check sent message
	msg := <-ts.messages
	var req SubscribeRequest
	err = json.Unmarshal(msg, &req)
	assert.NoError(t, err)
	assert.Equal(t, OpSubscribe, req.Op)
	assert.Equal(t, ChannelTrades, req.Args[0].Channel)
	assert.Equal(t, "BTC-USDT", req.Args[0].InstId)

	// Subscribe to candles
	err = service.SubscribeCandles("BTC-USDT", CandleInterval1m, func(msg []byte) {})
	assert.NoError(t, err)

	// Check sent message
	msg = <-ts.messages
	err = json.Unmarshal(msg, &req)
	assert.NoError(t, err)
	assert.Equal(t, OpSubscribe, req.Op)
	assert.Equal(t, CandleInterval1m, req.Args[0].Channel)
	assert.Equal(t, "BTC-USDT", req.Args[0].InstId)

	// Subscribe to orderbook
	err = service.SubscribeOrderbook("BTC-USDT", func(msg []byte) {})
	assert.NoError(t, err)

	// Check sent message
	msg = <-ts.messages
	err = json.Unmarshal(msg, &req)
	assert.NoError(t, err)
	assert.Equal(t, OpSubscribe, req.Op)
	assert.Equal(t, ChannelOrderbook, req.Args[0].Channel)
	assert.Equal(t, "BTC-USDT", req.Args[0].InstId)

	// Subscribe to tickers
	err = service.SubscribeTickers("BTC-USDT", func(msg []byte) {})
	assert.NoError(t, err)

	// Check sent message
	msg = <-ts.messages
	err = json.Unmarshal(msg, &req)
	assert.NoError(t, err)
	assert.Equal(t, OpSubscribe, req.Op)
	assert.Equal(t, ChannelTickers, req.Args[0].Channel)
	assert.Equal(t, "BTC-USDT", req.Args[0].InstId)

	// Subscribe to funding rate
	err = service.SubscribeFundingRate("BTC-USDT", func(msg []byte) {})
	assert.NoError(t, err)

	// Check sent message
	msg = <-ts.messages
	err = json.Unmarshal(msg, &req)
	assert.NoError(t, err)
	assert.Equal(t, OpSubscribe, req.Op)
	assert.Equal(t, ChannelFundingRate, req.Args[0].Channel)
	assert.Equal(t, "BTC-USDT", req.Args[0].InstId)

	// Try to subscribe to the same channel
	err = service.SubscribeTrades("BTC-USDT", func(msg []byte) {})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already subscribed")

	service.Close()
}

func TestMarketService_Unsubscribe(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	config := &WsConfig{Endpoint: ts.URL()}
	service := NewMarketService(config)

	err := service.Connect()
	assert.NoError(t, err)

	// Subscribe to trades
	err = service.SubscribeTrades("BTC-USDT", func(msg []byte) {})
	assert.NoError(t, err)

	// Check sent message
	msg := <-ts.messages
	var req SubscribeRequest
	err = json.Unmarshal(msg, &req)
	assert.NoError(t, err)
	assert.Equal(t, OpSubscribe, req.Op)
	assert.Equal(t, ChannelTrades, req.Args[0].Channel)
	assert.Equal(t, "BTC-USDT", req.Args[0].InstId)

	// Send subscription response
	subResp := SubscribeResponse{
		Event: "subscribe",
		Arg: ChannelArgs{
			Channel: ChannelTrades,
			InstId:  "BTC-USDT",
		},
		Code: "0",
	}
	ts.conn.WriteJSON(subResp)

	// Unsubscribe from trades
	err = service.Unsubscribe(ChannelTrades, "BTC-USDT")
	assert.NoError(t, err)

	// Check sent message
	msg = <-ts.messages
	err = json.Unmarshal(msg, &req)
	assert.NoError(t, err)
	assert.Equal(t, OpUnsubscribe, req.Op)
	assert.Equal(t, ChannelTrades, req.Args[0].Channel)
	assert.Equal(t, "BTC-USDT", req.Args[0].InstId)

	// Try to unsubscribe from non-subscribed channel
	err = service.Unsubscribe(ChannelTrades, "BTC-USDT")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not subscribed")

	service.Close()
}

func TestMarketService_MessageHandling(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	config := &WsConfig{Endpoint: ts.URL()}
	service := NewMarketService(config)

	messages := make(chan []byte, 1)
	err := service.Connect()
	assert.NoError(t, err)

	// Subscribe to trades
	err = service.SubscribeTrades("BTC-USDT", func(msg []byte) {
		messages <- msg
	})
	assert.NoError(t, err)

	// Send subscription response
	subResp := SubscribeResponse{
		Event: "subscribe",
		Arg: ChannelArgs{
			Channel: ChannelTrades,
			InstId:  "BTC-USDT",
		},
		Code: "0",
	}
	ts.conn.WriteJSON(subResp)

	// Send test message
	testMsg := []byte(`{"arg":{"channel":"trades","instId":"BTC-USDT"},"data":[{"instId":"BTC-USDT","tradeId":"123","price":"50000","size":"1","side":"buy","ts":"1234567890"}]}`)
	ts.conn.WriteMessage(websocket.TextMessage, testMsg)

	// Check received message
	msg := <-messages
	assert.Equal(t, testMsg, msg)

	// Parse trade message
	var tradeMsg TradeMessage
	err = json.Unmarshal(msg, &tradeMsg)
	assert.NoError(t, err)
	assert.Equal(t, ChannelTrades, tradeMsg.Arg.Channel)
	assert.Equal(t, "BTC-USDT", tradeMsg.Arg.InstId)
	assert.Len(t, tradeMsg.Data, 1)
	assert.Equal(t, "BTC-USDT", tradeMsg.Data[0].InstId)
	assert.Equal(t, "123", tradeMsg.Data[0].TradeId)
	assert.Equal(t, "50000", tradeMsg.Data[0].Price)
	assert.Equal(t, "1", tradeMsg.Data[0].Size)
	assert.Equal(t, "buy", tradeMsg.Data[0].Side)
	assert.Equal(t, "1234567890", tradeMsg.Data[0].Ts)

	service.Close()
}

func TestMarketService_CandleHandling(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	config := &WsConfig{Endpoint: ts.URL()}
	service := NewMarketService(config)

	messages := make(chan []byte, 1)
	err := service.Connect()
	assert.NoError(t, err)

	// Subscribe to candles
	err = service.SubscribeCandles("BTC-USDT", CandleInterval1d, func(msg []byte) {
		messages <- msg
	})
	assert.NoError(t, err)

	// Check sent message
	msg := <-ts.messages
	var req SubscribeRequest
	err = json.Unmarshal(msg, &req)
	assert.NoError(t, err)
	assert.Equal(t, OpSubscribe, req.Op)
	assert.Equal(t, CandleInterval1d, req.Args[0].Channel)
	assert.Equal(t, "BTC-USDT", req.Args[0].InstId)

	// Send subscription response
	subResp := SubscribeResponse{
		Event: "subscribe",
		Arg: ChannelArgs{
			Channel: CandleInterval1d,
			InstId:  "BTC-USDT",
		},
		Code: "0",
	}
	ts.conn.WriteJSON(subResp)

	// Send test message
	testMsg := []byte(`{"arg":{"channel":"candle1D","instId":"BTC-USDT"},"data":[["1696636800000","27491.5","27495","27483","27489.5","95359","95.359","2621407.651","0"]]}`)
	ts.conn.WriteMessage(websocket.TextMessage, testMsg)

	// Check received message
	msg = <-messages
	assert.Equal(t, testMsg, msg)

	// Parse candle message
	var candleMsg CandleMessage
	err = json.Unmarshal(msg, &candleMsg)
	assert.NoError(t, err)
	assert.Equal(t, CandleInterval1d, candleMsg.Arg.Channel)
	assert.Equal(t, "BTC-USDT", candleMsg.Arg.InstId)
	assert.Len(t, candleMsg.Data, 1)
	assert.Len(t, candleMsg.Data[0], 9)
	assert.Equal(t, "1696636800000", candleMsg.Data[0][0]) // ts
	assert.Equal(t, "27491.5", candleMsg.Data[0][1])       // open
	assert.Equal(t, "27495", candleMsg.Data[0][2])         // high
	assert.Equal(t, "27483", candleMsg.Data[0][3])         // low
	assert.Equal(t, "27489.5", candleMsg.Data[0][4])       // close
	assert.Equal(t, "95359", candleMsg.Data[0][5])         // vol
	assert.Equal(t, "95.359", candleMsg.Data[0][6])        // volCurrency
	assert.Equal(t, "2621407.651", candleMsg.Data[0][7])   // volCurrencyQuote
	assert.Equal(t, "0", candleMsg.Data[0][8])             // confirm

	service.Close()
}

func TestMarketService_OrderBookHandling(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	config := &WsConfig{Endpoint: ts.URL()}
	service := NewMarketService(config)

	messages := make(chan []byte, 1)
	err := service.Connect()
	assert.NoError(t, err)

	// Subscribe to orderbook
	err = service.SubscribeOrderbook("BTC-USDT", func(msg []byte) {
		messages <- msg
	})
	assert.NoError(t, err)

	// Check sent message
	msg := <-ts.messages
	var req SubscribeRequest
	err = json.Unmarshal(msg, &req)
	assert.NoError(t, err)
	assert.Equal(t, OpSubscribe, req.Op)
	assert.Equal(t, ChannelOrderbook, req.Args[0].Channel)
	assert.Equal(t, "BTC-USDT", req.Args[0].InstId)

	// Send subscription response
	subResp := SubscribeResponse{
		Event: "subscribe",
		Arg: ChannelArgs{
			Channel: ChannelOrderbook,
			InstId:  "BTC-USDT",
		},
		Code: "0",
	}
	ts.conn.WriteJSON(subResp)

	// Send snapshot message
	snapshotMsg := []byte(`{"arg":{"channel":"books","instId":"BTC-USDT"},"action":"snapshot","data":{"asks":[[27491.5,392],[27495,541]],"bids":[[27489.5,6817],[27483,4744]],"ts":"1696670727520","prevSeqId":"0","seqId":"107600747"}}`)
	ts.conn.WriteMessage(websocket.TextMessage, snapshotMsg)

	// Check received message
	msg = <-messages
	assert.Equal(t, snapshotMsg, msg)

	// Parse orderbook message
	var orderbookMsg OrderBookMessage
	err = json.Unmarshal(msg, &orderbookMsg)
	assert.NoError(t, err)
	assert.Equal(t, ChannelOrderbook, orderbookMsg.Arg.Channel)
	assert.Equal(t, "BTC-USDT", orderbookMsg.Arg.InstId)
	assert.Equal(t, "snapshot", orderbookMsg.Action)
	assert.Len(t, orderbookMsg.Data.Asks, 2)
	assert.Len(t, orderbookMsg.Data.Bids, 2)
	assert.Equal(t, 27491.5, orderbookMsg.Data.Asks[0][0])
	assert.Equal(t, 392.0, orderbookMsg.Data.Asks[0][1])
	assert.Equal(t, 27495.0, orderbookMsg.Data.Asks[1][0])
	assert.Equal(t, 541.0, orderbookMsg.Data.Asks[1][1])
	assert.Equal(t, 27489.5, orderbookMsg.Data.Bids[0][0])
	assert.Equal(t, 6817.0, orderbookMsg.Data.Bids[0][1])
	assert.Equal(t, 27483.0, orderbookMsg.Data.Bids[1][0])
	assert.Equal(t, 4744.0, orderbookMsg.Data.Bids[1][1])
	assert.Equal(t, "1696670727520", orderbookMsg.Data.Ts)
	assert.Equal(t, "0", orderbookMsg.Data.PrevSeqId)
	assert.Equal(t, "107600747", orderbookMsg.Data.SeqId)

	// Send update message
	updateMsg := []byte(`{"arg":{"channel":"books","instId":"BTC-USDT"},"action":"update","data":{"asks":[[27495,2208],[27496,4605]],"bids":[[27489.5,7115],[27483,4791]],"ts":"1696670728525","prevSeqId":"107600747","seqId":"107600806"}}`)
	ts.conn.WriteMessage(websocket.TextMessage, updateMsg)

	// Check received message
	msg = <-messages
	assert.Equal(t, updateMsg, msg)

	// Parse orderbook message
	err = json.Unmarshal(msg, &orderbookMsg)
	assert.NoError(t, err)
	assert.Equal(t, ChannelOrderbook, orderbookMsg.Arg.Channel)
	assert.Equal(t, "BTC-USDT", orderbookMsg.Arg.InstId)
	assert.Equal(t, "update", orderbookMsg.Action)
	assert.Len(t, orderbookMsg.Data.Asks, 2)
	assert.Len(t, orderbookMsg.Data.Bids, 2)
	assert.Equal(t, 27495.0, orderbookMsg.Data.Asks[0][0])
	assert.Equal(t, 2208.0, orderbookMsg.Data.Asks[0][1])
	assert.Equal(t, 27496.0, orderbookMsg.Data.Asks[1][0])
	assert.Equal(t, 4605.0, orderbookMsg.Data.Asks[1][1])
	assert.Equal(t, 27489.5, orderbookMsg.Data.Bids[0][0])
	assert.Equal(t, 7115.0, orderbookMsg.Data.Bids[0][1])
	assert.Equal(t, 27483.0, orderbookMsg.Data.Bids[1][0])
	assert.Equal(t, 4791.0, orderbookMsg.Data.Bids[1][1])
	assert.Equal(t, "1696670728525", orderbookMsg.Data.Ts)
	assert.Equal(t, "107600747", orderbookMsg.Data.PrevSeqId)
	assert.Equal(t, "107600806", orderbookMsg.Data.SeqId)

	service.Close()
}

func TestMarketService_TickerHandling(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	config := &WsConfig{Endpoint: ts.URL()}
	service := NewMarketService(config)

	messages := make(chan []byte, 1)
	err := service.Connect()
	assert.NoError(t, err)

	// Subscribe to tickers
	err = service.SubscribeTickers("BTC-USDT", func(msg []byte) {
		messages <- msg
	})
	assert.NoError(t, err)

	// Check sent message
	msg := <-ts.messages
	var req SubscribeRequest
	err = json.Unmarshal(msg, &req)
	assert.NoError(t, err)
	assert.Equal(t, OpSubscribe, req.Op)
	assert.Equal(t, ChannelTickers, req.Args[0].Channel)
	assert.Equal(t, "BTC-USDT", req.Args[0].InstId)

	// Send subscription response
	subResp := SubscribeResponse{
		Event: "subscribe",
		Arg: ChannelArgs{
			Channel: ChannelTickers,
			InstId:  "BTC-USDT",
		},
		Code: "0",
	}
	ts.conn.WriteJSON(subResp)

	// Send test message
	testMsg := []byte(`{"arg":{"channel":"tickers","instId":"BTC-USDT"},"data":[{"instId":"BTC-USDT","last":"9999.99","lastSize":"0.1","askPrice":"9999.99","askSize":"11","bidPrice":"8888.88","bidSize":"5","open24h":"9000","high24h":"10000","low24h":"8888.88","volCurrency24h":"2222","vol24h":"2222","ts":"1597026383085"}]}`)
	ts.conn.WriteMessage(websocket.TextMessage, testMsg)

	// Check received message
	msg = <-messages
	assert.Equal(t, testMsg, msg)

	// Parse ticker message
	var tickerMsg TickerMessage
	err = json.Unmarshal(msg, &tickerMsg)
	assert.NoError(t, err)
	assert.Equal(t, ChannelTickers, tickerMsg.Arg.Channel)
	assert.Equal(t, "BTC-USDT", tickerMsg.Arg.InstId)
	assert.Len(t, tickerMsg.Data, 1)
	assert.Equal(t, "BTC-USDT", tickerMsg.Data[0].InstId)
	assert.Equal(t, "9999.99", tickerMsg.Data[0].Last)
	assert.Equal(t, "0.1", tickerMsg.Data[0].LastSize)
	assert.Equal(t, "9999.99", tickerMsg.Data[0].AskPrice)
	assert.Equal(t, "11", tickerMsg.Data[0].AskSize)
	assert.Equal(t, "8888.88", tickerMsg.Data[0].BidPrice)
	assert.Equal(t, "5", tickerMsg.Data[0].BidSize)
	assert.Equal(t, "9000", tickerMsg.Data[0].Open24h)
	assert.Equal(t, "10000", tickerMsg.Data[0].High24h)
	assert.Equal(t, "8888.88", tickerMsg.Data[0].Low24h)
	assert.Equal(t, "2222", tickerMsg.Data[0].VolCurrency24h)
	assert.Equal(t, "2222", tickerMsg.Data[0].Vol24h)
	assert.Equal(t, "1597026383085", tickerMsg.Data[0].Ts)

	service.Close()
}

func TestMarketService_FundingRateHandling(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	config := &WsConfig{Endpoint: ts.URL()}
	service := NewMarketService(config)

	messages := make(chan []byte, 1)
	err := service.Connect()
	assert.NoError(t, err)

	// Subscribe to funding rate
	err = service.SubscribeFundingRate("BTC-USDT", func(msg []byte) {
		messages <- msg
	})
	assert.NoError(t, err)

	// Check sent message
	msg := <-ts.messages
	var req SubscribeRequest
	err = json.Unmarshal(msg, &req)
	assert.NoError(t, err)
	assert.Equal(t, OpSubscribe, req.Op)
	assert.Equal(t, ChannelFundingRate, req.Args[0].Channel)
	assert.Equal(t, "BTC-USDT", req.Args[0].InstId)

	// Send subscription response
	subResp := SubscribeResponse{
		Event: "subscribe",
		Arg: ChannelArgs{
			Channel: ChannelFundingRate,
			InstId:  "BTC-USDT",
		},
		Code: "0",
	}
	ts.conn.WriteJSON(subResp)

	// Send test message
	testMsg := []byte(`{"arg":{"channel":"funding-rate","instId":"BTC-USDT"},"data":[{"fundingRate":"0.0001875391284828","fundingTime":"1700726400000","instId":"BTC-USDT"}]}`)
	ts.conn.WriteMessage(websocket.TextMessage, testMsg)

	// Check received message
	msg = <-messages
	assert.Equal(t, testMsg, msg)

	// Parse funding rate message
	var fundingRateMsg FundingRateMessage
	err = json.Unmarshal(msg, &fundingRateMsg)
	assert.NoError(t, err)
	assert.Equal(t, ChannelFundingRate, fundingRateMsg.Arg.Channel)
	assert.Equal(t, "BTC-USDT", fundingRateMsg.Arg.InstId)
	assert.Len(t, fundingRateMsg.Data, 1)
	assert.Equal(t, "BTC-USDT", fundingRateMsg.Data[0].InstId)
	assert.Equal(t, "0.0001875391284828", fundingRateMsg.Data[0].FundingRate)
	assert.Equal(t, "1700726400000", fundingRateMsg.Data[0].FundingTime)

	service.Close()
}

func TestMarketService_SubscribeMulti(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	config := &WsConfig{Endpoint: ts.URL()}
	service := NewMarketService(config)

	err := service.Connect()
	assert.NoError(t, err)

	// Subscribe to trades for multiple instruments
	err = service.SubscribeTradesMulti([]string{"BTC-USDT", "ETH-USDT"}, func(msg []byte) {})
	assert.NoError(t, err)

	// Check sent message
	msg := <-ts.messages
	var req SubscribeRequest
	err = json.Unmarshal(msg, &req)
	assert.NoError(t, err)
	assert.Equal(t, OpSubscribe, req.Op)
	assert.Len(t, req.Args, 2)
	assert.Equal(t, ChannelTrades, req.Args[0].Channel)
	assert.Equal(t, "BTC-USDT", req.Args[0].InstId)
	assert.Equal(t, ChannelTrades, req.Args[1].Channel)
	assert.Equal(t, "ETH-USDT", req.Args[1].InstId)

	// Try to subscribe to already subscribed instruments
	err = service.SubscribeTradesMulti([]string{"BTC-USDT", "ETH-USDT"}, func(msg []byte) {})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already subscribed")

	service.Close()
}
