package httpx

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	cfg := ClientConfig{
		Timeout:      3 * time.Second,
		DebugEnabled: true,
	}

	c := NewClient(cfg)

	ctx := context.TODO()
	uri := "https://webhook.site/280adc37-95f5-4c7b-a44b-18c8f9993bc1"
	params := make(url.Values)
	params.Add("testKey1", "value1")
	header := make(http.Header)
	header.Add("testHeader", "testValue")
	if err := c.GetJSON(ctx, uri, params, header, nil); err != nil {
		t.Fatal(err)
	}

	t.Log("TestGet Ok!")
}
