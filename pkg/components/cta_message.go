package components

import (
	"encoding/json"
	"fmt"

	"github.com/wapikit/wapi.go/internal"
)

type CallToAction struct {
	DisplayText string `json:"display_text" validate:"required"`
	Url         string `json:"url" validate:"required"`
}

func (p *CallToAction) SetDisplayText(text string) {
	p.DisplayText = text
}

func (p *CallToAction) SetUrl(url string) {
	p.Url = url
}

type CtaMessageBody struct {
	Text string `json:"text" validate:"required"`
}

type CtaMessageFooter struct {
	Text string `json:"text"`
}

type CtaMessageHeaderType string

const (
	CtaMessageHeaderTypeText CtaMessageHeaderType = "text"
)

// ! TODO: support more header types
type CtaMessageHeader struct {
	Type CtaMessageHeaderType `json:"type" validate:"required"`
	Text string               `json:"text"`
}

type CtaMessageActionParams struct {
	DisplayText string `json:"display_text" validate:"required"`
	Url         string `json:"url" validate:"required"`
}

type CtaMessageAction struct {
	Name       string                 `json:"name" validate:"required"`
	Parameters CtaMessageActionParams `json:"parameters" validate:"required"`
}

// CtaMessage represents a cta message.
type CtaMessage struct {
	Action CtaMessageAction       `json:"action" validate:"required"`
	Body   CtaMessageBody         `json:"body" validate:"required"`
	Footer *CtaMessageFooter      `json:"footer,omitempty"`
	Header CtaMessageHeader       `json:"header,omitempty"`
	Type   InteractiveMessageType `json:"type" validate:"required"`
}

func (message *CtaMessage) SetBody(text string) {
	message.Body = CtaMessageBody{
		Text: text,
	}
}

func (message *CtaMessage) SetFooter(text string) {
	message.Footer = &CtaMessageFooter{
		Text: text,
	}
}

func (message *CtaMessage) SetHeader(text string) {
	message.Header = CtaMessageHeader{
		Type: CtaMessageHeaderTypeText,
		Text: text,
	}
}

func (message *CtaMessage) SetAction(params CtaMessageActionParams) {
	message.Action = CtaMessageAction{
		Name:       "cta_url",
		Parameters: params,
	}
}

// CtaMessageParams represents the parameters for creating a cta message.
type CtaMessageParams struct {
	BodyText string  `validate:"required"`
	Footer   *string `validate:"required"`
}

// CtaMessageApiPayload represents the API payload for a cta message.
type CtaMessageApiPayload struct {
	BaseMessagePayload
	Interactive CtaMessage `json:"interactive" validate:"required"`
}

// NewCtaMessage creates a new cta message.
func NewCtaMessage(params CtaMessageParams) (*CtaMessage, error) {
	if err := internal.GetValidator().Struct(params); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}

	return &CtaMessage{
		Type: InteractiveMessageTypeCta,
		Body: CtaMessageBody{
			Text: params.BodyText,
		},
	}, nil
}

// ToJson converts the cta message to JSON.
func (m *CtaMessage) ToJson(configs ApiCompatibleJsonConverterConfigs) ([]byte, error) {
	if err := internal.GetValidator().Struct(configs); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}

	jsonData := CtaMessageApiPayload{
		BaseMessagePayload: NewBaseMessagePayload(configs.SendToPhoneNumber, MessageTypeInteractive),
		Interactive:        *m,
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
