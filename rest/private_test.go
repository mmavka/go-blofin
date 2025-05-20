package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mmavka/go-blofin"
	"github.com/stretchr/testify/assert"
)

func newTestClient(t *testing.T) *RestClient {
	client := NewRestClient("http://localhost:8080")
	client.SetAuth("test-key", "test-secret", "test-passphrase")
	return client
}

func newTestPrivateRestClient(handler http.HandlerFunc) *RestClient {
	server := httptest.NewServer(handler)
	client := NewRestClient(server.URL)
	client.SetAuth("test-key", "test-secret", "test-passphrase")
	return client
}

func TestGetBalanceService(t *testing.T) {
	mockResp := `{"code":"0","data":[{"currency":"USDT","balance":"1000"}]}`
	client := newTestPrivateRestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResp))
	})
	resp, err := client.NewGetBalanceService().Do(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].Currency != "USDT" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestGetAccountBalanceService(t *testing.T) {
	client := newTestClient(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/account/balance", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		// Check authentication headers
		assert.NotEmpty(t, r.Header.Get("ACCESS-KEY"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-SIGN"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-TIMESTAMP"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-NONCE"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-PASSPHRASE"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"code": "0",
			"msg": "success",
			"data": {
				"ts": "1697021343571",
				"totalEquity": "10011254.077985990315787910",
				"isolatedEquity": "861.763132108800000000",
				"details": [
					{
						"currency": "USDT",
						"equity": "10014042.988958415234430699548",
						"balance": "10013119.885958415234430699",
						"ts": "1697021343571",
						"isolatedEquity": "862.003200000000000000048",
						"available": "9996399.4708691159703362725",
						"availableEquity": "9996399.4708691159703362725",
						"frozen": "15805.149672632597427761",
						"orderFrozen": "14920.994472632597427761",
						"equityUsd": "10011254.077985990315787910",
						"isolatedUnrealizedPnl": "-22.151999999999999999952",
						"bonus": "0"
					}
				]
			}
		}`))
	}))
	defer server.Close()

	client.SetBaseURL(server.URL)

	// Execute request
	balance, err := client.GetAccountBalance().Do(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, balance)

	// Check response data
	assert.Equal(t, "0", balance.Code)
	assert.Equal(t, "success", balance.Msg)
	assert.Equal(t, "1697021343571", balance.Data.Ts)
	assert.Equal(t, "10011254.077985990315787910", balance.Data.TotalEquity)
	assert.Equal(t, "861.763132108800000000", balance.Data.IsolatedEquity)

	// Check details
	assert.Len(t, balance.Data.Details, 1)
	detail := balance.Data.Details[0]
	assert.Equal(t, "USDT", detail.Currency)
	assert.Equal(t, "10014042.988958415234430699548", detail.Equity)
	assert.Equal(t, "10013119.885958415234430699", detail.Balance)
	assert.Equal(t, "862.003200000000000000048", detail.IsolatedEquity)
	assert.Equal(t, "9996399.4708691159703362725", detail.Available)
	assert.Equal(t, "15805.149672632597427761", detail.Frozen)
	assert.Equal(t, "14920.994472632597427761", detail.OrderFrozen)
	assert.Equal(t, "-22.151999999999999999952", detail.IsolatedUnrealizedPnl)
	assert.Equal(t, "0", detail.Bonus)
}

func TestGetPositionsService(t *testing.T) {
	client := newTestClient(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/account/positions", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		// Check authentication headers
		assert.NotEmpty(t, r.Header.Get("ACCESS-KEY"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-SIGN"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-TIMESTAMP"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-NONCE"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-PASSPHRASE"))

		// Check request parameters
		instId := r.URL.Query().Get("instId")
		if instId != "" {
			assert.Equal(t, "BTC-USDT", instId)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"code": "0",
			"msg": "success",
			"data": [
				{
					"positionId": "7982",
					"instId": "ETH-USDT",
					"instType": "SWAP",
					"marginMode": "isolated",
					"positionSide": "net",
					"adl": "5",
					"positions": "1",
					"availablePositions": "1",
					"averagePrice": "1591.800000000000000000",
					"margin": "53.060000000000000000",
					"markPrice": "1591.69",
					"marginRatio": "72.453752328329172684",
					"liquidationPrice": "1066.104078762306610407",
					"unrealizedPnl": "-0.011",
					"unrealizedPnlRatio": "-0.000207312476441764",
					"maintenanceMargin": "0.636676",
					"createTime": "1695352782370",
					"updateTime": "1695352782372",
					"leverage": "3"
				}
			]
		}`))
	}))
	defer server.Close()

	client.SetBaseURL(server.URL)

	// Test without parameters
	positions, err := client.GetPositions().Do(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, positions)

	// Check response data
	assert.Equal(t, "0", positions.Code)
	assert.Equal(t, "success", positions.Msg)
	assert.Len(t, positions.Data, 1)

	position := positions.Data[0]
	assert.Equal(t, "7982", position.PositionId)
	assert.Equal(t, "ETH-USDT", position.InstId)
	assert.Equal(t, "SWAP", position.InstType)
	assert.Equal(t, "isolated", position.MarginMode)
	assert.Equal(t, "1591.800000000000000000", position.AveragePrice)
	assert.Equal(t, "-0.011", position.UnrealizedPnl)
	assert.Equal(t, "3", position.Leverage)

	// Test with instId parameter
	positions, err = client.GetPositions().InstId("BTC-USDT").Do(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, positions)
}

func TestGetMarginModeService(t *testing.T) {
	client := newTestClient(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/account/margin-mode", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		// Check authentication headers
		assert.NotEmpty(t, r.Header.Get("ACCESS-KEY"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-SIGN"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-TIMESTAMP"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-NONCE"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-PASSPHRASE"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"code": "0",
			"msg": "success",
			"data": {
				"marginMode": "isolated"
			}
		}`))
	}))
	defer server.Close()

	client.SetBaseURL(server.URL)

	// Execute request
	marginMode, err := client.GetMarginMode().Do(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, marginMode)

	// Check response data
	assert.Equal(t, "0", marginMode.Code)
	assert.Equal(t, "success", marginMode.Msg)
	assert.Equal(t, blofin.MarginModeIsolated, marginMode.Data.MarginMode)
}

func TestSetMarginModeService(t *testing.T) {
	client := newTestClient(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/account/set-margin-mode", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		// Check authentication headers
		assert.NotEmpty(t, r.Header.Get("ACCESS-KEY"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-SIGN"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-TIMESTAMP"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-NONCE"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-PASSPHRASE"))

		// Check request body
		var req SetMarginModeRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, blofin.MarginModeIsolated, req.MarginMode)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"code": "0",
			"msg": "success",
			"data": {
				"marginMode": "isolated"
			}
		}`))
	}))
	defer server.Close()

	client.SetBaseURL(server.URL)

	// Test success
	marginMode, err := client.SetMarginMode().MarginMode(blofin.MarginModeIsolated).Do(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, marginMode)
	assert.Equal(t, "0", marginMode.Code)
	assert.Equal(t, "success", marginMode.Msg)
	assert.Equal(t, blofin.MarginModeIsolated, marginMode.Data.MarginMode)

	// Test with empty marginMode
	_, err = client.SetMarginMode().Do(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "marginMode required")
}

func TestGetPositionModeService(t *testing.T) {
	client := newTestClient(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/account/position-mode", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		// Check authentication headers
		assert.NotEmpty(t, r.Header.Get("ACCESS-KEY"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-SIGN"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-TIMESTAMP"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-NONCE"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-PASSPHRASE"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"code": "0",
			"msg": "success",
			"data": {
				"positionMode": "net_mode"
			}
		}`))
	}))
	defer server.Close()

	client.SetBaseURL(server.URL)

	// Execute request
	positionMode, err := client.GetPositionMode().Do(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, positionMode)

	// Check response data
	assert.Equal(t, "0", positionMode.Code)
	assert.Equal(t, "success", positionMode.Msg)
	assert.Equal(t, blofin.PositionModeNet, positionMode.Data.PositionMode)
}

func TestSetPositionModeService(t *testing.T) {
	client := newTestClient(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/account/set-position-mode", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		// Check authentication headers
		assert.NotEmpty(t, r.Header.Get("ACCESS-KEY"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-SIGN"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-TIMESTAMP"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-NONCE"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-PASSPHRASE"))

		// Check request body
		var req SetPositionModeRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, blofin.PositionModeNet, req.PositionMode)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"code": "0",
			"msg": "success",
			"data": {
				"positionMode": "net_mode"
			}
		}`))
	}))
	defer server.Close()

	client.SetBaseURL(server.URL)

	// Test successful request
	positionMode, err := client.SetPositionMode().PositionMode(blofin.PositionModeNet).Do(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, positionMode)
	assert.Equal(t, "0", positionMode.Code)
	assert.Equal(t, "success", positionMode.Msg)
	assert.Equal(t, blofin.PositionModeNet, positionMode.Data.PositionMode)

	// Test with empty positionMode
	_, err = client.SetPositionMode().Do(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "positionMode required")
}

