package components

import (
	"encoding/json"
	"fmt"

	"github.com/gTahidi/wapi.go/internal"
)

type locationMessageAction struct {
	Name string `json:"name" validate:"required"` // Name of the action.
}

type LocationRequestMessageBody struct {
	Text string `json:"text" validate:"required"` // Text of the body.
}
type LocationRequestMessage struct {
	Type   InteractiveMessageType     `json:"type" validate:"required"`   // Type of the message.
	Action locationMessageAction      `json:"action" validate:"required"` // Action of the message.
	Body   LocationRequestMessageBody `json:"body,omitempty"`             // Body of the message.
}

type LocationRequestMessageParams struct {
	BodyText string `json:"-" validate:"required"` // Text of the body.
}

func NewLocationRequestMessage(params LocationRequestMessageParams) (*LocationRequestMessage, error) {
	if err := internal.GetValidator().Struct(params); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}

	return &LocationRequestMessage{
		Type: InteractiveMessageTypeLocationRequest,
		Body: LocationRequestMessageBody{
			Text: params.BodyText,
		},
		Action: locationMessageAction{
			Name: "send_location",
		},
	}, nil
}

type LocationRequestMessageApiPayload struct {
	BaseMessagePayload
	Interactive LocationRequestMessage `json:"interactive" validate:"required"` // Interactive message.
}

// SetBodyText sets the body text of the list message.
func (m *LocationRequestMessage) SetBodyText(bodyText string) {
	m.Body.Text = bodyText

}

// ToJson converts the list message to JSON.
func (m *LocationRequestMessage) ToJson(configs ApiCompatibleJsonConverterConfigs) ([]byte, error) {
	if err := internal.GetValidator().Struct(configs); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}

	jsonData := LocationRequestMessageApiPayload{
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
