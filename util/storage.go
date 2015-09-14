package util

import "google.golang.org/appengine"

var FileNames = map[string]string{
	"py":   "app.py",
	"c":    "app.c",
	"cpp":  "app.cpp",
	"java": "Application.java",
}

const TemplateBucket = "coduno-templates"

func SubmissionBucket() string {
	if appengine.IsDevAppServer() {
		return "coduno-dev"
	}
	return "coduno"
}
