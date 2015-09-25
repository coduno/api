package controllers

import (
	"net/http"

	"google.golang.org/cloud/compute/metadata"
)

func init() {
	http.HandleFunc("/ip", hsts(ip))
}

func ip(w http.ResponseWriter, r *http.Request) {
	ip := "127.0.0.1"
	if metadata.OnGCE() {
		var err error
		ip, err = metadata.InternalIP()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	w.Write([]byte(ip))
}
