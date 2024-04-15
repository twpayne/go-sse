package sse_test

import (
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-sse"
)

func TestEvent_MarshalText(t *testing.T) {
	for _, tc := range []struct {
		name     string
		event    sse.Event
		expected string
	}{
		{
			name:     "empty",
			expected: "data: \n\n",
		},
		{
			name: "id",
			event: sse.Event{
				Name: "eventName",
				ID:   "eventID",
			},
			expected: joinLines(
				"event: eventName",
				"id: eventID",
				"data: ",
				"",
			),
		},
		{
			name: "data",
			event: sse.Event{
				Name: "eventName",
				Data: []byte("eventData"),
			},
			expected: joinLines(
				"event: eventName",
				"data: eventData",
				"",
			),
		},
		{
			name: "multiline_data",
			event: sse.Event{
				Name: "eventName",
				Data: []byte("eventDataLine1\neventDataLine2"),
			},
			expected: joinLines(
				"event: eventName",
				"data: eventDataLine1",
				"data: eventDataLine2",
				"",
			),
		},
		{
			name: "multiline_data_with_trailing_newline",
			event: sse.Event{
				Name: "eventName",
				Data: []byte("eventDataLine1\neventDataLine2\n"),
			},
			expected: joinLines(
				"event: eventName",
				"data: eventDataLine1",
				"data: eventDataLine2",
				"data: ",
				"",
			),
		},
		{
			name: "retry_1s",
			event: sse.Event{
				Name:  "eventName",
				Retry: 1 * time.Second,
			},
			expected: joinLines(
				"event: eventName",
				"data: ",
				"retry: 1000",
				"",
			),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := tc.event.MarshalText()
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, string(actual))
		})
	}
}

func TestEvent_Write(t *testing.T) {
	var sseEvent sse.Event
	_, err := sseEvent.Write([]byte("eventData"))
	assert.NoError(t, err)
	assert.Equal(t, string("eventData"), string(sseEvent.Data))
}
