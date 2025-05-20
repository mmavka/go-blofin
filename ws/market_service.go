package ws

import (
	"encoding/json"
	"fmt"
	"sync"
)

// MarketService handles market data websocket connections
type MarketService struct {
	client     *WsClient
	config     *WsConfig
	mu         sync.RWMutex
	subs       map[string]struct{} // track subscriptions
	handlers   map[string]WsHandler
	errHandler ErrHandler
}

// NewMarketService creates a new market service
func NewMarketService(config *WsConfig) *MarketService {
	if config.Endpoint == "" {
		config.Endpoint = PublicWebSocketURL
	}
	return &MarketService{
		config:   config,
		subs:     make(map[string]struct{}),
		handlers: make(map[string]WsHandler),
	}
}

// Connect connects to the websocket server
func (s *MarketService) Connect() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.client != nil {
		return fmt.Errorf("already connected")
	}

	client := NewWsClient(s.config)
	client.SetHandler(s.handleMessage)
	client.SetErrHandler(s.handleError)

	err := client.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	s.client = client
	return nil
}

// Close closes the websocket connection
func (s *MarketService) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.client == nil {
		return fmt.Errorf("not connected")
	}

	s.client.Close()
	s.client = nil
	s.subs = make(map[string]struct{})
	return nil
}

// handleMessage handles incoming messages
func (s *MarketService) handleMessage(message []byte) {
	var msg struct {
		Event string          `json:"event"`
		Arg   ChannelArgs     `json:"arg"`
		Data  json.RawMessage `json:"data"`
		Code  string          `json:"code"`
		Msg   string          `json:"msg"`
	}

	err := json.Unmarshal(message, &msg)
	if err != nil {
		if s.errHandler != nil {
			s.errHandler(fmt.Errorf("failed to unmarshal message: %w", err))
		}
		return
	}

	// Handle subscription response
	if msg.Event == "subscribe" || msg.Event == "unsubscribe" {
		// Skip subscription messages as they are handled by the client
		return
	}

	// Handle data message
	key := fmt.Sprintf("%s:%s", msg.Arg.Channel, msg.Arg.InstId)
	s.mu.RLock()
	handler, ok := s.handlers[key]
	s.mu.RUnlock()

	if ok && handler != nil {
		handler(message)
	}
}

// handleError handles errors
func (s *MarketService) handleError(err error) {
	if s.errHandler != nil {
		s.errHandler(err)
	}
}

// SetErrHandler sets the error handler
func (s *MarketService) SetErrHandler(errHandler ErrHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.errHandler = errHandler
}

// Subscribe subscribes to a channel
func (s *MarketService) Subscribe(channel string, instId string, handler WsHandler) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.client == nil {
		return fmt.Errorf("not connected")
	}

	key := fmt.Sprintf("%s:%s", channel, instId)
	if _, ok := s.subs[key]; ok {
		return fmt.Errorf("already subscribed to %s", key)
	}

	// Check subscription length
	req := SubscribeRequest{
		Op: OpSubscribe,
		Args: []ChannelArgs{
			{
				Channel: channel,
				InstId:  instId,
			},
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal subscribe request: %w", err)
	}

	if len(data) > MaxSubscriptionLength {
		return fmt.Errorf("subscription length exceeds maximum allowed size of %d bytes", MaxSubscriptionLength)
	}

	if err := s.client.Subscribe(channel, instId, handler); err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	s.subs[key] = struct{}{}
	s.handlers[key] = handler
	return nil
}

// Unsubscribe unsubscribes from a channel
func (s *MarketService) Unsubscribe(channel string, instId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.client == nil {
		return fmt.Errorf("not connected")
	}

	key := fmt.Sprintf("%s:%s", channel, instId)
	if _, ok := s.subs[key]; !ok {
		return fmt.Errorf("not subscribed to %s", key)
	}

	if err := s.client.Unsubscribe(channel, instId); err != nil {
		return fmt.Errorf("failed to unsubscribe: %w", err)
	}

	delete(s.subs, key)
	delete(s.handlers, key)
	return nil
}