func TestGetLeverageInfoService(t *testing.T) {
	client := newTestClient(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/account/leverage-info", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		// Check authentication headers
		assert.NotEmpty(t, r.Header.Get("ACCESS-KEY"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-SIGN"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-TIMESTAMP"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-NONCE"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-PASSPHRASE"))

		// Check request parameters
		assert.Equal(t, "BTC-USDT", r.URL.Query().Get("instId"))
		assert.Equal(t, blofin.MarginModeCross, r.URL.Query().Get("marginMode"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"code": "0",
			"msg": "success",
			"data": {
				"instId": "BTC-USDT",
				"leverage": "3",
				"marginMode": "cross"
			}
		}`))
	}))
	defer server.Close()

	client.SetBaseURL(server.URL)

	// Test successful request
	leverage, err := client.GetLeverageInfo().
		InstId("BTC-USDT").
		MarginMode(blofin.MarginModeCross).
		Do(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, leverage)
	assert.Equal(t, "0", leverage.Code)
	assert.Equal(t, "success", leverage.Msg)
	assert.Equal(t, "BTC-USDT", leverage.Data.InstId)
	assert.Equal(t, "3", leverage.Data.Leverage)
	assert.Equal(t, blofin.MarginModeCross, leverage.Data.MarginMode)

	// Test with empty instId
	_, err = client.GetLeverageInfo().MarginMode(blofin.MarginModeCross).Do(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "instId required")

	// Test with empty marginMode
	_, err = client.GetLeverageInfo().InstId("BTC-USDT").Do(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "marginMode required")
}

func TestGetBatchLeverageInfoService(t *testing.T) {
	client := newTestClient(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/account/batch-leverage-info", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		// Check authentication headers
		assert.NotEmpty(t, r.Header.Get("ACCESS-KEY"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-SIGN"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-TIMESTAMP"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-NONCE"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-PASSPHRASE"))

		// Check request parameters
		assert.Equal(t, "BTC-USDT,ETH-USDT", r.URL.Query().Get("instId"))
		assert.Equal(t, blofin.MarginModeCross, r.URL.Query().Get("marginMode"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"code": "0",
			"msg": "success",
			"data": [
				{
					"leverage": "50",
					"marginMode": "cross",
					"instId": "BTC-USDT",
					"positionSide": "net"
				},
				{
					"leverage": "3",
					"marginMode": "cross",
					"instId": "ETH-USDT",
					"positionSide": "net"
				}
			]
		}`))
	}))
	defer server.Close()

	client.SetBaseURL(server.URL)

	// Test successful request
	leverage, err := client.NewGetBatchLeverageInfoService().
		InstIds([]string{"BTC-USDT", "ETH-USDT"}).
		MarginMode(blofin.MarginModeCross).
		Do(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, leverage)
	assert.Equal(t, "0", leverage.Code)
	assert.Equal(t, "success", leverage.Msg)
	assert.Len(t, leverage.Data, 2)

	// Check first instrument
	assert.Equal(t, "BTC-USDT", leverage.Data[0].InstId)
	assert.Equal(t, "50", leverage.Data[0].Leverage)
	assert.Equal(t, blofin.MarginModeCross, leverage.Data[0].MarginMode)
	assert.Equal(t, blofin.PositionSideNet, leverage.Data[0].PositionSide)

	// Check second instrument
	assert.Equal(t, "ETH-USDT", leverage.Data[1].InstId)
	assert.Equal(t, "3", leverage.Data[1].Leverage)
	assert.Equal(t, blofin.MarginModeCross, leverage.Data[1].MarginMode)
	assert.Equal(t, blofin.PositionSideNet, leverage.Data[1].PositionSide)

	// Test with empty list of instruments
	_, err = client.NewGetBatchLeverageInfoService().
		MarginMode(blofin.MarginModeCross).
		Do(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "instIds required")

	// Test with exceeding instrument limit
	instIds := make([]string, 21)
	for i := range instIds {
		instIds[i] = fmt.Sprintf("BTC-USDT-%d", i)
	}
	_, err = client.NewGetBatchLeverageInfoService().
		InstIds(instIds).
		MarginMode(blofin.MarginModeCross).
		Do(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "too many instIds (max 20)")

	// Test with empty marginMode
	_, err = client.NewGetBatchLeverageInfoService().
		InstIds([]string{"BTC-USDT"}).
		Do(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "marginMode required")
}

func TestSetLeverageService(t *testing.T) {
	client := newTestClient(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/account/set-leverage", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		// Check authentication headers
		assert.NotEmpty(t, r.Header.Get("ACCESS-KEY"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-SIGN"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-TIMESTAMP"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-NONCE"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-PASSPHRASE"))

		// Check request body
		var req SetLeverageRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "BTC-USDT", req.InstId)
		assert.Equal(t, "100", req.Leverage)
		assert.Equal(t, blofin.MarginModeCross, req.MarginMode)
		assert.Equal(t, blofin.PositionSideLong, req.PositionSide)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"code": "0",
			"msg": "success",
			"data": {
				"instId": "BTC-USDT",
				"leverage": "100",
				"marginMode": "cross",
				"positionSide": "long"
			}
		}`))
	}))
	defer server.Close()

	client.SetBaseURL(server.URL)

	// Test successful request
	leverage, err := client.NewSetLeverageService().
		InstId("BTC-USDT").
		Leverage("100").
		MarginMode(blofin.MarginModeCross).
		PositionSide(blofin.PositionSideLong).
		Do(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, leverage)
	assert.Equal(t, "0", leverage.Code)
	assert.Equal(t, "success", leverage.Msg)
	assert.Equal(t, "BTC-USDT", leverage.Data.InstId)
	assert.Equal(t, "100", leverage.Data.Leverage)
	assert.Equal(t, blofin.MarginModeCross, leverage.Data.MarginMode)
	assert.Equal(t, blofin.PositionSideLong, leverage.Data.PositionSide)

	// Test with empty instId
	_, err = client.NewSetLeverageService().
		Leverage("100").
		MarginMode(blofin.MarginModeCross).
		Do(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "instId required")

	// Test with empty leverage
	_, err = client.NewSetLeverageService().
		InstId("BTC-USDT").
		MarginMode(blofin.MarginModeCross).
		Do(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "leverage required")

	// Test with empty marginMode
	_, err = client.NewSetLeverageService().
		InstId("BTC-USDT").
		Leverage("100").
		Do(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "marginMode required")
}

func TestPlaceOrderService(t *testing.T) {
	client := newTestClient(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/trade/order", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		// Check authentication headers
		assert.NotEmpty(t, r.Header.Get("ACCESS-KEY"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-SIGN"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-TIMESTAMP"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-NONCE"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-PASSPHRASE"))

		// Check request body
		var req PlaceOrderRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "BTC-USDT", req.InstID)
		assert.Equal(t, blofin.MarginModeCross, req.MarginMode)
		assert.Equal(t, blofin.PositionSideLong, req.PositionSide)
		assert.Equal(t, blofin.OrderSideSell, req.Side)
		assert.Equal(t, blofin.OrderTypeLimit, req.OrderType)
		assert.Equal(t, "23212.2", req.Price)
		assert.Equal(t, "2", req.Size)
		assert.Equal(t, "test1597321", req.ClientOrderId)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"code": "0",
			"msg": "",
			"data": [
				{
					"orderId": "28150801",
					"clientOrderId": "test1597321",
					"msg": "",
					"code": "0"
				}
			]
		}`))
	}))
	defer server.Close()

	client.SetBaseURL(server.URL)

	// Test successful request
	order, err := client.NewPlaceOrderService().
		InstId("BTC-USDT").
		MarginMode(blofin.MarginModeCross).
		PositionSide(blofin.PositionSideLong).
		Side(blofin.OrderSideSell).
		OrderType(blofin.OrderTypeLimit).
		Price("23212.2").
		Size("2").
		ClientOrderId("test1597321").
		Do(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, "0", order.Code)
	assert.Equal(t, "", order.Msg)
	assert.Len(t, order.Data, 1)
	assert.Equal(t, "28150801", order.Data[0].OrderId)
	assert.Equal(t, "test1597321", order.Data[0].ClientOrderId)
	assert.Equal(t, "0", order.Data[0].Code)
	assert.Equal(t, "", order.Data[0].Msg)

	// Test with empty instId
	_, err = client.NewPlaceOrderService().
		MarginMode(blofin.MarginModeCross).
		PositionSide(blofin.PositionSideLong).
		Side(blofin.OrderSideSell).
		OrderType(blofin.OrderTypeLimit).
		Price("23212.2").
		Size("2").
		Do(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "instId required")

	// Test with empty marginMode
	_, err = client.NewPlaceOrderService().
		InstId("BTC-USDT").
		PositionSide(blofin.PositionSideLong).
		Side(blofin.OrderSideSell).
		OrderType(blofin.OrderTypeLimit).
		Price("23212.2").
		Size("2").
		Do(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "marginMode required")

	// Test with empty positionSide
	_, err = client.NewPlaceOrderService().
		InstId("BTC-USDT").
		MarginMode(blofin.MarginModeCross).
		Side(blofin.OrderSideSell).
		OrderType(blofin.OrderTypeLimit).
		Price("23212.2").
		Size("2").
		Do(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "positionSide required")

	// Test with empty side
	_, err = client.NewPlaceOrderService().
		InstId("BTC-USDT").
		MarginMode(blofin.MarginModeCross).
		PositionSide(blofin.PositionSideLong).
		OrderType(blofin.OrderTypeLimit).
		Price("23212.2").
		Size("2").
		Do(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "side required")

	// Test with empty orderType
	_, err = client.NewPlaceOrderService().
		InstId("BTC-USDT").
		MarginMode(blofin.MarginModeCross).
		PositionSide(blofin.PositionSideLong).
		Side(blofin.OrderSideSell).
		Price("23212.2").
		Size("2").
		Do(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "orderType required")

	// Test with empty size
	_, err = client.NewPlaceOrderService().
		InstId("BTC-USDT").
		MarginMode(blofin.MarginModeCross).
		PositionSide(blofin.PositionSideLong).
		Side(blofin.OrderSideSell).
		OrderType(blofin.OrderTypeLimit).
		Price("23212.2").
		Do(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "size required")

	// Test with empty price for limit order
	_, err = client.NewPlaceOrderService().
		InstId("BTC-USDT").
		MarginMode(blofin.MarginModeCross).
		PositionSide(blofin.PositionSideLong).
		Side(blofin.OrderSideSell).
		OrderType(blofin.OrderTypeLimit).
		Size("2").
		Do(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "price required for non-market orders")
}

