package ws

import (
	"encoding/json"
	"testing"
)

func TestWSModelUnmarshal(t *testing.T) {
	// TODO: реализовать тесты для проверки сериализации/десериализации моделей WebSocket (например, TradeWS, CandleWSMessage и др.)
	_ = json.Unmarshal
}

func TestTradeWSMessageUnmarshal(t *testing.T) {
	jsonStr := `{"arg":{"channel":"trades","instId":"BTC-USDT"},"data":[{"tradeId":"1","price":"100","size":"0.1","side":"buy","ts":"123456"}]}`
	var msg TradeWSMessage
	err := json.Unmarshal([]byte(jsonStr), &msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Arg.Channel != "trades" || len(msg.Data) != 1 || msg.Data[0].Price != "100" {
		t.Errorf("unexpected result: %+v", msg)
	}
}

func TestCandleWSMessageUnmarshal(t *testing.T) {
	jsonStr := `{"arg":{"channel":"candle1m","instId":"BTC-USDT"},"data":[["123456","100","110","90","105","10","1000","1000","1"]]}`
	var msg CandleWSMessage
	err := json.Unmarshal([]byte(jsonStr), &msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Arg.Channel != "candle1m" || len(msg.Data) != 1 || msg.Data[0].Open != "100" {
		t.Errorf("unexpected result: %+v", msg)
	}
}

func TestFundingRateWSMessageUnmarshal(t *testing.T) {
	jsonStr := `{"arg":{"channel":"funding-rate","instId":"BTC-USDT"},"data":[{"fundingRate":"0.0001","fundingTime":"1700726400000","instId":"BTC-USDT"}]}`
	var msg FundingRateWSMessage
	err := json.Unmarshal([]byte(jsonStr), &msg)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if msg.Arg.Channel != "funding-rate" || len(msg.Data) != 1 || msg.Data[0].FundingRate != "0.0001" {
		t.Errorf("unexpected result: %+v", msg)
	}
}
