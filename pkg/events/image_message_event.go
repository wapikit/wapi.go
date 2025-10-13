package events

import "github.com/gTahidi/wapi.go/pkg/components"

// ImageMessageEvent represents an event for an image message.
type ImageMessageEvent struct {
	BaseMediaMessageEvent `json:",inline"`
	Image                 components.ImageMessage `json:"image"`
}

// NewImageMessageEvent creates a new ImageMessageEvent instance.
func NewImageMessageEvent(baseMessageEvent BaseMessageEvent, image components.ImageMessage, mimeType, sha256, mediaId string) *ImageMessageEvent {
	return &ImageMessageEvent{
		BaseMediaMessageEvent: BaseMediaMessageEvent{
			MediaId:          mediaId,
			Sha256:           sha256,
			MimeType:         mimeType,
			BaseMessageEvent: baseMessageEvent,
		},
		Image: image,
	}
}
