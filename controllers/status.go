package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"
	"time"

	"google.golang.org/appengine"
)

type status struct {
	Init      time.Time
	Appengine appengineStatus
	Runtime   runtimeStatus
	Environ   []string
}

type appengineStatus struct {
	InstanceID,
	AppID,
	Datacenter,
	DefaultVersionHostname,
	ModuleName string
	IsDevAppServer bool
}

type runtimeStatus struct {
	GOMAXPROCS int
	GOARCH,
	GOOS,
	GOROOT string
	NumCPU     int
	NumCgoCall int64
	Version    string
	MemStats   runtime.MemStats
}

var initTime time.Time

func init() {
	initTime = time.Now()
	router.HandleFunc("/status", Status)
}

// Status gathers a quick overview of the system state
// and dumps it in JSON format.
func Status(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	ctx := appengine.NewContext(r)

	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)

	s := &status{
		initTime,
		appengineStatus{
			appengine.InstanceID(),
			appengine.AppID(ctx),
			appengine.Datacenter(ctx),
			appengine.DefaultVersionHostname(ctx),
			appengine.ModuleName(ctx),
			appengine.IsDevAppServer(),
		},
		runtimeStatus{
			runtime.GOMAXPROCS(0),
			runtime.GOARCH,
			runtime.GOOS,
			runtime.GOROOT(),
			runtime.NumCPU(),
			runtime.NumCgoCall(),
			runtime.Version(),
			*m,
		},
		os.Environ(),
	}

	b, err := json.Marshal(s)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; encoding=utf-8")
	w.Write(b)
}
