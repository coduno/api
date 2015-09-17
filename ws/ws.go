package ws

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write messages to the client.
	// If a write takes longer than this duration, the connection
	// is declared dead and removed from the global pool.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	// If there is no pong received after this duration, the connection
	// is declared dead and removed from the global pool.
	pongWait = 30 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// ErrNoSocket is returned if there could be no WebSocket found for some
// key.
var ErrNoSocket = errors.New("ws: could not find associated WebSocket")

var upgrader = websocket.Upgrader{HandshakeTimeout: 10 * time.Second}

var conns = struct {
	sync.RWMutex
	m map[id]*websocket.Conn
}{
	m: make(map[id]*websocket.Conn),
}

// ID is a a shallow copy of datastore.Key. It has no references to
// parent keys.
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
		// TODO(flowlo): Remove magic constant.
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
		// NOTE: Returning an HTTP error is done by upgrader.
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
			die(id, err)
			return
		}
	}
}

// ping will loop infinitely and send ping messages over ws. If a write takes longer than
// writeWait, it will remove ws from the connection pool.
func ping(ws *websocket.Conn, id id) {
	for t := range time.Tick(pingPeriod) {
		m := []byte(fmt.Sprint("Tick", t))
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		if err := ws.WriteMessage(websocket.PingMessage, m); err != nil {
			die(id, err)
			return
		}
	}
}

// lookup is a shorthand for synchronized read access on conns.m.
func lookup(id id) *websocket.Conn {
	conns.RLock()
	ws := conns.m[id]
	conns.RUnlock()
	return ws
}

// die removes the WebSocket references by id from the global pool and logs
// that it did so.
func die(id id, err error) {
	log.Debugf(appengine.BackgroundContext(), "ws: connection for %v died: %s", id, err)
	conns.Lock()
	delete(conns.m, id)
	conns.Unlock()
}

// ktoi makes a shallow copy of key. Links to parent keys are not followed.
func ktoi(key *datastore.Key) id {
	return id{
		kind:      key.Kind(),
		stringID:  key.StringID(),
		intID:     key.IntID(),
		appID:     key.AppID(),
		namespace: key.Namespace(),
	}
}

// Write picks the right WebSocket to communicate with the user owning
// the entity references by the passed key and writes data to the socket.
func Write(key *datastore.Key, buf []byte) error {
	id := ktoi(key)
	ws := lookup(id)

	if ws == nil {
		return ErrNoSocket
	}

	if err := ws.SetReadDeadline(time.Now().Add(writeWait)); err != nil {
		die(id, err)
		return err
	}

	if err := ws.WriteMessage(websocket.TextMessage, buf); err != nil {
		die(id, err)
		return err
	}

	return nil
}

// Close closes the WebSocket associated with the given key. It does not take care of
// messages currently being sent, and just roughly closes the underlying connection.
func Close(key *datastore.Key) error {
	id := ktoi(key)
	ws := lookup(id)

	if ws == nil {
		return ErrNoSocket
	}

	if err := ws.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
		die(id, err)
		return err
	}

	return nil
}
