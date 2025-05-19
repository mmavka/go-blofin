package ws

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

// Handler handles WebSocket connections and messages
type Handler struct {
	conn               *websocket.Conn
	positionsCallback  func(Position)
	ordersCallback     func(Order)
	algoOrdersCallback func(AlgoOrder)
	accountCallback    func(Account)
}

// NewHandler creates a new Handler
func NewHandler(conn *websocket.Conn) *Handler {
	return &Handler{
		conn: conn,
	}
}

// handlePositions handles positions channel messages
func (h *Handler) handlePositions(msg []byte) error {
	var positionsData PositionsData
	if err := json.Unmarshal(msg, &positionsData); err != nil {
		return fmt.Errorf("failed to unmarshal positions data: %w", err)
	}

	for _, position := range positionsData.Data {
		h.positionsCallback(position)
	}
	return nil
}

// handleOrders handles orders channel messages
func (h *Handler) handleOrders(msg []byte) error {
	var ordersData OrdersData
	if err := json.Unmarshal(msg, &ordersData); err != nil {
		return fmt.Errorf("failed to unmarshal orders data: %w", err)
	}

	for _, order := range ordersData.Data {
		h.ordersCallback(order)
	}
	return nil
}

// handleAlgoOrders handles algo orders channel messages
func (h *Handler) handleAlgoOrders(msg []byte) error {
	var algoOrdersData AlgoOrdersData
	if err := json.Unmarshal(msg, &algoOrdersData); err != nil {
		return fmt.Errorf("failed to unmarshal algo orders data: %w", err)
	}

	for _, order := range algoOrdersData.Data {
		h.algoOrdersCallback(order)
	}
	return nil
}

// handleAccount handles account channel messages
func (h *Handler) handleAccount(msg []byte) error {
	var accountData AccountData
	if err := json.Unmarshal(msg, &accountData); err != nil {
		return fmt.Errorf("failed to unmarshal account data: %w", err)
	}

	h.accountCallback(accountData.Data)
	return nil
}

// SubscribePositions subscribes to positions channel
func (h *Handler) SubscribePositions(instId string, callback func(Position)) error {
	h.positionsCallback = callback

	req := PositionsRequest{
		Op: "subscribe",
		Args: []struct {
			Channel string `json:"channel"`
			InstId  string `json:"instId,omitempty"`
		}{
			{
				Channel: "positions",
				InstId:  instId,
			},
		},
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal positions request: %w", err)
	}

	if err := h.conn.WriteMessage(websocket.TextMessage, reqBytes); err != nil {
		return fmt.Errorf("failed to send positions request: %w", err)
	}

	return nil
}

// SubscribeOrders subscribes to orders channel
func (h *Handler) SubscribeOrders(instId string, callback func(Order)) error {
	h.ordersCallback = callback

	req := OrdersRequest{
		Op: "subscribe",
		Args: []struct {
			Channel string `json:"channel"`
			InstId  string `json:"instId,omitempty"`
		}{
			{
				Channel: "orders",
				InstId:  instId,
			},
		},
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal orders request: %w", err)
	}

	if err := h.conn.WriteMessage(websocket.TextMessage, reqBytes); err != nil {
		return fmt.Errorf("failed to send orders request: %w", err)
	}

	return nil
}

// SubscribeAlgoOrders subscribes to algo orders channel
func (h *Handler) SubscribeAlgoOrders(instId string, callback func(AlgoOrder)) error {
	h.algoOrdersCallback = callback

	req := AlgoOrdersRequest{
		Op: "subscribe",
		Args: []struct {
			Channel string `json:"channel"`
			InstId  string `json:"instId,omitempty"`
		}{
			{
				Channel: "orders-algo",
				InstId:  instId,
			},
		},
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal algo orders request: %w", err)
	}

	if err := h.conn.WriteMessage(websocket.TextMessage, reqBytes); err != nil {
		return fmt.Errorf("failed to send algo orders request: %w", err)
	}

	return nil
}

// SubscribeAccount subscribes to account channel
func (h *Handler) SubscribeAccount(callback func(Account)) error {
	h.accountCallback = callback

	req := AccountRequest{
		Op: "subscribe",
		Args: []struct {
			Channel string `json:"channel"`
		}{
			{
				Channel: "account",
			},
		},
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal account request: %w", err)
	}

	if err := h.conn.WriteMessage(websocket.TextMessage, reqBytes); err != nil {
		return fmt.Errorf("failed to send account request: %w", err)
	}

	return nil
}

// HandleMessages handles incoming WebSocket messages
func (h *Handler) HandleMessages() error {
	for {
		_, msg, err := h.conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}

		// Try to unmarshal as base message first
		var baseMsg struct {
			Arg struct {
				Channel string `json:"channel"`
			} `json:"arg"`
			Event string `json:"event"`
		}
		if err := json.Unmarshal(msg, &baseMsg); err != nil {
			return fmt.Errorf("failed to unmarshal base message: %w", err)
		}

		// Handle subscription response
		if baseMsg.Event == "subscribe" {
			continue
		}

		// Handle data messages
		switch baseMsg.Arg.Channel {
		case "positions":
			if err := h.handlePositions(msg); err != nil {
				return fmt.Errorf("failed to handle positions message: %w", err)
			}
		case "orders":
			if err := h.handleOrders(msg); err != nil {
				return fmt.Errorf("failed to handle orders message: %w", err)
			}
		case "orders-algo":
			if err := h.handleAlgoOrders(msg); err != nil {
				return fmt.Errorf("failed to handle algo orders message: %w", err)
			}
		case "account":
			if err := h.handleAccount(msg); err != nil {
				return fmt.Errorf("failed to handle account message: %w", err)
			}
		}
	}
}
