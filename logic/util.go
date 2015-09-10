package logic

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/coduno/api/model"
	"github.com/coduno/api/util"
	"google.golang.org/appengine/datastore"
	"google.golang.org/cloud/storage"
)

type InsDel struct {
	Inserted,
	Deleted int
}

func (id *InsDel) Add(insDel InsDel) {
	id.Inserted += insDel.Inserted
	id.Deleted += insDel.Deleted
}

func computeInsertedDeletedLines(oldCodeR, newCodeR io.Reader) (id InsDel, err error) {
	var i, d int
	// TODO(flowlo): get rid of ReadAll
	var oldCode, newCode []byte
	if oldCode, err = ioutil.ReadAll(oldCodeR); err != nil {
		return
	}

	if newCode, err = ioutil.ReadAll(newCodeR); err != nil {
		return
	}
	currentFields := bytes.Split(newCode, []byte("\n"))
	oldFields := bytes.Split(oldCode, []byte("\n"))
	for _, val := range currentFields {
		if !bytes.Contains(oldCode, val) {
			i++
		}
	}
	for _, val := range oldFields {
		if !bytes.Contains(oldCode, val) {
			d++
		}
	}
	return InsDel{i, d}, nil
}

func getTaskIndex(c model.Challenge, task *datastore.Key) int {
	for i, val := range c.Tasks {
		if val.Equal(task) {
			return i
		}
	}
	return -1
}

func readFromGCS(so model.StoredObject) (io.ReadCloser, error) {
	return storage.NewReader(util.CloudContext(nil), so.Bucket, so.Name)
}
