// Package engine handles HTTP request execution and response processing.
package engine

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/EdgarOrtegaRamirez/httpforge/models"
)

// Client wraps http.Client with HttpForge-specific functionality.
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new HttpForge HTTP client.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// Allow up to 10 redirects
				if len(via) >= 10 {
					return fmt.Errorf("too many redirects")
				}
				return nil
			},
		},
	}
}

// Execute sends an HTTP request and returns the response.
func (c *Client) Execute(req *models.Request) (*models.Response, error) {
	// Build the full URL with query params
	fullURL, err := buildURL(req.URL, req.QueryParams)
	if err != nil {
		return nil, fmt.Errorf("building URL: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest(string(req.Method), fullURL, strings.NewReader(req.Body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Set headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Set content type based on body_type
	if req.Body != "" && req.BodyType != "" {
		switch req.BodyType {
		case "json":
			httpReq.Header.Set("Content-Type", "application/json")
		case "form":
			httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		case "xml":
			httpReq.Header.Set("Content-Type", "application/xml")
		}
	}

	// Set timeout if specified
	if req.Timeout > 0 {
		c.httpClient.Timeout = req.Timeout
	}

	// Execute request
	start := time.Now()
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	timing := time.Since(start)

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	// Build response headers map
	headers := make(map[string]string)
	for key, values := range resp.Header {
		headers[key] = strings.Join(values, ", ")
	}

	return &models.Response{
		StatusCode: resp.StatusCode,
		Headers:    headers,
		Body:       string(bodyBytes),
		Timing:     timing,
		Size:       len(bodyBytes),
		Request:    req,
	}, nil
}

// buildURL constructs the full URL with query parameters.
func buildURL(baseURL string, params map[string]string) (string, error) {
	if len(params) == 0 {
		return baseURL, nil
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// ExecuteChain executes a chain of requests, passing variables between them.
func ExecuteChain(requests []*models.Request, initialVars map[string]string) ([]*models.Response, error) {
	client := NewClient()
	vars := make(map[string]string)
	for k, v := range initialVars {
		vars[k] = v
	}

	var responses []*models.Response
	for i, req := range requests {
		resp, err := client.Execute(req)
		if err != nil {
			return responses, fmt.Errorf("request %d failed: %w", i+1, err)
		}
		responses = append(responses, resp)

		// Store response data as variables for next request
		vars["response.status_code"] = fmt.Sprintf("%d", resp.StatusCode)
		vars["response.body"] = resp.Body
		vars["response.size"] = fmt.Sprintf("%d", resp.Size)
		vars["response.timing"] = resp.Timing.String()
	}

	return responses, nil
}
