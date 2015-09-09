package util

import (
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine"
	"google.golang.org/cloud"
)

var FileNames = map[string]string{
	"py":   "app.py",
	"c":    "app.c",
	"cpp":  "app.cpp",
	"java": "Application.java",
}

var cloudClient *http.Client

const projID = "coduno"
const TemplateBucket = "coduno-templates"

func init() {
	var err error
	cloudClient, err = google.DefaultClient(context.Background())
	if err != nil {
		panic(err)
	}
}

func CloudContext(parent context.Context) context.Context {
	if parent == nil {
		return cloud.NewContext(projID, cloudClient)
	}
	return cloud.WithContext(parent, projID, cloudClient)
}

func SubmissionBucket() string {
	if appengine.IsDevAppServer() {
		return "coduno-dev"
	}
	return "coduno"
}
