package client

import (
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	c := New("http://127.0.0.1:8081")

	// Get stream by hash
	s, e := c.GetStream("6BF8AA3990AEAC71D3CD0476095F08F6")
	fmt.Printf("Stream: %v, Error: %v\n", s, e)

	// Update stream's info
	s.URL = fmt.Sprintf("http://127.0.0.1:1935/%d", time.Now().Unix())
	e = c.UpdateStream(s)
	fmt.Printf("Stream: %v, Error: %v\n", s, e)
}
