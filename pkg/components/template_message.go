package components

import (
	"encoding/json"
	"fmt"

	"github.com/wapikit/wapi.go/internal"
)

// TemplateMessageComponentType represents the type of a template message component.
type TemplateMessageComponentType string

const (
	TemplateMessageComponentTypeHeader TemplateMessageComponentType = "header"
	TemplateMessageComponentTypeBody   TemplateMessageComponentType = "body"
	TemplateMessageComponentTypeButton TemplateMessageComponentType = "button"
)

// TemplateMessageButtonComponentType represents the subtype of a button component.
type TemplateMessageButtonComponentType string

const (
	TemplateMessageButtonComponentTypeQuickReply TemplateMessageButtonComponentType = "quick_reply"
	TemplateMessageButtonComponentTypeUrl        TemplateMessageButtonComponentType = "url"
	TemplateMessageButtonComponentTypeCatalog    TemplateMessageButtonComponentType = "catalog"
)

// TemplateMessageComponent is an interface for all template message components.
type TemplateMessageComponent interface {
	GetComponentType() string
}

// TemplateMessageComponentButtonType represents a button component.
type TemplateMessageComponentButtonType struct {
	Type       TemplateMessageComponentType       `json:"type" validate:"required"`                 // e.g., "button"
	SubType    TemplateMessageButtonComponentType `json:"sub_type" validate:"required"`             // e.g., "quick_reply", "url", etc.
	Index      int                                `json:"index" validate:"required"`                // Position index of the button (0 to 9)
	Parameters *[]TemplateMessageParameter        `json:"parameters,omitempty" validate:"required"` // Parameters for the button component.
}

// GetComponentType returns the component type.
func (t TemplateMessageComponentButtonType) GetComponentType() string {
	return string(t.Type)
}

// TemplateMessageComponentHeaderType represents a header component.
type TemplateMessageComponentHeaderType struct {
	Type       TemplateMessageComponentType `json:"type" validate:"required"`                 // "header"
	Parameters *[]TemplateMessageParameter  `json:"parameters,omitempty" validate:"required"` // Parameters for the header component.
}

// GetComponentType returns the component type.
func (t TemplateMessageComponentHeaderType) GetComponentType() string {
	return string(t.Type)
}

// TemplateMessageComponentBodyType represents a body component.
type TemplateMessageComponentBodyType struct {
	Type       TemplateMessageComponentType `json:"type" validate:"required"`       // "body"
	Parameters []TemplateMessageParameter   `json:"parameters" validate:"required"` // Parameters for the body component.
}

// GetComponentType returns the component type.
func (t TemplateMessageComponentBodyType) GetComponentType() string {
	return string(t.Type)
}

// TemplateMessageParameterType represents the type of a parameter.
type TemplateMessageParameterType string

const (
	TemplateMessageParameterTypeCurrency TemplateMessageParameterType = "currency"
	TemplateMessageParameterTypeDateTime TemplateMessageParameterType = "date_time"
	TemplateMessageParameterTypeDocument TemplateMessageParameterType = "document"
	TemplateMessageParameterTypeImage    TemplateMessageParameterType = "image"
	TemplateMessageParameterTypeText     TemplateMessageParameterType = "text"
	TemplateMessageParameterTypeVideo    TemplateMessageParameterType = "video"
	TemplateMessageParameterTypeLocation TemplateMessageParameterType = "location"
)

// TemplateMessageParameterCurrency represents a currency parameter.
type TemplateMessageParameterCurrency struct {
	FallbackValue string `json:"fallback_value" validate:"required"` // Default text if localization fails.
	Code          string `json:"code" validate:"required"`           // ISO 4217 currency code.
	Amount1000    int    `json:"amount_1000" validate:"required"`    // Amount multiplied by 1000.
}

// TemplateMessageParameterDateTime represents a date-time parameter.
type TemplateMessageParameterDateTime struct {
	FallbackValue string `json:"fallback_value" validate:"required"` // Default text if localization fails.
}

// TemplateMessageParameterMedia represents a media parameter (for document, image, video).
type TemplateMessageParameterMedia struct {
	Link string `json:"link" validate:"required"` // URL link of the media.
}

// TemplateMessageParameterLocation represents a location parameter.
type TemplateMessageParameterLocation struct {
	Latitude  string `json:"latitude" validate:"required"`  // Latitude.
	Longitude string `json:"longitude" validate:"required"` // Longitude.
	Name      string `json:"name" validate:"required"`      // Location name.
	Address   string `json:"address" validate:"required"`   // Address.
}

// TemplateMessageParameter is an interface for all parameter types.
type TemplateMessageParameter interface {
	GetParameterType() string
}