func TestPlaceBatchOrdersService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1/trade/batch-orders", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var req []PlaceOrderRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Len(t, req, 2)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":"0","msg":"","data":[{"orderId":"123","clientOrderId":"456"}]}`))
	}))
	defer server.Close()

	client := NewRestClient(server.URL)
	client.SetAuth("test-key", "test-secret", "test-passphrase")

	// Create test orders
	orders := []PlaceOrderRequest{
		{
			InstID:       "BTC-USDT",
			MarginMode:   "cross",
			PositionSide: "long",
			Side:         "buy",
			OrderType:    "limit",
			Price:        "50000",
			Size:         "0.1",
		},
		{
			InstID:       "ETH-USDT",
			MarginMode:   "cross",
			PositionSide: "long",
			Side:         "buy",
			OrderType:    "limit",
			Price:        "3000",
			Size:         "1",
		},
	}

	resp, err := client.NewPlaceBatchOrdersService().Orders(orders).Do(nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "0", resp.Code)
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, "123", resp.Data[0].OrderId)
	assert.Equal(t, "456", resp.Data[0].ClientOrderId)
}

func TestPlaceTPSLOrderService(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		client := newTestClient(t)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/trade/order-tpsl", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			// Check authentication headers
			assert.NotEmpty(t, r.Header.Get("ACCESS-KEY"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-SIGN"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-TIMESTAMP"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-NONCE"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-PASSPHRASE"))

			// Check request body
			var req PlaceTPSLOrderRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)
			assert.Equal(t, "ETH-USDT", req.InstID)
			assert.Equal(t, blofin.MarginModeCross, req.MarginMode)
			assert.Equal(t, blofin.PositionSideShort, req.PositionSide)
			assert.Equal(t, blofin.OrderSideSell, req.Side)
			assert.Equal(t, "1661.1", req.TpTriggerPrice)
			assert.Equal(t, "-1", req.TpOrderPrice)
			assert.Equal(t, "2", req.Size)
			assert.Equal(t, "true", req.ReduceOnly)
			assert.Equal(t, "test123", req.ClientOrderId)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"code": "0",
				"msg": "success",
				"data": {
					"tpslId": "1012",
					"clientOrderId": "test123",
					"code": "0",
					"msg": "Order placed successfully"
				}
			}`))
		}))
		defer server.Close()

		client.SetBaseURL(server.URL)

		// Test successful request
		order, err := client.NewPlaceTPSLOrderService().
			InstId("ETH-USDT").
			MarginMode(blofin.MarginModeCross).
			PositionSide(blofin.PositionSideShort).
			Side(blofin.OrderSideSell).
			TakeProfitParams("1661.1", "-1").
			Size("2").
			ReduceOnly("true").
			ClientOrderId("test123").
			Do(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, "0", order.Code)
		assert.Equal(t, "success", order.Msg)
		assert.Equal(t, "1012", order.Data.TpslId)
		assert.Equal(t, "test123", order.Data.ClientOrderId)
		assert.Equal(t, "0", order.Data.Code)
		assert.Equal(t, "Order placed successfully", order.Data.Msg)
	})

	t.Run("error response", func(t *testing.T) {
		client := newTestClient(t)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{
				"code": "50001",
				"msg": "Parameter error",
				"data": {
					"tpslId": "",
					"clientOrderId": "test123",
					"code": "50001",
					"msg": "Invalid trigger price"
				}
			}`))
		}))
		defer server.Close()

		client.SetBaseURL(server.URL)

		// Test with error
		order, err := client.NewPlaceTPSLOrderService().
			InstId("ETH-USDT").
			MarginMode(blofin.MarginModeCross).
			PositionSide(blofin.PositionSideShort).
			Side(blofin.OrderSideSell).
			TakeProfitParams("invalid", "-1").
			Size("2").
			ClientOrderId("test123").
			Do(context.Background())
		assert.NoError(t, err) // API returns error in response body
		assert.NotNil(t, order)
		assert.Equal(t, "50001", order.Code)
		assert.Equal(t, "Parameter error", order.Msg)
		assert.Empty(t, order.Data.TpslId)
		assert.Equal(t, "test123", order.Data.ClientOrderId)
		assert.Equal(t, "50001", order.Data.Code)
		assert.Equal(t, "Invalid trigger price", order.Data.Msg)
	})

	t.Run("validation errors", func(t *testing.T) {
		client := &RestClient{}

		// Test with empty instId
		_, err := client.NewPlaceTPSLOrderService().
			MarginMode(blofin.MarginModeCross).
			PositionSide(blofin.PositionSideShort).
			Side(blofin.OrderSideSell).
			TakeProfitParams("1661.1", "-1").
			Size("2").
			Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "instId required")

		// Test with empty marginMode
		_, err = client.NewPlaceTPSLOrderService().
			InstId("ETH-USDT").
			PositionSide(blofin.PositionSideShort).
			Side(blofin.OrderSideSell).
			TakeProfitParams("1661.1", "-1").
			Size("2").
			Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "marginMode required")

		// Test with empty positionSide
		_, err = client.NewPlaceTPSLOrderService().
			InstId("ETH-USDT").
			MarginMode(blofin.MarginModeCross).
			Side(blofin.OrderSideSell).
			TakeProfitParams("1661.1", "-1").
			Size("2").
			Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "positionSide required")

		// Test with empty side
		_, err = client.NewPlaceTPSLOrderService().
			InstId("ETH-USDT").
			MarginMode(blofin.MarginModeCross).
			PositionSide(blofin.PositionSideShort).
			Side(blofin.OrderSideSell).
			TakeProfitParams("1661.1", "-1").
			Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "side required")

		// Test with empty size
		_, err = client.NewPlaceTPSLOrderService().
			InstId("ETH-USDT").
			MarginMode(blofin.MarginModeCross).
			PositionSide(blofin.PositionSideShort).
			Side(blofin.OrderSideSell).
			TakeProfitParams("1661.1", "-1").
			Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "size required")

		// Test without TP and SL
		_, err = client.NewPlaceTPSLOrderService().
			InstId("ETH-USDT").
			MarginMode(blofin.MarginModeCross).
			PositionSide(blofin.PositionSideShort).
			Side(blofin.OrderSideSell).
			Size("2").
			Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least one of tpTriggerPrice or slTriggerPrice required")

		// Test with TP trigger without TP order price
		_, err = client.NewPlaceTPSLOrderService().
			InstId("ETH-USDT").
			MarginMode(blofin.MarginModeCross).
			PositionSide(blofin.PositionSideShort).
			Side(blofin.OrderSideSell).
			Size("2").
			TakeProfitParams("1661.1", "").
			Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tpOrderPrice required when tpTriggerPrice is set")
	})
}

func TestPlaceAlgoOrderService(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		client := newTestClient(t)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/trade/order-algo", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			// Check authentication headers
			assert.NotEmpty(t, r.Header.Get("ACCESS-KEY"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-SIGN"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-TIMESTAMP"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-NONCE"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-PASSPHRASE"))

			// Check request body
			var req PlaceAlgoOrderRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)
			assert.Equal(t, "ETH-USDT", req.InstID)
			assert.Equal(t, blofin.MarginModeCross, req.MarginMode)
			assert.Equal(t, blofin.PositionSideShort, req.PositionSide)
			assert.Equal(t, blofin.OrderSideSell, req.Side)
			assert.Equal(t, "1", req.Size)
			assert.Equal(t, "-1", req.OrderPrice)
			assert.Equal(t, "trigger", req.OrderType)
			assert.Equal(t, "3000", req.TriggerPrice)
			assert.Equal(t, "last", req.TriggerPriceType)
			assert.Equal(t, "test123", req.ClientOrderId)

			// Check attached TP/SL orders
			assert.Len(t, req.AttachAlgoOrders, 1)
			assert.Equal(t, "3500", req.AttachAlgoOrders[0].TpTriggerPrice)
			assert.Equal(t, "3600", req.AttachAlgoOrders[0].TpOrderPrice)
			assert.Equal(t, "last", req.AttachAlgoOrders[0].TpTriggerPriceType)
			assert.Equal(t, "2600", req.AttachAlgoOrders[0].SlTriggerPrice)
			assert.Equal(t, "2500", req.AttachAlgoOrders[0].SlOrderPrice)
			assert.Equal(t, "last", req.AttachAlgoOrders[0].SlTriggerPriceType)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"code": "0",
				"msg": "success",
				"data": {
					"algoId": "1012",
					"clientOrderId": "test123",
					"code": "0",
					"msg": "Order placed successfully"
				}
			}`))
		}))
		defer server.Close()

		client.SetBaseURL(server.URL)

		// Создаем прикрепленные TP/SL ордера
		attachOrders := []AttachAlgoOrder{
			{
				TpTriggerPrice:     "3500",
				TpOrderPrice:       "3600",
				TpTriggerPriceType: "last",
				SlTriggerPrice:     "2600",
				SlOrderPrice:       "2500",
				SlTriggerPriceType: "last",
			},
		}

		// Test successful request
		order, err := client.NewPlaceAlgoOrderService().
			InstId("ETH-USDT").
			MarginMode(blofin.MarginModeCross).
			PositionSide(blofin.PositionSideShort).
			Side(blofin.OrderSideSell).
			Size("1").
			OrderPrice("-1").
			OrderType("trigger").
			TriggerPrice("3000").
			TriggerPriceType("last").
			ClientOrderId("test123").
			AttachAlgoOrders(attachOrders).
			Do(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, "0", order.Code)
		assert.Equal(t, "success", order.Msg)
		assert.Equal(t, "1012", order.Data.AlgoId)
		assert.Equal(t, "test123", order.Data.ClientOrderId)
		assert.Equal(t, "0", order.Data.Code)
		assert.Equal(t, "Order placed successfully", order.Data.Msg)
	})

	t.Run("validation errors", func(t *testing.T) {
		client := &RestClient{}

		// Test with empty instId
		_, err := client.NewPlaceAlgoOrderService().
			MarginMode(blofin.MarginModeCross).
			PositionSide(blofin.PositionSideShort).
			Side(blofin.OrderSideSell).
			Size("1").
			OrderType("trigger").
			TriggerPrice("3000").
			Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "instId required")

		// Test with empty marginMode
		_, err = client.NewPlaceAlgoOrderService().
			InstId("ETH-USDT").
			PositionSide(blofin.PositionSideShort).
			Side(blofin.OrderSideSell).
			Size("1").
			OrderType("trigger").
			TriggerPrice("3000").
			Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "marginMode required")

		// Test with empty positionSide
		_, err = client.NewPlaceAlgoOrderService().
			InstId("ETH-USDT").
			MarginMode(blofin.MarginModeCross).
			Side(blofin.OrderSideSell).
			Size("1").
			OrderType("trigger").
			TriggerPrice("3000").
			Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "positionSide required")

		// Test with empty side
		_, err = client.NewPlaceAlgoOrderService().
			InstId("ETH-USDT").
			MarginMode(blofin.MarginModeCross).
			PositionSide(blofin.PositionSideShort).
			Size("1").
			OrderType("trigger").
			TriggerPrice("3000").
			Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "side required")

		// Test with empty size
		_, err = client.NewPlaceAlgoOrderService().
			InstId("ETH-USDT").
			MarginMode(blofin.MarginModeCross).
			PositionSide(blofin.PositionSideShort).
			Side(blofin.OrderSideSell).
			OrderType("trigger").
			TriggerPrice("3000").
			Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "size required")

		// Test with empty orderType
		_, err = client.NewPlaceAlgoOrderService().
			InstId("ETH-USDT").
			MarginMode(blofin.MarginModeCross).
			PositionSide(blofin.PositionSideShort).
			Side(blofin.OrderSideSell).
			Size("1").
			TriggerPrice("3000").
			Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "orderType required")

		// Test with empty triggerPrice
		_, err = client.NewPlaceAlgoOrderService().
			InstId("ETH-USDT").
			MarginMode(blofin.MarginModeCross).
			PositionSide(blofin.PositionSideShort).
			Side(blofin.OrderSideSell).
			Size("1").
			OrderType("trigger").
			Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "triggerPrice required")
	})
}

