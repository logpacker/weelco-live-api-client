package client

import (
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	c := New("https://as.weelco.com/liveapi")

	// Get stream by hash
	s, e := c.GetStream("75AB8BF9947A9E7C23A2B026CF623581")
	fmt.Printf("Stream: %v, Error: %v\n", s, e)

	// Update stream's info
	s.OutputURL = fmt.Sprintf("http://127.0.0.1:1935/%d", time.Now().Unix())
	s.OutputNodeAddr = fmt.Sprintf("http://127.0.0.1:1936/%d", time.Now().Unix())
	e = c.UpdateStream(s)
	fmt.Printf("Stream: %v, Error: %v\n", s, e)
}
