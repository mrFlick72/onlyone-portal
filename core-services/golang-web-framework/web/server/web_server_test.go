package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// runServerUntilShutdown mirrors the lifecycle in StartEngine: spin up the
// server in a goroutine, block until the trigger context is cancelled (the
// test's stand-in for SIGTERM), then drain via srv.Shutdown.
func runServerUntilShutdown(t *testing.T, srv *http.Server, trigger context.Context, shutdownTimeout time.Duration) error {
	t.Helper()

	serverErr := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()

	select {
	case err := <-serverErr:
		return err
	case <-trigger.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}

// freePort grabs an OS-assigned port and immediately releases it; the returned
// port is racy in theory but fine for serial tests.
func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := l.Addr().(*net.TCPAddr).Port
	require.NoError(t, l.Close())
	return port
}

// Proves the shutdown signal triggers a graceful drain: an in-flight slow
// request started before shutdown completes successfully (200 OK) rather than
// being aborted, even though shutdown was triggered while it was running.
func TestStartEngine_GracefulShutdownDrainsInFlightRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	requestStarted := make(chan struct{})
	router := gin.New()
	router.GET("/slow", func(c *gin.Context) {
		close(requestStarted)
		time.Sleep(300 * time.Millisecond)
		c.String(http.StatusOK, "drained")
	})

	port := freePort(t)
	srv := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Handler: router,
	}

	trigger, cancelTrigger := context.WithCancel(context.Background())
	defer cancelTrigger()

	done := make(chan error, 1)
	go func() {
		done <- runServerUntilShutdown(t, srv, trigger, 5*time.Second)
	}()

	// Wait for the server to bind before issuing the request.
	require.Eventually(t, func() bool {
		conn, err := net.DialTimeout("tcp", srv.Addr, 50*time.Millisecond)
		if err != nil {
			return false
		}
		conn.Close()
		return true
	}, 2*time.Second, 10*time.Millisecond)

	// In-flight slow request.
	respCh := make(chan *http.Response, 1)
	errCh := make(chan error, 1)
	go func() {
		resp, err := http.Get(fmt.Sprintf("http://%s/slow", srv.Addr))
		if err != nil {
			errCh <- err
			return
		}
		respCh <- resp
	}()

	// Wait until the handler has started, then trigger shutdown.
	<-requestStarted
	cancelTrigger()

	select {
	case resp := <-respCh:
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode, "in-flight request should have drained")
		assert.Equal(t, "drained", string(body))
	case err := <-errCh:
		t.Fatalf("in-flight request was aborted by shutdown: %v", err)
	case <-time.After(3 * time.Second):
		t.Fatal("in-flight request never returned")
	}

	select {
	case err := <-done:
		assert.NoError(t, err, "graceful shutdown should not error")
	case <-time.After(2 * time.Second):
		t.Fatal("server did not finish shutdown")
	}
}

// Proves a request that arrives AFTER shutdown is initiated is rejected
// (server has stopped accepting new connections).
func TestStartEngine_ShutdownStopsAcceptingNewConnections(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	port := freePort(t)
	srv := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Handler: router,
	}

	trigger, cancelTrigger := context.WithCancel(context.Background())
	defer cancelTrigger()

	done := make(chan error, 1)
	go func() {
		done <- runServerUntilShutdown(t, srv, trigger, 2*time.Second)
	}()

	require.Eventually(t, func() bool {
		conn, err := net.DialTimeout("tcp", srv.Addr, 50*time.Millisecond)
		if err != nil {
			return false
		}
		conn.Close()
		return true
	}, 2*time.Second, 10*time.Millisecond)

	cancelTrigger()

	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(3 * time.Second):
		t.Fatal("server did not finish shutdown")
	}

	_, err := http.Get(fmt.Sprintf("http://%s/ping", srv.Addr))
	assert.Error(t, err, "post-shutdown request should fail")
}

// Proves that when a slow handler exceeds the shutdown timeout, srv.Shutdown
// returns context.DeadlineExceeded — the framework propagates it as an error
// rather than silently swallowing.
func TestStartEngine_ShutdownTimeoutExpires(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var requestsStarted atomic.Int32
	requestStarted := make(chan struct{}, 1)
	router := gin.New()
	router.GET("/very-slow", func(c *gin.Context) {
		if requestsStarted.Add(1) == 1 {
			requestStarted <- struct{}{}
		}
		time.Sleep(2 * time.Second)
		c.String(http.StatusOK, "late")
	})

	port := freePort(t)
	srv := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Handler: router,
	}

	trigger, cancelTrigger := context.WithCancel(context.Background())
	defer cancelTrigger()

	done := make(chan error, 1)
	go func() {
		done <- runServerUntilShutdown(t, srv, trigger, 100*time.Millisecond)
	}()

	require.Eventually(t, func() bool {
		conn, err := net.DialTimeout("tcp", srv.Addr, 50*time.Millisecond)
		if err != nil {
			return false
		}
		conn.Close()
		return true
	}, 2*time.Second, 10*time.Millisecond)

	go func() {
		http.Get(fmt.Sprintf("http://%s/very-slow", srv.Addr))
	}()

	<-requestStarted
	cancelTrigger()

	select {
	case err := <-done:
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	case <-time.After(3 * time.Second):
		t.Fatal("shutdown never returned")
	}
}
