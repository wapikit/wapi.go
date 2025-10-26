package events

import "github.com/gTahidi/wapi.go/pkg/components"

// VideoMessageEvent represents a WhatsApp video message event.
type VideoMessageEvent struct {
	BaseMediaMessageEvent `json:",inline"`
	Video                 components.VideoMessage `json:"video"`
}

// NewVideoMessageEvent creates a new VideoMessageEvent instance.
func NewVideoMessageEvent(baseMessageEvent BaseMessageEvent, video components.VideoMessage, mimeType, sha256, mediaId string) *VideoMessageEvent {
	return &VideoMessageEvent{
		BaseMediaMessageEvent: BaseMediaMessageEvent{
			MediaId:          mediaId,
			Sha256:           sha256,
			MimeType:         mimeType,
			BaseMessageEvent: baseMessageEvent,
		},
		Video: video,
	}
}
