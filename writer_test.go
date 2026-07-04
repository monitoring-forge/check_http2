package main

import "testing"

func TestCapWriterWriteWithinCap(t *testing.T) {
	w := &CapWriter{Cap: 8}

	n, err := w.Write([]byte("hello"))
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if n != 5 {
		t.Fatalf("Write() n = %d, want 5", n)
	}
	if got := string(w.Bytes()); got != "hello" {
		t.Fatalf("Bytes() = %q, want %q", got, "hello")
	}
	if got := w.Size(); got != 5 {
		t.Fatalf("Size() = %d, want 5", got)
	}
}

func TestCapWriterWriteSameSizeCapDiscard(t *testing.T) {
	w := &CapWriter{Cap: 4}

	n, err := w.Write([]byte("abcd"))
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if n != 4 {
		t.Fatalf("Write() n = %d, want 4", n)
	}
	if got := string(w.Bytes()); got != "abcd" {
		t.Fatalf("Bytes() = %q, want %q", got, "abcd")
	}
	if got := w.Size(); got != 4 {
		t.Fatalf("Size() = %d, want 4", got)
	}
}

func TestCapWriterWriteOverCapDiscard(t *testing.T) {
	w := &CapWriter{Cap: 4}

	n, err := w.Write([]byte("abcdef"))
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if n != 6 {
		t.Fatalf("Write() n = %d, want 6", n)
	}
	if got := string(w.Bytes()); got != "abcd" {
		t.Fatalf("Bytes() = %q, want %q", got, "abcd")
	}
	if got := w.Size(); got != 6 {
		t.Fatalf("Size() = %d, want 6", got)
	}
}

func TestCapWriterWriteOverCapWithExistingBufferDiscard(t *testing.T) {
	w := &CapWriter{Cap: 5}
	if _, err := w.Write([]byte("abc")); err != nil {
		t.Fatalf("first Write() error = %v", err)
	}

	n, err := w.Write([]byte("defg"))
	if err != nil {
		t.Fatalf("second Write() error = %v", err)
	}
	if n != 4 {
		t.Fatalf("second Write() n = %d, want 4", n)
	}
	if got := string(w.Bytes()); got != "abcde" {
		t.Fatalf("Bytes() = %q, want %q", got, "abcde")
	}
	if got := w.Size(); got != 7 {
		t.Fatalf("Size() = %d, want 7", got)
	}
}

func TestCapWriterWriteOverCapNoDiscard(t *testing.T) {
	w := &CapWriter{Cap: 4, NoDiscard: true}

	n, err := w.Write([]byte("abcdef"))
	if err == nil {
		t.Fatalf("Write() error = nil, want non-nil")
	}
	if n != 0 {
		t.Fatalf("Write() n = %d, want 0", n)
	}
	if got := len(w.Bytes()); got != 0 {
		t.Fatalf("len(Bytes()) = %d, want 0", got)
	}
}
