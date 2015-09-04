package runner

import (
	"bytes"
	"errors"
	"os"
	"path"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine/log"

	"github.com/coduno/api/model"
	"github.com/fsouza/go-dockerclient"
)

type waitResult struct {
	ExitCode int
	Err      error
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

func Simple(ctx context.Context, sub model.KeyedSubmission) (stdout, stderr *bytes.Buffer, err error) {
	image := newImage(sub.Language)

	if err = prepareImage(ctx, image); err != nil {
		return nil, nil, err
	}

	v, err := dc.CreateVolume(docker.CreateVolumeOptions{
		Driver: "gcs",
		Name:   sub.Code.Bucket + "/" + path.Dir(sub.Code.Name),
	})
	if err != nil {
		return nil, nil, err
	}

	c, err := dc.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: newImage(sub.Language),
		},
		HostConfig: &docker.HostConfig{
			Privileged: false,
			Memory:     0, // TODO(flowlo): Limit memory
			Binds:      []string{v.Name + ":/run"},
		},
	})
	if err != nil {
		return nil, nil, err
	}

	err = dc.StartContainer(c.ID, c.HostConfig)
	if err != nil {
		return nil, nil, err
	}

	waitc := make(chan waitResult)
	go func() {
		exitCode, err := dc.WaitContainer(c.ID)
		waitc <- waitResult{exitCode, err}
	}()

	var res waitResult
	select {
	case res = <-waitc:
	case <-time.After(time.Minute):
		return nil, nil, errors.New("execution timed out")
	}

	if res.Err != nil {
		return nil, nil, res.Err
	}

	stdout, stderr = new(bytes.Buffer), new(bytes.Buffer)
	err = dc.Logs(docker.LogsOptions{
		Container:    c.ID,
		OutputStream: stdout,
		ErrorStream:  stderr,
		Stdout:       true,
		Stderr:       true,
	})
	if err != nil {
		return nil, nil, err
	}

	return
}
