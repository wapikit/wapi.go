package wapi

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wapikit/wapi.go/pkg/components"
)

func TestAsGraphAPIErrorThroughSend(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = io.WriteString(w, `{"error":{"message":"rate limited","code":130429}}`)
	}))
	defer srv.Close()

	client := New(&ClientConfig{BusinessAccountId: "waba", ApiAccessToken: "t", WebhookSecret: "s", GraphBaseURL: srv.URL})
	text, err := components.NewTextMessage(components.TextMessageConfigs{Text: "hi"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.NewMessagingClient("pn-1").Message.Send(text, "15551234567")

	apiErr, ok := AsGraphAPIError(fmt.Errorf("wrapped: %w", err))
	if !ok {
		t.Fatalf("err = %v, want a GraphAPIError", err)
	}
	if apiErr.StatusCode != 429 || apiErr.Code != 130429 || !apiErr.IsRetryable() {
		t.Fatalf("apiErr = %+v", apiErr)
	}
	if _, ok := AsGraphAPIError(fmt.Errorf("plain")); ok {
		t.Fatalf("plain error reported as GraphAPIError")
	}
}
