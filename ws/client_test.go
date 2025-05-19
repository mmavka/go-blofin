package ws

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWSClientBasic(t *testing.T) {
	// TODO: реализовать тесты для login, subscribe, маршрутизации сообщений с использованием мок-соединения
}

func TestWSClientPing(t *testing.T) {
	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
		}
		defer conn.Close()
		for {
			mt, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if string(msg) == "ping" {
				conn.WriteMessage(mt, []byte("pong"))
			}
		}
	}))
	defer server.Close()

	url := "ws" + server.URL[len("http"):] // http:// -> ws://
	client := NewClient(url)
	err := client.Connect()
	if err != nil {
		t.Fatalf("connect error: %v", err)
	}
	defer client.Close()

	err = client.Ping()
	if err != nil {
		t.Errorf("ping error: %v", err)
	}
}

func TestWSClientSubscribe(t *testing.T) {
	upgrader := websocket.Upgrader{}
	ch := make(chan []byte, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
		}
		defer conn.Close()
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("read error: %v", err)
		}
		ch <- msg
		// не отвечаем, просто завершаем
	}))
	defer server.Close()

	url := "ws" + server.URL[len("http"):]
	client := NewClient(url)
	err := client.Connect()
	if err != nil {
		t.Fatalf("connect error: %v", err)
	}
	defer client.Close()

	channels := []ChannelArgs{{Channel: "trades", InstId: "BTC-USDT"}}
	err = client.Subscribe(channels)
	if err != nil {
		t.Errorf("subscribe error: %v", err)
	}

	select {
	case msg := <-ch:
		if len(msg) == 0 {
			t.Error("empty subscribe message")
		}
	case <-time.After(time.Second):
		t.Error("timeout waiting for subscribe message")
	}
}

func TestWSClientPushTrades(t *testing.T) {
	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
		}
		defer conn.Close()
		// Ждём подписки, затем отправляем push trades
		_, _, _ = conn.ReadMessage()
		push := `{"arg":{"channel":"trades","instId":"BTC-USDT"},"data":[{"tradeId":"1","price":"100","size":"0.1","side":"buy","ts":"123456"}]}`
		conn.WriteMessage(websocket.TextMessage, []byte(push))
		// ждём немного, чтобы клиент успел обработать
		time.Sleep(100 * time.Millisecond)
	}))
	defer server.Close()

	url := "ws" + server.URL[len("http"):]
	client := NewClient(url)
	err := client.Connect()
	if err != nil {
		t.Fatalf("connect error: %v", err)
	}
	defer client.Close()

	channels := []ChannelArgs{{Channel: "trades", InstId: "BTC-USDT"}}
	err = client.Subscribe(channels)
	if err != nil {
		t.Fatalf("subscribe error: %v", err)
	}

	select {
	case trade := <-client.Trades():
		if len(trade.Data) != 1 || trade.Data[0].Price != "100" {
			t.Errorf("unexpected trade: %+v", trade)
		}
	case <-time.After(time.Second):
		t.Error("timeout waiting for trade push")
	}
}

func TestWSClientPushCandles(t *testing.T) {
	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
		}
		defer conn.Close()
		// Ждём подписки, затем отправляем push candles
		_, _, _ = conn.ReadMessage()
		push := `{"arg":{"channel":"candle1m","instId":"BTC-USDT"},"data":[["123456","100","110","90","105","10","1000","1000","1"]]}`
		conn.WriteMessage(websocket.TextMessage, []byte(push))
		// ждём немного, чтобы клиент успел обработать
		time.Sleep(100 * time.Millisecond)
	}))
	defer server.Close()

	url := "ws" + server.URL[len("http"):]
	client := NewClient(url)
	err := client.Connect()
	if err != nil {
		t.Fatalf("connect error: %v", err)
	}
	defer client.Close()

	channels := []ChannelArgs{{Channel: "candle1m", InstId: "BTC-USDT"}}
	err = client.Subscribe(channels)
	if err != nil {
		t.Fatalf("subscribe error: %v", err)
	}

	select {
	case candle := <-client.Candles():
		if len(candle.Data) != 1 || candle.Data[0].Open != "100" {
			t.Errorf("unexpected candle: %+v", candle)
		}
	case <-time.After(time.Second):
		t.Error("timeout waiting for candle push")
	}
}

func TestWSClientPushTickers(t *testing.T) {
	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
		}
		defer conn.Close()
		_, _, _ = conn.ReadMessage()
		push := `{"arg":{"channel":"tickers","instId":"BTC-USDT"},"data":[{"instId":"BTC-USDT","last":"100"}]}`
		conn.WriteMessage(websocket.TextMessage, []byte(push))
		time.Sleep(100 * time.Millisecond)
	}))
	defer server.Close()

	url := "ws" + server.URL[len("http"):]
	client := NewClient(url)
	err := client.Connect()
	if err != nil {
		t.Fatalf("connect error: %v", err)
	}
	defer client.Close()

	channels := []ChannelArgs{{Channel: "tickers", InstId: "BTC-USDT"}}
	err = client.Subscribe(channels)
	if err != nil {
		t.Fatalf("subscribe error: %v", err)
	}

	select {
	case ticker := <-client.Tickers():
		if len(ticker.Data) != 1 || ticker.Data[0].Last != "100" {
			t.Errorf("unexpected ticker: %+v", ticker)
		}
	case <-time.After(time.Second):
		t.Error("timeout waiting for ticker push")
	}
}

func TestWSClientPushFundingRate(t *testing.T) {
	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
		}
		defer conn.Close()
		_, _, _ = conn.ReadMessage()
		push := `{"arg":{"channel":"funding-rate","instId":"BTC-USDT"},"data":[{"fundingRate":"0.0001","fundingTime":"1700726400000","instId":"BTC-USDT"}]}`
		conn.WriteMessage(websocket.TextMessage, []byte(push))
		time.Sleep(100 * time.Millisecond)
	}))
	defer server.Close()

	url := "ws" + server.URL[len("http"):]
	client := NewClient(url)
	err := client.Connect()
	if err != nil {
		t.Fatalf("connect error: %v", err)
	}
	defer client.Close()

	channels := []ChannelArgs{{Channel: "funding-rate", InstId: "BTC-USDT"}}
	err = client.Subscribe(channels)
	if err != nil {
		t.Fatalf("subscribe error: %v", err)
	}

	select {
	case fr := <-client.FundingRates():
		if len(fr.Data) != 1 || fr.Data[0].FundingRate != "0.0001" {
			t.Errorf("unexpected funding rate: %+v", fr)
		}
	case <-time.After(time.Second):
		t.Error("timeout waiting for funding-rate push")
	}
}
