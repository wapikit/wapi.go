package components

import (
	"encoding/json"
	"fmt"

	"github.com/gTahidi/wapi.go/internal"
)

// StickerMessage represents a sticker message.
type StickerMessage struct {
	Id   string `json:"id,omitempty"`
	Link string `json:"link,omitempty"`
}

// StickerMessageApiPayload represents the API payload for a sticker message.
type StickerMessageApiPayload struct {
	BaseMessagePayload
	Sticker StickerMessage `json:"sticker" validate:"required"`
}

// StickerMessageConfigs represents the configurations for a sticker message.
type StickerMessageConfigs struct {
	Id   string `json:"id,omitempty"`
	Link string `json:"link,omitempty"`
}

// NewStickerMessage creates a new sticker message based on the provided configurations.
func NewStickerMessage(params *StickerMessageConfigs) (*StickerMessage, error) {
	if err := internal.GetValidator().Struct(params); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}
	idSet := params.Id != ""
	linkSet := params.Link != ""

	if idSet && linkSet {
		return nil, fmt.Errorf("only one of ID or Link can be provided")
	}
	if !idSet && !linkSet {
		return nil, fmt.Errorf("either ID or Link must be provided")
	}

	return &StickerMessage{
		Id:   params.Id,
		Link: params.Link,
	}, nil
}

// ToJson converts the sticker message to JSON based on the provided configurations.
func (sticker *StickerMessage) ToJson(configs ApiCompatibleJsonConverterConfigs) ([]byte, error) {
	if err := internal.GetValidator().Struct(configs); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}

	jsonData := StickerMessageApiPayload{
		BaseMessagePayload: NewBaseMessagePayload(configs.SendToPhoneNumber, MessageTypeSticker),
		Sticker:            *sticker,
	}

	if configs.ReplyToMessageId != "" {
		jsonData.Context = &Context{
			MessageId: configs.ReplyToMessageId,
		}
	}

	jsonToReturn, err := json.Marshal(jsonData)

	if err != nil {
		return nil, fmt.Errorf("error marshalling json: %v", err)
	}

	return jsonToReturn, nil
}
