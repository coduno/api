package runner

import (
	"archive/tar"
	"bytes"
	"encoding/xml"
	"io"
	"path"
	"strings"
	"time"

	"github.com/coduno/api/model"
	"github.com/coduno/api/util"
	"github.com/fsouza/go-dockerclient"
	"golang.org/x/net/context"
)

func JUnit(ctx context.Context, testFle string, sub model.KeyedSubmission) (testResults model.JunitTestResult, err error) {
	image := newImage("javaut")

	if err = prepareImage(image); err != nil {
		return
	}

	var v *docker.Volume
	if v, err = createDockerVolume(sub.Code.Bucket + "/" + path.Dir(sub.Code.Name)); err != nil {
		return
	}

	var testV *docker.Volume
	if testV, err = createDockerVolume(util.TestsBucket + "/" + testFle); err != nil {
		return
	}

	binds := []string{v.Name + ":/run/src/main/java", testV.Name + ":/run/src/test/java"}
	var c *docker.Container
	if c, err = createDockerContainer(image, binds); err != nil {
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

	pr, pw := io.Pipe()

	err = dc.CopyFromContainer(docker.CopyFromContainerOptions{
		Container:    c.ID,
		Resource:     util.JUnitResultsPath,
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
		testResults = model.JunitTestResult{
			Stdout:  stdout.String(),
			Results: utr,
			Stderr:  stderr.String(),
			Start:   start,
			End:     end}
	}

	return
}
