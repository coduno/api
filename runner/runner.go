package runner

import (
	"bytes"
	"errors"
	"net/url"
	"os"
	"time"

	"github.com/fsouza/go-dockerclient"
	"google.golang.org/appengine"
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

func createDockerContainer(image string, binds []string) (*docker.Container, error) {
	// TODO(victorbalan): Pass the memory limit from test
	return dc.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: image,
		},
		HostConfig: &docker.HostConfig{
			Privileged: false,
			Memory:     0, // TODO(flowlo): Limit memory
			Binds:      binds,
		},
	})
}

func waitForContainer(cID string) (err error) {
	waitc := make(chan waitResult)
	go func() {
		exitCode, err := dc.WaitContainer(cID)
		waitc <- waitResult{exitCode, err}
	}()

	var res waitResult
	select {
	case res = <-waitc:
	case <-time.After(time.Minute):
		err = errors.New("execution timed out")
		return
	}

	return res.Err
}

func getLogs(cID string) (stdout, stderr *bytes.Buffer, err error) {
	stdout = new(bytes.Buffer)
	stderr = new(bytes.Buffer)
	err = dc.Logs(docker.LogsOptions{
		Container:    cID,
		OutputStream: stdout,
		ErrorStream:  stderr,
		Stdout:       true,
		Stderr:       true,
	})
	return
}

func prepareImage(name string) (err error) {
	if _, err = dc.InspectImage(name); err == nil {
		return nil
	}

	if err != docker.ErrNoSuchImage {
		return
	}

	err = dc.PullImage(docker.PullImageOptions{
		Repository:   name,
		OutputStream: os.Stderr,
	}, docker.AuthConfiguration{})
	return
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
