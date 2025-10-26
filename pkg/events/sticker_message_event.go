package events

import "github.com/gTahidi/wapi.go/pkg/components"

// StickerMessageEvent represents an event for a sticker message.
type StickerMessageEvent struct {
	BaseMediaMessageEvent `json:",inline"`
	Sticker               components.StickerMessage
}

// NewStickerMessageEvent creates a new StickerMessageEvent instance.
func NewStickerMessageEvent(baseMessageEvent BaseMessageEvent, sticker components.StickerMessage, mediaId, sha256, mimeType string) *StickerMessageEvent {
	return &StickerMessageEvent{
		BaseMediaMessageEvent: BaseMediaMessageEvent{
			MediaId:          mediaId,
			Sha256:           sha256,
			MimeType:         mimeType,
			BaseMessageEvent: baseMessageEvent,
		},
		Sticker: sticker,
	}
}
