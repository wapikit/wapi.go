package components

import (
	"encoding/json"
	"fmt"

	"github.com/wapikit/wapi.go/internal"
)

// AudioMessage represents an audio message.
type AudioMessage struct {
	Id   string `json:"id,omitempty"`
	Link string `json:"link,omitempty"`
}

// AudioMessageApiPayload represents the payload for an audio message API request.
type AudioMessageApiPayload struct {
	BaseMessagePayload
	Audio AudioMessage `json:"audio" validate:"required"`
}

// AudioMessageConfigs is an alias for AudioMessage.
type AudioMessageConfigs = AudioMessage

// NewAudioMessage creates a new AudioMessage object.
func NewAudioMessage(params AudioMessageConfigs) (*AudioMessage, error) {
	if err := internal.GetValidator().Struct(params); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}

	idSet := params.Id != ""
	linkSet := params.Link != ""

	if !idSet && !linkSet {
		return nil, fmt.Errorf("either ID or Link must be provided")
	}

	return &AudioMessage{
		Id:   params.Id,
		Link: params.Link,
	}, nil
}

// ToJson converts the AudioMessage object to JSON.
func (audio *AudioMessage) ToJson(configs ApiCompatibleJsonConverterConfigs) ([]byte, error) {
	if err := internal.GetValidator().Struct(configs); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}

	jsonData := AudioMessageApiPayload{
		BaseMessagePayload: NewBaseMessagePayload(configs.SendToPhoneNumber, MessageTypeAudio),
		Audio:              *audio,
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
