package astisparkpost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	astihttp "github.com/asticode/go-astitools/http"
	"github.com/pkg/errors"
)

const baseURL = "https://api.eu.sparkpost.com/api"

// Client represents the client
type Client struct {
	c Configuration
	s *astihttp.Sender
}

// New creates a new client
func New(c Configuration) *Client {
	return &Client{
		c: c,
		s: astihttp.NewSender(c.Sender),
	}
}

type ErrorPayload struct {
	Errors []Error `json:"errors"`
}

type Error struct {
	Code        int    `json:"code"`
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
			err = errors.Wrapf(err, "astisparkpost: marshaling payload of %s request to %s failed", method, url)
			return
		}

		// Set body
		body = buf
	}

	// Create request
	var req *http.Request
	if req, err = http.NewRequest(method, baseURL+url, body); err != nil {
		err = errors.Wrapf(err, "astisparkpost: creating %s request to %s failed", method, url)
		return
	}

	// Add headers
	req.Header.Set("Authorization", c.c.Key)
	req.Header.Set("Content-Type", "application/json")

	// Send
	var resp *http.Response
	if resp, err = c.s.Send(req); err != nil {
		err = errors.Wrapf(err, "astisparkpost: sending %s request to %s failed", req.Method, req.URL.Path)
		return
	}
	defer resp.Body.Close()

	// Process error
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		// Unmarshal
		var e ErrorPayload
		if err = json.NewDecoder(resp.Body).Decode(&e); err != nil {
			err = errors.Wrap(err, "astisparkpost: unmarshaling error failed")
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
			err = errors.Wrap(err, "astisparkpost: unmarshaling response failed")
			return
		}
	}
	return
}
