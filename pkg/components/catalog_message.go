package components

import (
    "encoding/json"
    "fmt"

    "github.com/wapikit/wapi.go/internal"
)

type CatalogMessageActionParameter struct {
	ThumbnailProductRetailerId string `json:"thumbnail_product_retailer_id" validate:"required"`
}

type CatalogMessageAction struct {
	Name       string                        `json:"name" validate:"required"`
	Parameters CatalogMessageActionParameter `json:"parameters,omitempty"`
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

func NewCatalogMessage(name, thumbnailProductRetailerId string) (*CatalogMessage, error) {
    if thumbnailProductRetailerId == "" {
        return nil, fmt.Errorf("thumbnail_product_retailer_id is required for catalog_message")
    }
    if name == "" {
        name = "catalog_message"
    }
    return &CatalogMessage{
        Type: InteractiveMessageTypeCatalog,
        Action: CatalogMessageAction{
            Name: name,
            Parameters: CatalogMessageActionParameter{
                ThumbnailProductRetailerId: thumbnailProductRetailerId,
            },
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
    // Validate message structure and required fields as well
    if err := internal.GetValidator().Struct(m); err != nil {
        return nil, fmt.Errorf("error validating catalog message: %v", err)
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