// TemplateMessageBodyAndHeaderParameter represents parameters for body and header components.
type TemplateMessageBodyAndHeaderParameter struct {
	Type          TemplateMessageParameterType      `json:"type" validate:"required"` // e.g., "text", "currency", etc.
	ParameterName *string                           `json:"parameter_name,omitempty"` // Optional: name of the parameter (for named parameters).
	Currency      *TemplateMessageParameterCurrency `json:"currency,omitempty"`       // Currency details (if type is currency).
	DateTime      *TemplateMessageParameterDateTime `json:"date_time,omitempty"`      // Date/time details (if type is date_time).
	Document      *TemplateMessageParameterMedia    `json:"document,omitempty"`       // Document details (if type is document).
	Image         *TemplateMessageParameterMedia    `json:"image,omitempty"`          // Image details (if type is image).
	Text          *string                           `json:"text,omitempty"`           // Text content (if type is text).
	Video         *TemplateMessageParameterMedia    `json:"video,omitempty"`          // Video details (if type is video).
	Location      *TemplateMessageParameterLocation `json:"location,omitempty"`       // Location details (if type is location).
}

// GetParameterType returns the parameter type as a string.
func (t TemplateMessageBodyAndHeaderParameter) GetParameterType() string {
	return string(t.Type)
}

// TemplateMessageButtonParameterType represents the type for button parameters.
type TemplateMessageButtonParameterType string

const (
	TemplateMessageButtonParameterTypePayload TemplateMessageButtonParameterType = "payload"
	TemplateMessageButtonParameterTypeText    TemplateMessageButtonParameterType = "text"
)

// TemplateMessageButtonParameter represents a parameter for a button component.
type TemplateMessageButtonParameter struct {
	Type    TemplateMessageButtonParameterType `json:"type" validate:"required"` // e.g., "payload" or "text"
	Payload string                             `json:"payload,omitempty"`        // Required for quick_reply buttons.
	Text    string                             `json:"text,omitempty"`           // Required for URL buttons.
}

// GetParameterType returns the button parameter type as a string.
func (t TemplateMessageButtonParameter) GetParameterType() string {
	return string(t.Type)
}

// TemplateMessageLanguage represents the language configuration.
type TemplateMessageLanguage struct {
	Code   string `json:"code" validate:"required"`   // e.g., "en_US"
	Policy string `json:"policy" validate:"required"` // e.g., "deterministic"
}

// TemplateMessage represents a template message.
type TemplateMessage struct {
	Name       string                     `json:"name" validate:"required"`       // Template name.
	Language   TemplateMessageLanguage    `json:"language" validate:"required"`   // Language configuration.
	Components []TemplateMessageComponent `json:"components" validate:"required"` // Array of components.
}

// TemplateMessageApiPayload represents the API payload for sending a template message.
type TemplateMessageApiPayload struct {
	BaseMessagePayload
	Template TemplateMessage `json:"template" validate:"required"`
}

// TemplateMessageConfigs represents basic configurations for a template message.
type TemplateMessageConfigs struct {
	Name     string `json:"name" validate:"required"`     // Template name.
	Language string `json:"language" validate:"required"` // Language code.
}

// NewTemplateMessage creates a new TemplateMessage instance.
func NewTemplateMessage(params *TemplateMessageConfigs) (*TemplateMessage, error) {
	return &TemplateMessage{
		Name: params.Name,
		Language: TemplateMessageLanguage{
			Code:   params.Language,
			Policy: "deterministic",
		},
	}, nil
}

// AddHeader adds (or overrides) a header component in the template message.
// Only one header is allowed.
func (tm *TemplateMessage) AddHeader(params TemplateMessageComponentHeaderType) {
	var existingHeaderIndex int
	var found bool

	for i, component := range tm.Components {
		if TemplateMessageComponentType(component.GetComponentType()) == TemplateMessageComponentTypeHeader {
			existingHeaderIndex = i
			found = true
			break
		}
	}

	if found {
		// Override the existing header.
		tm.Components[existingHeaderIndex] = params
	} else {
		// Append the new header.
		tm.Components = append(tm.Components, params)
	}
}

// AddBody adds (or overrides) a body component in the template message.
// Only one body component is allowed.
func (tm *TemplateMessage) AddBody(params TemplateMessageComponentBodyType) {
	var existingBodyIndex int
	var found bool

	for i, component := range tm.Components {
		if TemplateMessageComponentType(component.GetComponentType()) == TemplateMessageComponentTypeBody {
			existingBodyIndex = i
			found = true
			break
		}
	}

	if found {
		// Override the existing body.
		tm.Components[existingBodyIndex] = params
	} else {
		// Append the new body.
		tm.Components = append(tm.Components, params)
	}
}

// AddButton adds a button component to the template message.
// A maximum of 10 buttons is allowed.
func (tm *TemplateMessage) AddButton(params TemplateMessageComponentButtonType) error {
	// Count existing button components.
	numberOfButtons := 0
	for _, component := range tm.Components {
		if TemplateMessageComponentType(component.GetComponentType()) == TemplateMessageComponentTypeButton {
			numberOfButtons++
		}
	}

	if numberOfButtons >= 10 {
		return fmt.Errorf("maximum number of buttons reached")
	}
	tm.Components = append(tm.Components, params)
	return nil
}

// ToJson converts the template message into a JSON payload compatible with the API.
func (m *TemplateMessage) ToJson(configs ApiCompatibleJsonConverterConfigs) ([]byte, error) {
	if err := internal.GetValidator().Struct(configs); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}

	jsonData := TemplateMessageApiPayload{
		BaseMessagePayload: NewBaseMessagePayload(configs.SendToPhoneNumber, MessageTypeTemplate),
		Template:           *m,
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
