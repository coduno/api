package runner

import (
	"bytes"
	"io"
	"time"

	"golang.org/x/net/context"

	"github.com/coduno/api/model"
	"github.com/fsouza/go-dockerclient"
)

type waitResult struct {
	ExitCode int
	Err      error
}

func Simple(ctx context.Context, sub model.Submission, ball io.Reader) (testResult model.SimpleTestResult, err error) {
	image := newImage(sub.Language)

	if err = prepareImage(image); err != nil {
		return
	}

	var c *docker.Container
	if c, err = itoc(image); err != nil {
		return
	}

	err = dc.UploadToContainer(c.ID, docker.UploadToContainerOptions{
		Path:        "/run",
		InputStream: ball,
	})
	if err != nil {
		return
	}

	start := time.Now()
	if err = dc.StartContainer(c.ID, c.HostConfig); err != nil {
		return
	}

	if err = waitForContainer(c.ID); err != nil {
		return
	}
	end := time.Now()

	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	if stdout, stderr, err = getLogs(c.ID); err != nil {
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
