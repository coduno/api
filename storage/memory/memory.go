package google

import (
	"bytes"
	"time"

	"golang.org/x/net/context"

	"github.com/coduno/api/storage"
)

type provider map[string]*object

func (p provider) Create(ctx context.Context, name string, maxAge time.Duration, contentType string) storage.Object {
	o := &object{
		b: new(bytes.Buffer),
		n: name,
	}
	p[o.n] = o
	return o
}

type object struct {
	n string
	b *bytes.Buffer
}

func (o *object) Write(p []byte) (n int, err error) {
	return o.b.Write(p)
}

func (o *object) Close() error {
	return nil
}

func (o *object) Read(p []byte) (n int, err error) {
	return o.b.Read(p)
}

func (o *object) Name() string {
	return o.n
}

func NewProvider() storage.Provider {
	return provider{}
}