func TestCancelOrderService(t *testing.T) {
	t.Run("successful request with orderId", func(t *testing.T) {
		client := newTestClient(t)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/trade/cancel-order", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			// Check authentication headers
			assert.NotEmpty(t, r.Header.Get("ACCESS-KEY"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-SIGN"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-TIMESTAMP"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-NONCE"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-PASSPHRASE"))

			// Check request body
			var req CancelOrderRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)
			assert.Equal(t, "23209016", req.OrderId)
			assert.Equal(t, "BTC-USDT", req.InstID)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"code": "0",
				"msg": "success",
				"data": {
					"orderId": "23209016",
					"clientOrderId": null,
					"code": "0",
					"msg": null
				}
			}`))
		}))
		defer server.Close()

		client.SetBaseURL(server.URL)

		// Test successful request with orderId
		order, err := client.NewCancelOrderService().
			InstId("BTC-USDT").
			OrderId("23209016").
			Do(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, "0", order.Code)
		assert.Equal(t, "success", order.Msg)
		assert.Equal(t, "23209016", order.Data.OrderId)
		assert.Equal(t, "0", order.Data.Code)
	})

	t.Run("successful request with clientOrderId", func(t *testing.T) {
		client := newTestClient(t)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req CancelOrderRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)
			assert.Equal(t, "test123", req.ClientOrderId)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"code": "0",
				"msg": "success",
				"data": {
					"orderId": "23209016",
					"clientOrderId": "test123",
					"code": "0",
					"msg": null
				}
			}`))
		}))
		defer server.Close()

		client.SetBaseURL(server.URL)

		// Test successful request with clientOrderId
		order, err := client.NewCancelOrderService().
			ClientOrderId("test123").
			Do(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, "test123", order.Data.ClientOrderId)
	})

	t.Run("validation errors", func(t *testing.T) {
		client := &RestClient{}

		// Test without orderId and clientOrderId
		_, err := client.NewCancelOrderService().Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "either orderId or clientOrderId required")
	})
}

func TestCancelBatchOrdersService(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		client := newTestClient(t)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/trade/cancel-batch-orders", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			// Check authentication headers
			assert.NotEmpty(t, r.Header.Get("ACCESS-KEY"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-SIGN"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-TIMESTAMP"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-NONCE"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-PASSPHRASE"))

			// Check request body
			var orders []CancelOrderRequest
			err := json.NewDecoder(r.Body).Decode(&orders)
			assert.NoError(t, err)
			assert.Len(t, orders, 2)
			assert.Equal(t, "ETH-USDT", orders[0].InstID)
			assert.Equal(t, "22619976", orders[0].OrderId)
			assert.Equal(t, "ETH-USDT", orders[1].InstID)
			assert.Equal(t, "22619977", orders[1].OrderId)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"code": "0",
				"msg": "success",
				"data": [
					{
						"orderId": "22619976",
						"clientOrderId": "eeeeee112231121",
						"code": "0",
						"msg": null
					},
					{
						"orderId": "22619977",
						"clientOrderId": "eeeeee11223211",
						"code": "0",
						"msg": null
					},
					{
						"orderId": "22619977111",
						"clientOrderId": null,
						"msg": "Cancel failed as the order has been filled, triggered, canceled or does not exist.",
						"code": "1000"
					}
				]
			}`))
		}))
		defer server.Close()

		client.SetBaseURL(server.URL)

		// Create list of orders to cancel
		orders := []CancelOrderRequest{
			{
				InstID:  "ETH-USDT",
				OrderId: "22619976",
			},
			{
				InstID:  "ETH-USDT",
				OrderId: "22619977",
			},
		}

		// Test successful request
		response, err := client.NewCancelBatchOrdersService().
			Orders(orders).
			Do(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "0", response.Code)
		assert.Equal(t, "success", response.Msg)
		assert.Len(t, response.Data, 3)

		// Check results
		assert.Equal(t, "22619976", response.Data[0].OrderId)
		assert.Equal(t, "eeeeee112231121", response.Data[0].ClientOrderId)
		assert.Equal(t, "0", response.Data[0].Code)

		assert.Equal(t, "22619977", response.Data[1].OrderId)
		assert.Equal(t, "eeeeee11223211", response.Data[1].ClientOrderId)
		assert.Equal(t, "0", response.Data[1].Code)

		assert.Equal(t, "22619977111", response.Data[2].OrderId)
		assert.Equal(t, "1000", response.Data[2].Code)
		assert.Equal(t, "Cancel failed as the order has been filled, triggered, canceled or does not exist.", response.Data[2].Msg)
	})

	t.Run("validation errors", func(t *testing.T) {
		client := &RestClient{}

		// Test with empty list of orders
		_, err := client.NewCancelBatchOrdersService().Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least one order required")

		// Test with exceeding order limit
		orders := make([]CancelOrderRequest, 21)
		for i := range orders {
			orders[i] = CancelOrderRequest{InstID: "ETH-USDT", OrderId: fmt.Sprintf("%d", i)}
		}
		_, err = client.NewCancelBatchOrdersService().Orders(orders).Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "too many orders (max 20)")

		// Test with different instId
		orders = []CancelOrderRequest{
			{InstID: "ETH-USDT", OrderId: "1"},
			{InstID: "BTC-USDT", OrderId: "2"},
		}
		_, err = client.NewCancelBatchOrdersService().Orders(orders).Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "all orders must have the same instId")

		// Test without orderId and clientOrderId
		orders = []CancelOrderRequest{
			{InstID: "ETH-USDT"},
		}
		_, err = client.NewCancelBatchOrdersService().Orders(orders).Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "either orderId or clientOrderId required")
	})

	t.Run("add order method", func(t *testing.T) {
		service := &CancelBatchOrdersService{}
		order := CancelOrderRequest{
			InstID:  "ETH-USDT",
			OrderId: "22619976",
		}

		service.AddOrder(order)
		assert.Len(t, service.orders, 1)
		assert.Equal(t, order, service.orders[0])
	})
}

func TestCancelTPSLOrderService(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		client := newTestClient(t)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/trade/cancel-tpsl", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			// Check authentication headers
			assert.NotEmpty(t, r.Header.Get("ACCESS-KEY"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-SIGN"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-TIMESTAMP"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-NONCE"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-PASSPHRASE"))

			// Check request body
			var orders []CancelTPSLOrderRequest
			err := json.NewDecoder(r.Body).Decode(&orders)
			assert.NoError(t, err)
			assert.Len(t, orders, 2)
			assert.Equal(t, "ETH-USDT", orders[0].InstID)
			assert.Equal(t, "22619976", orders[0].TpslId)
			assert.Equal(t, "ETH-USDT", orders[1].InstID)
			assert.Equal(t, "22619977", orders[1].TpslId)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"code": "0",
				"msg": "success",
				"data": [
					{
						"tpslId": "1009",
						"clientOrderId": null,
						"code": "500",
						"msg": "Cancel failed as the order has been filled, triggered, canceled or does not exist."
					},
					{
						"tpslId": "1010",
						"clientOrderId": null,
						"code": "500",
						"msg": "Cancel failed as the order has been filled, triggered, canceled or does not exist."
					}
				]
			}`))
		}))
		defer server.Close()

		client.SetBaseURL(server.URL)

		// Создаем список ордеров для отмены
		orders := []CancelTPSLOrderRequest{
			{
				InstID: "ETH-USDT",
				TpslId: "22619976",
			},
			{
				InstID: "ETH-USDT",
				TpslId: "22619977",
			},
		}

		// Test successful request
		response, err := client.NewCancelTPSLOrderService().
			Orders(orders).
			Do(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "0", response.Code)
		assert.Equal(t, "success", response.Msg)
		assert.Len(t, response.Data, 2)

		// Check results
		assert.Equal(t, "1009", response.Data[0].TpslId)
		assert.Equal(t, "500", response.Data[0].Code)
		assert.Equal(t, "Cancel failed as the order has been filled, triggered, canceled or does not exist.", response.Data[0].Msg)

		assert.Equal(t, "1010", response.Data[1].TpslId)
		assert.Equal(t, "500", response.Data[1].Code)
		assert.Equal(t, "Cancel failed as the order has been filled, triggered, canceled or does not exist.", response.Data[1].Msg)
	})

	t.Run("validation errors", func(t *testing.T) {
		client := &RestClient{}

		// Test with empty list of orders
		_, err := client.NewCancelTPSLOrderService().Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least one order required")

		// Test with exceeding order limit
		orders := make([]CancelTPSLOrderRequest, 21)
		for i := range orders {
			orders[i] = CancelTPSLOrderRequest{InstID: "ETH-USDT", TpslId: fmt.Sprintf("%d", i)}
		}
		_, err = client.NewCancelTPSLOrderService().Orders(orders).Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "too many orders (max 20)")

		// Test with different instId
		orders = []CancelTPSLOrderRequest{
			{InstID: "ETH-USDT", TpslId: "1"},
			{InstID: "BTC-USDT", TpslId: "2"},
		}
		_, err = client.NewCancelTPSLOrderService().Orders(orders).Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "all orders must have the same instId")

		// Test without tpslId and clientOrderId
		orders = []CancelTPSLOrderRequest{
			{InstID: "ETH-USDT"},
		}
		_, err = client.NewCancelTPSLOrderService().Orders(orders).Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "either tpslId or clientOrderId required")
	})

	t.Run("add order method", func(t *testing.T) {
		service := &CancelTPSLOrderService{}
		order := CancelTPSLOrderRequest{
			InstID: "ETH-USDT",
			TpslId: "22619976",
		}

		service.AddOrder(order)
		assert.Len(t, service.orders, 1)
		assert.Equal(t, order, service.orders[0])
	})
}

