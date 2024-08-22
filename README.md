# go-sse

[![PkgGoDev](https://pkg.go.dev/badge/github.com/twpayne/go-sse)](https://pkg.go.dev/github.com/twpayne/go-sse)

Package sse implements a [Server-Sent
Events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events)
server.

In short, it makes sending events to a client just like writing to a Go channel.

## Example

```go
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
```

## License

MIT
