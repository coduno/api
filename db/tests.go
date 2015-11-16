package db

import (
	"io"

	"github.com/coduno/api/model"
)

func LoadFile(name string) io.ReadCloser {
	return nil
}

func LoadTestsForTask(taskId int64) ([]model.Test, error) {
	return nil, nil
}