func TestCancelAlgoOrderService(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		client := newTestPrivateRestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/api/v1/trade/cancel-algo", r.URL.Path)
			assert.NotEmpty(t, r.Header.Get("ACCESS-KEY"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-SIGN"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-TIMESTAMP"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-NONCE"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-PASSPHRASE"))

			var req CancelAlgoOrderRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)
			assert.Equal(t, "ETH-USDT", req.InstID)
			assert.Equal(t, "22619976", req.AlgoId)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"code": "0",
				"msg": "success",
				"data": {
					"algoId": "22619976",
					"clientOrderId": null,
					"code": "0",
					"msg": null
				}
			}`))
		}))

		resp, err := client.NewCancelAlgoOrderService().
			InstId("ETH-USDT").
			AlgoId("22619976").
			Do(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "0", resp.Code)
		assert.Equal(t, "success", resp.Msg)
		assert.Equal(t, "22619976", resp.Data.AlgoId)
		assert.Equal(t, "0", resp.Data.Code)
	})

	t.Run("validation errors", func(t *testing.T) {
		client := newTestClient(t)

		// Check that either algoId or clientOrderId is required
		_, err := client.NewCancelAlgoOrderService().Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "either algoId or clientOrderId required")
	})

	t.Run("error response", func(t *testing.T) {
		client := newTestPrivateRestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"code": "0",
				"msg": "success",
				"data": {
					"algoId": "1009",
					"clientOrderId": null,
					"code": "500",
					"msg": "Cancel failed as the order has been filled, triggered, canceled or does not exist."
				}
			}`))
		}))

		resp, err := client.NewCancelAlgoOrderService().
			InstId("ETH-USDT").
			AlgoId("1009").
			Do(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "0", resp.Code)
		assert.Equal(t, "success", resp.Msg)
		assert.Equal(t, "1009", resp.Data.AlgoId)
		assert.Equal(t, "500", resp.Data.Code)
		assert.Equal(t, "Cancel failed as the order has been filled, triggered, canceled or does not exist.", resp.Data.Msg)
	})
}

func TestGetPendingOrdersService(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		client := newTestPrivateRestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/api/v1/trade/orders-pending", r.URL.Path)
			assert.NotEmpty(t, r.Header.Get("ACCESS-KEY"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-SIGN"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-TIMESTAMP"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-NONCE"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-PASSPHRASE"))

			// Check request parameters
			assert.Equal(t, "ETH-USDT", r.URL.Query().Get("instId"))
			assert.Equal(t, "limit", r.URL.Query().Get("orderType"))
			assert.Equal(t, "live", r.URL.Query().Get("state"))
			assert.Equal(t, "20", r.URL.Query().Get("limit"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"code": "0",
				"msg": "success",
				"data": [
					{
						"orderId": "29531103",
						"clientOrderId": "",
						"instId": "ETH-USDT",
						"marginMode": "isolated",
						"positionSide": "net",
						"side": "buy",
						"orderType": "limit",
						"price": "1514.150000000000000000",
						"size": "1.000000000000000000",
						"reduceOnly": "false",
						"leverage": "3",
						"state": "live",
						"filledSize": "0.000000000000000000",
						"filled_amount": "0.000000000000000000",
						"averagePrice": "0.000000000000000000",
						"fee": "0.000000000000000000",
						"pnl": "0.000000000000000000",
						"createTime": "1697031292762",
						"updateTime": "1697031292788",
						"orderCategory": "normal",
						"tpTriggerPrice": "1688.000000000000000000",
						"slTriggerPrice": "1299.000000000000000000",
						"slOrderPrice": null,
						"tpOrderPrice": null,
						"algoClientOrderId": "aaa",
						"algoId": "11756185",
						"brokerId": ""
					}
				]
			}`))
		}))

		resp, err := client.NewGetPendingOrdersService().
			InstId("ETH-USDT").
			OrderType("limit").
			State("live").
			Limit("20").
			Do(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "0", resp.Code)
		assert.Equal(t, "success", resp.Msg)
		assert.Len(t, resp.Data, 1)

		order := resp.Data[0]
		assert.Equal(t, "29531103", order.OrderId)
		assert.Equal(t, "ETH-USDT", order.InstId)
		assert.Equal(t, "isolated", order.MarginMode)
		assert.Equal(t, "net", order.PositionSide)
		assert.Equal(t, "buy", order.Side)
		assert.Equal(t, "limit", order.OrderType)
		assert.Equal(t, "1514.150000000000000000", order.Price)
		assert.Equal(t, "1.000000000000000000", order.Size)
		assert.Equal(t, "false", order.ReduceOnly)
		assert.Equal(t, "3", order.Leverage)
		assert.Equal(t, "live", order.State)
		assert.Equal(t, "1688.000000000000000000", order.TpTriggerPrice)
		assert.Equal(t, "1299.000000000000000000", order.SlTriggerPrice)
		assert.Equal(t, "aaa", order.AlgoClientOrderId)
		assert.Equal(t, "11756185", order.AlgoId)
	})

	t.Run("pagination", func(t *testing.T) {
		client := newTestPrivateRestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/api/v1/trade/orders-pending", r.URL.Path)

			// Check pagination parameters
			assert.Equal(t, "29531103", r.URL.Query().Get("after"))
			assert.Equal(t, "50", r.URL.Query().Get("limit"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"code": "0",
				"msg": "success",
				"data": []
			}`))
		}))

		resp, err := client.NewGetPendingOrdersService().
			After("29531103").
			Limit("50").
			Do(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "0", resp.Code)
		assert.Equal(t, "success", resp.Msg)
		assert.Empty(t, resp.Data)
	})
}

