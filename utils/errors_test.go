package utils

import (
	"testing"
)

func TestErrAPIRequest(t *testing.T) {
	if ErrAPIRequest == nil {
		t.Error("ErrAPIRequest is nil")
	}
	if ErrAPIRequest.Error() != "blofin: api request error" {
		t.Errorf("unexpected error string: %s", ErrAPIRequest.Error())
	}
}
