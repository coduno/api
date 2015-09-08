package runner

import (
	"bytes"
	"errors"
	"path"
	"time"

	"golang.org/x/net/context"

	"github.com/coduno/api/model"
	"github.com/fsouza/go-dockerclient"
)

type waitResult struct {
	ExitCode int
	Err      error
}

func Simple(ctx context.Context, sub model.KeyedSubmission) (testResult model.SimpleTestResult, err error) {
	image := newImage(sub.Language)

	if err = prepareImage(ctx, image); err != nil {
		return
	}

	var v *docker.Volume
	if v, err = createDockerVolume(sub.Code.Bucket + "/" + path.Dir(sub.Code.Name)); err != nil {
		return
	}

	var c *docker.Container
	c, err = dc.CreateContainer(docker.CreateContainerOptions{
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
		return
	}

	start := time.Now()
	if err = dc.StartContainer(c.ID, c.HostConfig); err != nil {
		return
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
		err = errors.New("execution timed out")
		return
	}

	if res.Err != nil {
		err = res.Err
		return
	}
	end := time.Now()

	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	err = dc.Logs(docker.LogsOptions{
		Container:    c.ID,
		OutputStream: stdout,
		ErrorStream:  stderr,
		Stdout:       true,
		Stderr:       true,
	})
	if err != nil {
		return
	}

	testResult = model.SimpleTestResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
		Start:  start,
		End:    end,
	}

	return
}
