package turso

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// Client manages communication with the Turso API
type Client struct {
	cfg    *Config
	client HTTPRequestDoer

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// Services
	Organization   organizationService
	Group          groupService
	Database       databaseService
	DatabaseTokens databaseTokensService
}

type service struct {
	client *Client
}

// HTTPRequestDoer implements the standard http.Client interface.
type HTTPRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client returns the http client
func (c *Client) Client() HTTPRequestDoer {
	return c.client
}

// NewClient creates a new client for interacting with the Turso API.
func NewClient(c Config) (*Client, error) {
	if c.Token == "" {
		return nil, ErrAPITokenNotSet
	}

	client := &Client{
		cfg: &c,
	}

	if client.client == nil {
		client.client = http.DefaultClient
	}

	client.common.client = client

	// initialize services
	client.Organization = (*OrganizationService)(&client.common)
	client.Database = (*DatabaseService)(&client.common)
	client.Group = (*GroupService)(&client.common)
	client.DatabaseTokens = (*DatabaseTokensService)(&client.common)

	return client, nil
}

// DoRequest performs an HTTP request and returns the response
func (c *Client) DoRequest(ctx context.Context, method string, url string, data interface{}) (*http.Response, error) {
	var bodyReader io.Reader

	buf, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	bodyReader = bytes.NewReader(buf)

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	// Add Headers
	req.Header.Add("Authorization", "Bearer "+c.cfg.Token)
	req.Header.Add("Content-Type", "application/json")

	return c.client.Do(req.WithContext(ctx))
}
