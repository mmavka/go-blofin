/**
 * @file: signature.go
 * @description: Генерация подписи для BloFin
 * @dependencies: -
 * @created: 2025-05-19
 */

package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// SignRequest generates signature for BloFin
func SignRequest(secret, method, path, timestamp, nonce, body string) string {
	prehash := path + method + timestamp + nonce + body
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(prehash))
	hexSignature := hex.EncodeToString(h.Sum(nil))
	signature := base64.StdEncoding.EncodeToString([]byte(hexSignature))
	return signature
}
