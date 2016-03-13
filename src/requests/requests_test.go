package requests

import (
	"testing"
)

func TestRequests(t *testing.T) {

	pool := New(2, true)
	pool.send(&Request{
		Method: "GET",
		Url:    "https://google.fr",
	})
}
