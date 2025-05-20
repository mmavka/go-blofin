/**
 * @file: signature.go
 * @description: Генерация подписи для WebSocket login BloFin
 * @dependencies: -
 * @created: 2025-05-19
 */

package ws

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// SignWebSocketLogin generates signature for WebSocket login
func SignWebSocketLogin(secret, timestamp, nonce string) string {
	path := "/users/self/verify"
	method := "GET"
	prehash := path + method + timestamp + nonce
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(prehash))
	hexSignature := hex.EncodeToString(h.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(hexSignature))
}
