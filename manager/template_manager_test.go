package manager

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/wapikit/wapi.go/internal/request_client"
)

type loggedRequest struct {
	method, path, body string
	query              url.Values
}

// newRoutedGraphStub answers each "METHOD /path" with its response and logs
// every request in order.
func newRoutedGraphStub(t *testing.T, responses map[string]string) (*request_client.RequestClient, func() []loggedRequest) {
	t.Helper()
	var mu sync.Mutex
	var log []loggedRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		mu.Lock()
		log = append(log, loggedRequest{method: r.Method, path: r.URL.Path, body: string(body), query: r.URL.Query()})
		mu.Unlock()
		response, ok := responses[r.Method+" "+r.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			_, _ = io.WriteString(w, `{"error":{"message":"unexpected request","code":100}}`)
			return
		}
		_, _ = io.WriteString(w, response)
	}))
	t.Cleanup(srv.Close)
	return request_client.NewRequestClientWithOptions("t", request_client.Options{BaseURL: srv.URL}), func() []loggedRequest {
		mu.Lock()
		defer mu.Unlock()
		return append([]loggedRequest(nil), log...)
	}
}

func newTestTemplateManager(requester *request_client.RequestClient) *TemplateManager {
	return NewTemplateManager(&TemplateManagerConfig{BusinessAccountId: "waba", ApiAccessToken: "t", Requester: requester})
}

func TestTemplateDeleteLooksUpNameAndDeletesById(t *testing.T) {
	requester, requests := newRoutedGraphStub(t, map[string]string{
		"GET /v20.0/tpl-1":                     `{"name":"order_update","id":"tpl-1"}`,
		"DELETE /v20.0/waba/message_templates": `{"success":true}`,
	})
	if err := newTestTemplateManager(requester).Delete("tpl-1"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	log := requests()
	if len(log) != 2 {
		t.Fatalf("requests = %+v", log)
	}
	del := log[1]
	if del.method != http.MethodDelete || del.query.Get("name") != "order_update" || del.query.Get("hsm_id") != "tpl-1" {
		t.Fatalf("delete request = %+v", del)
	}
}

func TestTemplateDeleteByNameAndFailure(t *testing.T) {
	requester, requests := newRoutedGraphStub(t, map[string]string{
		"DELETE /v20.0/waba/message_templates": `{"success":false}`,
	})
	tm := newTestTemplateManager(requester)
	if err := tm.DeleteByName("promo"); err == nil {
		t.Fatalf("expected an error when Meta returns success=false")
	}
	if q := requests()[0].query; q.Get("name") != "promo" || q.Has("hsm_id") {
		t.Fatalf("query = %v", q)
	}
	if err := tm.DeleteByName(""); err == nil {
		t.Fatalf("expected an error for an empty name")
	}
}

func TestMigrateTemplatesUsesPOST(t *testing.T) {
	requester, requests := newRoutedGraphStub(t, map[string]string{
		"POST /v20.0/waba/migrate_message_templates": `{"migrated_templates":["1"],"failed_templates":{"2":"reason"}}`,
	})
	res, err := newTestTemplateManager(requester).MigrateFromOtherBusinessAccount(0, 987)
	if err != nil {
		t.Fatalf("Migrate: %v", err)
	}
	var body map[string]string
	_ = json.Unmarshal([]byte(requests()[0].body), &body)
	if body["source_waba_id"] != "987" || body["page_number"] != "0" {
		t.Fatalf("body = %v", body)
	}
	if len(res.MigratedTemplates) != 1 || res.FailedTemplates["2"] != "reason" {
		t.Fatalf("res = %+v", res)
	}
}

func TestTemplateUpdateReportsSuccess(t *testing.T) {
	requester, _ := newRoutedGraphStub(t, map[string]string{"POST /v20.0/tpl-1": `{"success":true}`})
	res, err := newTestTemplateManager(requester).Update("tpl-1", WhatsAppBusinessAccountMessageTemplateUpdateRequestBody{})
	if err != nil || !res.Success {
		t.Fatalf("res=%+v err=%v", res, err)
	}
}

func TestDeleteMediaUsesMediaIdPath(t *testing.T) {
	requester, requests := newRoutedGraphStub(t, map[string]string{"DELETE /v20.0/media-1": `{"success":true}`})
	if _, err := NewMediaManager(*requester).DeleteMedia("media-1"); err != nil {
		t.Fatalf("DeleteMedia: %v", err)
	}
	if path := requests()[0].path; path != "/v20.0/media-1" {
		t.Fatalf("path = %q", path)
	}
}
