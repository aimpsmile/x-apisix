package apisix

import (
	"testing"
)

func TestApisix(t *testing.T) {
	client := NewClient()
	t.Log(client)
}
