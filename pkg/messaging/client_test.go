package messaging

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wapikit/wapi.go/internal/request_client"
)

func TestRegisterSendsPinInBodyToRegisterEndpoint(t *testing.T) {
	var method, path string
	var body map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		if len(r.URL.Query()) != 0 {
			t.Errorf("unexpected query params %v", r.URL.Query())
		}
		_, _ = io.WriteString(w, `{"success":true}`)
	}))
	defer srv.Close()

	client := &MessagingClient{
		PhoneNumberId: "pn-1",
		Requester:     request_client.NewRequestClientWithOptions("t", request_client.Options{BaseURL: srv.URL}),
	}
	res, err := client.Register("123456")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if !res.Success {
		t.Fatalf("Success = false")
	}
	if method != http.MethodPost || path != "/v20.0/pn-1/register" {
		t.Fatalf("sent %s %s", method, path)
	}
	if body["messaging_product"] != "whatsapp" || body["pin"] != "123456" {
		t.Fatalf("body = %v", body)
	}
}
