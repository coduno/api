package runner

import (
	"bytes"

	"github.com/coduno/api/model"
	"github.com/fsouza/go-dockerclient"
	"golang.org/x/net/context"
)

func diffRunner(ctx context.Context, test *model.Test, sub model.KeyedSubmission) error {
	c, err := dc.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			// TODO(flowlo): Check if the language is known.
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

	stdout, stderr, stdin := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)

	err = dc.AttachToContainer(docker.AttachToContainerOptions{
		Container:    c.ID,
		OutputStream: stdout,
		Stdout:       true,
		Stream:       true,
	})
	if err != nil {
		return err
	}

	err = dc.AttachToContainer(docker.AttachToContainerOptions{
		Container:    c.ID,
		OutputStream: stderr,
		Stderr:       true,
		Stream:       true,
	})
	if err != nil {
		return err
	}

	err = dc.AttachToContainer(docker.AttachToContainerOptions{
		Container:    c.ID,
		OutputStream: stdin,
		Stdin:        true,
		Stream:       true,
	})
	if err != nil {
		return err
	}

	// TODO(flowlo): Save result.
	return nil
}
