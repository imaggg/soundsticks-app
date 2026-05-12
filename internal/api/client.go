package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	http    *http.Client
	baseURL string
}

// NewClient creates an mTLS HTTPS client for the SoundSticks API.
// HTTP/2 is explicitly disabled — the device only supports HTTP/1.1.
func NewClient(tlsCert tls.Certificate, ip string) *Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{tlsCert},
			InsecureSkipVerify: true,
			NextProtos:         []string{"http/1.1"}, // device rejects h2
		},
		// Disable Go's automatic HTTP/2 upgrade
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}
	return &Client{
		http:    &http.Client{Transport: transport, Timeout: 10 * time.Second},
		baseURL: fmt.Sprintf("https://%s", ip),
	}
}

// GetRaw calls a GET command and returns the decoded JSON as a raw map.
func (c *Client) GetRaw(command string) (map[string]interface{}, error) {
	var out map[string]interface{}
	return out, c.get(command, &out)
}

func (c *Client) get(command string, result interface{}) error {
	u := fmt.Sprintf("%s/httpapi.asp?command=%s", c.baseURL, url.QueryEscape(command))
	resp, err := c.http.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if result == nil {
		return nil
	}
	return json.Unmarshal(body, result)
}

// rawGet calls a GET command without URL-encoding the command string.
// Use for commands that contain literal colons, e.g. setPlayerCmd:vol:70.
func (c *Client) rawGet(command string) error {
	u := fmt.Sprintf("%s/httpapi.asp?command=%s", c.baseURL, command)
	resp, err := c.http.Get(u)
	if err != nil {
		return err
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	return nil
}

// post sends command + payload using application/x-www-form-urlencoded,
// matching the format the device actually accepts:
//
//	command=setLightInfo&payload={"brightness":"80"}
func (c *Client) post(command string, payload interface{}) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	// Device expects raw (unencoded) JSON after "payload=" — same as Python sample
	body := fmt.Sprintf("command=%s&payload=%s", command, string(payloadJSON))
	req, err := http.NewRequest("POST", c.baseURL+"/httpapi.asp", strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	return nil
}
