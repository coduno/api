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

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// NOTE: Returning a HTTP error is done by upgrader.
		return
	}

	id := ktoi(key)

	conns.Lock()
	conns.m[ktoi(key)] = ws
	conns.Unlock()

	go ping(ws, id)

	h := func(_ string) error {
		return ws.SetReadDeadline(time.Now().Add(pongWait))
	}

	ws.SetReadLimit(512)
	h("")
	ws.SetPongHandler(h)
	for {
		// Ignore all incoming messages.
		_, _, err := ws.ReadMessage()
		if err != nil {
			die(id)
			return
		}
	}
}

func ping(ws *websocket.Conn, id id) {
	for t := range time.Tick(pingPeriod) {
		m := []byte(fmt.Sprint("Tick", t))
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		if err := ws.WriteMessage(websocket.PingMessage, m); err != nil {
			die(id)
			return
		}
	}
}

// Write picks the right WebSocket to communicate with the user owning
// the entity references by the passed key and writes data to the socket.
func Write(key *datastore.Key, buf []byte) error {
	id := ktoi(key)
	ws := lookup(id)

	if ws == nil {
		return errors.New("ws: could not find associated WebSocket")
	}

	if err := ws.SetReadDeadline(time.Now().Add(writeWait)); err != nil {
		die(id)
		return err
	}

	if err := ws.WriteMessage(websocket.TextMessage, buf); err != nil {
		die(id)
		return err
	}

	return nil
}

func Close(key *datastore.Key) error {
	id := ktoi(key)
	ws := lookup(id)

	if ws == nil {
		return errors.New("ws: could not find associated WebSocket")
	}

	die(id)

	return ws.WriteMessage(websocket.CloseMessage, []byte{})
}

func lookup(id id) *websocket.Conn {
	conns.RLock()
	ws := conns.m[id]
	conns.RUnlock()
	return ws
}

func die(id id) {
	conns.Lock()
	delete(conns.m, id)
	conns.Unlock()
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
