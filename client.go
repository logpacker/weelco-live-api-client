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
	ID             int64  `json:"id"`
	Hash           string `json:"hash"`
	Name           string `json:"name"`
	OwnerID        uint64 `json:"owner_id"`
	OutputURL      string `json:"output_url"`
	OutputNodeAddr string `json:"output_node_addr"`
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

// UpdateStream updates Stream
func (c *Client) UpdateStream(s *Stream) error {
	streamBytes, _ := json.Marshal(s)
	return c.api(fmt.Sprintf("/streams/update?hash=%s", s.Hash), "POST", string(streamBytes), nil)
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
