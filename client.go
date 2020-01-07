package astisparkpost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/asticode/go-astikit"
)

const baseURL = "https://api.eu.sparkpost.com/api"

// Client represents the client
type Client struct {
	c Configuration
	s *astikit.HTTPSender
}

// New creates a new client
func New(c Configuration) *Client {
	return &Client{
		c: c,
		s: astikit.NewHTTPSender(c.Sender),
	}
}

type ErrorPayload struct {
	Errors []Error `json:"errors"`
}

type Error struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Message     string `json:"message"`
}

// Unfortunately there's not a constant way of indicating the key :(
func (c *Client) send(method, url string, reqPayload, respPayload interface{}) (err error) {
	// Create body
	var body io.Reader
	if reqPayload != nil {
		// Marshal
		buf := &bytes.Buffer{}
		if err = json.NewEncoder(buf).Encode(reqPayload); err != nil {
			err = fmt.Errorf("astisparkpost: marshaling payload of %s request to %s failed: %w", method, url, err)
			return
		}

		// Set body
		body = buf
	}

	// Create request
	var req *http.Request
	if req, err = http.NewRequest(method, baseURL+url, body); err != nil {
		err = fmt.Errorf("astisparkpost: creating %s request to %s failed: %w", method, url, err)
		return
	}

	// Add headers
	req.Header.Set("Authorization", c.c.Key)
	req.Header.Set("Content-Type", "application/json")

	// Send
	var resp *http.Response
	if resp, err = c.s.Send(req); err != nil {
		err = fmt.Errorf("astisparkpost: sending %s request to %s failed: %w", req.Method, req.URL.Path, err)
		return
	}
	defer resp.Body.Close()

	// Process error
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		// Unmarshal
		var e ErrorPayload
		if err = json.NewDecoder(resp.Body).Decode(&e); err != nil {
			err = fmt.Errorf("astisparkpost: unmarshaling error failed: %w", err)
			return
		}

		// Set error
		err = fmt.Errorf("astisparkpost: invalid status code %d: %+v", resp.StatusCode, e.Errors)
		return
	}

	// Parse response
	if respPayload != nil {
		// Unmarshal
		if err = json.NewDecoder(resp.Body).Decode(&respPayload); err != nil {
			err = fmt.Errorf("astisparkpost: unmarshaling response failed: %w", err)
			return
		}
	}
	return
}
