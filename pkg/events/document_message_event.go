package events

import "github.com/gTahidi/wapi.go/pkg/components"

// DocumentMessageEvent represents an event that occurs when a document message is received.
type DocumentMessageEvent struct {
	BaseMediaMessageEvent
	Document components.DocumentMessage
}

// NewDocumentMessageEvent creates a new DocumentMessageEvent instance.
func NewDocumentMessageEvent(baseMessageEvent BaseMessageEvent, document components.DocumentMessage, mediaId, sha256, mimeType string) *DocumentMessageEvent {
	return &DocumentMessageEvent{
		BaseMediaMessageEvent: BaseMediaMessageEvent{
			MediaId:          mediaId,
			Sha256:           sha256,
			MimeType:         mimeType,
			BaseMessageEvent: baseMessageEvent,
		},
		Document: document,
	}
}
