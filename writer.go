package main

import "fmt"

type CapWriter struct {
	Cap       uint64
	NoDiscard bool
	size      uint64
	buffer    []byte
}

func (w *CapWriter) Write(p []byte) (int, error) {
	w.size += uint64(len(p))
	if w.size > w.Cap && w.NoDiscard {
		return 0, fmt.Errorf("could not write body buffer. buffer is full")
	}

	if w.size > w.Cap {
		q := w.Cap - uint64(len(w.buffer))
		if q != 0 {
			w.buffer = append(w.buffer, p[0:q-1]...)
		}
	} else {
		w.buffer = append(w.buffer, p...)
	}

	return len(p), nil
}

func (w *CapWriter) Size() uint64 {
	return w.size
}

func (w *CapWriter) Bytes() []byte {
	return w.buffer
}
