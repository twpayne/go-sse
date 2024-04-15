// Package sse implements a Server Sent Events server.
//
// See https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events.
package sse

import (
	"context"
	"net/http"
)

// A Server is a Server Sent Events Server.
type Server struct {
	connectFunc     ConnectFunc
	establishedFunc EstablishedFunc
	errorFunc       ErrorFunc
}

// A ConnectFunc is called when a client connects. If it returns true then the
// connection is established. If it returns false then it should write the
// response to w.
type ConnectFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) bool

// An EstablishedFunc is called when a connection with the client is
// established. It receives a channel to which it should write events to be sent
// to the client. The function must close ch before returning and terminate when
// r's context is done.
type EstablishedFunc func(ctx context.Context, ch chan<- Event, r *http.Request)

// A ErrorFunc is called when an error is encountered while sending an event. If
// it returns non-nil then the connection to the client is terminated.
type ErrorFunc func(ctx context.Context, err error, r *http.Request) error

// A ServerOption sets an option on a Server.
type ServerOption func(*Server)

// WithConnectFunc sets the connection function.
func WithConnectFunc(connectFunc ConnectFunc) ServerOption {
	return func(s *Server) {
		s.connectFunc = connectFunc
	}
}

// WithEstablishedFunc sets the established function.
func WithEstablishedFunc(establishedFunc EstablishedFunc) ServerOption {
	return func(s *Server) {
		s.establishedFunc = establishedFunc
	}
}

// WithErrorFunc sets the error function.
func WithErrorFunc(errorFunc ErrorFunc) ServerOption {
	return func(s *Server) {
		s.errorFunc = errorFunc
	}
}

// NewServer returns a new Server with the given options.
func NewServer(options ...ServerOption) *Server {
	s := &Server{
		connectFunc:     DefaultConnectFunc,
		establishedFunc: DefaultEstablishedFunc,
		errorFunc:       DefaultErrorFunc,
	}
	for _, option := range options {
		option(s)
	}
	return s
}

// ServeHTTP implements net/http.Handler.ServeHTTP.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if !s.connectFunc(ctx, w, r) {
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	ch := make(chan Event)
	go s.establishedFunc(ctx, ch, r)

	var flush func()
	if flusher, ok := w.(http.Flusher); ok {
		flush = flusher.Flush
	}

	done := ctx.Done()
	for {
		select {
		case <-done:
			return
		case event, ok := <-ch:
			if !ok {
				return
			}
			data, err := event.MarshalText()
			if err == nil {
				_, err = w.Write(data)
				if flush != nil {
					flush()
				}
			}
			if err != nil {
				err = s.errorFunc(ctx, err, r)
			}
			if err != nil {
				return
			}
		}
	}
}

// DefaultConnectFunc is the default connect function. It always returns true.
func DefaultConnectFunc(ctx context.Context, w http.ResponseWriter, r *http.Request) bool {
	return true
}

// DefaultEstablishedFunc is the default established function. It immediately
// closes ch, i.e. it terminates the connection immediately.
func DefaultEstablishedFunc(ctx context.Context, ch chan<- Event, r *http.Request) {
	close(ch)
}

// DefaultErrorFunc is the default error function. It returns err unchanged,
// i.e. it terminates the connection as soon as any error is encountered.
func DefaultErrorFunc(ctx context.Context, err error, r *http.Request) error {
	return err
}
