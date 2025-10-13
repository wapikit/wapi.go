package events

import (
	"github.com/gTahidi/wapi.go/pkg/components"
)

// AudioMessageEvent represents an event for an audio message.
type AudioMessageEvent struct {
	BaseMediaMessageEvent `json:",inline"`
	Audio                 components.AudioMessage `json:"audio"`
}

// NewAudioMessageEvent creates a new AudioMessageEvent instance.
func NewAudioMessageEvent(baseMessageEvent BaseMessageEvent, audio components.AudioMessage, mimeType, sha256, mediaId string) *AudioMessageEvent {
	return &AudioMessageEvent{
		BaseMediaMessageEvent: BaseMediaMessageEvent{
			BaseMessageEvent: baseMessageEvent,
			MediaId:          mediaId,
			Sha256:           sha256,
			MimeType:         mimeType,
		},
		Audio: audio,
	}
}
