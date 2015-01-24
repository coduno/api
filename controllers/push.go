package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"

	"github.com/coduno/app/gitlab"
)

type ContainerCreation struct {
	Id string `json:"Id"`
}

func Push(w http.ResponseWriter, req *http.Request, c context.Context) {
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.Warningf(c, err.Error())
	} else if len(body) < 1 {
		log.Warningf(c, "Received empty body.")
	} else {
		push, err := gitlab.NewPush(body)

		if err != nil {
			log.Warningf(c, err.Error())
		} else {
			commit := push.Commits[0]

			docker, _ := http.NewRequest("POST", "http://docker.cod.uno:2375/v1.15/containers/create", strings.NewReader(`
				{
					"Image": "coduno/git:experimental",
					"Cmd": ["/start.sh"],
					"Env": [
						"CODUNO_REPOSITORY_NAME=`+push.Repository.Name+`",
						"CODUNO_REPOSITORY_URL=`+push.Repository.URL+`",
						"CODUNO_REPOSITORY_HOMEPAGE=`+push.Repository.Homepage+`",
						"CODUNO_REF=`+push.Ref+`"
					]
				}
			`))

			docker.Header = map[string][]string{
				"Content-Type": {"application/json"},
				"Accept":       {"application/json"},
			}

			client := urlfetch.Client(c)
			res, err := client.Do(docker)

			if err != nil {
				log.Debugf(c, "Docker API response:", err.Error())
				return
			}

			var result ContainerCreation
			body, err := ioutil.ReadAll(res.Body)
			err = json.Unmarshal(body, &result)

			if err != nil {
				log.Debugf(c, "Received body %d: %s", res.StatusCode, string(body))
				log.Debugf(c, "Unmarshalling API response: %s", err.Error())
				return
			}

			docker, _ = http.NewRequest("POST", "http://docker.cod.uno:2375/v1.15/containers/"+result.Id+"/start", nil)

			docker.Header = map[string][]string{
				"Content-Type": {"application/json"},
				"Accept":       {"application/json"},
			}

			res, err = client.Do(docker)

			if err != nil {
				log.Debugf(c, "Docker API response 2: %s", err.Error())
			} else {
				result, err := ioutil.ReadAll(res.Body)
				if err != nil {
					log.Debugf(c, err.Error())
				} else {
					log.Debugf(c, string(result))
				}
			}
			defer res.Body.Close()

			log.Infof(c, "Received push from %s", commit.Author.Email)
		}
	}
}
