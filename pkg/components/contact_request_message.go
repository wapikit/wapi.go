package components

import (
	"encoding/json"
	"fmt"

	"github.com/wapikit/wapi.go/internal"
)

// ContactRequestMessage builds an interactive `contact_request` message whose
// action is `request_contact_info` — the high-intent prompt that asks a user to
// share their contact info (BSUID/phone) with the business. The button label is
// Meta-controlled; only the body text is business-set.

type contactRequestAction struct {
	Name string `json:"name" validate:"required"` // "request_contact_info"
}

type ContactRequestMessageBody struct {
	Text string `json:"text" validate:"required"`
}

type ContactRequestMessage struct {
	Type   InteractiveMessageType    `json:"type" validate:"required"`
	Action contactRequestAction      `json:"action" validate:"required"`
	Body   ContactRequestMessageBody `json:"body,omitempty"`
}

type ContactRequestMessageParams struct {
	BodyText string `json:"-" validate:"required"`
}

func NewContactRequestMessage(params ContactRequestMessageParams) (*ContactRequestMessage, error) {
	if err := internal.GetValidator().Struct(params); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}
	return &ContactRequestMessage{
		Type: InteractiveMessageTypeContactRequest,
		Body: ContactRequestMessageBody{Text: params.BodyText},
		Action: contactRequestAction{
			Name: "request_contact_info",
		},
	}, nil
}

func (m *ContactRequestMessage) SetBodyText(bodyText string) {
	m.Body.Text = bodyText
}

type ContactRequestMessageApiPayload struct {
	BaseMessagePayload
	Interactive ContactRequestMessage `json:"interactive" validate:"required"`
}

// ToJson converts the contact-request message to WhatsApp API compatible JSON.
func (m *ContactRequestMessage) ToJson(configs ApiCompatibleJsonConverterConfigs) ([]byte, error) {
	if err := internal.GetValidator().Struct(configs); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}

	jsonData := ContactRequestMessageApiPayload{
		BaseMessagePayload: NewBaseMessagePayload(configs, MessageTypeInteractive),
		Interactive:        *m,
	}

	if configs.ReplyToMessageId != "" {
		jsonData.Context = &Context{MessageId: configs.ReplyToMessageId}
	}

	jsonToReturn, err := json.Marshal(jsonData)
	if err != nil {
		return nil, fmt.Errorf("error marshalling json: %v", err)
	}
	return jsonToReturn, nil
}
