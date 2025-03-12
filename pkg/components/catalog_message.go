package components

import (
	"encoding/json"
	"fmt"

	"github.com/wapikit/wapi.go/internal"
)

type CatalogMessageAction struct {
	CatalogId string `json:"catalog_id" validate:"required"`
}

type CatalogMessageBody struct {
	Text string `json:"text" validate:"required"`
}

type CatalogMessageFooter struct {
	Text string `json:"text" validate:"required"`
}

type CatalogMessageHeader struct {
	Text string `json:"text" validate:"required"`
}

// ProductMessage represents a product message.
type CatalogMessage struct {
	Type   InteractiveMessageType `json:"type" validate:"required"`
	Action CatalogMessageAction   `json:"action" validate:"required"`
	Body   CatalogMessageBody     `json:"body" validate:"required"`
	Footer *CatalogMessageFooter  `json:"footer,omitempty"`
	Header *CatalogMessageHeader  `json:"header,omitempty"`
}

func NewCatalogMessage(catalogId string) (*CatalogMessage, error) {
	return &CatalogMessage{
		Type: InteractiveMessageTypeCatalog,
		Action: CatalogMessageAction{
			CatalogId: catalogId,
		},
	}, nil
}

func (m *CatalogMessage) SetHeader(text string) {
	m.Header = &CatalogMessageHeader{
		Text: text,
	}
}

func (m *CatalogMessage) SetBody(text string) {
	m.Body = CatalogMessageBody{
		Text: text,
	}
}

func (m *CatalogMessage) SetFooter(text string) {
	m.Footer = &CatalogMessageFooter{
		Text: text,
	}
}

// ProductMessageApiPayload represents the API payload for a product message.
type CatalogMessageApiPayload struct {
	BaseMessagePayload
	Interactive CatalogMessage `json:"interactive" validate:"required"`
}

// ToJson converts the product message to JSON with the given configurations.
func (m *CatalogMessage) ToJson(configs ApiCompatibleJsonConverterConfigs) ([]byte, error) {
	if err := internal.GetValidator().Struct(configs); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}

	jsonData := CatalogMessageApiPayload{
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
