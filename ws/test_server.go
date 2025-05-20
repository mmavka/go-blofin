package ws

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// TestServer is a test websocket server
type TestServer struct {
	t        *testing.T
	server   *httptest.Server
	conn     *websocket.Conn
	messages chan []byte
}

// newTestServer creates a new test server
func newTestServer(t *testing.T) *TestServer {
	ts := &TestServer{
		t:        t,
		messages: make(chan []byte, 100),
	}

	ts.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatal(err)
		}
		ts.conn = conn

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return
			}
			ts.messages <- message
		}
	}))

	return ts
}

// URL returns the server URL
func (ts *TestServer) URL() string {
	return "ws" + ts.server.URL[4:]
}

// Close closes the server
func (ts *TestServer) Close() {
	if ts.conn != nil {
		ts.conn.Close()
	}
	ts.server.Close()
}
