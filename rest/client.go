/**
 * @file: client.go
 * @description: RestClient для работы с REST API BloFin
 * @dependencies: rest/public.go, auth/signature.go, utils/errors.go
 * @created: 2025-05-19
 */

package rest

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/mmavka/go-blofin"
	"github.com/mmavka/go-blofin/auth"
)

// UseTestnet use testnet
var UseTestnet = false

// RestClient основной клиент для работы с REST API BloFin
// Позже будет расширен приватными методами и middleware для подписи

type RestClient struct {
	baseURL    string
	httpClient *resty.Client
	apiKey     string
	apiSecret  string
	passphrase string
}

// NewRestClient создает новый экземпляр RestClient
func NewRestClient(baseURL string) *RestClient {
	client := resty.New().
		SetTimeout(10 * time.Second).
		SetBaseURL(baseURL)
	return &RestClient{
		baseURL:    baseURL,
		httpClient: client,
	}
}

// SetAuth устанавливает ключи для приватных методов
func (c *RestClient) SetAuth(apiKey, apiSecret, passphrase string) {
	c.apiKey = apiKey
	c.apiSecret = apiSecret
	c.passphrase = passphrase
}

// addAuthHeaders добавляет заголовки для приватных запросов
func (c *RestClient) addAuthHeaders(req *resty.Request, method, path, body string) {
	timestamp := time.Now().UTC().Format("20060102150405")
	nonce := uuid.New().String()
	sign := auth.SignRequest(c.apiSecret, method, path, timestamp, nonce, body)

	req.SetHeader("ACCESS-KEY", c.apiKey)
	req.SetHeader("ACCESS-SIGN", sign)
	req.SetHeader("ACCESS-TIMESTAMP", timestamp)
	req.SetHeader("ACCESS-NONCE", nonce)
	req.SetHeader("ACCESS-PASSPHRASE", c.passphrase)
}

// NewDefaultRestClient creates a new REST client with default settings
func NewDefaultRestClient() *RestClient {
	return NewRestClient(blofin.DefaultBaseURL)
}

func (c *RestClient) SetBaseURL(url string) {
	c.baseURL = url
}

// GetBaseURL returns the base URL
func (c *RestClient) GetBaseURL() string {
	return c.baseURL
}

// GetAPIKey returns the API key
func (c *RestClient) GetAPIKey() string {
	return c.apiKey
}

// GetSecretKey returns the secret key
func (c *RestClient) GetSecretKey() string {
	return c.apiSecret
}

// GetHTTPClient returns the HTTP client
func (c *RestClient) GetHTTPClient() *resty.Client {
	return c.httpClient
}

// Request sends a request to the API
func (c *RestClient) Request(ctx context.Context, method, path string, body io.Reader) (*resty.Response, error) {
	req := c.httpClient.R().SetContext(ctx)
	if body != nil {
		req.SetBody(body)
	}

	if c.apiKey != "" {
		req.SetHeader("X-API-KEY", c.apiKey)
	}

	resp, err := req.Execute(method, path)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return resp, nil
}
