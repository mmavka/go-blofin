package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestPrivateRestClient(handler http.HandlerFunc) *RestClient {
	ts := httptest.NewServer(handler)
	client := NewRestClient(ts.URL)
	client.SetAuth("key", "secret", "pass")
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

func TestGetPositionsService(t *testing.T) {
	mockResp := `{"code":"0","data":[{"instId":"BTC-USDT","pos":"1","side":"long"}]}`
	client := newTestPrivateRestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResp))
	})
	resp, err := client.NewGetPositionsService().Do(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].InstID != "BTC-USDT" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestPlaceOrderService(t *testing.T) {
	mockResp := `{"code":"0","data":{"orderId":"123"}}`
	client := newTestPrivateRestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResp))
	})
	order := OrderRequest{InstID: "BTC-USDT", Side: "buy", OrderType: "limit", Price: "100", Size: "1"}
	resp, err := client.NewPlaceOrderService().Order(order).Do(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Data == nil {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestCancelOrderService(t *testing.T) {
	mockResp := `{"code":"0","data":{"orderId":"123"}}`
	client := newTestPrivateRestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResp))
	})
	resp, err := client.NewCancelOrderService().OrderId("123").Do(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Data == nil {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestGetActiveOrdersService(t *testing.T) {
	mockResp := `{"code":"0","data":{"orders":[]}}`
	client := newTestPrivateRestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResp))
	})
	resp, err := client.NewGetActiveOrdersService().Do(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Data == nil {
		t.Errorf("unexpected response: %+v", resp)
	}
}
