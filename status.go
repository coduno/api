package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/coduno/app/util"

	"google.golang.org/appengine"
)

type Status struct {
	InstanceID,
	AppID,
	Datacenter,
	DefaultVersionHostname,
	ModuleName string
	IsDevAppServer bool
	Init           time.Time
}

var initTime time.Time

func init() {
	initTime = time.Now()
}

func status(w http.ResponseWriter, r *http.Request) {
	if !util.CheckMethod(w, r, "GET") {
		return
	}

	ctx := appengine.NewContext(r)

	s := &Status{
		appengine.InstanceID(),
		appengine.AppID(ctx),
		appengine.Datacenter(ctx),
		appengine.DefaultVersionHostname(ctx),
		appengine.ModuleName(ctx),
		appengine.IsDevAppServer(),
		initTime,
	}

	b, err := json.Marshal(s)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
