package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"testing"
)

func TestExpectedStatusCodeMatched(t *testing.T) {
	opt := Opt{Expect: "HTTP/1.1 200,HTTP/2.0 200"}

	got := opt.ExpectedStatusCode("HTTP/2.0 200 OK")
	if got != "HTTP/2.0 200" {
		t.Fatalf("expectedStatusCode() = %q, want %q", got, "HTTP/2.0 200")
	}
}

func TestExpectedStatusCodeNoMatch(t *testing.T) {
	opt := Opt{Expect: "HTTP/1.1 200,HTTP/2.0 200"}

	got := opt.ExpectedStatusCode("HTTP/1.1 500 Internal Server Error")
	if got != "" {
		t.Fatalf("expectedStatusCode() = %q, want empty", got)
	}
}

func TestExpectedStatusCodeReturnsFirstMatch(t *testing.T) {
	opt := Opt{Expect: "HTTP/,HTTP/2.0 200"}

	got := opt.ExpectedStatusCode("HTTP/2.0 200 OK")
	if got != "HTTP/" {
		t.Fatalf("expectedStatusCode() = %q, want %q", got, "HTTP/")
	}
}

// BuildRequest tests
func TestBuildRequest(t *testing.T) {
	opt := Opt{
		Hostname:      "example.com",
		URI:           "/path",
		Method:        "POST",
		UserAgent:     "test-agent",
		Authorization: "user:pass",
		Port:          80,
	}

	ctx := context.Background()
	req, err := opt.BuildRequest(ctx)
	if err != nil {
		t.Fatalf("BuildRequest() error = %v", err)
	}

	if req.Method != "POST" {
		t.Fatalf("BuildRequest() Method = %q, want %q", req.Method, "POST")
	}

	if req.URL.String() != "http://example.com/path" {
		t.Fatalf("BuildRequest() URL = %q, want %q", req.URL.String(), "http://example.com/path")
	}

	if req.Header.Get("User-Agent") != "test-agent" {
		t.Fatalf("BuildRequest() User-Agent = %q, want %q", req.Header.Get("User-Agent"), "test-agent")
	}

	username, password, ok := req.BasicAuth()
	if !ok || username != "user" || password != "pass" {
		t.Fatalf("BuildRequest() BasicAuth = (%q, %q), want (%q, %q)", username, password, "user", "pass")
	}
}

// BuildRequestTest with SSL and SNI
func TestBuildRequestWithSSLAndSNI(t *testing.T) {
	opt := Opt{
		Hostname:  "example.com",
		URI:       "/path",
		Method:    "GET",
		UserAgent: "test-agent",
		SSL:       true,
		SNI:       true,
		Port:      443,
	}

	ctx := context.Background()
	req, err := opt.BuildRequest(ctx)
	if err != nil {
		t.Fatalf("BuildRequest() error = %v", err)
	}

	if req.URL.String() != "https://example.com/path" {
		t.Fatalf("BuildRequest() URL = %q, want %q", req.URL.String(), "https://example.com/path")
	}
}

// BuildRequestTest with invalid authorization
func TestBuildRequestWithInvalidAuthorization(t *testing.T) {
	opt := Opt{
		Hostname:      "example.com",
		URI:           "/path",
		Method:        "GET",
		UserAgent:     "test-agent",
		Authorization: "invalid-authorization",
	}

	ctx := context.Background()
	_, err := opt.BuildRequest(ctx)
	if err == nil {
		t.Fatalf("BuildRequest() error = nil, want non-nil")
	}
	if err.Error() != "invalid authorization args" {
		t.Fatalf("BuildRequest() error = %q, want %q", err.Error(), "invalid authorization args")
	}
}

// MakeTransport tests
func TestMakeTransport(t *testing.T) {
	opt := Opt{
		SSL:      true,
		SNI:      true,
		Hostname: "example.com:443",
	}

	tripper := opt.MakeTransport()
	if tripper == nil {
		t.Fatalf("MakeTransport() returned nil, want non-nil")
	}
	transport, ok := tripper.(*http.Transport)
	if !ok {
		t.Fatalf("MakeTransport() returned non-http.Transport, want http.Transport")
	}
	if transport.TLSClientConfig.ServerName != "example.com" {
		t.Fatalf("MakeTransport() TLSClientConfig.ServerName = %q, want %q", transport.TLSClientConfig.ServerName, "example.com")
	}
}

// MakeTransport tests with TLSMaxVersion
func TestMakeTransportWithTLSMaxVersion(t *testing.T) {
	{
		opt := Opt{
			SSL:           true,
			TLSMaxVersion: "1.2",
		}
		tripper := opt.MakeTransport()
		transport, ok := tripper.(*http.Transport)
		if !ok {
			t.Fatalf("MakeTransport() returned non-http.Transport, want http.Transport")
		}
		if transport == nil {
			t.Fatalf("MakeTransport() returned nil, want non-nil")
		}
		if transport.TLSClientConfig.MaxVersion != tls.VersionTLS12 {
			t.Fatalf("MakeTransport() TLSClientConfig.MaxVersion = %d, want %d", transport.TLSClientConfig.MaxVersion, tls.VersionTLS12)
		}
	}
	{
		opt := Opt{
			SSL:           true,
			TLSMaxVersion: "1.3",
		}
		tripper := opt.MakeTransport()
		transport, ok := tripper.(*http.Transport)
		if !ok {
			t.Fatalf("MakeTransport() returned non-http.Transport, want http.Transport")
		}
		if transport == nil {
			t.Fatalf("MakeTransport() returned nil, want non-nil")
		}
		if transport.TLSClientConfig.MaxVersion != tls.VersionTLS13 {
			t.Fatalf("MakeTransport() TLSClientConfig.MaxVersion = %d, want %d", transport.TLSClientConfig.MaxVersion, tls.VersionTLS13)
		}
	}
	{
		opt := Opt{
			SSL:           true,
			TLSMaxVersion: "1.1",
		}
		tripper := opt.MakeTransport()
		transport, ok := tripper.(*http.Transport)
		if !ok {
			t.Fatalf("MakeTransport() returned non-http.Transport, want http.Transport")
		}
		if transport == nil {
			t.Fatalf("MakeTransport() returned nil, want non-nil")
		}
		if transport.TLSClientConfig.MaxVersion != tls.VersionTLS11 {
			t.Fatalf("MakeTransport() TLSClientConfig.MaxVersion = %d, want %d", transport.TLSClientConfig.MaxVersion, tls.VersionTLS11)
		}
		if transport.TLSClientConfig.MinVersion != tls.VersionTLS11 {
			t.Fatalf("MakeTransport() TLSClientConfig.MinVersion = %d, want %d", transport.TLSClientConfig.MinVersion, tls.VersionTLS11)
		}
	}
}
