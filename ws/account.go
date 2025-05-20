package ws

import (
	"encoding/json"
	"fmt"

	"github.com/mmavka/go-blofin"
)

// WsAccountHandler handle websocket account event
type WsAccountHandler func(event *WsAccountEvent)

// WsAccountEvent define websocket account event
type WsAccountEvent struct {
	EventType string `json:"e"`
	Time      int64  `json:"E"`
	Account   struct {
		Balances []WsBalance `json:"B"`
	} `json:"a"`
}

// WsBalance define websocket balance
type WsBalance struct {
	Asset                  string `json:"a"`
	Free                   string `json:"f"`
	Locked                 string `json:"l"`
	UnrealizedProfit       string `json:"up"`
	MarginBalance          string `json:"mb"`
	MaintenanceMargin      string `json:"mm"`
	InitialMargin          string `json:"im"`
	PositionInitialMargin  string `json:"pim"`
	OpenOrderInitialMargin string `json:"oim"`
	CrossWalletBalance     string `json:"cw"`
	AvailableBalance       string `json:"ab"`
	MaxWithdrawAmount      string `json:"mwa"`
}

// WsAccountServe serve websocket that pushes account information
func WsAccountServe(apiKey, secretKey string, handler WsAccountHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/account/%s", blofin.DefaultWSPrivate, apiKey)
	config := &WsConfig{
		Endpoint: endpoint,
	}
	wsHandler := func(message []byte) {
		var event WsAccountEvent
		if err := json.Unmarshal(message, &event); err != nil {
			errHandler(err)
			return
		}
		handler(&event)
	}
	return WsServe(config, wsHandler, errHandler)
}

// WsPositionHandler handle websocket position event
type WsPositionHandler func(event *WsPositionEvent)

// WsPositionEvent define websocket position event
type WsPositionEvent struct {
	EventType string `json:"e"`
	Time      int64  `json:"E"`
	Position  struct {
		Symbol              string `json:"s"`
		PositionAmount      string `json:"pa"`
		EntryPrice          string `json:"ep"`
		AccumulatedRealized string `json:"cr"`
		UnrealizedProfit    string `json:"up"`
		MarginType          string `json:"mt"`
		IsolatedWallet      string `json:"iw"`
		PositionSide        string `json:"ps"`
	} `json:"p"`
}

// WsPositionServe serve websocket that pushes position information
func WsPositionServe(apiKey, secretKey string, handler WsPositionHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/positions/%s", blofin.DefaultWSPrivate, apiKey)
	config := &WsConfig{
		Endpoint: endpoint,
	}
	wsHandler := func(message []byte) {
		var event WsPositionEvent
		if err := json.Unmarshal(message, &event); err != nil {
			errHandler(err)
			return
		}
		handler(&event)
	}
	return WsServe(config, wsHandler, errHandler)
}

// WsOrderHandler handle websocket order event
type WsOrderHandler func(event *WsOrderEvent)

// WsOrderEvent define websocket order event
type WsOrderEvent struct {
	EventType string `json:"e"`
	Time      int64  `json:"E"`
	Order     struct {
		Symbol               string `json:"s"`
		ClientOrderID        string `json:"c"`
		Side                 string `json:"S"`
		Type                 string `json:"o"`
		TimeInForce          string `json:"f"`
		Quantity             string `json:"q"`
		Price                string `json:"p"`
		StopPrice            string `json:"sp"`
		ExecutionType        string `json:"x"`
		OrderStatus          string `json:"X"`
		RejectReason         string `json:"r"`
		OrderID              string `json:"i"`
		LastFilledQty        string `json:"l"`
		FilledAccQty         string `json:"z"`
		LastFilledPrice      string `json:"L"`
		Commission           string `json:"n"`
		CommissionAsset      string `json:"N"`
		OrderTradeTime       int64  `json:"T"`
		TradeID              string `json:"t"`
		IsOnBook             bool   `json:"w"`
		IsMaker              bool   `json:"m"`
		Ignore               bool   `json:"M"`
		CreateTime           int64  `json:"O"`
		FilledQuoteQty       string `json:"Z"`
		QuoteQty             string `json:"Q"`
		WorkingTime          int64  `json:"W"`
		SelfTradePreventMode string `json:"V"`
		GoodTillDate         int64  `json:"g"`
	} `json:"o"`
}

// WsOrderServe serve websocket that pushes order information
func WsOrderServe(apiKey, secretKey string, handler WsOrderHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/orders/%s", blofin.DefaultWSPrivate, apiKey)
	config := &WsConfig{
		Endpoint: endpoint,
	}
	wsHandler := func(message []byte) {
		var event WsOrderEvent
		if err := json.Unmarshal(message, &event); err != nil {
			errHandler(err)
			return
		}
		handler(&event)
	}
	return WsServe(config, wsHandler, errHandler)
}
