package components

import (
	"encoding/json"
	"fmt"

	"github.com/gTahidi/wapi.go/internal"
)

// textMessage represents a text message.
type textMessage struct {
	Text         string // The text content of the message.
	AllowPreview bool   // Whether to allow preview of the message.
}

// TextMessageConfigs represents the configuration options for a text message.
type TextMessageConfigs struct {
	Text         string `json:"text" validate:"required"` // The text content of the message.
	AllowPreview bool   `json:"allowPreview,omitempty"`   // Whether to allow preview of the message.
}

// TextMessageApiPayloadText represents the text payload for the WhatsApp API.
type TextMessageApiPayloadText struct {
	Body         string `json:"body" validate:"required"` // The text content of the message.
	AllowPreview bool   `json:"preview_url,omitempty"`    // Whether to allow preview of the message.
}

// SetText sets the text content of the message.
func (m *textMessage) SetText(text string) {
	m.Text = text
}

// TextMessageApiPayload represents the payload for the WhatsApp API.
type TextMessageApiPayload struct {
	BaseMessagePayload `json:",inline"`
	Text               TextMessageApiPayloadText `json:"text" validate:"required"` // The text content of the message.
}

// NewTextMessage creates a new text message with the given configurations.
func NewTextMessage(configs TextMessageConfigs) (*textMessage, error) {
	err := internal.GetValidator().Struct(configs)
	if err != nil {
		return nil, fmt.Errorf("error validating text message config: %v", err)
	}
	return &textMessage{
		Text:         configs.Text,
		AllowPreview: configs.AllowPreview,
	}, nil
}

// ToJson converts the text message struct to WhatsApp API compatible JSON.
func (m *textMessage) ToJson(configs ApiCompatibleJsonConverterConfigs) ([]byte, error) {
	if err := internal.GetValidator().Struct(configs); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}

	jsonData := TextMessageApiPayload{
		BaseMessagePayload: NewBaseMessagePayload(configs.SendToPhoneNumber, MessageTypeText),
		Text: TextMessageApiPayloadText{
			Body:         m.Text,
			AllowPreview: m.AllowPreview,
		},
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
