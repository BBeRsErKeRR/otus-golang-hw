package main

import (
	"io"
	"time"
)

type Reader struct {
	io.Reader
	bar            *Progress
	freezeDuration time.Duration
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	r.bar.Add(n)
	// sleep to see copy
	time.Sleep(r.freezeDuration)
	return
}

func (r *Reader) Close() (err error) {
	// Close the reader when it implements io.Closer
	r.bar.Finish()
	if closer, ok := r.Reader.(io.Closer); ok {
		return closer.Close()
	}
	return
}
