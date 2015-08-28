package runner

import (
	"net/url"

	"github.com/fsouza/go-dockerclient"
	"google.golang.org/appengine"
)

var compute *url.URL

var dc *docker.Client

func init() {
	var err error
	if appengine.IsDevAppServer() {
		dc, err = docker.NewClientFromEnv()
		if err != nil {
			panic(err)
		}
	}
}

// newImage returns the correct docker image name for a
// specific language.
func newImage(language string) string {
	const imagePattern string = "coduno/fingerprint-"
	return imagePattern + language
}
