package logic

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/coduno/api/cache"
	"github.com/coduno/api/model"
	"github.com/coduno/api/util"
	"google.golang.org/appengine/datastore"
)

type InsDel struct {
	Inserted,
	Deleted int
}

func (id *InsDel) Add(insDel InsDel) {
	id.Inserted += insDel.Inserted
	id.Deleted += insDel.Deleted
}

func getInsertedDeleted(submissions []model.KeyedSubmission) (insDel *InsDel, err error) {
	var oldCode io.Reader
	oldCode, err = readFromGCS(submissions[0].Code)
	if err != nil {
		return
	}
	var initialCode []byte
	if initialCode, err = ioutil.ReadAll(oldCode); err != nil {
		return
	}
	insDel = &InsDel{len(bytes.Split(initialCode, []byte("\n"))), 0}
	// Iterate all submissions
	for i := 1; i < len(submissions); i++ {
		var newCode io.Reader
		newCode, err = readFromGCS(submissions[i].Code)
		if err != nil {
			return
		}
		var id InsDel
		if id, err = computeInsertedDeletedLines(newCode, oldCode); err != nil {
			return
		}
		insDel.Add(id)
		oldCode = newCode
	}
	return
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

func readFromGCS(so model.StoredObject) (io.Reader, error) {
	return cache.PutGCS(util.CloudContext(nil), so.Bucket, so.Name)
}
