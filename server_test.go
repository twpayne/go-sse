package sse_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
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
	assert.Equal(t, "no-cache", resp.Header.Get("Cache-Control"))
	assert.Equal(t, "keep-alive", resp.Header.Get("Connection"))
	assert.Equal(t, "text/event-stream", resp.Header.Get("Content-Type"))
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, joinLines(
		"event: eventName",
		"id: eventID",
		"data: eventData",
		"retry: 1000",
		"",
	), string(body))
}

func TestConnectFunc(t *testing.T) {
	server := newTestSSEServer(
		sse.WithConnectFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) bool {
			statusCode, _ := strconv.Atoi(r.URL.Query().Get("statusCode"))
			w.WriteHeader(statusCode)
			result, _ := strconv.ParseBool(r.URL.Query().Get("result"))
			return result
		}),
	)

	{
		resp, err := server.Client().Get(server.URL + "/events?statusCode=200&result=true")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}

	{
		resp, err := server.Client().Get(server.URL + "/events?statusCode=429&result=false")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
	}
}

func newTestSSEServer(options ...sse.ServerOption) *httptest.Server {
	sseServer := sse.NewServer(options...)
	serveMux := http.NewServeMux()
	serveMux.Handle("GET /events", sseServer)
	return httptest.NewServer(serveMux)
}
