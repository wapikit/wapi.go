package wapi

import (
	"errors"

	"github.com/wapikit/wapi.go/internal/request_client"
	"github.com/wapikit/wapi.go/manager"
)

// GraphAPIError is the error returned for a non-2xx Graph API response. It
// carries the HTTP status, Meta's error code/subcode/message and the raw body.
//
// Use AsGraphAPIError (or errors.As) to inspect it instead of matching on the
// error string:
//
//	if apiErr, ok := wapi.AsGraphAPIError(err); ok && apiErr.IsRetryable() { ... }
type GraphAPIError = request_client.GraphAPIError

// AsGraphAPIError reports whether err is, or wraps, a *GraphAPIError and
// returns it.
func AsGraphAPIError(err error) (*GraphAPIError, bool) {
	var apiErr *GraphAPIError
	if errors.As(err, &apiErr) {
		return apiErr, true
	}
	return nil, false
}

// ErrInvalidWebhookSignature is returned by VerifyWebhookSignature when the
// signature header is missing, malformed or does not match the body.
var ErrInvalidWebhookSignature = manager.ErrInvalidWebhookSignature

// VerifyWebhookSignature checks a Meta X-Hub-Signature-256 header value
// ("sha256=<hex>") against the HMAC-SHA256 of the raw request body keyed with
// the app secret. It is what the webhook POST handler runs when
// ClientConfig.WebhookAppSecret is set, exposed for callers that verify
// signatures themselves.
func VerifyWebhookSignature(appSecret string, body []byte, signatureHeader string) error {
	return manager.VerifyWebhookSignature(appSecret, body, signatureHeader)
}
