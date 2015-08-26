package runner

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/coduno/api/model"
	"github.com/fsouza/go-dockerclient"
	"golang.org/x/net/context"
)

type waitResult struct {
	ExitCode int
	Err      error
}

func simpleRunner(ctx context.Context, test *model.Test, sub model.KeyedSubmission) error {
	c, err := dc.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: newImage(sub.Language),
		},
		HostConfig: &docker.HostConfig{
			Privileged: false,
			Memory:     0, // TODO(flowlo): Limit memory
		},
	})
	if err != nil {
		return err
	}

	err = dc.StartContainer(c.ID, c.HostConfig)
	if err != nil {
		return err
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
		return errors.New("execution timed out")
	}

	if res.Err != nil {
		return err
	}

	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	err = dc.Logs(docker.LogsOptions{
		OutputStream: stdout,
		ErrorStream:  stderr,
		Stdout:       true,
		Stderr:       true,
	})
	if err != nil {
		return err
	}

	var result = struct {
		Stdout, Stderr string
	}{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	fmt.Printf("%v", result)

	// TODO(flowlo): Store result in Datastore.
	return nil
}
