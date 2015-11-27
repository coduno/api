package main

import (
	"net/http"

	"github.com/coduno/api/controllers"

	"google.golang.org/appengine"
)

func main() {
	http.Handle("/", controllers.Handler())
	appengine.Main()
}
