package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/twpayne/go-sse"
)

func main() {
	sseServer := sse.NewServer(
		sse.WithConnectFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) bool {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			return true
		}),
		sse.WithEstablishedFunc(func(ctx context.Context, ch chan<- sse.Event, r *http.Request) {
			defer close(ch)
			for i := 0; i < 4; i++ {
				time.Sleep(time.Second)
				ch <- sse.Event{
					Name: fmt.Sprintf("event-%d", i),
					Data: []byte(strconv.Itoa(i)),
				}
			}
		}),
	)

	serverMux := http.NewServeMux()
	serverMux.Handle("GET /events", sseServer)

	http.ListenAndServe(":8080", serverMux)
}
