package main

import (
	"encoding/json"
	"net/http"
)

// newMux serves GET / with the identity of this build and GET /healthz for
// probes. Both endpoints exist so a promotion is observable (version) and
// rollouts are gated (health).
func newMux(version, env string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response{App: "promo-demo", Version: version, Env: env})
	})
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return mux
}