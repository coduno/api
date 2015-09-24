package runner

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"path"
	"time"

	"golang.org/x/net/context"

	"github.com/coduno/api/model"
	"github.com/coduno/api/util"
	"github.com/fsouza/go-dockerclient"
)

type waitResult struct {
	ExitCode int
	Err      error
}

func Simple(ctx context.Context, sub model.KeyedSubmission) (testResult model.SimpleTestResult, err error) {
	image := newImage(sub.Language)

	if err = prepareImage(image); err != nil {
		return
	}

	var v *docker.Volume
	if v, err = createDockerVolume(sub.Code.Bucket + "/" + path.Dir(sub.Code.Name)); err != nil {
		return
	}

	var c *docker.Container
	if c, err = createDockerContainer(image, []string{v.Name + ":/run"}); err != nil {
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

	errc := make(chan error)
	pr, pw := io.Pipe()

	go func() {
		errc <- dc.DownloadFromContainer(c.ID, docker.DownloadFromContainerOptions{
			Path:         util.StatsPath,
			OutputStream: pw,
		})
	}()

	var buf []byte
	buf, err = ioutil.ReadAll(pr)
	if err != nil {
		return
	}

	tr := tar.NewReader(bytes.NewReader(buf))
	d := json.NewDecoder(tr)
	var h *tar.Header
	var rusage model.Rusage
	for {
		h, err = tr.Next()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return
		}
		if h.Name != "stats.log" {
			continue
		}
		if err = d.Decode(&rusage); err != nil {
			return
		}
	}
	testResult.Rusage = rusage

	return
}
