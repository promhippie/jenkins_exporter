package jenkins

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

const (
	// UserAgent defines the used user ganet for request.
	UserAgent = "go-jenkins/" + Version
)

// Client is a client for the Jenkins API.
type Client struct {
	httpClient *http.Client
	httpDumper Dumper
	endpoint   string
	username   string
	password   string

	Job JobClient
}

// A ClientOption is used to configure a Client.
type ClientOption func(*Client) error

// WithHTTPClient configures a Client to use the specified HTTP client.
func WithHTTPClient(value *http.Client) ClientOption {
	return func(client *Client) error {
		client.httpClient = value
		return nil
	}
}

// WithHTTPDumper configures a Client to use the specified debug dumper.
func WithHTTPDumper(value Dumper) ClientOption {
	return func(client *Client) error {
		client.httpDumper = value
		return nil
	}
}

// WithEndpoint configures a Client to use the specified API endpoint.
func WithEndpoint(endpoint string) ClientOption {
	return func(client *Client) error {
		client.endpoint = strings.TrimRight(endpoint, "/")
		return nil
	}
}

// WithUsername configures a Client to use the specified username for authentication.
func WithUsername(username string) ClientOption {
	return func(client *Client) error {
		client.username = username
		return nil
	}
}

// WithPassword configures a Client to use the specified password for authentication.
func WithPassword(password string) ClientOption {
	return func(client *Client) error {
		client.password = password
		return nil
	}
}

// NewClient creates a new client.
func NewClient(options ...ClientOption) (*Client, error) {
	client := &Client{
		// httpDumper: StandardDumper(true),
	}

	for _, option := range options {
		if err := option(client); err != nil {
			return nil, err
		}
	}

	if client.httpClient == nil {
		pool, err := x509.SystemCertPool()

		if err != nil {
			return nil, err
		}

		client.httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				TLSClientConfig: &tls.Config{
					RootCAs: pool,
				},
			},
		}
	}

	client.Job = JobClient{client: client}

	return client, nil
}

// NewRequest creates an HTTP request against the Jenkins API.
func (c *Client) NewRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, path, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)

	if c.username != "" && c.password != "" {
		req.SetBasicAuth(
			c.username,
			c.password,
		)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req.WithContext(ctx), nil
}

// Do performs an HTTP request against the Jenkins API.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	if c.httpDumper != nil {
		c.httpDumper.DumpRequest(req)
	}

	res, err := c.httpClient.Do(req)

	if res != nil {
		defer func() { _ = res.Body.Close() }()
	}

	if err != nil {
		return nil, err
	}

	if c.httpDumper != nil {
		c.httpDumper.DumpResponse(res)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return &Response{Response: res}, err
	}

	res.Body = io.NopCloser(bytes.NewReader(body))

	if res.StatusCode >= 400 && res.StatusCode <= 599 {
		return &Response{Response: res}, errors.New(http.StatusText(res.StatusCode))
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, bytes.NewReader(body))
		} else {
			err = json.Unmarshal(body, v)
		}
	}

	return &Response{Response: res}, err
}

// Response simply wraps the standard response type.
type Response struct {
	*http.Response
}
