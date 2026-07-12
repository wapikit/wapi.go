package request_client

import "testing"

func TestNewGraphAPIErrorParsesMetaEnvelope(t *testing.T) {
	body := `{"error":{"message":"(#131056) Rate limit hit","type":"OAuthException","code":131056,"error_subcode":2494055,"fbtrace_id":"Abc123"}}`
	e := newGraphAPIError(429, body)
	if e.StatusCode != 429 {
		t.Fatalf("status=%d", e.StatusCode)
	}
	if e.Message != "(#131056) Rate limit hit" || e.Code != 131056 || e.Subcode != 2494055 {
		t.Fatalf("envelope not parsed: %+v", e)
	}
	if e.Body != body {
		t.Fatalf("raw body not preserved")
	}
	if !e.IsRetryable() {
		t.Fatalf("429 should be retryable")
	}
}

// The Error() string must embed the status code so string-based classifiers can
// detect rate limits / server errors.
func TestGraphAPIErrorStringEmbedsStatus(t *testing.T) {
	if got := (&GraphAPIError{StatusCode: 429, Message: "rate limited", Code: 4}).Error(); !contains(got, "429") {
		t.Fatalf("expected status in error string, got %q", got)
	}
	if got := (&GraphAPIError{StatusCode: 500, Body: "oops"}).Error(); !contains(got, "500") {
		t.Fatalf("expected status in error string, got %q", got)
	}
}

func TestGraphAPIErrorRetryable(t *testing.T) {
	cases := map[int]bool{400: false, 401: false, 404: false, 429: true, 500: true, 503: true}
	for status, want := range cases {
		if got := (&GraphAPIError{StatusCode: status}).IsRetryable(); got != want {
			t.Fatalf("status %d retryable=%v, want %v", status, got, want)
		}
	}
}

// non-JSON / empty bodies must not panic and still carry status + raw body.
func TestNewGraphAPIErrorNonJSONBody(t *testing.T) {
	e := newGraphAPIError(502, "<html>Bad Gateway</html>")
	if e.StatusCode != 502 || e.Message != "" || e.Body != "<html>Bad Gateway</html>" {
		t.Fatalf("unexpected: %+v", e)
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
