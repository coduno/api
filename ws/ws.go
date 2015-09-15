package ws

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write messages to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{HandshakeTimeout: 10 * time.Second}

type conn struct {
	ws   *websocket.Conn
	send chan []byte
}

func newConn(ws *websocket.Conn) *conn {
	return &conn{
		ws:   ws,
		send: make(chan []byte),
	}
}

func (c *conn) writeMessage(data []byte) {
	c.send <- data
}

func (c *conn) handlePong(_ string) error {
	return c.ws.SetReadDeadline(time.Now().Add(pongWait))
}

func (c *conn) readLoop() {
	defer c.ws.Close()
	c.ws.SetReadLimit(512)
	c.handlePong("")
	c.ws.SetPongHandler(c.handlePong)
	for {
		_, _, err := c.ws.ReadMessage()
		if err != nil {
			// Ignore all incoming messages.
			return
		}
	}
}

func (c *conn) write(messageType int, data []byte) error {
	if err := c.ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return err
	}
	return c.ws.WriteMessage(messageType, data)
}

func (c *conn) writeLoop() {
	defer c.ws.Close()

	ticker := time.NewTicker(pingPeriod)
	// Make sure ticker is stopped if we end writing for some reason.
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				m := []byte("Upstream channel was closed.")
				c.write(websocket.CloseMessage, m)
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case t := <-ticker.C:
			m := []byte(fmt.Sprintf("Tick %s", t.Format(time.RFC3339)))
			if err := c.write(websocket.PingMessage, m); err != nil {
				return
			}
		}
	}
}

var conns = struct {
	sync.RWMutex
	m map[id]*conn
}{
	m: make(map[id]*conn),
}

type id struct {
	kind      string
	stringID  string
	intID     int64
	appID     string
	namespace string
}

func init() {
	if appengine.IsDevAppServer() {
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
		return
	}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		origin, ok := r.Header["Origin"]
		if !ok {
			return false
		}
		return origin[0] == "https://app.cod.uno"
	}
}

// Handle performs an upgrade from HTTP to WebSocket. It hijacks the connection
// and registers it in the global pool.
func Handle(w http.ResponseWriter, r *http.Request) {
	m := r.URL.Query()

	result := m.Get("result")
	if result == "" {
		http.Error(w, "Bad Request: missing parameter \"result\"", http.StatusBadRequest)
		return
	}

	key, err := datastore.DecodeKey(result)
	if err != nil {
		http.Error(w, "Bad Request: cannot decode \"result\"", http.StatusBadRequest)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// NOTE: Returning a HTTP error is done by upgrader.
		return
	}

	conn := newConn(ws)

	conns.Lock()
	conns.m[ktoi(key)] = conn
	conns.Unlock()

	go conn.writeLoop()
	conn.readLoop()
}

// Write picks the right WebSocket to communicate with the user owning
// the entity references by the passed key and writes data to the socket.
func Write(key *datastore.Key, buf []byte) error {
	conn := lookup(key)

	if conn == nil {
		return errors.New("ws: could not find associated WebSocket")
	}

	conn.writeMessage(buf)
	return nil
}

func lookup(key *datastore.Key) *conn {
	conns.RLock()
	defer conns.RUnlock()

	for key != nil {
		conn, ok := conns.m[ktoi(key)]
		if ok {
			return conn
		}
		key = key.Parent()
	}

	return nil
}

func ktoi(key *datastore.Key) id {
	return id{
		kind:      key.Kind(),
		stringID:  key.StringID(),
		intID:     key.IntID(),
		appID:     key.AppID(),
		namespace: key.Namespace(),
	}
}
