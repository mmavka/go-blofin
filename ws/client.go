package ws

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Connection settings
	pingInterval = 20 * time.Second
	pongTimeout  = 10 * time.Second
)

// WsHandler handle raw websocket message
type WsHandler func(message []byte)

// ErrHandler handles errors
type ErrHandler func(err error)

// WsConfig represents WebSocket configuration
type WsConfig struct {
	Endpoint string
	APIKey   string
	Secret   string
}

// WsClient represents WebSocket client
type WsClient struct {
	config          *WsConfig
	conn            *websocket.Conn
	mu              sync.Mutex
	subs            map[string]func([]byte)
	stopCh          chan struct{}
	shouldReconnect bool
	lastPong        time.Time
	pingTicker      *time.Ticker
	handler         WsHandler
	errHandler      ErrHandler
}

// NewWsClient creates a new WebSocket client
func NewWsClient(config *WsConfig) *WsClient {
	return &WsClient{
		config: config,
		subs:   make(map[string]func([]byte)),
		stopCh: make(chan struct{}),
	}
}

// SetHandler sets the message handler
func (c *WsClient) SetHandler(handler WsHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handler = handler
}

// SetErrHandler sets the error handler
func (c *WsClient) SetErrHandler(errHandler ErrHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.errHandler = errHandler
}

// Connect establishes WebSocket connection
func (c *WsClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return nil
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, _, err := dialer.Dial(c.config.Endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	c.conn = conn
	c.shouldReconnect = true
	c.lastPong = time.Now()

	// Start ping/pong
	c.pingTicker = time.NewTicker(PingInterval * time.Second)
	go c.pingPong()

	// Start message handler
	go c.handleMessages()

	return nil
}

// Close closes WebSocket connection
func (c *WsClient) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		c.shouldReconnect = false
		c.pingTicker.Stop()
		c.conn.Close()
		c.conn = nil
	}
}

// Subscribe subscribes to a channel
func (c *WsClient) Subscribe(channel string, instId string, handler func([]byte)) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return fmt.Errorf("not connected")
	}

	key := fmt.Sprintf("%s:%s", channel, instId)
	if _, ok := c.subs[key]; ok {
		return fmt.Errorf("already subscribed to %s", key)
	}

	req := SubscribeRequest{
		Op: OpSubscribe,
		Args: []ChannelArgs{
			{
				Channel: channel,
				InstId:  instId,
			},
		},
	}

	// Check subscription length
	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal subscribe request: %v", err)
	}

	if len(data) > MaxSubscriptionLength {
		return fmt.Errorf("subscription length exceeds maximum allowed size of %d bytes", MaxSubscriptionLength)
	}

	if err := c.conn.WriteJSON(req); err != nil {
		return fmt.Errorf("failed to send subscribe request: %v", err)
	}

	c.subs[key] = handler
	return nil
}

// Unsubscribe unsubscribes from a channel
func (c *WsClient) Unsubscribe(channel string, instId string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return fmt.Errorf("not connected")
	}

	key := fmt.Sprintf("%s:%s", channel, instId)
	if _, ok := c.subs[key]; !ok {
		return fmt.Errorf("not subscribed to %s", key)
	}

	req := SubscribeRequest{
		Op: OpUnsubscribe,
		Args: []ChannelArgs{
			{
				Channel: channel,
				InstId:  instId,
			},
		},
	}

	if err := c.conn.WriteJSON(req); err != nil {
		return fmt.Errorf("failed to send unsubscribe request: %v", err)
	}

	delete(c.subs, key)
	return nil
}

