package ws

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestHandler_SubscribePositions(t *testing.T) {
	positionReceived := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatal(err)
		}

		// Read subscription request
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Fatal(err)
		}

		var req PositionsRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			t.Fatal(err)
		}

		// Verify request
		if req.Op != "subscribe" {
			t.Errorf("expected op=subscribe, got %s", req.Op)
		}
		if len(req.Args) != 1 {
			t.Errorf("expected 1 arg, got %d", len(req.Args))
		}
		if req.Args[0].Channel != "positions" {
			t.Errorf("expected channel=positions, got %s", req.Args[0].Channel)
		}
		if req.Args[0].InstId != "ETH-USDT" {
			t.Errorf("expected instId=ETH-USDT, got %s", req.Args[0].InstId)
		}

		// Send subscription response
		resp := PositionsResponse{
			Event: "subscribe",
			Arg: struct {
				Channel string `json:"channel"`
				InstId  string `json:"instId,omitempty"`
			}{
				Channel: "positions",
				InstId:  "ETH-USDT",
			},
		}
		respBytes, _ := json.Marshal(resp)
		if err := conn.WriteMessage(websocket.TextMessage, respBytes); err != nil {
			t.Fatal(err)
		}

		// Send position data
		positionsData := PositionsData{
			Arg: struct {
				Channel string `json:"channel"`
			}{
				Channel: "positions",
			},
			Data: []Position{
				{
					InstType:           "SWAP",
					InstId:             "ETH-USDT",
					MarginMode:         "cross",
					PositionId:         "8138",
					PositionSide:       "net",
					Positions:          "-100",
					AvailablePositions: "-100",
					AveragePrice:       "130.06",
					UnrealizedPnl:      "-77.1",
					UnrealizedPnlRatio: "-1.778409964631708442",
					Leverage:           "3",
					LiquidationPrice:   "107929.699398660166170462",
					MarkPrice:          "207.16",
					InitialMargin:      "69.053333333333333333",
					Margin:             "",
					MarginRatio:        "131.337873621866389829",
					MaintenanceMargin:  "1.0358",
					Adl:                "3",
					CreateTime:         "1695795726481",
					UpdateTime:         "1695795726484",
				},
			},
		}
		positionsBytes, _ := json.Marshal(positionsData)
		if err := conn.WriteMessage(websocket.TextMessage, positionsBytes); err != nil {
			t.Fatal(err)
		}

		// Wait for position to be received before closing
		<-positionReceived
		conn.Close()
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := NewHandler(conn)

	var receivedPosition Position
	err = handler.SubscribePositions("ETH-USDT", func(position Position) {
		receivedPosition = position
		close(positionReceived)
	})
	if err != nil {
		conn.Close()
		t.Fatal(err)
	}

	// Start message handling in a goroutine
	go func() {
		if err := handler.HandleMessages(); err != nil {
			// Ignore connection closed error
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("HandleMessages error: %v", err)
			}
		}
	}()

	// Wait for position data with timeout
	select {
	case <-positionReceived:
		// Position received successfully
	case <-time.After(1 * time.Second):
		conn.Close()
		t.Fatal("timeout waiting for position data")
	}

	// Verify received position
	if receivedPosition.InstId != "ETH-USDT" {
		t.Errorf("expected instId=ETH-USDT, got %s", receivedPosition.InstId)
	}
	if receivedPosition.Positions != "-100" {
		t.Errorf("expected positions=-100, got %s", receivedPosition.Positions)
	}

	// Close connection after all checks
	conn.Close()
}

