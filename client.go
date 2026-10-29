package ipclues

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const DefaultBaseURL = "https://api.ipclues.com/v1"

// NewClient creates an IPClues client. An empty API key is allowed for tests.
func NewClient(apiKey string, options ...ClientOption) *Client {
	c := &Client{apiKey: apiKey, baseURL: DefaultBaseURL, httpClient: &http.Client{Timeout: 10 * time.Second}}
	for _, option := range options {
		if option != nil {
			option(c)
		}
	}
	return c
}

// WithBaseURL overrides the API base URL, useful for local testing.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		if strings.TrimSpace(baseURL) != "" {
			c.baseURL = strings.TrimRight(baseURL, "/")
		}
	}
}

// WithHTTPClient supplies a custom HTTP client.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		if client != nil {
			c.httpClient = client
		}
	}
}

// WithTimeout sets the timeout on the default HTTP client. It has no effect
// when a custom client is supplied with WithHTTPClient.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		if timeout > 0 {
			c.httpClient.Timeout = timeout
		}
	}
}

// Client calls the IPClues API and is safe for concurrent use.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// Lookup fetches IP intelligence. HTTP error statuses (401, 402, 429, etc.)
// are returned in APIResponse with a nil error. Network and transport failures
// return an error.
func (c *Client) Lookup(ctx context.Context, ip string, options *LookupOptions) (*Result, error) {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return nil, ErrEmptyIPAddress
	}

	u, err := url.Parse(strings.Join([]string{c.baseURL, "lookup/ip", url.PathEscape(ip)}, "/"))
	if err != nil {
		return nil, fmt.Errorf("ipclues: build request URL: %w", err)
	}
	if options != nil && len(options.Databases) > 0 {
		dbs := make([]string, len(options.Databases))
		for _, db := range options.Databases {
			if !db.Valid() {
				return nil, fmt.Errorf("ipclues: unsupported database %q", db)
			}
			dbs = append(dbs, string(db))
		}
		q := u.Query()
		q.Set("db", strings.Join(dbs, ","))
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("ipclues: create request: %w", err)
	}
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", GetUserAgent())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return nil, fmt.Errorf("ipclues: request timeout: %w", err)
		}
		return nil, fmt.Errorf("ipclues: request failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ipclues: read response body: %w", err)
	}

	var data any
	if len(body) > 0 {
		if err := json.Unmarshal(body, &data); err != nil {
			data = string(body)
		}
	}
	return &Result{
		Data:     data,
		Metadata: parseMetadata(resp.StatusCode, http.StatusText(resp.StatusCode), resp.Header),
	}, nil
}

// ClientOption configures a Client.
type ClientOption func(*Client)

// LookupOptions controls an IP lookup.
type LookupOptions struct {
	// Databases optionally selects database IDs. An empty slice omits the db parameter.
	Databases []Database
}
