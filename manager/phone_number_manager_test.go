package manager

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wapikit/wapi.go/internal/request_client"
)

type seenRequest struct {
	method, path, body string
}

func newGraphStub(t *testing.T, response string) (*request_client.RequestClient, *seenRequest) {
	t.Helper()
	seen := &seenRequest{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		*seen = seenRequest{method: r.Method, path: r.URL.Path, body: string(body)}
		_, _ = io.WriteString(w, response)
	}))
	t.Cleanup(srv.Close)
	return request_client.NewRequestClientWithOptions("t", request_client.Options{BaseURL: srv.URL}), seen
}

func TestGetQrCodeByIdUsesGET(t *testing.T) {
	requester, seen := newGraphStub(t, `{"data":[{"code":"ABC","prefilled_message":"hi"}]}`)
	pm := NewPhoneNumberManager(&PhoneNumberManagerConfig{BusinessAccountId: "waba", ApiAccessToken: "t", Requester: requester})

	res, err := pm.GetQrCodeById("pn-1", "ABC")
	if err != nil {
		t.Fatalf("GetQrCodeById: %v", err)
	}
	if seen.method != http.MethodGet || seen.path != "/v20.0/pn-1/message_qrdls/ABC" {
		t.Fatalf("sent %s %s; fetching a QR code must not DELETE it", seen.method, seen.path)
	}
	if len(res.Data) != 1 || res.Data[0].Code != "ABC" {
		t.Fatalf("res = %+v", res)
	}
}

func TestQrCodeResponseDecodeErrorsAreReturned(t *testing.T) {
	requester, _ := newGraphStub(t, `not json`)
	pm := NewPhoneNumberManager(&PhoneNumberManagerConfig{BusinessAccountId: "waba", ApiAccessToken: "t", Requester: requester})
	if _, err := pm.GetAllQrCodes("pn-1"); err == nil {
		t.Fatalf("expected a decode error instead of an empty result")
	}
}
