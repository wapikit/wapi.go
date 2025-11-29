package components

import (
	"encoding/json"
	"fmt"

	"github.com/wapikit/wapi.go/internal"
)

// VideoMessage represents a video message.
type VideoMessage struct {
	Id      string `json:"id,omitempty"`
	Link    string `json:"link,omitempty"`
	Caption string `json:"caption,omitempty"`
}

// VideoMessageApiPayload represents the API payload for a video message.
type VideoMessageApiPayload struct {
	BaseMessagePayload
	Video VideoMessage `json:"video" validate:"required"`
}

// VideoMessageConfigs is an alias for VideoMessage.
type VideoMessageConfigs = VideoMessage

// NewVideoMessage creates a new VideoMessage instance.
func NewVideoMessage(params VideoMessageConfigs) (*VideoMessage, error) {
	if err := internal.GetValidator().Struct(params); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}

	idSet := params.Id != ""
	linkSet := params.Link != ""

	if !idSet && !linkSet {
		return nil, fmt.Errorf("either ID or Link must be provided")
	}

	return &VideoMessage{
		Id:   params.Id,
		Link: params.Link,
	}, nil
}

// SetCaption sets the caption for the video message.
func (video *VideoMessage) SetCaption(params string) {
	video.Caption = params
}

// ToJson converts the video message to JSON.
func (video *VideoMessage) ToJson(configs ApiCompatibleJsonConverterConfigs) ([]byte, error) {
	if err := internal.GetValidator().Struct(configs); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}

	jsonData := VideoMessageApiPayload{
		BaseMessagePayload: NewBaseMessagePayload(configs.SendToPhoneNumber, MessageTypeVideo),
		Video:              *video,
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