func TestGetPendingTPSLOrdersService(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		client := newTestPrivateRestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/api/v1/trade/orders-tpsl-pending", r.URL.Path)
			assert.NotEmpty(t, r.Header.Get("ACCESS-KEY"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-SIGN"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-TIMESTAMP"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-NONCE"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-PASSPHRASE"))

			// Check request parameters
			assert.Equal(t, "ETH-USDT", r.URL.Query().Get("instId"))
			assert.Equal(t, "20", r.URL.Query().Get("limit"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"code": "0",
				"msg": "success",
				"data": [
					{
						"tpslId": "2411",
						"instId": "ETH-USDT",
						"marginMode": "cross",
						"positionSide": "net",
						"side": "sell",
						"tpTriggerPrice": "1666.000000000000000000",
						"tpOrderPrice": null,
						"slTriggerPrice": "1222.000000000000000000",
						"slOrderPrice": null,
						"size": "1",
						"state": "live",
						"leverage": "3",
						"reduceOnly": "false",
						"actualSize": null,
						"clientOrderId": "aabbc",
						"createTime": "1697016700775",
						"brokerId": ""
					}
				]
			}`))
		}))

		resp, err := client.NewGetPendingTPSLOrdersService().
			InstId("ETH-USDT").
			Limit("20").
			Do(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "0", resp.Code)
		assert.Equal(t, "success", resp.Msg)
		assert.Len(t, resp.Data, 1)

		order := resp.Data[0]
		assert.Equal(t, "2411", order.TpslId)
		assert.Equal(t, "ETH-USDT", order.InstId)
		assert.Equal(t, "cross", order.MarginMode)
		assert.Equal(t, "net", order.PositionSide)
		assert.Equal(t, "sell", order.Side)
		assert.Equal(t, "1666.000000000000000000", order.TpTriggerPrice)
		assert.Equal(t, "1222.000000000000000000", order.SlTriggerPrice)
		assert.Equal(t, "1", order.Size)
		assert.Equal(t, "live", order.State)
		assert.Equal(t, "3", order.Leverage)
		assert.Equal(t, "false", order.ReduceOnly)
		assert.Equal(t, "aabbc", order.ClientOrderId)
		assert.Equal(t, "1697016700775", order.CreateTime)
	})

	t.Run("pagination", func(t *testing.T) {
		client := newTestPrivateRestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/api/v1/trade/orders-tpsl-pending", r.URL.Path)

			// Check pagination parameters
			assert.Equal(t, "2411", r.URL.Query().Get("after"))
			assert.Equal(t, "50", r.URL.Query().Get("limit"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"code": "0",
				"msg": "success",
				"data": []
			}`))
		}))

		resp, err := client.NewGetPendingTPSLOrdersService().
			After("2411").
			Limit("50").
			Do(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "0", resp.Code)
		assert.Equal(t, "success", resp.Msg)
		assert.Empty(t, resp.Data)
	})
}

func TestGetPendingAlgoOrdersService(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		client := newTestPrivateRestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/api/v1/trade/orders-algo-pending", r.URL.Path)
			assert.NotEmpty(t, r.Header.Get("ACCESS-KEY"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-SIGN"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-TIMESTAMP"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-NONCE"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-PASSPHRASE"))

			// Check request parameters
			assert.Equal(t, "ETH-USDT", r.URL.Query().Get("instId"))
			assert.Equal(t, "trigger", r.URL.Query().Get("orderType"))
			assert.Equal(t, "20", r.URL.Query().Get("limit"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"code": "0",
				"msg": "success",
				"data": [
					{
						"algoId": "2101",
						"clientOrderId": "BBBBqqqq",
						"instId": "ETH-USDT",
						"marginMode": "cross",
						"positionSide": "net",
						"side": "sell",
						"orderType": "trigger",
						"size": "1",
						"leverage": "3",
						"state": "canceled",
						"triggerPrice": "1661.100000000000000000",
						"triggerPriceType": "last",
						"brokerId": "",
						"attachAlgoOrders": [
							{
								"tpTriggerPrice": "1666.000000000000000000",
								"tpOrderPrice": "-1",
								"tpTriggerPriceType": "last",
								"slTriggerPrice": "1222.000000000000000000",
								"slOrderPrice": "-1",
								"slTriggerPriceType": "last"
							}
						]
					}
				]
			}`))
		}))

		resp, err := client.NewGetPendingAlgoOrdersService().
			InstId("ETH-USDT").
			OrderType("trigger").
			Limit("20").
			Do(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "0", resp.Code)
		assert.Equal(t, "success", resp.Msg)
		assert.Len(t, resp.Data, 1)

		order := resp.Data[0]
		assert.Equal(t, "2101", order.AlgoId)
		assert.Equal(t, "BBBBqqqq", order.ClientOrderId)
		assert.Equal(t, "ETH-USDT", order.InstId)
		assert.Equal(t, "cross", order.MarginMode)
		assert.Equal(t, "net", order.PositionSide)
		assert.Equal(t, "sell", order.Side)
		assert.Equal(t, "trigger", order.OrderType)
		assert.Equal(t, "1", order.Size)
		assert.Equal(t, "3", order.Leverage)
		assert.Equal(t, "canceled", order.State)
		assert.Equal(t, "1661.100000000000000000", order.TriggerPrice)
		assert.Equal(t, "last", order.TriggerPriceType)
		assert.Empty(t, order.BrokerId)
		assert.Len(t, order.AttachAlgoOrders, 1)

		attachOrder := order.AttachAlgoOrders[0]
		assert.Equal(t, "1666.000000000000000000", attachOrder.TpTriggerPrice)
		assert.Equal(t, "-1", attachOrder.TpOrderPrice)
		assert.Equal(t, "last", attachOrder.TpTriggerPriceType)
		assert.Equal(t, "1222.000000000000000000", attachOrder.SlTriggerPrice)
		assert.Equal(t, "-1", attachOrder.SlOrderPrice)
		assert.Equal(t, "last", attachOrder.SlTriggerPriceType)
	})

	t.Run("validation errors", func(t *testing.T) {
		client := newTestPrivateRestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		// Check required orderType parameter
		_, err := client.NewGetPendingAlgoOrdersService().
			InstId("ETH-USDT").
			Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "orderType required")

		// Check simultaneous use of after and before
		_, err = client.NewGetPendingAlgoOrdersService().
			OrderType("trigger").
			After("2101").
			Before("2102").
			Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "after and before cannot be used simultaneously")
	})

	t.Run("pagination", func(t *testing.T) {
		client := newTestPrivateRestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/api/v1/trade/orders-algo-pending", r.URL.Path)

			// Check pagination parameters
			assert.Equal(t, "2101", r.URL.Query().Get("after"))
			assert.Equal(t, "50", r.URL.Query().Get("limit"))
			assert.Equal(t, "trigger", r.URL.Query().Get("orderType"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"code": "0",
				"msg": "success",
				"data": []
			}`))
		}))

		resp, err := client.NewGetPendingAlgoOrdersService().
			OrderType("trigger").
			After("2101").
			Limit("50").
			Do(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "0", resp.Code)
		assert.Equal(t, "success", resp.Msg)
		assert.Empty(t, resp.Data)
	})
}

func TestClosePositionService(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		client := newTestPrivateRestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/api/v1/trade/close-position", r.URL.Path)
			assert.NotEmpty(t, r.Header.Get("ACCESS-KEY"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-SIGN"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-TIMESTAMP"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-NONCE"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-PASSPHRASE"))

			// Check request body
			var req ClosePositionRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)
			assert.Equal(t, "BTC-USDT", req.InstID)
			assert.Equal(t, "cross", req.MarginMode)
			assert.Equal(t, "long", req.PositionSide)
			assert.Empty(t, req.ClientOrderId)
			assert.Empty(t, req.BrokerId)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"code": "0",
				"msg": "success",
				"data": {
					"instId": "BTC-USDT",
					"positionSide": "long",
					"clientOrderId": ""
				}
			}`))
		}))

		resp, err := client.NewClosePositionService().
			InstId("BTC-USDT").
			MarginMode("cross").
			PositionSide("long").
			Do(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "0", resp.Code)
		assert.Equal(t, "success", resp.Msg)
		assert.Equal(t, "BTC-USDT", resp.Data.InstID)
		assert.Equal(t, "long", resp.Data.PositionSide)
		assert.Empty(t, resp.Data.ClientOrderId)
	})

	t.Run("validation errors", func(t *testing.T) {
		client := newTestPrivateRestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		// Check required parameters
		_, err := client.NewClosePositionService().
			Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "instId required")

		_, err = client.NewClosePositionService().
			InstId("BTC-USDT").
			Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "marginMode required")

		_, err = client.NewClosePositionService().
			InstId("BTC-USDT").
			MarginMode("cross").
			Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "positionSide required")
	})

	t.Run("with client order id", func(t *testing.T) {
		client := newTestPrivateRestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req ClosePositionRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)
			assert.Equal(t, "BTC-USDT", req.InstID)
			assert.Equal(t, "cross", req.MarginMode)
			assert.Equal(t, "long", req.PositionSide)
			assert.Equal(t, "test-order-id", req.ClientOrderId)
			assert.Empty(t, req.BrokerId)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"code": "0",
				"msg": "success",
				"data": {
					"instId": "BTC-USDT",
					"positionSide": "long",
					"clientOrderId": "test-order-id"
				}
			}`))
		}))

		resp, err := client.NewClosePositionService().
			InstId("BTC-USDT").
			MarginMode("cross").
			PositionSide("long").
			ClientOrderId("test-order-id").
			Do(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "0", resp.Code)
		assert.Equal(t, "success", resp.Msg)
		assert.Equal(t, "BTC-USDT", resp.Data.InstID)
		assert.Equal(t, "long", resp.Data.PositionSide)
		assert.Equal(t, "test-order-id", resp.Data.ClientOrderId)
	})
}

