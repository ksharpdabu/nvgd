package tail

import (
	"bytes"
	"io"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/ringbuf"
)

// Tail is "tail" like filter.
type Tail struct {
	filter.Base

	rf bool
	rb *ringbuf.Buffer
}

// NewTail creates an instance of tail filter.
func NewTail(r io.ReadCloser, limit int) *Tail {
	if limit <= 0 {
		limit = 10
	}
	t := &Tail{
		rb: ringbuf.New(limit),
	}
	t.Base.Init(r, t.readNext)
	return t
}

func (t *Tail) readNext(buf *bytes.Buffer) error {
	if !t.rf {
		t.rf = true
		if err := t.readAll(); err != nil {
			return err
		}
	}
	if t.rb.Empty() {
		return io.EOF
	}
	for {
		v, ok := t.rb.Get()
		if !ok {
			return io.EOF
		}
		_, err := buf.Write(v.([]byte))
		if err != nil {
			return err
		}
	}
}

func (t *Tail) readAll() error {
	for {
		b, err := t.ReadLine()
		if err == io.EOF {
			if len(b) > 0 {
				t.rb.Put(b)
			}
			return nil
		} else if err != nil {
			return err
		}
		t.rb.Put(b)
	}
}
