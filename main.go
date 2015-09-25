package main

import (
	"net/http"

	"github.com/coduno/api/controllers"
	"github.com/coduno/api/ws"

	"google.golang.org/appengine"
)

func init() {
	controllers.InvitationTemplatePath = "./mail/template.invitation"
	controllers.SubTemplatePath = "./mail/template.subscription"
}

func main() {
	go http.ListenAndServe(":8090", http.HandlerFunc(ws.Handle))
	http.Handle("/", controllers.Handler())
	appengine.Main()
}