func TestGetOrderHistoryService(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		client := newTestPrivateRestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/api/v1/trade/orders-history", r.URL.Path)
			assert.NotEmpty(t, r.Header.Get("ACCESS-KEY"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-SIGN"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-TIMESTAMP"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-NONCE"))
			assert.NotEmpty(t, r.Header.Get("ACCESS-PASSPHRASE"))

			// Check request parameters
			assert.Equal(t, "ETH-USDT", r.URL.Query().Get("instId"))
			assert.Equal(t, "limit", r.URL.Query().Get("orderType"))
			assert.Equal(t, "canceled", r.URL.Query().Get("state"))
			assert.Equal(t, "20", r.URL.Query().Get("limit"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"code": "0",
				"msg": "success",
				"data": [
					{
						"orderId": "29419717",
						"clientOrderId": "aabbc",
						"instId": "ETH-USDT",
						"marginMode": "cross",
						"positionSide": "net",
						"side": "buy",
						"orderType": "limit",
						"price": "1523.000000000000000000",
						"size": "1.000000000000000000",
						"reduceOnly": "false",
						"leverage": "3",
						"state": "canceled",
						"filledSize": "0.000000000000000000",
						"pnl": "0.000000000000000000",
						"averagePrice": "0.000000000000000000",
						"fee": "0.000000000000000000",
						"createTime": "1697010303781",
						"updateTime": "1697014607770",
						"orderCategory": "normal",
						"tpTriggerPrice": null,
						"tpOrderPrice": null,
						"slTriggerPrice": null,
						"slOrderPrice": null,
						"cancelSource": "user_canceled",
						"cancelSourceReason": "Order canceled by user",
						"algoClientOrderId": "aaa",
						"algoId": "11756185",
						"brokerId": ""
					}
				]
			}`))
		}))

		resp, err := client.NewGetOrderHistoryService().
			InstId("ETH-USDT").
			OrderType("limit").
			State("canceled").
			Limit("20").
			Do(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "0", resp.Code)
		assert.Equal(t, "success", resp.Msg)
		assert.Len(t, resp.Data, 1)
		assert.Equal(t, "29419717", resp.Data[0].OrderId)
		assert.Equal(t, "aabbc", resp.Data[0].ClientOrderId)
		assert.Equal(t, "ETH-USDT", resp.Data[0].InstId)
		assert.Equal(t, "cross", resp.Data[0].MarginMode)
		assert.Equal(t, "net", resp.Data[0].PositionSide)
		assert.Equal(t, "buy", resp.Data[0].Side)
		assert.Equal(t, "limit", resp.Data[0].OrderType)
		assert.Equal(t, "1523.000000000000000000", resp.Data[0].Price)
		assert.Equal(t, "1.000000000000000000", resp.Data[0].Size)
		assert.Equal(t, "false", resp.Data[0].ReduceOnly)
		assert.Equal(t, "3", resp.Data[0].Leverage)
		assert.Equal(t, "canceled", resp.Data[0].State)
		assert.Equal(t, "0.000000000000000000", resp.Data[0].FilledSize)
		assert.Equal(t, "0.000000000000000000", resp.Data[0].Pnl)
		assert.Equal(t, "0.000000000000000000", resp.Data[0].AveragePrice)
		assert.Equal(t, "0.000000000000000000", resp.Data[0].Fee)
		assert.Equal(t, "1697010303781", resp.Data[0].CreateTime)
		assert.Equal(t, "1697014607770", resp.Data[0].UpdateTime)
		assert.Equal(t, "normal", resp.Data[0].OrderCategory)
		assert.Equal(t, "user_canceled", resp.Data[0].CancelSource)
		assert.Equal(t, "Order canceled by user", resp.Data[0].CancelSourceReason)
		assert.Equal(t, "aaa", resp.Data[0].AlgoClientOrderId)
		assert.Equal(t, "11756185", resp.Data[0].AlgoId)
		assert.Equal(t, "", resp.Data[0].BrokerId)
	})

	t.Run("validation errors", func(t *testing.T) {
		client := newTestPrivateRestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		// Check that after and before cannot be used simultaneously
		_, err := client.NewGetOrderHistoryService().
			After("123").
			Before("456").
			Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "after and before cannot be used simultaneously")
	})

	t.Run("pagination", func(t *testing.T) {
		client := newTestPrivateRestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "123", r.URL.Query().Get("after"))
			assert.Equal(t, "50", r.URL.Query().Get("limit"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"code": "0",
				"msg": "success",
				"data": []
			}`))
		}))

		resp, err := client.NewGetOrderHistoryService().
			After("123").
			Limit("50").
			Do(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "0", resp.Code)
		assert.Equal(t, "success", resp.Msg)
		assert.Empty(t, resp.Data)
	})

	t.Run("time filter", func(t *testing.T) {
		client := newTestPrivateRestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "1697010303781", r.URL.Query().Get("begin"))
			assert.Equal(t, "1697014607770", r.URL.Query().Get("end"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"code": "0",
				"msg": "success",
				"data": []
			}`))
		}))

		resp, err := client.NewGetOrderHistoryService().
			Begin("1697010303781").
			End("1697014607770").
			Do(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "0", resp.Code)
		assert.Equal(t, "success", resp.Msg)
		assert.Empty(t, resp.Data)
	})
}

func TestGetTPSLOrderHistoryService(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		service        *GetTPSLOrderHistoryService
		expectedParams map[string]string
		expectedError  string
	}{
		{
			name: "Basic request",
			service: newTestPrivateRestClient(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"code": "0",
					"msg": "success",
					"data": [
						{
							"tpslId": "2101",
							"clientOrderId": "BBBBqqqq",
							"instId": "ETH-USDT",
							"marginMode": "cross",
							"positionSide": "net",
							"side": "sell",
							"orderType": null,
							"size": "1",
							"reduceOnly": "true",
							"leverage": "3",
							"state": "canceled",
							"actualSize": null,
							"triggerType": null,
							"orderCategory": "normal",
							"tpTriggerPrice": "1661.100000000000000000",
							"tpOrderPrice": null,
							"slTriggerPrice": null,
							"slOrderPrice": null,
							"brokerId": ""
						}
					]
				}`))
			}).NewGetTPSLOrderHistoryService(),
			expectedParams: map[string]string{},
		},
		{
			name: "With all parameters",
			service: newTestPrivateRestClient(func(w http.ResponseWriter, r *http.Request) {
				// Check query parameters
				query := r.URL.Query()
				expected := map[string]string{
					"instId":        "BTC-USDT",
					"tpslId":        "123",
					"clientOrderId": "client-123",
					"state":         "live",
					"limit":         "20",
				}
				for key, value := range expected {
					if query.Get(key) != value {
						t.Errorf("expected %s=%s, got %s=%s", key, value, key, query.Get(key))
					}
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"code": "0",
					"msg": "success",
					"data": [
						{
							"tpslId": "2101",
							"clientOrderId": "BBBBqqqq",
							"instId": "ETH-USDT",
							"marginMode": "cross",
							"positionSide": "net",
							"side": "sell",
							"orderType": null,
							"size": "1",
							"reduceOnly": "true",
							"leverage": "3",
							"state": "canceled",
							"actualSize": null,
							"triggerType": null,
							"orderCategory": "normal",
							"tpTriggerPrice": "1661.100000000000000000",
							"tpOrderPrice": null,
							"slTriggerPrice": null,
							"slOrderPrice": null,
							"brokerId": ""
						}
					]
				}`))
			}).NewGetTPSLOrderHistoryService().
				InstId("BTC-USDT").
				TpslId("123").
				ClientOrderId("client-123").
				State("live").
				Limit("20"),
			expectedParams: map[string]string{
				"instId":        "BTC-USDT",
				"tpslId":        "123",
				"clientOrderId": "client-123",
				"state":         "live",
				"limit":         "20",
			},
		},
		{
			name: "With after parameter",
			service: newTestPrivateRestClient(func(w http.ResponseWriter, r *http.Request) {
				// Check query parameters
				query := r.URL.Query()
				if query.Get("after") != "123" {
					t.Errorf("expected after=123, got after=%s", query.Get("after"))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"code": "0",
					"msg": "success",
					"data": [
						{
							"tpslId": "2101",
							"clientOrderId": "BBBBqqqq",
							"instId": "ETH-USDT",
							"marginMode": "cross",
							"positionSide": "net",
							"side": "sell",
							"orderType": null,
							"size": "1",
							"reduceOnly": "true",
							"leverage": "3",
							"state": "canceled",
							"actualSize": null,
							"triggerType": null,
							"orderCategory": "normal",
							"tpTriggerPrice": "1661.100000000000000000",
							"tpOrderPrice": null,
							"slTriggerPrice": null,
							"slOrderPrice": null,
							"brokerId": ""
						}
					]
				}`))
			}).NewGetTPSLOrderHistoryService().
				After("123"),
			expectedParams: map[string]string{
				"after": "123",
			},
		},
		{
			name: "With before parameter",
			service: newTestPrivateRestClient(func(w http.ResponseWriter, r *http.Request) {
				// Check query parameters
				query := r.URL.Query()
				if query.Get("before") != "123" {
					t.Errorf("expected before=123, got before=%s", query.Get("before"))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"code": "0",
					"msg": "success",
					"data": [
						{
							"tpslId": "2101",
							"clientOrderId": "BBBBqqqq",
							"instId": "ETH-USDT",
							"marginMode": "cross",
							"positionSide": "net",
							"side": "sell",
							"orderType": null,
							"size": "1",
							"reduceOnly": "true",
							"leverage": "3",
							"state": "canceled",
							"actualSize": null,
							"triggerType": null,
							"orderCategory": "normal",
							"tpTriggerPrice": "1661.100000000000000000",
							"tpOrderPrice": null,
							"slTriggerPrice": null,
							"slOrderPrice": null,
							"brokerId": ""
						}
					]
				}`))
			}).NewGetTPSLOrderHistoryService().
				Before("123"),
			expectedParams: map[string]string{
				"before": "123",
			},
		},
		{
			name: "Error: after and before used simultaneously",
			service: newTestPrivateRestClient(nil).NewGetTPSLOrderHistoryService().
				After("123").
				Before("456"),
			expectedError: "after and before cannot be used simultaneously",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute request
			resp, err := tt.service.Do(context.Background())

			// Check error
			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("expected error %s, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError {
					t.Errorf("expected error %s, got %s", tt.expectedError, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Check response
			if resp.Code != "0" {
				t.Errorf("expected code 0, got %s", resp.Code)
			}
			if resp.Msg != "success" {
				t.Errorf("expected msg success, got %s", resp.Msg)
			}
			if len(resp.Data) != 1 {
				t.Errorf("expected 1 order, got %d", len(resp.Data))
			}

			order := resp.Data[0]
			if order.TpslId != "2101" {
				t.Errorf("expected tpslId 2101, got %s", order.TpslId)
			}
			if order.ClientOrderId != "BBBBqqqq" {
				t.Errorf("expected clientOrderId BBBBqqqq, got %s", order.ClientOrderId)
			}
			if order.InstId != "ETH-USDT" {
				t.Errorf("expected instId ETH-USDT, got %s", order.InstId)
			}
		})
	}
}

func TestGetAlgoOrderHistoryService(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/api/v1/trade/orders-algo-history", r.URL.Path)
		assert.Equal(t, "trigger", r.URL.Query().Get("orderType"))
		assert.Equal(t, "ETH-USDT", r.URL.Query().Get("instId"))
		assert.Equal(t, "2101", r.URL.Query().Get("algoId"))
		assert.Equal(t, "BBBBqqqq", r.URL.Query().Get("clientOrderId"))
		assert.Equal(t, "canceled", r.URL.Query().Get("state"))
		assert.Equal(t, "20", r.URL.Query().Get("limit"))

		response := `{
			"code": "0",
			"msg": "success",
			"data": [
				{
					"algoId": "2101",
					"clientOrderId": "BBBBqqqq",
					"instId": "ETH-USDT",
					"marginMode": "cross",
					"positionSide": "net",
					"side": "sell",
					"orderType": "trigger",
					"size": "1",
					"actualSize": "1",
					"leverage": "3",
					"state": "canceled",
					"triggerPrice": "1661.100000000000000000",
					"triggerPriceType": "last",
					"brokerId": "",
					"attachAlgoOrders": [
						{
							"tpTriggerPrice": "1666.000000000000000000",
							"tpOrderPrice": "-1",
							"tpTriggerPriceType": "last",
							"slTriggerPrice": "1222.000000000000000000",
							"slOrderPrice": "-1",
							"slTriggerPriceType": "last"
						}
					]
				}
			]
		}`
		w.Write([]byte(response))
	}

	client := newTestPrivateRestClient(handler)
	service := client.NewGetAlgoOrderHistoryService()

	service.OrderType("trigger").
		InstId("ETH-USDT").
		AlgoId("2101").
		ClientOrderId("BBBBqqqq").
		State("canceled").
		Limit("20")

	resp, err := service.Do(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "0", resp.Code)
	assert.Equal(t, "success", resp.Msg)
	assert.Len(t, resp.Data, 1)

	order := resp.Data[0]
	assert.Equal(t, "2101", order.AlgoId)
	assert.Equal(t, "BBBBqqqq", order.ClientOrderId)
	assert.Equal(t, "ETH-USDT", order.InstId)
	assert.Equal(t, "cross", order.MarginMode)
	assert.Equal(t, "net", order.PositionSide)
	assert.Equal(t, "sell", order.Side)
	assert.Equal(t, "trigger", order.OrderType)
	assert.Equal(t, "1", order.Size)
	assert.Equal(t, "1", order.ActualSize)
	assert.Equal(t, "3", order.Leverage)
	assert.Equal(t, "canceled", order.State)
	assert.Equal(t, "1661.100000000000000000", order.TriggerPrice)
	assert.Equal(t, "last", order.TriggerPriceType)
	assert.Equal(t, "", order.BrokerId)
	assert.Len(t, order.AttachAlgoOrders, 1)

	attachOrder := order.AttachAlgoOrders[0]
	assert.Equal(t, "1666.000000000000000000", attachOrder.TpTriggerPrice)
	assert.Equal(t, "-1", attachOrder.TpOrderPrice)
	assert.Equal(t, "last", attachOrder.TpTriggerPriceType)
	assert.Equal(t, "1222.000000000000000000", attachOrder.SlTriggerPrice)
	assert.Equal(t, "-1", attachOrder.SlOrderPrice)
	assert.Equal(t, "last", attachOrder.SlTriggerPriceType)
}

func TestGetAlgoOrderHistoryService_Validation(t *testing.T) {
	client := newTestClient(t)
	service := client.NewGetAlgoOrderHistoryService()

	// Test missing orderType
	_, err := service.Do(context.Background())
	assert.Error(t, err)
	assert.Equal(t, "orderType required", err.Error())

	// Test after and before used simultaneously
	service.OrderType("trigger").After("123").Before("456")
	_, err = service.Do(context.Background())
	assert.Error(t, err)
	assert.Equal(t, "after and before cannot be used simultaneously", err.Error())
}

func TestGetTradeHistoryService(t *testing.T) {
	client := newTestPrivateRestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/api/v1/trade/fills-history", r.URL.Path)
		assert.NotEmpty(t, r.Header.Get("ACCESS-KEY"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-SIGN"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-TIMESTAMP"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-NONCE"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-PASSPHRASE"))

		// Проверяем параметры запроса
		assert.Equal(t, "BTC-USDT", r.URL.Query().Get("instId"))
		assert.Equal(t, "100", r.URL.Query().Get("limit"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"code": "0",
			"msg": "",
			"data": [{
				"instId": "BTC-USDT",
				"tradeId": "123456",
				"orderId": "789012",
				"fillPrice": "50000.0",
				"fillSize": "0.1",
				"fillPnl": "100.0",
				"positionSide": "long",
				"side": "buy",
				"fee": "2.5",
				"ts": "1620000000000",
				"brokerId": "9999"
			}]
		}`))
	}))

	// Тест успешного получения истории сделок
	t.Run("Success", func(t *testing.T) {
		resp, err := client.NewGetTradeHistoryService().
			InstId("BTC-USDT").
			Limit("100").
			Do(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "0", resp.Code)
		assert.Len(t, resp.Data, 1)
		assert.Equal(t, "BTC-USDT", resp.Data[0].InstId)
		assert.Equal(t, "123456", resp.Data[0].TradeId)
		assert.Equal(t, "789012", resp.Data[0].OrderId)
		assert.Equal(t, "50000.0", resp.Data[0].FillPrice)
		assert.Equal(t, "0.1", resp.Data[0].FillSize)
		assert.Equal(t, "100.0", resp.Data[0].FillPnl)
		assert.Equal(t, "long", resp.Data[0].PositionSide)
		assert.Equal(t, "buy", resp.Data[0].Side)
		assert.Equal(t, "2.5", resp.Data[0].Fee)
		assert.Equal(t, "1620000000000", resp.Data[0].Ts)
		assert.Equal(t, "9999", resp.Data[0].BrokerId)
	})

	// Тест валидации параметров
	t.Run("Validation", func(t *testing.T) {
		_, err := client.NewGetTradeHistoryService().
			After("100").
			Before("200").
			Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "after and before cannot be used simultaneously")
	})
}

