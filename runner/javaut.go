package runner

import (
	"archive/tar"
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/coduno/api/model"
	"github.com/coduno/api/util"
	"github.com/fsouza/go-dockerclient"
	"golang.org/x/net/context"
)

func JUnit(ctx context.Context, testFile string, sub model.KeyedSubmission) (testResults model.JunitTestResult, err error) {
	image := newImage("javaut")

	if err = prepareImage(image); err != nil {
		return
	}

	var v *docker.Volume
	if v, err = createDockerVolume(sub.Code.Bucket + "/" + path.Dir(sub.Code.Name)); err != nil {
		return
	}

	var testV *docker.Volume
	if testV, err = createDockerVolume(util.TestsBucket + "/" + testFile); err != nil {
		return
	}

	binds := []string{v.Name + ":/run/src/main/java", testV.Name + ":/run/src/test/java/" + testFile}
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

	errc := make(chan error)
	pr, pw := io.Pipe()

	go func() {
		errc <- dc.DownloadFromContainer(c.ID, docker.DownloadFromContainerOptions{
			Path:         util.JUnitResultsPath,
			OutputStream: pw,
		})
	}()

	// TODO(flowlo): encoding/xml can only parse XML 1.0, but JUnit
	// reports are XML 1.1. It appears the reports are valid XML 1.0
	// too, so this replaces the version attribute.
	// As soon as encoding/xml can parse XML 1.1 we can remove this
	// and directly stream without buffering.
	var buf []byte
	buf, err = ioutil.ReadAll(pr)
	if err != nil {
		return
	}
	buf = bytes.Replace(buf, []byte(`version="1.1"`), []byte(`version="1.0"`), 1)

	tr := tar.NewReader(bytes.NewReader(buf))
	d := xml.NewDecoder(tr)
	var h *tar.Header
	for {
		h, err = tr.Next()
		if err == io.EOF {
			err = nil
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

	err = <-errc
	return
}