func TestHandler_SubscribeOrders(t *testing.T) {
	orderReceived := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatal(err)
		}

		// Read subscription request
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Fatal(err)
		}

		var req OrdersRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			t.Fatal(err)
		}

		// Verify request
		if req.Op != "subscribe" {
			t.Errorf("expected op=subscribe, got %s", req.Op)
		}
		if len(req.Args) != 1 {
			t.Errorf("expected 1 arg, got %d", len(req.Args))
		}
		if req.Args[0].Channel != "orders" {
			t.Errorf("expected channel=orders, got %s", req.Args[0].Channel)
		}
		if req.Args[0].InstId != "ETH-USDT" {
			t.Errorf("expected instId=ETH-USDT, got %s", req.Args[0].InstId)
		}

		// Send subscription response
		resp := OrdersResponse{
			Event: "subscribe",
			Arg: struct {
				Channel string `json:"channel"`
				InstId  string `json:"instId,omitempty"`
			}{
				Channel: "orders",
				InstId:  "ETH-USDT",
			},
		}
		respBytes, _ := json.Marshal(resp)
		if err := conn.WriteMessage(websocket.TextMessage, respBytes); err != nil {
			t.Fatal(err)
		}

		// Send order data
		ordersData := OrdersData{
			Action: "snapshot",
			Arg: struct {
				Channel string `json:"channel"`
			}{
				Channel: "orders",
			},
			Data: []Order{
				{
					InstType:           "SWAP",
					InstId:             "ETH-USDT",
					OrderId:            "28334314",
					ClientOrderId:      "",
					Price:              "28000.000000000000000000",
					Size:               "10",
					OrderType:          "limit",
					Side:               "sell",
					PositionSide:       "net",
					MarginMode:         "cross",
					FilledSize:         "0",
					FilledAmount:       "0.000000000000000000",
					AveragePrice:       "0.000000000000000000",
					State:              "live",
					Leverage:           "2",
					TpTriggerPrice:     "27000.000000000000000000",
					TpTriggerPriceType: "last",
					TpOrderPrice:       "-1",
					SlTriggerPrice:     "",
					SlTriggerPriceType: "",
					SlOrderPrice:       "",
					Fee:                "0.000000000000000000",
					Pnl:                "0.000000000000000000",
					CancelSource:       "",
					OrderCategory:      "pre_tp_sl",
					CreateTime:         "1696760245931",
					UpdateTime:         "1696760245973",
					ReduceOnly:         "false",
					BrokerId:           "",
				},
			},
		}
		ordersBytes, _ := json.Marshal(ordersData)
		if err := conn.WriteMessage(websocket.TextMessage, ordersBytes); err != nil {
			t.Fatal(err)
		}

		// Wait for order to be received before closing
		<-orderReceived
		conn.Close()
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := NewHandler(conn)

	var receivedOrder Order
	err = handler.SubscribeOrders("ETH-USDT", func(order Order) {
		receivedOrder = order
		close(orderReceived)
	})
	if err != nil {
		conn.Close()
		t.Fatal(err)
	}

	// Start message handling in a goroutine
	go func() {
		if err := handler.HandleMessages(); err != nil {
			// Ignore connection closed error
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("HandleMessages error: %v", err)
			}
		}
	}()

	// Wait for order data with timeout
	select {
	case <-orderReceived:
		// Order received successfully
	case <-time.After(1 * time.Second):
		conn.Close()
		t.Fatal("timeout waiting for order data")
	}

	// Verify received order
	if receivedOrder.InstId != "ETH-USDT" {
		t.Errorf("expected instId=ETH-USDT, got %s", receivedOrder.InstId)
	}
	if receivedOrder.OrderId != "28334314" {
		t.Errorf("expected orderId=28334314, got %s", receivedOrder.OrderId)
	}
	if receivedOrder.State != "live" {
		t.Errorf("expected state=live, got %s", receivedOrder.State)
	}

	// Close connection after all checks
	conn.Close()
}

