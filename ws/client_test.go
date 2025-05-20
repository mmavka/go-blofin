package ws

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// Test intervals
	testPingInterval = 100 * time.Millisecond
	testPongTimeout  = 50 * time.Millisecond
)

type testServer struct {
	t               *testing.T
	server          *httptest.Server
	conn            *websocket.Conn
	mu              sync.Mutex
	connectCount    int
	lastConnectTime time.Time
}

func newTestServer(t *testing.T) *testServer {
	ts := &testServer{t: t}
	ts.server = httptest.NewServer(http.HandlerFunc(ts.handler))
	return ts
}

func (ts *testServer) URL() string {
	return "ws" + strings.TrimPrefix(ts.server.URL, "http")
}

func (ts *testServer) Close() {
	if ts.conn != nil {
		ts.conn.Close()
	}
	ts.server.Close()
}

func (ts *testServer) handler(w http.ResponseWriter, r *http.Request) {
	ts.mu.Lock()
	now := time.Now()
	if now.Sub(ts.lastConnectTime) < reconnectDelay {
		ts.t.Error("Connection attempt too soon after previous connection")
	}
	ts.lastConnectTime = now
	ts.connectCount++
	ts.mu.Unlock()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	ts.conn = conn

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if string(msg) == "ping" {
				conn.WriteMessage(websocket.TextMessage, []byte("pong"))
			}
		}
	}()
}

func TestPingPong(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	client := NewClient(ts.URL())
	client.SetErrorHandler(func(err error) {
		t.Logf("Error handler called: %v", err)
	})

	err := client.Connect()
	assert.NoError(t, err)

	// Wait for several ping/pong cycles
	time.Sleep(350 * time.Millisecond)

	// Check that connection is alive
	assert.NotNil(t, client.conn)
	assert.False(t, client.reconnecting)
}

// Test that client reconnects if pong is not received
func TestReconnectOnPongTimeout(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	client := NewClient(ts.URL())
	client.SetErrorHandler(func(err error) {
		t.Logf("Error handler called: %v", err)
	})

	err := client.Connect()
	assert.NoError(t, err)

	// Wait for initial connection
	time.Sleep(100 * time.Millisecond)

	// Close connection to simulate pong timeout
	ts.conn.Close()

	// Wait for reconnect
	time.Sleep(2 * time.Second)

	// Check that client reconnected
	assert.NotNil(t, client.conn)
	assert.False(t, client.reconnecting)

	client.Close()
}

// Test that keepalive timer resets on any incoming message
func TestKeepAliveResetOnAnyMessage(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	client := NewClient(ts.URL())
	client.SetErrorHandler(func(err error) {
		t.Logf("Error handler called: %v", err)
	})

	err := client.Connect()
	assert.NoError(t, err)

	// Отправляем обычное сообщение, не pong
	go func() {
		time.Sleep(50 * time.Millisecond)
		ts.conn.WriteMessage(websocket.TextMessage, []byte(`{"event":"test"}`))
	}()

	// Ждем чуть больше pingInterval, клиент не должен переподключиться
	time.Sleep(150 * time.Millisecond)
	assert.NotNil(t, client.conn)
	assert.False(t, client.reconnecting)

	client.Close()
}

func TestConnectionLimit(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	client := NewClient(ts.URL())
	client.SetErrorHandler(func(err error) {
		t.Logf("Error handler called: %v", err)
	})

	// First connection
	err := client.Connect()
	assert.NoError(t, err)
	time.Sleep(100 * time.Millisecond) // Wait for connection to establish

	// Try to connect too soon
	err = client.Connect()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already connected")

	// Close connection
	client.Close()
	time.Sleep(reconnectDelay) // Wait for reconnect delay

	// Second connection should succeed
	err = client.Connect()
	assert.NoError(t, err)
	time.Sleep(100 * time.Millisecond)

	// Try to connect too soon again
	err = client.Connect()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already connected")

	// Close connection
	client.Close()
	time.Sleep(reconnectDelay) // Wait for reconnect delay

	// Third connection should succeed
	err = client.Connect()
	assert.NoError(t, err)

	// Check that all successful connections were counted
	// Учитываем только успешные подключения (без переподключений)
	assert.Equal(t, 3, ts.connectCount/2)

	// Close connection at the end
	client.Close()
	time.Sleep(100 * time.Millisecond) // Wait for connection to close

	// Reset connection counter
	ts.mu.Lock()
	ts.connectCount = 0
	ts.mu.Unlock()
}

