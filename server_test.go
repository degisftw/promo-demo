package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// The identity endpoint is the entire point of the demo app: if it stops
// reporting version/env, promotions become invisible.
func TestRootReportsVersionAndEnv(t *testing.T) {
	srv := httptest.NewServer(newMux("v1.2.3", "non-prod"))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET / = %d, want 200", resp.StatusCode)
	}
	var got response
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	want := response{App: "promo-demo", Version: "v1.2.3", Env: "non-prod"}
	if got != want {
		t.Fatalf("GET / = %+v, want %+v", got, want)
	}
}

func TestHealthz(t *testing.T) {
	srv := httptest.NewServer(newMux("v1.2.3", "prod"))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /healthz = %d, want 200", resp.StatusCode)
	}
}

func TestUnknownPathNotFound(t *testing.T) {
	srv := httptest.NewServer(newMux("v1.2.3", "prod"))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/nope")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("GET /nope = %d, want 404", resp.StatusCode)
	}
}