package cache

import (
	"bytes"
	"io"
	"sync"
	"time"

	"google.golang.org/cloud/storage"

	"golang.org/x/net/context"
)

type item struct {
	// Holds the parts of the resource that were already
	// read from r. Once buf is exhausted, further contents
	// should be read from r.
	buf []byte

	// The upstream source of the resource.
	r io.Reader

	// Guards reading on both buf and r.
	sync.Mutex

	// Set once r was read till the end. Inspected for
	// cleaning the cache.
	closed time.Time
}

func newItem(r io.Reader, c int) *item {
	return &item{
		buf:   make([]byte, 0, c),
		r:     r,
		Mutex: sync.Mutex{},
	}
}

type reader struct {
	// Holds the buffer to read from and write to
	// if exhausted.
	*item

	// Read from buf[off].
	off int
}

func (r *reader) Read(p []byte) (n int, err error) {
	n = len(p)

	r.Lock()
	defer r.Unlock()

	// Underlying source is already gone, so serve all
	// reads from buf. Basically act like a bytes.Reader.
	if r.r == nil {
		if r.off+n >= len(r.buf) {
			n = len(r.buf) - r.off
			err = io.EOF
		}
		n = copy(p, r.buf[r.off:])
		r.off += n
		return
	}

	// Underlying source is still there, but we might have
	// everything in buf already.
	if r.off+n < len(r.buf) {
		n = copy(p, r.buf[r.off:])
		r.off += n
		return n, nil
	}

	// We're out of luck, advance the underlying reader.
	n, err = r.r.Read(p)
	if err != nil && err != io.EOF {
		return 0, err
	}

	// If we read something, append it to buf.
	if n != 0 {
		r.buf = append(r.buf, p[:n]...)
	}

	n = copy(p, r.buf[r.off:r.off+n])
	r.off += n

	// Read from underlying reader reached io.EOF,
	// so we clean up and close the underlying
	// reader if possible.
	if err == io.EOF {
		c, ok := r.r.(io.Closer)
		if ok {
			c.Close()
		}
		r.r = nil
		r.closed = time.Now()
	}

	return n, err
}

// Cache is a map of strings to items that is
type Cache struct {
	items map[string]*item
	sync.RWMutex
}

var global = New()

// New makes an empty Cache.
func New() *Cache {
	return &Cache{
		items:   map[string]*item{},
		RWMutex: sync.RWMutex{},
	}
}

// Item is an object to be cached. It is referred to by it's key,
// returned by Key(). If the item is not found in the cache, Fetch()
// will be called to obtain a copy.
type Item interface {
	Key() string
	Fetch() (io.Reader, error)
}

// Sizer let's an Item additionaly hint the size of the object to
// be cached so that the buffer inside the cache will be sized to
// this capacity, reducing overhead from appending to it.
type Sizer interface {
	Size() int
}

type gcsf struct {
	Context context.Context
	Bucket  string
	Name    string
}

func (f gcsf) Key() string {
	return f.Bucket + "/" + f.Name
}

func (f gcsf) Fetch() (io.Reader, error) {
	return storage.NewReader(f.Context, f.Bucket, f.Name)
}

// Put adds an Item to c. If the Item is not already present in the
// cache, it.Fetch will be called to pull it's contents.
func (c *Cache) Put(it Item) (io.Reader, error) {
	k := it.Key()

	if r, ok := c.Get(k); ok {
		return r, nil
	}

	c.Lock()
	defer c.Unlock()
	if r, ok := c.Get(k); ok {
		return r, nil
	}

	r, err := it.Fetch()
	if err != nil {
		return nil, err
	}

	cap := 0
	if s, ok := it.(Sizer); ok {
		cap = s.Size()
	}

	i := newItem(r, cap)
	c.items[k] = i
	return &reader{i, 0}, err
}

// PutRaw is a lower-level version of Put, with the ability to directly
// Specify a reader for the item in case it's not cached yet. Use with
// caution.
func (c *Cache) PutRaw(key string, r io.Reader) io.Reader {
	k := key

	if r, ok := c.Get(k); ok {
		return r
	}

	c.Lock()
	defer c.Unlock()
	if r, ok := c.Get(k); ok {
		return r
	}

	i := newItem(r, 0)
	c.items[k] = i
	return &reader{i, 0}
}

// Get looks something up in the cache.
func (c *Cache) Get(key string) (io.Reader, bool) {
	c.RLock()
	i, ok := c.items[key]
	c.RUnlock()
	if !ok {
		return nil, false
	}
	if i.r == nil {
		return bytes.NewReader(i.buf), true
	}
	return &reader{i, 0}, true
}

// PutGCS adds a Google Cloud Storage object to the cache.
func (c *Cache) PutGCS(ctx context.Context, bucket, name string) (io.Reader, error) {
	return c.Put(gcsf{Name: name, Bucket: bucket, Context: ctx})
}

func Get(key string) (io.Reader, bool) {
	return global.Get(key)
}

func PutGCS(ctx context.Context, bucket, name string) (io.Reader, error) {
	return global.PutGCS(ctx, bucket, name)
}
