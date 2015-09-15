package ws

import (
	"errors"
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

var conns = struct {
	sync.RWMutex
	m map[id]*websocket.Conn
}{
	m: make(map[id]*websocket.Conn),
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

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// NOTE: Returning a HTTP error is done by upgrader.
		return
	}

	conns.Lock()
	conns.m[ktoi(key)] = conn
	conns.Unlock()

	reader(conn)
}

// Write picks the right WebSocket to communicate with the user owning
// the entity references by the passed key and writes data to the socket.
func Write(key *datastore.Key, buf []byte) error {
	conn := lookup(key)

	if conn == nil {
		return errors.New("ws: could not find associated WebSocket")
	}

	conn.SetWriteDeadline(time.Now().Add(writeWait))
	return conn.WriteMessage(websocket.TextMessage, buf)
}

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func lookup(key *datastore.Key) *websocket.Conn {
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
