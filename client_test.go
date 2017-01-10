package client

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	c := New("https://as.weelco.com/liveapi")

	// Get stream by hash
	s, e := c.GetStream("75AB8BF9947A9E7C23A2B026CF623581")
	fmt.Printf("Stream: %v, Error: %v\n", s, e)

	// Update stream's info
	s.OutputURL = "http://sample.vodobox.net/skate_phantom_flex_4k/skate_phantom_flex_4k.m3u8"
	s.OutputNodeAddr = "localhost"
	s.Status = 1
	e = c.UpdateStream(s)
	fmt.Printf("Stream: %v, Error: %v\n", s, e)
}
