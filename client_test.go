package client

import "testing"

func TestGetStreams(t *testing.T) {
	c := New("https://as.weelco.com/liveapi")

	_, e := c.GetStreams()
	if e != nil {
		t.Errorf("Unable to get streams. %s", e.Error())
	}
}

func TestGetStream(t *testing.T) {
	c := New("https://as.weelco.com/liveapi")

	_, e := c.GetStream("invalidhash")
	if e == nil {
		t.Errorf("Expected error")
	}
}
