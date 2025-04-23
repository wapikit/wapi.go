package components

import (
	"encoding/json"
	"fmt"

	"github.com/wapikit/wapi.go/internal"
)

// quickReplyButtonMessageButtonReply represents the reply structure of a quick reply button.
type quickReplyButtonMessageButtonReply struct {
	Title string `json:"title" validate:"required"` // Title of the quick reply button.
	Id    string `json:"id" validate:"required"`    // ID of the quick reply button.
}

// quickReplyButtonMessageButton represents a quick reply button.
type quickReplyButtonMessageButton struct {
	Type  string                             `json:"type" validate:"required"`  // Type of the quick reply button.
	Reply quickReplyButtonMessageButtonReply `json:"reply" validate:"required"` // Reply structure of the quick reply button.
}

// NewQuickReplyButton creates a new quick reply button with the given ID and title.
func NewQuickReplyButton(id, title string) (*quickReplyButtonMessageButton, error) {
	return &quickReplyButtonMessageButton{
		Type: "reply",
		Reply: quickReplyButtonMessageButtonReply{
			Title: title,
			Id:    id,
		},
	}, nil
}

// QuickReplyButtonMessageAction represents the action of a quick reply button message.
type QuickReplyButtonMessageAction struct {
	Buttons []quickReplyButtonMessageButton `json:"buttons" validate:"required"` // List of quick reply buttons.
}

// QuickReplyButtonMessageBody represents the body of a quick reply button message.
type QuickReplyButtonMessageBody struct {
	Text string `json:"text" validate:"required"` // Text of the quick reply button message.
}

// QuickReplyButtonMessage represents a quick reply button message.
type QuickReplyButtonMessage struct {
	Type   InteractiveMessageType        `json:"type" validate:"required"`   // Type of the quick reply button message.
	Body   QuickReplyButtonMessageBody   `json:"body" validate:"required"`   // Body of the quick reply button message.
	Action QuickReplyButtonMessageAction `json:"action" validate:"required"` // Action of the quick reply button message.
}

// QuickReplyButtonMessageApiPayload represents the API payload for a quick reply button message.
type QuickReplyButtonMessageApiPayload struct {
	BaseMessagePayload
	Interactive QuickReplyButtonMessage `json:"interactive" validate:"required"` // Interactive part of the API payload.
}

// NewQuickReplyButtonMessage creates a new quick reply button message with the given body text.
func NewQuickReplyButtonMessage(bodyText string) (*QuickReplyButtonMessage, error) {
	return &QuickReplyButtonMessage{
		Type: InteractiveMessageTypeButton,
		Body: QuickReplyButtonMessageBody{
			Text: bodyText,
		},
		Action: QuickReplyButtonMessageAction{
			Buttons: []quickReplyButtonMessageButton{},
		},
	}, nil
}

func (m *QuickReplyButtonMessage) AddButton(id, title string) error {
	button, err := NewQuickReplyButton(id, title)
	if err != nil {
		return fmt.Errorf("error creating quick reply button: %v", err)
	}
	m.Action.Buttons = append(m.Action.Buttons, *button)
	return nil
}

// ToJson converts the quick reply button message to JSON.
func (m *QuickReplyButtonMessage) ToJson(configs ApiCompatibleJsonConverterConfigs) ([]byte, error) {
	if err := internal.GetValidator().Struct(configs); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}

	jsonData := QuickReplyButtonMessageApiPayload{
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
