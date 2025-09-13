package hostman

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const exampleResponse = `{
	"hosts": {
		"a.example.com": "127.0.0.1",
		"b.example.com": "127.0.0.1",
		"c.example.com": "127.0.0.2"
	}
}
`

func TestHTTPSource_GetFromRemote_Success(t *testing.T) {
	//
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(exampleResponse))
	}))
	defer ts.Close()

	src := &HTTPSource{Endpoint: ts.URL}

	// Act
	cfg, err := src.GetFromRemote()

	// Assert
	if err != nil {
		t.Fatalf("GetFromRemote unexpected error: %v", err)
	}

	if cfg == nil || cfg.Hosts == nil {
		t.Fatalf("GetFromRemote returned nil or missing hosts: %#v", cfg)
	}
	if got := cfg.Hosts["c.example.com"]; got != "127.0.0.2" {
		t.Fatalf("unexpected host mapping: got %q want %q", got, "127.0.0.2")
	}

	m, err := src.GetMapping()
	if err != nil {
		t.Fatalf("GetMapping unexpected error: %v", err)
	}
	if m["c.example.com"] != "127.0.0.2" {
		t.Fatalf("GetMapping value mismatch: got %q want %q", m["c.example.com"], "127.0.0.2")
	}
}

func TestHTTPSource_GetFromRemote_HTTPError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	defer ts.Close()

	src := &HTTPSource{Endpoint: ts.URL}
	if _, err := src.GetFromRemote(); err == nil {
		t.Fatalf("expected error for non-2xx status")
	}
}

func TestHTTPSource_GetFromRemote_InvalidJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{ not-json }"))
	}))
	defer ts.Close()

	src := &HTTPSource{Endpoint: ts.URL}
	if _, err := src.GetFromRemote(); err == nil {
		t.Fatalf("expected error for invalid JSON")
	}
}
