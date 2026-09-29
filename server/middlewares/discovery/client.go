package discovery

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

// DefaultTimeout is the default timeout for the Perplexity API.
const DefaultTimeout = 30 * time.Second

// DefaultClient is the default HTTP client for the Perplexity API.
var DefaultClient = &http.Client{
	Timeout: DefaultTimeout,
}

// Client is an interface for discovering services.
type Client interface {
	// Discover returns the well-known configuration for the given base URL.
	Discover(ctx context.Context, baseURL string) (*WellKnownConfig, error)
}

type discoveryClient struct {
	client *http.Client
}

// Option is a function that configures the discovery client.
type Option func(*discoveryClient)

// NewClient returns a new discovery client.
func NewClient(opts ...Option) Client {
	c := &discoveryClient{
		client: DefaultClient,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Discover returns the well-known configuration for the given base URL.
func (d *discoveryClient) Discover(ctx context.Context, baseURL string) (*WellKnownConfig, error) {
	u, err := url.JoinPath(baseURL, WellKnownConfigurationURL)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	config := &WellKnownConfig{}
	if err := json.NewDecoder(resp.Body).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
