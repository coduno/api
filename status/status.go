package status

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/coduno/app/util"

	"google.golang.org/appengine"
)

type Status struct {
	Init      time.Time
	Appengine Appengine
	Runtime   Runtime
	Environ   []string
}

type Appengine struct {
	InstanceID,
	AppID,
	Datacenter,
	DefaultVersionHostname,
	ModuleName string
	IsDevAppServer bool
}

type Runtime struct {
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
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if !util.CheckMethod(w, r, "GET") {
		return
	}

	ctx := appengine.NewContext(r)

	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)

	s := &Status{
		initTime,
		Appengine{
			appengine.InstanceID(),
			appengine.AppID(ctx),
			appengine.Datacenter(ctx),
			appengine.DefaultVersionHostname(ctx),
			appengine.ModuleName(ctx),
			appengine.IsDevAppServer(),
		},
		Runtime{
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

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
