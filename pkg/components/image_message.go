package components

import (
	"encoding/json"
	"fmt"

	"github.com/wapikit/wapi.go/internal"
)

// ImageMessage represents a message with an image.
type ImageMessage struct {
	Id      string `json:"id,omitempty"`
	Link    string `json:"link,omitempty"`
	Caption string `json:"caption,omitempty"`
}

// ImageMessageApiPayload represents the API payload for an image message.
type ImageMessageApiPayload struct {
	BaseMessagePayload
	Image ImageMessage `json:"image" validate:"required"`
}

// ImageMessageConfigs is an alias for ImageMessage.
type ImageMessageConfigs = ImageMessage

// SetCaption sets the caption for the image message.
func (image *ImageMessage) SetCaption(params string) {
	image.Caption = params
}

// NewImageMessage creates a new ImageMessage instance.
func NewImageMessage(params ImageMessageConfigs) (*ImageMessage, error) {
	if err := internal.GetValidator().Struct(params); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}

	idSet := params.Id != ""
	linkSet := params.Link != ""

	if !idSet && !linkSet {
		return nil, fmt.Errorf("either ID or Link must be provided")
	}

	return &ImageMessage{
		Id:      params.Id,
		Link:    params.Link,
		Caption: params.Caption,
	}, nil
}

// ToJson converts the ImageMessage to JSON.
func (image *ImageMessage) ToJson(configs ApiCompatibleJsonConverterConfigs) ([]byte, error) {
	if err := internal.GetValidator().Struct(configs); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}

	jsonData := ImageMessageApiPayload{
		BaseMessagePayload: NewBaseMessagePayload(configs.SendToPhoneNumber, MessageTypeImage),
		Image:              *image,
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
