package ws

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gorilla/websocket"
)

// Handler handles WebSocket connections and messages
type Handler struct {
	conn               *websocket.Conn
	positionsCallback  func(Position)
	ordersCallback     func(Order)
	algoOrdersCallback func(AlgoOrder)
	accountCallback    func(Account)
	orderBook          *OrderBook
	ticker             *Ticker
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

// handleOrderBook handles order book channel messages
func (h *Handler) handleOrderBook(data []byte) error {
	var msg struct {
		Data struct {
			Bids [][]string `json:"bids"`
			Asks [][]string `json:"asks"`
		} `json:"data"`
	}

	if err := json.Unmarshal(data, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal order book: %w", err)
	}

	// Очищаем стакан если пришли пустые массивы
	if len(msg.Data.Bids) == 0 && len(msg.Data.Asks) == 0 {
		h.orderBook = &OrderBook{
			Bids: make([]OrderBookEntry, 0),
			Asks: make([]OrderBookEntry, 0),
		}
		return nil
	}

	// Обновляем стакан
	h.orderBook = &OrderBook{
		Bids: make([]OrderBookEntry, len(msg.Data.Bids)),
		Asks: make([]OrderBookEntry, len(msg.Data.Asks)),
	}

	for i, bid := range msg.Data.Bids {
		price, _ := strconv.ParseFloat(bid[0], 64)
		size, _ := strconv.ParseFloat(bid[1], 64)
		h.orderBook.Bids[i] = OrderBookEntry{Price: price, Size: size}
	}

	for i, ask := range msg.Data.Asks {
		price, _ := strconv.ParseFloat(ask[0], 64)
		size, _ := strconv.ParseFloat(ask[1], 64)
		h.orderBook.Asks[i] = OrderBookEntry{Price: price, Size: size}
	}

	return nil
}

// handleTicker handles ticker channel messages
func (h *Handler) handleTicker(data []byte) error {
	var msg struct {
		Data struct {
			LastPrice          string `json:"lastPrice"`
			LastSize           string `json:"lastSize"`
			BestBidPrice       string `json:"bestBidPrice"`
			BestBidSize        string `json:"bestBidSize"`
			BestAskPrice       string `json:"bestAskPrice"`
			BestAskSize        string `json:"bestAskSize"`
			Volume24h          string `json:"volume24h"`
			PriceChange        string `json:"priceChange"`
			PriceChangePercent string `json:"priceChangePercent"`
			Timestamp          int64  `json:"timestamp"`
		} `json:"data"`
	}

	if err := json.Unmarshal(data, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal ticker: %w", err)
	}

	// Обрабатываем null значения
	lastPrice, _ := strconv.ParseFloat(msg.Data.LastPrice, 64)
	lastSize, _ := strconv.ParseFloat(msg.Data.LastSize, 64)
	bestBidPrice, _ := strconv.ParseFloat(msg.Data.BestBidPrice, 64)
	bestBidSize, _ := strconv.ParseFloat(msg.Data.BestBidSize, 64)
	bestAskPrice, _ := strconv.ParseFloat(msg.Data.BestAskPrice, 64)
	bestAskSize, _ := strconv.ParseFloat(msg.Data.BestAskSize, 64)
	volume24h, _ := strconv.ParseFloat(msg.Data.Volume24h, 64)
	priceChange, _ := strconv.ParseFloat(msg.Data.PriceChange, 64)
	priceChangePercent, _ := strconv.ParseFloat(msg.Data.PriceChangePercent, 64)

	h.ticker = &Ticker{
		LastPrice:          lastPrice,
		LastSize:           lastSize,
		BestBidPrice:       bestBidPrice,
		BestBidSize:        bestBidSize,
		BestAskPrice:       bestAskPrice,
		BestAskSize:        bestAskSize,
		Volume24h:          volume24h,
		PriceChange:        priceChange,
		PriceChangePercent: priceChangePercent,
		Timestamp:          msg.Data.Timestamp,
	}

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
		case "order-book":
			if err := h.handleOrderBook(msg); err != nil {
				return fmt.Errorf("failed to handle order book message: %w", err)
			}
		case "ticker":
			if err := h.handleTicker(msg); err != nil {
				return fmt.Errorf("failed to handle ticker message: %w", err)
			}
		}
	}
}

// GetOrderBook returns the current order book state
func (h *Handler) GetOrderBook() *OrderBook {
	return h.orderBook
}

// GetTicker returns the current ticker state
func (h *Handler) GetTicker() *Ticker {
	return h.ticker
}