func TestHandler_SubscribeAlgoOrders(t *testing.T) {
	algoOrderReceived := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatal(err)
		}

		// Read subscription request
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Fatal(err)
		}

		var req AlgoOrdersRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			t.Fatal(err)
		}

		// Verify request
		if req.Op != "subscribe" {
			t.Errorf("expected op=subscribe, got %s", req.Op)
		}
		if len(req.Args) != 1 {
			t.Errorf("expected 1 arg, got %d", len(req.Args))
		}
		if req.Args[0].Channel != "orders-algo" {
			t.Errorf("expected channel=orders-algo, got %s", req.Args[0].Channel)
		}
		if req.Args[0].InstId != "ETH-USDT" {
			t.Errorf("expected instId=ETH-USDT, got %s", req.Args[0].InstId)
		}

		// Send subscription response
		resp := AlgoOrdersResponse{
			Event: "subscribe",
			Arg: struct {
				Channel string `json:"channel"`
				InstId  string `json:"instId,omitempty"`
			}{
				Channel: "orders-algo",
				InstId:  "ETH-USDT",
			},
		}
		respBytes, _ := json.Marshal(resp)
		if err := conn.WriteMessage(websocket.TextMessage, respBytes); err != nil {
			t.Fatal(err)
		}

		// Send algo order data
		algoOrdersData := AlgoOrdersData{
			Action: "snapshot",
			Arg: struct {
				Channel string `json:"channel"`
			}{
				Channel: "orders-algo",
			},
			Data: []AlgoOrder{
				{
					InstType:         "SWAP",
					InstId:           "ETH-USDT",
					TpslId:           "11779982",
					AlgoId:           "11779982",
					ClientOrderId:    "",
					Size:             "100",
					OrderType:        "conditional",
					Side:             "buy",
					PositionSide:     "long",
					MarginMode:       "cross",
					Leverage:         "10",
					State:            "live",
					TpTriggerPrice:   "73000.000000000000000000",
					TpOrderPrice:     "-1",
					SlTriggerPrice:   "",
					SlOrderPrice:     "",
					TriggerPrice:     "",
					TriggerPriceType: "last",
					OrderPrice:       "",
					ActualSize:       "",
					ActualSide:       "",
					ReduceOnly:       "false",
					CancelType:       "not_canceled",
					CreateTime:       "1731056529341",
					UpdateTime:       "1731056529341",
					Tag:              "",
					BrokerId:         "",
					AttachAlgoOrders: []AttachAlgoOrder{
						{
							TpTriggerPrice:     "75000",
							TpTriggerPriceType: "market",
							TpOrderPrice:       "-1",
							SlTriggerPriceType: "",
							SlTriggerPrice:     "",
							SlOrderPrice:       "",
						},
					},
				},
			},
		}
		algoOrdersBytes, _ := json.Marshal(algoOrdersData)
		if err := conn.WriteMessage(websocket.TextMessage, algoOrdersBytes); err != nil {
			t.Fatal(err)
		}

		// Wait for algo order to be received before closing
		<-algoOrderReceived
		conn.Close()
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := NewHandler(conn)

	var receivedAlgoOrder AlgoOrder
	err = handler.SubscribeAlgoOrders("ETH-USDT", func(order AlgoOrder) {
		receivedAlgoOrder = order
		close(algoOrderReceived)
	})
	if err != nil {
		conn.Close()
		t.Fatal(err)
	}

	// Start message handling in a goroutine
	go func() {
		if err := handler.HandleMessages(); err != nil {
			// Ignore connection closed error
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("HandleMessages error: %v", err)
			}
		}
	}()

	// Wait for algo order data with timeout
	select {
	case <-algoOrderReceived:
		// Algo order received successfully
	case <-time.After(1 * time.Second):
		conn.Close()
		t.Fatal("timeout waiting for algo order data")
	}

	// Verify received algo order
	if receivedAlgoOrder.InstId != "ETH-USDT" {
		t.Errorf("expected instId=ETH-USDT, got %s", receivedAlgoOrder.InstId)
	}
	if receivedAlgoOrder.AlgoId != "11779982" {
		t.Errorf("expected algoId=11779982, got %s", receivedAlgoOrder.AlgoId)
	}
	if receivedAlgoOrder.State != "live" {
		t.Errorf("expected state=live, got %s", receivedAlgoOrder.State)
	}
	if len(receivedAlgoOrder.AttachAlgoOrders) != 1 {
		t.Errorf("expected 1 attached algo order, got %d", len(receivedAlgoOrder.AttachAlgoOrders))
	}
	if receivedAlgoOrder.AttachAlgoOrders[0].TpTriggerPrice != "75000" {
		t.Errorf("expected tpTriggerPrice=75000, got %s", receivedAlgoOrder.AttachAlgoOrders[0].TpTriggerPrice)
	}

	// Close connection after all checks
	conn.Close()
}

