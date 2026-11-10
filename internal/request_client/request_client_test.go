package request_client

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

// recorded is what the test server saw for one request.
type recorded struct {
	method string
	path   string
	query  url.Values
	auth   string
	ctype  string
	body   string
}

func newTestServer(t *testing.T, status int, response string) (*httptest.Server, *recorded) {
	t.Helper()
	rec := &recorded{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		*rec = recorded{
			method: r.Method,
			path:   r.URL.Path,
			query:  r.URL.Query(),
			auth:   r.Header.Get("Authorization"),
			ctype:  r.Header.Get("Content-Type"),
			body:   string(body),
		}
		w.WriteHeader(status)
		_, _ = io.WriteString(w, response)
	}))
	t.Cleanup(srv.Close)
	return srv, rec
}

func TestBuildURLEncodesEveryQueryParamOnce(t *testing.T) {
	client := NewRequestClient("token")
	got := client.buildURL("123/messages", map[string]string{
		"limit":     "100",
		"after":     "QVFIUk==",
		"filtering": `[{"field":"account_mode","operator":"EQUAL","value":"LIVE"}]`,
		"name":      "a b&c",
	})

	u, err := url.Parse(got)
	if err != nil {
		t.Fatalf("parse %q: %v", got, err)
	}
	if u.Scheme != "https" || u.Host != "graph.facebook.com" || u.Path != "/v20.0/123/messages" {
		t.Fatalf("unexpected URL %q", got)
	}
	want := map[string]string{
		"limit":     "100",
		"after":     "QVFIUk==",
		"filtering": `[{"field":"account_mode","operator":"EQUAL","value":"LIVE"}]`,
		"name":      "a b&c",
	}
	q := u.Query()
	if len(q) != len(want) {
		t.Fatalf("got %d query keys (%v), want %d", len(q), q, len(want))
	}
	for key, value := range want {
		if q[key] == nil || len(q[key]) != 1 || q[key][0] != value {
			t.Fatalf("query %q = %v, want [%q]", key, q[key], value)
		}
	}
	// The old builder duplicated the string on every parameter after the first.
	if strings.Count(got, "?") != 1 {
		t.Fatalf("URL has more than one '?': %q", got)
	}
}

func TestBuildURLWithoutQuery(t *testing.T) {
	got := NewRequestClient("token").buildURL("123", nil)
	if got != "https://graph.facebook.com/v20.0/123" {
		t.Fatalf("got %q", got)
	}
}

func TestExecuteSendsFieldsFiltersAndQueryParams(t *testing.T) {
	srv, rec := newTestServer(t, http.StatusOK, `{"ok":true}`)
	client := NewRequestClientWithOptions("secret-token", Options{BaseURL: srv.URL, ApiVersion: "v99.0"})

	req := client.NewApiRequest("waba-1", http.MethodGet)
	req.AddField(ApiRequestQueryParamField{Name: "id"})
	analytics := req.AddField(ApiRequestQueryParamField{Name: "analytics"})
	analytics.AddFilter("start", "1")
	analytics.AddFilter("end", "2")
	req.AddQueryParam("limit", "10")
	req.AddQueryParam("after", "cursor")

	body, err := req.Execute()
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if body != `{"ok":true}` {
		t.Fatalf("body = %q", body)
	}
	if rec.method != http.MethodGet || rec.path != "/v99.0/waba-1" {
		t.Fatalf("got %s %s", rec.method, rec.path)
	}
	if rec.auth != "Bearer secret-token" {
		t.Fatalf("auth header = %q", rec.auth)
	}
	if got := rec.query.Get("fields"); got != "id,analytics.end(2).start(1)" {
		t.Fatalf("fields = %q", got)
	}
	if rec.query.Get("limit") != "10" || rec.query.Get("after") != "cursor" {
		t.Fatalf("query = %v", rec.query)
	}
}