// handleMessages handles incoming WebSocket messages
func (c *WsClient) handleMessages() {
	for {
		select {
		case <-c.stopCh:
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if c.errHandler != nil {
					c.errHandler(fmt.Errorf("read message error: %v", err))
				}
				if c.shouldReconnect {
					c.doReconnect()
				}
				return
			}

			// Update last pong time
			c.mu.Lock()
			c.lastPong = time.Now()
			c.mu.Unlock()

			// Handle pong message
			if string(message) == OpPong {
				continue
			}

			// Parse message
			var msg map[string]interface{}
			if err := json.Unmarshal(message, &msg); err != nil {
				if c.errHandler != nil {
					c.errHandler(fmt.Errorf("unmarshal error: %v", err))
				}
				continue
			}

			// Handle error message
			if event, ok := msg["event"].(string); ok && event == "error" {
				if code, ok := msg["code"].(string); ok {
					if c.errHandler != nil {
						c.errHandler(fmt.Errorf("websocket error: %s", code))
					}
					switch code {
					case ErrCodeRateLimit:
						// Wait for rate limit to reset
						time.Sleep(5 * time.Minute)
					case ErrCodeForbidden:
						// Wait for firewall restriction to reset
						time.Sleep(5 * time.Minute)
					}
				}
				continue
			}

			// Handle subscription response
			if event, ok := msg["event"].(string); ok && (event == "subscribe" || event == "unsubscribe") {
				if arg, ok := msg["arg"].(map[string]interface{}); ok {
					channel, _ := arg["channel"].(string)
					instId, _ := arg["instId"].(string)
					if c.handler != nil {
						subscriptionMsg := map[string]interface{}{
							"event": event,
							"arg": map[string]string{
								"channel": channel,
								"instId":  instId,
							},
						}
						if data, err := json.Marshal(subscriptionMsg); err == nil {
							c.handler(data)
						}
					}
				}
				continue
			}

			// Handle data message
			if arg, ok := msg["arg"].(map[string]interface{}); ok {
				channel, _ := arg["channel"].(string)
				instId, _ := arg["instId"].(string)
				key := fmt.Sprintf("%s:%s", channel, instId)

				c.mu.Lock()
				if handler, ok := c.subs[key]; ok {
					handler(message)
				} else if c.handler != nil {
					c.handler(message)
				} else {
					if c.errHandler != nil {
						c.errHandler(fmt.Errorf("no handler for %s", key))
					}
				}
				c.mu.Unlock()
			}
		}
	}
}

// pingPong sends ping messages and checks pong responses
func (c *WsClient) pingPong() {
	for {
		select {
		case <-c.stopCh:
			return
		case <-c.pingTicker.C:
			c.mu.Lock()
			if c.conn == nil {
				c.mu.Unlock()
				return
			}

			// Check if pong was received
			if time.Since(c.lastPong) > PongTimeout*time.Second {
				c.mu.Unlock()
				if c.shouldReconnect {
					c.doReconnect()
				}
				return
			}

			// Send ping
			if err := c.conn.WriteMessage(websocket.TextMessage, []byte(OpPing)); err != nil {
				c.mu.Unlock()
				if c.shouldReconnect {
					c.doReconnect()
				}
				return
			}
			c.mu.Unlock()
		}
	}
}

// doReconnect attempts to reconnect to WebSocket server
func (c *WsClient) doReconnect() {
	for i := 0; i < MaxReconnects; i++ {
		time.Sleep(ReconnectDelay * time.Second)

		c.mu.Lock()
		if !c.shouldReconnect {
			c.mu.Unlock()
			return
		}

		if c.conn != nil {
			c.conn.Close()
			c.conn = nil
		}

		dialer := websocket.Dialer{
			HandshakeTimeout: 10 * time.Second,
		}

		conn, _, err := dialer.Dial(c.config.Endpoint, nil)
		if err != nil {
			c.mu.Unlock()
			continue
		}

		c.conn = conn
		c.lastPong = time.Now()

		// Resubscribe to all channels
		for key := range c.subs {
			channel, instId := parseKey(key)
			req := SubscribeRequest{
				Op: OpSubscribe,
				Args: []ChannelArgs{
					{
						Channel: channel,
						InstId:  instId,
					},
				},
			}

			if err := c.conn.WriteJSON(req); err != nil {
				c.mu.Unlock()
				continue
			}
		}

		c.mu.Unlock()
		return
	}
}

// parseKey parses channel key into channel and instrument ID
func parseKey(key string) (string, string) {
	for i := 0; i < len(key); i++ {
		if key[i] == ':' {
			return key[:i], key[i+1:]
		}
	}
	return "", ""
}
