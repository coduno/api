package runner

import (
	"bytes"
	"errors"
	"time"

	"github.com/coduno/api/model"
	"github.com/fsouza/go-dockerclient"
)

type waitResult struct {
	ExitCode int
	Err      error
}

func Simple(sub model.KeyedSubmission) (stdout, stderr *bytes.Buffer, err error) {
	c, err := dc.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: newImage(sub.Language),
		},
		HostConfig: &docker.HostConfig{
			Privileged: false,
			Memory:     0, // TODO(flowlo): Limit memory
			Binds:      []string{"/tmp/submissions:/run"},
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