func TestAddFilterOnFieldWithoutFiltersDoesNotPanic(t *testing.T) {
	req := NewRequestClient("t").NewApiRequest("x", http.MethodGet)
	field := req.AddField(ApiRequestQueryParamField{Name: "f"})
	field.AddFilter("k", "v")
	if got := req.fieldsParam(); got != "f.k(v)" {
		t.Fatalf("fieldsParam = %q", got)
	}

	var bare ApiRequestQueryParamField
	bare.AddFilter("k", "v") // nil map must not panic
	if bare.Filters["k"] != "v" {
		t.Fatalf("filter not stored")
	}
}

func TestExecuteReturnsGraphAPIErrorAndBody(t *testing.T) {
	const envelope = `{"error":{"message":"(#131056) pair rate limit hit","type":"OAuthException","code":131056,"fbtrace_id":"abc"}}`
	srv, rec := newTestServer(t, http.StatusBadRequest, envelope)
	client := NewRequestClientWithOptions("t", Options{BaseURL: srv.URL})

	req := client.NewApiRequest("123/messages", http.MethodPost)
	req.SetBody(`{"to":"1"}`)
	body, err := req.Execute()

	var apiErr *GraphAPIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("err = %v, want *GraphAPIError", err)
	}
	if apiErr.StatusCode != 400 || apiErr.Code != 131056 {
		t.Fatalf("apiErr = %+v", apiErr)
	}
	// Callers may match on this prefix; keep the error text stable.
	if !strings.HasPrefix(err.Error(), "graph api error (status 400, code 131056): ") {
		t.Fatalf("error text changed: %q", err.Error())
	}
	if body != envelope {
		t.Fatalf("body = %q", body)
	}
	if rec.ctype != "application/json" || rec.body != `{"to":"1"}` {
		t.Fatalf("request ctype=%q body=%q", rec.ctype, rec.body)
	}
}

func TestExecuteContextHonoursCancellation(t *testing.T) {
	block := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-block
	}))
	defer srv.Close()
	defer close(block)

	client := NewRequestClientWithOptions("t", Options{BaseURL: srv.URL})
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := client.NewApiRequest("x", http.MethodGet).ExecuteContext(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want context.DeadlineExceeded", err)
	}
	if time.Since(start) > 5*time.Second {
		t.Fatalf("request was not cancelled promptly")
	}
}

func TestMultipartKeepsGraphAPIErrorUnwrapped(t *testing.T) {
	srv, rec := newTestServer(t, http.StatusInternalServerError, `{"error":{"message":"boom","code":1}}`)
	client := NewRequestClientWithOptions("t", Options{BaseURL: srv.URL})

	_, err := client.RequestMultipart(http.MethodPost, "123/media", strings.NewReader("data"), "multipart/form-data; boundary=x")
	var apiErr *GraphAPIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 500 {
		t.Fatalf("err = %v", err)
	}
	if !strings.HasPrefix(err.Error(), "graph api error (status 500") {
		t.Fatalf("error text changed: %q", err.Error())
	}
	if rec.ctype != "multipart/form-data; boundary=x" || rec.path != "/v20.0/123/media" {
		t.Fatalf("ctype=%q path=%q", rec.ctype, rec.path)
	}
}

func TestDefaultClientHasNoOverallTimeoutButBoundsResponseHeaders(t *testing.T) {
	client := NewRequestClient("t")
	if client.httpClient != defaultHTTPClient {
		t.Fatalf("expected shared default client")
	}
	if defaultHTTPClient.Timeout != 0 {
		t.Fatalf("overall timeout would cut off large uploads")
	}
	transport := defaultHTTPClient.Transport.(*http.Transport)
	if transport.ResponseHeaderTimeout != DefaultResponseHeaderTimeout {
		t.Fatalf("ResponseHeaderTimeout = %v", transport.ResponseHeaderTimeout)
	}

	custom := &http.Client{}
	if NewRequestClientWithOptions("t", Options{HTTPClient: custom}).httpClient != custom {
		t.Fatalf("custom HTTP client ignored")
	}
}
