package components

import (
	"encoding/json"
	"fmt"

	"github.com/wapikit/wapi.go/internal"
)

// DocumentMessage represents a document message.
type DocumentMessage struct {
	Id       string  `json:"id,omitempty"`
	Link     *string `json:"link,omitempty"`
	Caption  *string `json:"caption,omitempty"`
	FileName string  `json:"filename" validate:"required"`
}

// DocumentMessageApiPayload represents the API payload for a document message.
type DocumentMessageApiPayload struct {
	BaseMessagePayload
	Document DocumentMessage `json:"document" validate:"required"`
}

// DocumentMessageConfigs represents the configurations for a document message.
type DocumentMessageConfigs struct {
	Id       string  `json:"id" validate:"required"`
	Link     *string `json:"link,omitempty"`
	Caption  *string `json:"caption,omitempty"`
	FileName string  `json:"filename" validate:"required"`
}

// NewDocumentMessage creates a new DocumentMessage instance.
func NewDocumentMessage(params DocumentMessageConfigs) (*DocumentMessage, error) {
	if err := internal.GetValidator().Struct(params); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}
	return &DocumentMessage{
		Id:       params.Id,
		Link:     params.Link,
		Caption:  params.Caption,
		FileName: params.FileName,
	}, nil
}

// ToJson converts the DocumentMessage instance to JSON.
func (dm *DocumentMessage) ToJson(configs ApiCompatibleJsonConverterConfigs) ([]byte, error) {
	if err := internal.GetValidator().Struct(configs); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}

	jsonData := DocumentMessageApiPayload{
		BaseMessagePayload: NewBaseMessagePayload(configs.SendToPhoneNumber, MessageTypeDocument),
		Document:           *dm,
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
