package ws

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestSignWebSocketLogin(t *testing.T) {
	secret := "testsecret"
	timestamp := "1234567890"
	nonce := "abcdef"
	path := "/users/self/verify"
	method := "GET"
	prehash := path + method + timestamp + nonce
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(prehash))
	hexSignature := hex.EncodeToString(h.Sum(nil))
	expected := base64.StdEncoding.EncodeToString([]byte(hexSignature))

	got := SignWebSocketLogin(secret, timestamp, nonce)
	if got != expected {
		t.Errorf("signature mismatch: got %s, want %s", got, expected)
	}
}
