package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Client struct
type Client struct {
	Addr string
}

// Stream struct
type Stream struct {
	ID                 int64      `json:"id"`
	Hash               string     `json:"hash"`
	Name               string     `json:"name"`
	Status             int64      `json:"status"`
	StatusString       string     `json:"status_str"`
	OwnerID            uint64     `json:"owner_id"`
	OutputURL          string     `json:"output_url"`
	AdvertisedHost     string     `json:"advertised_host"`
	ControlQueue       string     `json:"control_queue"`
	StreamingInputAddr string     `json:"stream_input_addr"`
	StartTime          *time.Time `json:"start_time"`
	StopTime           *time.Time `json:"stop_time"`
	CurrentWatchers    uint64     `json:"current_watchers"`
}

type errorResponse struct {
	Message string `json:"message"`
}

// New initializes client
func New(addr string) *Client {
	c := new(Client)
	c.Addr = addr
	return c
}

// CreateStream creates Stream, only Name and OwnerID is used here
func (c *Client) CreateStream(s *Stream) error {
	if s == nil {
		return fmt.Errorf("Stream is empty")
	}
	return c.api(fmt.Sprintf("/streams/new?name=%s&owner_id=%d", s.Name, s.OwnerID), "POST", "", nil)
}

// GetStreams returns all streams
func (c *Client) GetStreams() ([]*Stream, error) {
	s := []*Stream{}
	err := c.api("/streams", "GET", "", &s)
	return s, err
}

// GetStream returns stream by hash
func (c *Client) GetStream(hash string) (*Stream, error) {
	s := new(Stream)
	err := c.api(fmt.Sprintf("/streams/find?hash=%s", hash), "GET", "", &s)
	return s, err
}

// Start updates status to Running and sets OutputURL, AdvertisedHost, ControlQueue, StartTime
func (c *Client) Start(hash string, outputURL string, advertisedHost string, controlQueue string) error {
	streamBytes, _ := json.Marshal(Stream{
		OutputURL:      outputURL,
		AdvertisedHost: advertisedHost,
		ControlQueue:   controlQueue,
	})
	return c.api(fmt.Sprintf("/streams/start?hash=%s", hash), "POST", string(streamBytes), nil)
}

// Stop updates status to Stopped and sets StopTime
func (c *Client) Stop(hash string) error {
	return c.api(fmt.Sprintf("/streams/stop?hash=%s", hash), "POST", "", nil)
}

// StopTranscoding sends request to transcoder to stop stream. Transcoder will update status after that
func (c *Client) StopTranscoding(hash string) error {
	return c.api(fmt.Sprintf("/streams/transcoder/stop?hash=%s", hash), "POST", "", nil)
}

func (c *Client) api(endpoint string, method string, dataStr string, v interface{}) error {
	// Build request
	u := strings.TrimRight(c.Addr, "/") + "/" + strings.TrimLeft(endpoint, "/")
	request, err := http.NewRequest(method, u, bytes.NewBufferString(dataStr))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Content-Length", fmt.Sprintf("%d", len(dataStr)))
	if err != nil {
		return err
	}

	// Do request and get raw body
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		errResp := errorResponse{}
		err = json.Unmarshal(body, &errResp)
		if err == nil {
			return fmt.Errorf("%s %s. Code: %d. Details: %s", method, u, response.StatusCode, errResp.Message)
		}
	}

	if v != nil {
		err = json.Unmarshal(body, v)
		if err != nil {
			return fmt.Errorf("Unable to parse response from %s. Details: %s", u, err.Error())
		}
	}

	return nil
}