func TestGetOrderPriceRangeService(t *testing.T) {
	client := newTestPrivateRestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/api/v1/trade/order/price-range", r.URL.Path)
		assert.NotEmpty(t, r.Header.Get("ACCESS-KEY"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-SIGN"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-TIMESTAMP"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-NONCE"))
		assert.NotEmpty(t, r.Header.Get("ACCESS-PASSPHRASE"))

		// Проверяем параметры запроса
		assert.Equal(t, "BTC-USDT", r.URL.Query().Get("instId"))
		assert.Equal(t, "buy", r.URL.Query().Get("side"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"code": "0",
			"msg": "success",
			"data": {
				"maxPrice": "1587.800000000000000000",
				"minPrice": "1187.000000000000000000"
			}
		}`))
	}))

	// Тест успешного получения диапазона цен
	t.Run("Success", func(t *testing.T) {
		resp, err := client.NewGetOrderPriceRangeService().
			InstId("BTC-USDT").
			Side("buy").
			Do(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "0", resp.Code)
		assert.Equal(t, "success", resp.Msg)
		assert.Equal(t, "1587.800000000000000000", resp.Data.MaxPrice)
		assert.Equal(t, "1187.000000000000000000", resp.Data.MinPrice)
	})

	// Тест валидации параметров
	t.Run("Validation", func(t *testing.T) {
		// Тест без instId
		_, err := client.NewGetOrderPriceRangeService().
			Side("buy").
			Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "instId required")

		// Тест без side
		_, err = client.NewGetOrderPriceRangeService().
			InstId("BTC-USDT").
			Do(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "side required")
	})
}
