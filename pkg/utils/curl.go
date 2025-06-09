package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

// ConnectionSpecial performs an HTTP POST request with query parameters
// request: map of query parameters
// baseURL: base URL for the request
// debug: enable debug output
// Returns: decoded JSON response and error if any
func ConnectionSpecial(request map[string]string, baseURL string, debug bool) (interface{}, error) {
	// Create query parameters
	params := url.Values{}
	for key, value := range request {
		params.Add(key, value)
	}

	// Append query parameters to URL
	fullURL := baseURL
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create POST request
	req, err := http.NewRequest("POST", fullURL, nil)
	if err != nil {
		return nil, err
	}

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if debug {
		// In Go, you might want to use a logger instead of direct printing
		println("HTTP Status Code:", resp.StatusCode)
		println("Response Body:", string(body))
	}

	// Parse JSON response
	var result interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
