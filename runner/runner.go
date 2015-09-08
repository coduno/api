package runner

import (
	"net/url"
	"os"

	"golang.org/x/net/context"

	"github.com/fsouza/go-dockerclient"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

var compute *url.URL

var dc *docker.Client

func init() {
	var err error
	if appengine.IsDevAppServer() {
		dc, err = docker.NewClientFromEnv()
	} else {
		// FIXME(flowlo)
		dc, err = docker.NewClient("tcp://10.240.10.141:2375")
	}
	if err != nil {
		panic(err)
	}
}

func prepareImage(ctx context.Context, name string) error {
	_, err := dc.InspectImage(name)

	if err != nil {
		return nil
	}

	if err != docker.ErrNoSuchImage {
		return err
	}

	log.Warningf(ctx, "Missing image %s will be pulled. Expect severe delay!", name)

	err = dc.PullImage(docker.PullImageOptions{
		Repository:   name,
		OutputStream: os.Stderr,
	}, docker.AuthConfiguration{})

	if err != nil {
		log.Warningf(ctx, "Failed pulling image %s because of: %s", name, err)
	}

	return err
}

// newImage returns the correct docker image name for a
// specific language.
func newImage(language string) string {
	const imagePattern string = "coduno/fingerprint-"
	return imagePattern + language
}

func createDockerVolume(name string) (*docker.Volume, error) {
	return dc.CreateVolume(docker.CreateVolumeOptions{
		Driver: "gcs",
		Name:   name,
	})
}