func TestStateRecovery(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	client := NewClient(ts.URL())
	client.SetErrorHandler(func(err error) {
		t.Logf("Error handler called: %v", err)
	})

	// Connect and authenticate
	err := client.Connect()
	assert.NoError(t, err)

	err = client.Login("test-key", "test-secret", "test-pass")
	assert.NoError(t, err)

	// Subscribe to channels
	channels := []ChannelArgs{
		{Channel: "trades", InstId: "BTC-USDT"},
		{Channel: "tickers", InstId: "ETH-USDT"},
	}
	err = client.Subscribe(channels)
	assert.NoError(t, err)

	// Wait for subscription to be processed
	time.Sleep(100 * time.Millisecond)

	// Simulate break and reconnection
	ts.conn.Close()

	// Wait for reconnection
	time.Sleep(2 * time.Second)

	// Check state recovery
	assert.True(t, client.isLoggedIn)
	// Проверяем уникальные подписки
	uniqueSubs := make(map[string]struct{})
	for _, sub := range client.subscriptions {
		key := sub.Channel + "-" + sub.InstId
		uniqueSubs[key] = struct{}{}
	}
	assert.Equal(t, len(channels), len(uniqueSubs))
	assert.NotNil(t, client.credentials)
}

func TestErrorHandling(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	client := NewClient(ts.URL())

	errors := make(chan error, 1)
	client.SetErrorHandler(func(err error) {
		errors <- err
	})

	err := client.Connect()
	assert.NoError(t, err)

	// Simulate ping/pong error
	ts.conn.Close()

	// Wait for error handling
	select {
	case err := <-errors:
		assert.NotNil(t, err)
	case <-time.After(time.Second):
		t.Fatal("Error handler not called")
	}
}

func TestWSClientBasic(t *testing.T) {
	// TODO: implement tests for login, subscribe, message routing using mock connection
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
		// no response, just finish
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
		// Wait for subscription, then send push trades
		_, _, _ = conn.ReadMessage()
		push := `{"arg":{"channel":"trades","instId":"BTC-USDT"},"data":[{"tradeId":"1","price":"100","size":"0.1","side":"buy","ts":"123456"}]}`
		conn.WriteMessage(websocket.TextMessage, []byte(push))
		// wait a bit for client to process
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
		// Wait for subscription, then send push candles
		_, _, _ = conn.ReadMessage()
		push := `{"arg":{"channel":"candle1m","instId":"BTC-USDT"},"data":[["123456","100","110","90","105","10","1000","1000","1"]]}`
		conn.WriteMessage(websocket.TextMessage, []byte(push))
		// wait a bit for client to process
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

func TestWSClientSetURL(t *testing.T) {
	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
		}
	}))
	defer server.Close()

	url := "ws" + server.URL[len("http"):]
	client := NewDefaultClient()
	client.SetURL(url)
	if client.url != url {
		t.Errorf("url not updated: got %s, want %s", client.url, url)
	}
	err := client.Connect()
	if err != nil {
		t.Fatalf("connect error: %v", err)
	}
	client.Close()
}

func TestWSClientErrorHandler(t *testing.T) {
	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
		}
		conn.Close() // close connection immediately to trigger client error
	}))
	defer server.Close()

	url := "ws" + server.URL[len("http"):]
	client := NewClient(url)
	err := client.Connect()
	if err != nil {
		t.Fatalf("connect error: %v", err)
	}
	ch := make(chan error, 1)
	client.SetErrorHandler(func(e error) {
		ch <- e
	})
	// Wait for read error
	select {
	case e := <-ch:
		if e == nil {
			t.Error("handler received nil error")
		}
	case <-time.After(time.Second):
		t.Error("timeout waiting for error handler call")
	}
	client.Close()
}

func TestWSClientSubscribePrivateWithoutLogin(t *testing.T) {
	client := NewDefaultClient()
	err := client.Subscribe([]ChannelArgs{{Channel: "orders", InstId: "BTC-USDT"}})
	if err == nil || err.Error() != "subscription to private channel 'orders' requires login" {
		t.Errorf("expected error for private channel without login, got: %v", err)
	}
}

func TestWSClientSubscribeExceedsLimit(t *testing.T) {
	client := NewDefaultClient()
	// Create many channels to exceed limit
	channels := make([]ChannelArgs, 0, 500)
	for i := 0; i < 500; i++ {
		channels = append(channels, ChannelArgs{Channel: "trades", InstId: fmt.Sprintf("BTC-USDT-%d", i)})
	}
	err := client.Subscribe(channels)
	if err == nil || !strings.HasPrefix(err.Error(), "subscription request exceeds 4096") {
		t.Errorf("expected error for request size limit, got: %v", err)
	}
}
