package storage

import (
	"io"
	"time"

	"golang.org/x/net/context"
)

type Object interface {
	Name() string

	io.Writer
	io.Reader
	io.Closer
}

type Provider interface {
	Create(ctx context.Context, name string, maxAge time.Duration, contentType string) (Object, error)
}
