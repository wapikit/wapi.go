package components

import (
	"encoding/json"
	"fmt"

	"github.com/gTahidi/wapi.go/internal"
)

// LocationMessage represents a location message with latitude, longitude, address, and name.
type LocationMessage struct {
	Latitude  float64 `json:"latitude" validate:"required"`  // Latitude of the location
	Longitude float64 `json:"longitude" validate:"required"` // Longitude of the location
	Address   string  `json:"address,omitempty"`             // Address of the location (optional)
	Name      string  `json:"name,omitempty"`                // Name of the location (optional)
}

// LocationMessageApiPayload represents the API payload for a location message.
type LocationMessageApiPayload struct {
	BaseMessagePayload
	Location LocationMessage `json:"location" validate:"required"` // Location message
}

// NewLocationMessage creates a new LocationMessage with the given latitude and longitude.
func NewLocationMessage(latitude float64, longitude float64) (*LocationMessage, error) {
	return &LocationMessage{
		Latitude:  latitude,
		Longitude: longitude,
	}, nil
}

// SetAddress sets the address of the location.
func (location *LocationMessage) SetAddress(params string) {
	location.Address = params
}

// SetName sets the name of the location.
func (location *LocationMessage) SetName(params string) {
	location.Name = params
}

// ToJson converts the LocationMessage to JSON with the given configurations.
func (location *LocationMessage) ToJson(configs ApiCompatibleJsonConverterConfigs) ([]byte, error) {
	if err := internal.GetValidator().Struct(configs); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}

	jsonData := LocationMessageApiPayload{
		BaseMessagePayload: NewBaseMessagePayload(configs.SendToPhoneNumber, MessageTypeLocation),
		Location:           *location,
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
