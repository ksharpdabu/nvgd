package ltsv

import (
	"bufio"
	"bytes"
	"io"
)

// Reader reads LSTV values.
type Reader struct {
	rd *bufio.Reader
}

// NewReader creates a new LSTV reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		rd: bufio.NewReader(r),
	}
}

func (r *Reader) readLine() ([]byte, error) {
	d, err := r.rd.ReadSlice('\n')
	if err == nil {
		return d, nil
	} else if err != bufio.ErrBufferFull {
		return nil, err
	}
	bb := bytes.NewBuffer(d)
	for {
		d2, err := r.rd.ReadSlice('\n')
		if len(d2) > 0 {
			if _, err := bb.Write(d2); err != nil {
				return nil, err
			}
		}
		if err == nil || err == io.EOF {
			return bb.Bytes(), nil
		}
		if err != bufio.ErrBufferFull {
			return nil, err
		}
	}
}

// Read read a LSTV value.
func (r *Reader) Read() (*Set, error) {
	d, err := r.readLine()
	if err != nil {
		return nil, err
	}
	s := &Set{
		Index: make(map[string][]int),
	}
	for _, raw := range bytes.Split(d, []byte("\t")) {
		kv := bytes.SplitN(raw, []byte(":"), 2)
		if len(kv) != 2 {
			continue
		}
		s.Put(string(kv[0]), string(kv[1]))
	}
	return s, nil
}

// Set is a set of LSTV values in a line.
type Set struct {
	Properties []Property
	Index      map[string][]int
}

// Put puts a property to the set.
func (s *Set) Put(label, value string) {
	n := len(s.Properties)
	s.Properties = append(s.Properties, Property{Lavel: label, Value: value})
	s.Index[label] = append(s.Index[label], n)
}

func (s *Set) Get(label string) []string {
	indexes, ok := s.Index[label]
	if !ok {
		return nil
	}
	list := make([]string, len(indexes))
	for i, n := range indexes {
		list[i] = s.Properties[n].Value
	}
	return list
}

// Property is a pair of label and value.
type Property struct {
	Lavel string
	Value string
}
