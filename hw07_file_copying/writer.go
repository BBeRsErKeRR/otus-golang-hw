package main

import (
	"io"
)

type Writer struct {
	io.Writer
	bar *Progress
}

func (r *Writer) Write(p []byte) (n int, err error) {
	n, err = r.Writer.Write(p)
	r.bar.Add(n)
	return
}

func (r *Writer) Close() (err error) {
	// Close the reader when it implements io.Closer
	r.bar.Finish()
	if closer, ok := r.Writer.(io.Closer); ok {
		return closer.Close()
	}
	return
}