func TestHandler_SubscribeAccount(t *testing.T) {
	accountReceived := make(chan Account)
	done := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()

		// Read subscription request
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Fatal(err)
		}

		var req AccountRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			t.Fatal(err)
		}

		// Verify request
		if req.Op != "subscribe" {
			t.Errorf("expected op=subscribe, got %s", req.Op)
		}
		if len(req.Args) != 1 {
			t.Errorf("expected 1 arg, got %d", len(req.Args))
		}
		if req.Args[0].Channel != "account" {
			t.Errorf("expected channel=account, got %s", req.Args[0].Channel)
		}

		// Send subscription response
		resp := AccountResponse{
			Event: "subscribe",
			Arg: struct {
				Channel string `json:"channel"`
			}{
				Channel: "account",
			},
		}
		respBytes, _ := json.Marshal(resp)
		if err := conn.WriteMessage(websocket.TextMessage, respBytes); err != nil {
			t.Fatal(err)
		}

		// Wait a bit to ensure subscription response is processed
		time.Sleep(100 * time.Millisecond)

		// Send account data
		accountData := AccountData{
			Arg: struct {
				Channel string `json:"channel"`
			}{
				Channel: "account",
			},
			Data: Account{
				Ts:             "1597026383085",
				TotalEquity:    "41624.32",
				IsolatedEquity: "3624.32",
				Details: []AccountDetail{
					{
						Currency:              "USDT",
						Equity:                "1",
						Balance:               "1",
						Ts:                    "1617279471503",
						IsolatedEquity:        "0",
						EquityUsd:             "45078.3790756226851775",
						AvailableEquity:       "1",
						Available:             "0",
						Frozen:                "0",
						OrderFrozen:           "0",
						UnrealizedPnl:         "0",
						IsolatedUnrealizedPnl: "0",
						CoinUsdPrice:          "1",
						SpotAvailable:         "0",
						Liability:             "0",
						BorrowFrozen:          "0",
						MarginRatio:           "0",
					},
				},
			},
		}
		accountBytes, _ := json.Marshal(accountData)
		if err := conn.WriteMessage(websocket.TextMessage, accountBytes); err != nil {
			t.Fatal(err)
		}

		// Wait for account to be received before closing
		<-done
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	handler := NewHandler(conn)

	err = handler.SubscribeAccount(func(account Account) {
		accountReceived <- account
	})
	if err != nil {
		t.Fatal(err)
	}

	// Start message handling in a goroutine
	go func() {
		if err := handler.HandleMessages(); err != nil {
			// Ignore connection closed error
			if !strings.Contains(err.Error(), "use of closed network connection") &&
				!strings.Contains(err.Error(), "websocket: close 1006") {
				t.Errorf("HandleMessages error: %v", err)
			}
		}
	}()

	// Wait for account data with timeout
	var receivedAccount Account
	select {
	case receivedAccount = <-accountReceived:
		// Account received successfully
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for account data")
	}

	// Verify received account
	if receivedAccount.TotalEquity != "41624.32" {
		t.Errorf("expected totalEquity=41624.32, got %s", receivedAccount.TotalEquity)
	}
	if receivedAccount.IsolatedEquity != "3624.32" {
		t.Errorf("expected isolatedEquity=3624.32, got %s", receivedAccount.IsolatedEquity)
	}
	if len(receivedAccount.Details) != 1 {
		t.Errorf("expected 1 detail, got %d", len(receivedAccount.Details))
	}
	if receivedAccount.Details[0].Currency != "USDT" {
		t.Errorf("expected currency=USDT, got %s", receivedAccount.Details[0].Currency)
	}
	if receivedAccount.Details[0].EquityUsd != "45078.3790756226851775" {
		t.Errorf("expected equityUsd=45078.3790756226851775, got %s", receivedAccount.Details[0].EquityUsd)
	}

	// Signal that we're done
	close(done)
}
