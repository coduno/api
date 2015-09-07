package runner

import (
	"archive/tar"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"path"
	"strings"
	"time"

	"github.com/coduno/api/model"
	"github.com/fsouza/go-dockerclient"
	"golang.org/x/net/context"
)

func JUnit(ctx context.Context, params map[string]string, sub model.KeyedSubmission) (stdout, stderr *bytes.Buffer, testResults []model.UnitTestResults, err error) {
	image := newImage("javaut")

	if err = prepareImage(ctx, image); err != nil {
		return nil, nil, []model.UnitTestResults{}, err
	}

	var v *docker.Volume
	if v, err = createDockerVolume(sub.Code.Bucket + "/" + path.Dir(sub.Code.Name)); err != nil {
		return
	}

	var testV *docker.Volume
	if testV, err = createDockerVolume(params["tests"]); err != nil {
		return
	}

	var c *docker.Container
	c, err = dc.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: image,
		},
		HostConfig: &docker.HostConfig{
			Privileged: false,
			Memory:     0, // TODO(flowlo): Limit memory
			Binds:      []string{v.Name + ":/run/src/main/java", testV.Name + ":/run/src/test/java"},
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
		return nil, nil, []model.UnitTestResults{}, errors.New("execution timed out")
	}

	if res.Err != nil {
		return nil, nil, []model.UnitTestResults{}, res.Err
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
		return
	}

	pr, pw := io.Pipe()

	err = dc.CopyFromContainer(docker.CopyFromContainerOptions{
		Container:    c.ID,
		Resource:     params["resultPath"],
		OutputStream: pw,
	})
	if err != nil {
		return
	}

	tr := tar.NewReader(pr)
	d := xml.NewDecoder(tr)
	var h *tar.Header
	for {
		h, err = tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		if !strings.HasSuffix(h.Name, ".xml") {
			continue
		}

		var utr model.UnitTestResults
		if err = d.Decode(&utr); err != nil {
			return
		}
		testResults = append(testResults, utr)
	}
	// TODO(victorbalan, flowlo): Get a single result to have the correct
	// start and end time when we will do different runs for every file
	// in the GCS bucket.
	// Also, datastore.PutMulti could be used to insert as batch here.
	for _, val := range testResults {
		jtr := model.JunitTestResult{
			Stdout:  stdout.String(),
			Results: val,
			Stderr:  stderr.String(),
			Start:   start,
			End:     time.Now()}
		if _, err = jtr.Put(ctx, nil); err != nil {
			return
		}
	}

	return
}
