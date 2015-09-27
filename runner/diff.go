package runner

import (
	"bytes"
	"io"
	"io/ioutil"
	"path"
	"strings"
	"time"

	"google.golang.org/appengine/datastore"

	"github.com/coduno/api/cache"
	"github.com/coduno/api/model"
	"github.com/coduno/api/util"
	"github.com/fsouza/go-dockerclient"
	"golang.org/x/net/context"
)

func IODiffRun(ctx context.Context, t model.KeyedTest, sub model.KeyedSubmission) (ts model.TestStats, err error) {
	image := newImage(sub.Language)

	if err = prepareImage(image); err != nil {
		return
	}

	var v *docker.Volume
	if v, err = createDockerVolume(sub.Code.Bucket + "/" + path.Dir(sub.Code.Name)); err != nil {
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
			Privileged: false,
			Memory:     0, // TODO(flowlo): Limit memory
			Binds:      []string{v.Name + ":/run"},
		},
	})
	if err != nil {
		return
	}

	var stdin io.Reader
	stdin, err = cache.PutGCS(util.CloudContext(ctx), util.TestsBucket, t.Params["input"])
	if err != nil {
		return
	}

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
	}

	return processDiffResults(ctx, tr, util.TestsBucket, t.Params["output"], t.Key)
}

func OutMatchDiffRun(ctx context.Context, t model.KeyedTest, sub model.KeyedSubmission) (ts model.TestStats, err error) {
	var str model.SimpleTestResult
	str, err = Simple(ctx, sub)
	if err != nil {
		return
	}
	tr := model.DiffTestResult{
		SimpleTestResult: str,
	}

	return processDiffResults(ctx, tr, util.TestsBucket, t.Params["tests"], t.Key)
}

func processDiffResults(ctx context.Context, tr model.DiffTestResult, bucket, testFile string, test *datastore.Key) (ts model.TestStats, err error) {
	var want io.Reader
	want, err = cache.PutGCS(util.CloudContext(ctx), bucket, testFile)
	if err != nil {
		return
	}

	have := strings.NewReader(tr.Stdout)
	diffLines, ok, err := compare(want, have)
	if err != nil {
		return
	}
	tr.DiffLines = diffLines

	ts = model.TestStats{
		Stdout: tr.Stdout,
		Stderr: tr.Stderr,
		Test:   test,
		Failed: !ok,
	}
	_, err = tr.Put(ctx, nil)
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
	for i := 0; i < len(wb); i++ {
		if bytes.Compare(wb[i], hb[i]) != 0 {
			diff = append(diff, i)
		}
	}

	return diff, true, nil
}
