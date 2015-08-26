package runner

import (
	"net/http"
	"net/url"

	"github.com/coduno/api/model"
	"github.com/fsouza/go-dockerclient"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

var compute *url.URL

var dc *docker.Client

type Runner int

const (
	SimpleRunner Runner = iota + 1
	DiffRunner
	JunitRunner
)

type RunnerFunc func(context.Context, *model.Test, model.KeyedSubmission) error

var Runners = map[Runner]RunnerFunc{
	SimpleRunner: simpleRunner,
	DiffRunner:   diffRunner,
	JunitRunner:  junitRunner,
}

func init() {
	var err error
	if appengine.IsDevAppServer() {
		dc, err = docker.NewClientFromEnv()
		if err != nil {
			panic(err)
		}
	}
}

func HandleCodeSubmission(ctx context.Context, w http.ResponseWriter, r *http.Request, resultKey, taskKey *datastore.Key) (int, error) {
	return http.StatusOK, nil
}

// newImage returns the correct docker image name for a
// specific language.
func newImage(language string) string {
	const imagePattern string = "coduno/fingerprint-"
	return imagePattern + language
}
