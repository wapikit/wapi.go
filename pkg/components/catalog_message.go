package components

import (
	"encoding/json"
	"fmt"

	"github.com/wapikit/wapi.go/internal"
)

type CatalogMessageAction struct {
	CatalogId string `json:"catalog_id" validate:"required"`
}

// ProductMessage represents a product message.
type CatalogMessage struct {
	Type   InteractiveMessageType `json:"type" validate:"required"`
	Action CatalogMessageAction   `json:"action" validate:"required"`
}

func NewCatalogMessage(catalogId string) (*CatalogMessage, error) {
	return &CatalogMessage{
		Type: InteractiveMessageTypeCatalog,
		Action: CatalogMessageAction{
			CatalogId: catalogId,
		},
	}, nil
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
