package request_client

import (
	"encoding/json"
	"fmt"
)

// GraphAPIError is returned when the Graph API responds with a non-2xx status.
// It carries the HTTP status, the raw response body, and the parsed fields from
// Meta's `{"error": {...}}` envelope when present. Callers can type-assert the
// error to inspect the status/code or decide on a retry.
type GraphAPIError struct {
	StatusCode int    // HTTP status code of the response.
	Body       string // Raw response body (never dropped).
	Message    string // error.message, when present.
	Type       string // error.type, when present.
	Code       int    // error.code, when present.
	Subcode    int    // error.error_subcode, when present.
	FBTraceID  string // error.fbtrace_id, for support/debugging.
}

func (e *GraphAPIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("graph api error (status %d, code %d): %s", e.StatusCode, e.Code, e.Message)
	}
	body := e.Body
	if len(body) > 300 {
		body = body[:300]
	}
	return fmt.Sprintf("graph api error (status %d): %s", e.StatusCode, body)
}

// IsRetryable reports whether the error is worth retrying — rate limits (429)
// and transient server errors (>=500). Callers own any backoff policy.
func (e *GraphAPIError) IsRetryable() bool {
	return e.StatusCode == 429 || e.StatusCode >= 500
}

// newGraphAPIError builds a GraphAPIError from a non-2xx status and body,
// parsing Meta's error envelope when the body contains one.
func newGraphAPIError(statusCode int, body string) *GraphAPIError {
	apiErr := &GraphAPIError{StatusCode: statusCode, Body: body}
	var parsed struct {
		Error struct {
			Message   string `json:"message"`
			Type      string `json:"type"`
			Code      int    `json:"code"`
			Subcode   int    `json:"error_subcode"`
			FBTraceID string `json:"fbtrace_id"`
		} `json:"error"`
	}
	if json.Unmarshal([]byte(body), &parsed) == nil {
		apiErr.Message = parsed.Error.Message
		apiErr.Type = parsed.Error.Type
		apiErr.Code = parsed.Error.Code
		apiErr.Subcode = parsed.Error.Subcode
		apiErr.FBTraceID = parsed.Error.FBTraceID
	}
	return apiErr
}
