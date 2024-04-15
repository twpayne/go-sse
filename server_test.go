package sse_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-sse"
)

func TestOneEvent(t *testing.T) {
	server := newTestSSEServer(
		sse.WithEstablishedFunc(func(ctx context.Context, ch chan<- sse.Event, r *http.Request) {
			defer close(ch)
			event := sse.Event{
				Name:  "eventName",
				ID:    "eventID",
				Data:  []byte("eventData"),
				Retry: 1 * time.Second,
			}
			select {
			case <-ctx.Done():
				return
			case ch <- event:
				return
			}
		}),
	)

	resp, err := server.Client().Get(server.URL + "/events")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/event-stream", resp.Header.Get("Content-Type"))
	assert.Equal(t, "no-cache", resp.Header.Get("Cache-Control"))
	assert.Equal(t, "keep-alive", resp.Header.Get("Connection"))
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, joinLines(
		"event: eventName",
		"id: eventID",
		"data: eventData",
		"retry: 1000",
		"",
	), body)
}

func joinLines(lines ...string) []byte {
	return []byte(strings.Join(lines, "\n") + "\n")
}

func newTestSSEServer(options ...sse.ServerOption) *httptest.Server {
	sseServer := sse.NewServer(options...)
	serveMux := http.NewServeMux()
	serveMux.Handle("GET /events", sseServer)
	return httptest.NewServer(serveMux)
}
