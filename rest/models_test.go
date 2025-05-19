package rest

import (
	"encoding/json"
	"testing"
)

func TestModelUnmarshal(t *testing.T) {
	// TODO: реализовать тесты для проверки сериализации/десериализации моделей (например, Instrument, Ticker и др.)
	_ = json.Unmarshal
}

func TestCandleUnmarshalJSON_Valid(t *testing.T) {
	jsonStr := `["123456","100","110","90","105","10","1000","1000","1"]`
	var c Candle
	err := json.Unmarshal([]byte(jsonStr), &c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Open != "100" || c.Close != "105" {
		t.Errorf("unexpected values: %+v", c)
	}
}

func TestCandleUnmarshalJSON_Empty(t *testing.T) {
	jsonStr := `[]`
	var c Candle
	err := json.Unmarshal([]byte(jsonStr), &c)
	if err == nil {
		t.Error("expected error for empty array")
	}
}

func TestCandleUnmarshalJSON_InvalidLength(t *testing.T) {
	jsonStr := `["1","2"]`
	var c Candle
	err := json.Unmarshal([]byte(jsonStr), &c)
	if err == nil {
		t.Error("expected error for invalid array length")
	}
}

func TestCandleUnmarshalJSON_InvalidType(t *testing.T) {
	jsonStr := `{}`
	var c Candle
	err := json.Unmarshal([]byte(jsonStr), &c)
	if err == nil {
		t.Error("expected error for invalid type")
	}
}
