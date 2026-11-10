package manager

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/wapikit/wapi.go/internal/request_client"
	"github.com/wapikit/wapi.go/pkg/events"
)

// A batch with an image that has neither id nor link (NewImageMessage
// rejects it) followed by a valid text message.
const mixedBatchPayload = `{
  "object": "whatsapp_business_account",
  "entry": [{
    "id": "waba-1",
    "changes": [
      {"field": "messages", "value": "not an object"},
      {"field": "messages", "value": {
        "messaging_product": "whatsapp",
        "metadata": {"display_phone_number": "15550001111", "phone_number_id": "pn-1"},
        "contacts": [{"profile": {"name": "Ada"}, "wa_id": "15551234567"}],
        "messages": [
          {"from": "15551234567", "id": "wamid.bad", "timestamp": "1", "type": "image", "image": {}},
          {"from": "15551234567", "id": "wamid.good", "timestamp": "2", "type": "text", "text": {"body": "hello"}}
        ]
      }}
    ]
  }]
}`

func newTestWebhook(t *testing.T, appSecret string) *WebhookManager {
	t.Helper()
	wh := NewWebhook(&WebhookManagerConfig{
		Secret:       "verify-token",
		AppSecret:    appSecret,
		EventManager: NewEventManager(),
		Requester:    *request_client.NewRequestClient("token"),
	})
	if wh == nil {
		t.Fatal("NewWebhook returned nil")
	}
	return wh
}

func postWebhook(wh *WebhookManager, body string, headers map[string]string) (*httptest.ResponseRecorder, error) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(body))
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	rec := httptest.NewRecorder()
	err := wh.PostRequestHandler(e.NewContext(req, rec))
	return rec, err
}

func sign(secret, body string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(body))
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func TestBadMessageDoesNotFailTheRestOfTheBatch(t *testing.T) {
	wh := newTestWebhook(t, "")
	var mu sync.Mutex
	var got []*events.TextMessageEvent
	wh.EventManager.On(events.TextMessageEventType, func(event events.BaseEvent) {
		mu.Lock()
		got = append(got, event.(*events.TextMessageEvent))
		mu.Unlock()
	})
	imageEvents := make(chan struct{}, 1)
	wh.EventManager.On(events.ImageMessageEventType, func(events.BaseEvent) { imageEvents <- struct{}{} })

	rec, err := postWebhook(wh, mixedBatchPayload, nil)
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (a non-2xx makes Meta redeliver the whole batch)", rec.Code)
	}

	waitFor(t, "the text message event", func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(got) == 1
	})
	if got[0].MessageId != "wamid.good" || got[0].Text != "hello" || got[0].SenderName != "Ada" {
		t.Fatalf("unexpected event %+v", got[0])
	}
	select {
	case <-imageEvents:
		t.Fatalf("the invalid image should have been skipped")
	default:
	}
}

func TestInvalidJSONIsRejected(t *testing.T) {
	rec, err := postWebhook(newTestWebhook(t, ""), "{", nil)
	if err != nil || rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d err=%v", rec.Code, err)
	}
}

func TestSignatureIsCheckedOnlyWhenAppSecretIsSet(t *testing.T) {
	const body = `{"object":"whatsapp_business_account","entry":[]}`

	rec, _ := postWebhook(newTestWebhook(t, ""), body, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("no app secret: status = %d", rec.Code)
	}

	wh := newTestWebhook(t, "app-secret")
	for name, header := range map[string]string{
		"missing":   "",
		"malformed": "sha256=zz",
		"no prefix": strings.TrimPrefix(sign("app-secret", body), "sha256="),
		"wrong key": sign("other-secret", body),
	} {
		headers := map[string]string{}
		if header != "" {
			headers["X-Hub-Signature-256"] = header
		}
		rec, _ := postWebhook(wh, body, headers)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("%s signature: status = %d, want 401", name, rec.Code)
		}
	}

	rec, _ = postWebhook(wh, body, map[string]string{"X-Hub-Signature-256": sign("app-secret", body)})
	if rec.Code != http.StatusOK {
		t.Fatalf("valid signature: status = %d", rec.Code)
	}
}

func TestGetRequestHandlerVerifiesToken(t *testing.T) {
	wh := newTestWebhook(t, "")
	e := echo.New()
	for token, want := range map[string]int{"verify-token": 200, "wrong": 400} {
		req := httptest.NewRequest(http.MethodGet, "/webhook?hub.mode=subscribe&hub.challenge=42&hub.verify_token="+token, nil)
		rec := httptest.NewRecorder()
		_ = wh.GetRequestHandler(e.NewContext(req, rec))
		if rec.Code != want {
			t.Fatalf("token %q: status = %d, want %d", token, rec.Code, want)
		}
		if want == 200 && rec.Body.String() != "42" {
			t.Fatalf("challenge = %q", rec.Body.String())
		}
	}
}
