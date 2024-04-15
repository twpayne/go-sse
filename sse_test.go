package sse_test

import "strings"

func joinLines(lines ...string) string {
	return strings.Join(lines, "\n") + "\n"
}
