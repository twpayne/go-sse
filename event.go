package sse

import (
	"bytes"
	"strconv"
	"time"
)

// An Event is a server-sent event.
type Event struct {
	Name  string
	ID    string
	Data  []byte
	Retry time.Duration
}

// MarshalText implements encoding.TextMarshaler.MarshalText.
func (e *Event) MarshalText() ([]byte, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	if e.Name != "" {
		buffer.Write([]byte("event: "))
		buffer.WriteString(e.Name)
		buffer.WriteByte('\n')
	}
	if e.ID != "" {
		buffer.Write([]byte("id: "))
		buffer.WriteString(e.ID)
		buffer.WriteByte('\n')
	}
	for _, line := range bytes.Split(e.Data, []byte{'\n'}) {
		buffer.Write([]byte("data: "))
		buffer.Write(line)
		buffer.WriteByte('\n')
	}
	if e.Retry > 0 {
		buffer.Write([]byte("retry: "))
		buffer.WriteString(strconv.FormatInt(int64(e.Retry/time.Millisecond), 10))
		buffer.WriteByte('\n')
	}
	buffer.WriteByte('\n')
	return buffer.Bytes(), nil
}

func (e *Event) Write(data []byte) (int, error) {
	e.Data = append(e.Data, data...)
	return len(data), nil
}
