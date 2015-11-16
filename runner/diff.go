package runner

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/coduno/api/db"
	"github.com/coduno/api/model"
	"github.com/fsouza/go-dockerclient"
)

func IODiffRun(ctx context.Context, t model.Test, sub model.Submission, ball io.Reader) (ts model.TestStats, err error) {
	image := newImage(sub.Language)

	if err = prepareImage(image); err != nil {
		return
	}

	var c *docker.Container
	c, err = dc.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:     image,
			OpenStdin: true,
			StdinOnce: true,
		},
		HostConfig: &docker.HostConfig{
			Privileged:  false,
			NetworkMode: "none",
			Memory:      0, // TODO(flowlo): Limit memory
		},
	})
	if err != nil {
		return
	}

	err = dc.UploadToContainer(c.ID, docker.UploadToContainerOptions{
		Path:        "/run",
		InputStream: ball,
	})
	if err != nil {
		return
	}

	stdin := db.LoadFile(t.Params["input"])
	defer stdin.Close()

	start := time.Now()
	if err = dc.StartContainer(c.ID, c.HostConfig); err != nil {
		return
	}

	err = dc.AttachToContainer(docker.AttachToContainerOptions{
		Container:   c.ID,
		InputStream: stdin,
		Stdin:       true,
		Stream:      true,
	})
	if err != nil {
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

	tr := model.DiffTestResult{
		SimpleTestResult: model.SimpleTestResult{
			Stdout: stdout.String(),
			Stderr: stderr.String(),
			Start:  start,
			End:    end,
		},
		Endpoint: "diff-result",
	}

	return processDiffResults(sub, tr, t.Params["output"], t.ID)
}

func OutMatchDiffRun(ctx context.Context, t model.Test, sub model.Submission, ball io.Reader) (ts model.TestStats, err error) {
	var str model.SimpleTestResult
	str, err = Simple(ctx, sub, ball)
	if err != nil {
		return
	}
	tr := model.DiffTestResult{
		SimpleTestResult: str,
		Endpoint:         "diff-result",
	}

	return processDiffResults(sub, tr, t.Params["tests"], t.ID)
}

func processDiffResults(sub model.Submission, tr model.DiffTestResult, testFile string, test int64) (ts model.TestStats, err error) {
	want := db.LoadFile(testFile)
	defer want.Close()

	have := strings.NewReader(tr.Stdout)
	diffLines, ok, err := compare(want, have)
	if err != nil {
		return
	}
	tr.DiffLines = diffLines
	tr.Failed = !ok

	ts = model.TestStats{
		Stdout: tr.Stdout,
		Stderr: tr.Stderr,
		Test:   test,
		Failed: !ok,
	}

	// _, err = tr.PutWithParent(ctx, sub.Key)
	return
}

func compare(want, have io.Reader) ([]int, bool, error) {
	w, err := ioutil.ReadAll(want)
	if err != nil {
		return nil, false, err
	}
	h, err := ioutil.ReadAll(have)
	if err != nil {
		return nil, false, err
	}
	wb := bytes.Split(w, []byte("\n"))
	hb := bytes.Split(h, []byte("\n"))

	if len(wb) != len(hb) {
		return nil, false, nil
	}

	var diff []int
	ok := true
	for i := 0; i < len(wb); i++ {
		if bytes.Compare(wb[i], hb[i]) != 0 {
			diff = append(diff, i)
			ok = false
		}
	}

	return diff, ok, nil
}
