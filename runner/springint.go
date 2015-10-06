package runner

import (
	"archive/tar"
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/coduno/api/model"
	"github.com/fsouza/go-dockerclient"
	"golang.org/x/net/context"

	"google.golang.org/appengine/log"
)

func SpringInt(ctx context.Context, sub model.KeyedSubmission, ball io.Reader) (*model.JunitTestResult, error) {
	image := newImage("spring-integration")

	if err := prepareImage(image); err != nil {
		return nil, err
	}

	c, err := itoc(image)
	if err != nil {
		return nil, err
	}

	err = dc.UploadToContainer(c.ID, docker.UploadToContainerOptions{
		Path:        "/run/src/main/java/test/controller/",
		InputStream: ball,
	})
	if err != nil {
		return nil, err
	}

	start := time.Now()
	if err := dc.StartContainer(c.ID, c.HostConfig); err != nil {
		return nil, err
	}

	log.Debugf(ctx, "SpringInt: Waiting for container")
	if err := waitForContainer(c.ID); err != nil {
		return nil, err
	}
	end := time.Now()

	stdout, stderr, err := getLogs(c.ID)
	if err != nil {
		return nil, err
	}

	testResults := &model.JunitTestResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
		Start:  start,
		End:    end,
	}

	errc := make(chan error)
	dpr, dpw := io.Pipe()

	log.Debugf(ctx, "SpringInt: Download for container")
	go func() {
		errc <- dc.DownloadFromContainer(c.ID, docker.DownloadFromContainerOptions{
			Path:         "/run/target/surefire-reports/TEST-test.ControllerTestApplicationTests.xml",
			OutputStream: dpw,
		})
	}()

	// TODO(flowlo): encoding/xml can only parse XML 1.0, but JUnit
	// reports are XML 1.1. It appears the reports are valid XML 1.0
	// too, so this replaces the version attribute.
	// As soon as encoding/xml can parse XML 1.1 we can remove this
	// and directly stream without buffering.
	log.Debugf(ctx, "SpringInt: Before read")
	buf, err := ioutil.ReadAll(dpr)
	if err != nil {
		return nil, err
	}
	log.Debugf(ctx, "SpringInt: after read")

	buf = bytes.Replace(buf, []byte(`version="1.1"`), []byte(`version="1.0"`), 1)

	tr := tar.NewReader(bytes.NewReader(buf))
	d := xml.NewDecoder(tr)

	log.Debugf(ctx, "SpringInt: after decoder")
	for {
		h, err := tr.Next()
		if err == io.EOF {
			// We reached EOF, so this loop has no
			// chance to continue. Later reads from
			// err will EOF too, which is acceptable.
			break
		}
		if err != nil {
			return nil, err
		}
		if strings.HasSuffix(h.Name, ".xml") {
			// We're definitely looking for an XML
			// file, so skip everything else.
			break
		}
	}

	log.Debugf(ctx, "SpringInt: after decoding")
	if err := d.Decode(&testResults.Results); err != nil {
		// Decode might very well error, for
		// example with the EOF from above.
		// This at the same time indicates
		// that no XML file was found.
		log.Debugf(ctx, "SpringInt: error decoding %s", err)
		return nil, err
	}
	log.Debugf(ctx, "Spring runner done %+v", testResults)

	return testResults, <-errc
}
