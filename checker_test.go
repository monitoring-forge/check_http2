package main

import "testing"

func TestExpectedStatusCodeMatched(t *testing.T) {
	opts := commandOpts{Expect: "HTTP/1.1 200,HTTP/2.0 200"}

	got := expectedStatusCode(opts, "HTTP/2.0 200 OK")
	if got != "HTTP/2.0 200" {
		t.Fatalf("expectedStatusCode() = %q, want %q", got, "HTTP/2.0 200")
	}
}

func TestExpectedStatusCodeNoMatch(t *testing.T) {
	opts := commandOpts{Expect: "HTTP/1.1 200,HTTP/2.0 200"}

	got := expectedStatusCode(opts, "HTTP/1.1 500 Internal Server Error")
	if got != "" {
		t.Fatalf("expectedStatusCode() = %q, want empty", got)
	}
}

func TestExpectedStatusCodeReturnsFirstMatch(t *testing.T) {
	opts := commandOpts{Expect: "HTTP/,HTTP/2.0 200"}

	got := expectedStatusCode(opts, "HTTP/2.0 200 OK")
	if got != "HTTP/" {
		t.Fatalf("expectedStatusCode() = %q, want %q", got, "HTTP/")
	}
}
