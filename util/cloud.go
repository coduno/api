package util

import (
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/cloud"
)

var cloudClient *http.Client

func init() {
	var err error
	cloudClient, err = google.DefaultClient(context.Background())
	if err != nil {
		panic(err)
	}
}

const projID = "coduno"

func CloudContext(parent context.Context) context.Context {
	if parent == nil {
		return cloud.NewContext(projID, cloudClient)
	}
	return cloud.WithContext(parent, projID, cloudClient)
}