// SubscribeTrades subscribes to trades channel
func (s *MarketService) SubscribeTrades(instId string, handler WsHandler) error {
	return s.Subscribe(ChannelTrades, instId, handler)
}

// SubscribeCandles subscribes to candles channel
func (s *MarketService) SubscribeCandles(instId string, interval string, handler WsHandler) error {
	return s.Subscribe(interval, instId, handler)
}

// SubscribeOrderbook subscribes to orderbook channel
func (s *MarketService) SubscribeOrderbook(instId string, handler WsHandler) error {
	return s.Subscribe(ChannelOrderbook, instId, handler)
}

// SubscribeTickers subscribes to tickers channel
func (s *MarketService) SubscribeTickers(instId string, handler WsHandler) error {
	return s.Subscribe(ChannelTickers, instId, handler)
}

// SubscribeFundingRate subscribes to funding rate channel
func (s *MarketService) SubscribeFundingRate(instId string, handler WsHandler) error {
	return s.Subscribe(ChannelFundingRate, instId, handler)
}

// SubscribeMulti subscribes to a channel for multiple instruments
func (s *MarketService) SubscribeMulti(channel string, instIds []string, handler WsHandler) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.client == nil {
		return fmt.Errorf("not connected")
	}

	// Check if any of the instruments are already subscribed
	for _, instId := range instIds {
		key := fmt.Sprintf("%s:%s", channel, instId)
		if _, ok := s.subs[key]; ok {
			return fmt.Errorf("already subscribed to %s", key)
		}
	}

	// Create subscription request
	args := make([]ChannelArgs, len(instIds))
	for i, instId := range instIds {
		args[i] = ChannelArgs{
			Channel: channel,
			InstId:  instId,
		}
	}

	req := SubscribeRequest{
		Op:   OpSubscribe,
		Args: args,
	}

	// Check subscription length
	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal subscribe request: %w", err)
	}

	if len(data) > MaxSubscriptionLength {
		return fmt.Errorf("subscription length exceeds maximum allowed size of %d bytes", MaxSubscriptionLength)
	}

	// Send subscription request
	if err := s.client.conn.WriteJSON(req); err != nil {
		return fmt.Errorf("failed to send subscribe request: %w", err)
	}

	// Add subscriptions and register handlers
	for _, instId := range instIds {
		key := fmt.Sprintf("%s:%s", channel, instId)
		s.subs[key] = struct{}{}
		s.handlers[key] = handler
		s.client.subs[key] = handler
	}

	return nil
}

// SubscribeTradesMulti subscribes to trades channel for multiple instruments
func (s *MarketService) SubscribeTradesMulti(instIds []string, handler WsHandler) error {
	return s.SubscribeMulti(ChannelTrades, instIds, handler)
}

// SubscribeCandlesMulti subscribes to candles channel for multiple instruments
func (s *MarketService) SubscribeCandlesMulti(instIds []string, interval string, handler WsHandler) error {
	return s.SubscribeMulti(interval, instIds, handler)
}

// SubscribeOrderbookMulti subscribes to orderbook channel for multiple instruments
func (s *MarketService) SubscribeOrderbookMulti(instIds []string, handler WsHandler) error {
	return s.SubscribeMulti(ChannelOrderbook, instIds, handler)
}

// SubscribeTickersMulti subscribes to tickers channel for multiple instruments
func (s *MarketService) SubscribeTickersMulti(instIds []string, handler WsHandler) error {
	return s.SubscribeMulti(ChannelTickers, instIds, handler)
}

// SubscribeFundingRateMulti subscribes to funding rate channel for multiple instruments
func (s *MarketService) SubscribeFundingRateMulti(instIds []string, handler WsHandler) error {
	return s.SubscribeMulti(ChannelFundingRate, instIds, handler)
}
