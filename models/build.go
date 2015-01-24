package models

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/appengine/urlfetch"
)

type Build struct {
	Repository string `json:"Repository"`
	Hash       string `json:"Hash"`
	Container  string
}

// Inform Docker of the creation of this build, pass on everything
// needed to start execution.x
func (build Build) Create(c context.Context) error {
	r, _ := http.NewRequest("POST", "http://docker.cod.uno:2375/v1.15/containers/create", strings.NewReader(`
		{
			"Image": "coduno/git:experimental",
			"Cmd": ["/start.sh"],
			"Env": [
				"CODUNO_REPOSITORY_URL=`+build.Repository+`",
				"CODUNO_HASH=`+build.Hash+`"
			]
		}
	`))

	r.Header = map[string][]string{
		"Content-Type": {"application/json"},
		"Accept":       {"application/json"},
	}

	client := urlfetch.Client(c)
	res, err := client.Do(r)

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)

	return json.Unmarshal(body, &build.Container)
}

// Once a container was created for a Build, it has to be started.
func (build Build) Start(c context.Context) error {
	r, _ := http.NewRequest("POST", "http://docker.cod.uno:2375/v1.15/containers/"+build.Container+"/start", nil)

	r.Header = map[string][]string{
		"Content-Type": {"application/json"},
		"Accept":       {"application/json"},
	}

	client := urlfetch.Client(c)
	_, err := client.Do(r)
	return err
	// TODO(flowlo): Interpret response, did we actually start the container?
}
