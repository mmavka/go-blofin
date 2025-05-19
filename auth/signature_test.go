package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestSignRequest(t *testing.T) {
	secret := "testsecret"
	method := "GET"
	path := "/api/v1/test"
	timestamp := "1234567890"
	nonce := "abcdef"
	body := "{\"foo\":\"bar\"}"

	prehash := path + method + timestamp + nonce + body
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(prehash))
	hexSignature := hex.EncodeToString(h.Sum(nil))
	signature := base64.StdEncoding.EncodeToString([]byte(hexSignature))

	got := SignRequest(secret, method, path, timestamp, nonce, body)
	if got != signature {
		t.Errorf("signature mismatch: got %s, want %s", got, signature)
	}
}
