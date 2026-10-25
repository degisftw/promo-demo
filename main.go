// Command promo-demo is a deliberately small HTTP app used to exercise the
// Kargo promotion pipeline in talos-gitops: it reports its build version and
// the environment name it was deployed to, so a promotion is visible by
// curling the two namespaces.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// version is stamped at build time via -ldflags "-X main.version=...".
var version = "dev"

type response struct {
	App     string `json:"app"`
	Version string `json:"version"`
	Env     string `json:"env"`
}

func main() {
	addr := ":8080"
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "unknown"
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           newMux(version, env),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("promo-demo %s (env=%s) listening on %s", version, env, addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}