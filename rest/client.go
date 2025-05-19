/**
 * @file: client.go
 * @description: RestClient для работы с REST API BloFin
 * @dependencies: rest/public.go, auth/signature.go, utils/errors.go
 * @created: 2025-05-19
 */

package rest

import (
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/mmavka/go-blofin/auth"
)

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
	client := resty.New()
	client.SetBaseURL(baseURL)
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
