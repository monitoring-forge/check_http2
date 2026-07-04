package main

import "fmt"

type CapWriter struct {
	Cap       uint64
	NoDiscard bool
	size      uint64
	buffer    []byte
}

func (w *CapWriter) Write(p []byte) (int, error) {
	size := w.size + uint64(len(p))
	if size > w.Cap && w.NoDiscard {
		return 0, fmt.Errorf("could not write body buffer. buffer is full")
	}

	if size > w.Cap {
		q := w.Cap - uint64(len(w.buffer))
		if q != 0 {
			w.buffer = append(w.buffer, p[0:q]...)
		}
		w.size = w.Cap
		return int(q), nil
	}

	w.buffer = append(w.buffer, p...)
	w.size = size
	return len(p), nil
}

func (w *CapWriter) Size() uint64 {
	return w.size
}

func (w *CapWriter) Bytes() []byte {
	return w.buffer
}
