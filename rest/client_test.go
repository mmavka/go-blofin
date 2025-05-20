package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestRestClient(handler http.HandlerFunc) *RestClient {
	ts := httptest.NewServer(handler)
	return NewRestClient(ts.URL)
}

func TestGetInstrumentsService(t *testing.T) {
	mockResp := `{"code":"0","data":[{"instId":"BTC-USDT"}]}`
	client := newTestRestClient(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/market/instruments") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResp))
	})
	resp, err := client.NewGetInstrumentsService().Do(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].InstID != "BTC-USDT" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestGetTickersService(t *testing.T) {
	mockResp := `{"code":"0","data":[{"instId":"BTC-USDT","last":"100"}]}`
	client := newTestRestClient(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/market/tickers") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResp))
	})
	resp, err := client.NewGetTickersService().Do(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].InstID != "BTC-USDT" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestGetOrderBookService(t *testing.T) {
	mockResp := `{"code":"0","data":[{"asks":[["100","1"]],"bids":[["99","2"]]}]}`
	client := newTestRestClient(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/market/books") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResp))
	})
	resp, err := client.NewGetOrderBookService().InstId("BTC-USDT").Do(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || len(resp.Data[0].Asks) == 0 || len(resp.Data[0].Bids) == 0 {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestGetTradesService(t *testing.T) {
	mockResp := `{"code":"0","data":[{"tradeId":"1","instId":"BTC-USDT","price":"100","size":"0.1","side":"buy","ts":"123456"}]}`
	client := newTestRestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResp))
	})
	resp, err := client.NewGetTradesService().InstId("BTC-USDT").Do(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].InstID != "BTC-USDT" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestGetMarkPriceService(t *testing.T) {
	mockResp := `{"code":"0","data":[{"instId":"BTC-USDT","indexPrice":"100","markPrice":"101","ts":"123456"}]}`
	client := newTestRestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResp))
	})
	resp, err := client.NewGetMarkPriceService().InstId("BTC-USDT").Do(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].InstID != "BTC-USDT" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestGetFundingRateHistoryService(t *testing.T) {
	mockResp := `{"code":"0","data":[{"instId":"BTC-USDT","fundingRate":"0.01","fundingTime":"123456"}]}`
	client := newTestRestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResp))
	})
	resp, err := client.NewGetFundingRateHistoryService().InstId("BTC-USDT").Do(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].InstID != "BTC-USDT" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestGetCandlesService(t *testing.T) {
	mockResp := `{"code":"0","data":[["123456","100","110","90","105","10","1000","1000","1"]]}`
	client := newTestRestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResp))
	})
	resp, err := client.NewGetCandlesService().InstId("BTC-USDT").Do(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].Open != "100" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestRestClientSetBaseURL(t *testing.T) {
	called := false
	newURL := "http://localhost:12345"
	client := NewDefaultRestClient()
	client.SetBaseURL(newURL)
	if client.baseURL != newURL {
		t.Errorf("baseURL not updated: got %s, want %s", client.baseURL, newURL)
	}
	// Check that httpClient is also updated (via SetBaseURL)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.Write([]byte(`{"code":"0","data":[]}`))
	}))
	defer ts.Close()
	client.SetBaseURL(ts.URL)
	_, _ = client.NewGetInstrumentsService().Do(context.Background())
	if !called {
		t.Error("httpClient did not use new baseURL")
	}
}
