/*
This file contains tests for the graceful shutdown signal handler. It verifies that when the server receives a shutdown signal,
it stops accepting new requests and shuts down gracefully.
It is a mini-copy of the code in main.go to test the signal handling logic in isolation.
Probably not the cleanest way, one might in the future refactor all the HTTP bells and whistles into a separate module.
However for now this is sufficient to verify that the signal handling logic works as expected :)
*/
package signals

import (
	"context"
	"net"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"
)

// TestGracefulShutdownSignalHandler verifies that a server which shuts down
// in response to a signal stops accepting requests.
func TestGracefulShutdownSignalHandler(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()

	addr := ln.Addr().String()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	srv := &http.Server{Handler: mux}

	serverErrCh := make(chan error, 1)
	go func() { serverErrCh <- srv.Serve(ln) }()

	// wait for readiness
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := http.Get("http://" + addr + "/healthz")
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				break
			}
		}
		time.Sleep(50 * time.Millisecond)
	}

	resp, err := http.Get("http://" + addr + "/healthz")
	if err != nil {
		t.Fatalf("server did not become ready: %v", err)
	}
	_ = resp.Body.Close()

	// create a channel that the handler goroutine will listen on
	sigCh := make(chan os.Signal, 1)
	go func() {
		<-sigCh
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}()

	// simulate signal
	sigCh <- syscall.SIGTERM

	// wait for server to exit
	select {
	case err := <-serverErrCh:
		if err != nil && err != http.ErrServerClosed {
			t.Fatalf("server returned unexpected error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("server did not shut down in time")
	}

	// after shutdown, requests should fail
	_, err = http.Get("http://" + addr + "/healthz")
	if err == nil {
		t.Fatalf("expected connection error after shutdown, got success")
	}
}
