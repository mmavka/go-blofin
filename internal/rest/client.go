// Package rest provides REST API client for Blofin.
package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/mmavka/go-blofin/internal/models"
)

// Client is a REST API client for Blofin.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new REST API client. If no baseURL is provided, BaseURLProd is used.
func NewClient(baseURL ...string) *Client {
	url := BaseURLProd
	if len(baseURL) > 0 && baseURL[0] != "" {
		url = baseURL[0]
	}
	return &Client{
		baseURL:    url,
		httpClient: &http.Client{},
	}
}

// doGet performs a GET request to the given path with query params and decodes the response into result.
func (c *Client) doGet(ctx context.Context, path string, query url.Values, result interface{}) error {
	endpoint := c.baseURL + path
	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Unmarshal into a generic response
	var apiResp struct {
		Code string          `json:"code"`
		Msg  string          `json:"msg"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return err
	}

	if apiResp.Code != CodeSuccess {
		return &models.ApiError{Code: apiResp.Code, Message: apiResp.Msg}
	}

	if result != nil {
		if err := json.Unmarshal(apiResp.Data, result); err != nil {
			return fmt.Errorf("failed to decode data: %w", err)
		}
	}

	return nil
}
