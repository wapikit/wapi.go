package components

import (
	"encoding/json"
	"fmt"

	"github.com/gTahidi/wapi.go/internal"
)

// ListSection represents a section in the list message.
type ListSection struct {
	Title string           `json:"title" validate:"required"` // Title of the section.
	Rows  []listSectionRow `json:"rows" validate:"required"`  // Rows in the section.
}

// NewListSection creates a new list section with the given title.
func NewListSection(title string) (*ListSection, error) {
	return &ListSection{
		Title: title,
	}, nil
}

// AddRow adds a new row to the list section.
func (section *ListSection) AddRow(row *listSectionRow) {
	section.Rows = append(section.Rows, *row)
}

// SetTitle sets the title of the list section.
func (section *ListSection) SetTitle(title string) {
	section.Title = title
}

// listSectionRow represents a row in the list section.
type listSectionRow struct {
	Id          string `json:"id" validate:"required"`          // ID of the row.
	Description string `json:"description" validate:"required"` // Description of the row.
	Title       string `json:"title" validate:"required"`       // Title of the row.
}

// NewListSectionRow creates a new list section row with the given ID, title, and description.
func NewListSectionRow(id, title, description string) (*listSectionRow, error) {
	return &listSectionRow{
		Id:          id,
		Description: description,
		Title:       title,
	}, nil
}

// SetTitle sets the title of the list section row.
func (row *listSectionRow) SetTitle(title string) {
	row.Title = title
}

// SetDescription sets the description of the list section row.
func (row *listSectionRow) SetDescription(description string) {
	row.Description = description
}

// SetId sets the ID of the list section row.
func (row *listSectionRow) SetId(id string) {
	row.Id = id
}

// listMessageAction represents the action of the list message.
type listMessageAction struct {
	ButtonText string        `json:"button" validate:"required"`   // Text of the button.
	Sections   []ListSection `json:"sections" validate:"required"` // Sections in the list message.
}

// ListMessageBody represents the body of the list message.
type ListMessageBody struct {
	Text string `json:"text" validate:"required"` // Text of the body.
}

// ListMessageFooter represents the footer of the list message.
type ListMessageFooter struct {
	Text string `json:"text" validate:"required"` // Text of the footer.
}

type ListMessageHeader struct {
	Type string `json:"type" validate:"required"` // Type of the header.
	Text string `json:"text" validate:"required"` // Text of the header.
}

// listMessage represents an interactive list message.
type listMessage struct {
	Type   InteractiveMessageType `json:"type" validate:"required"`   // Type of the message.
	Action listMessageAction      `json:"action" validate:"required"` // Action of the message.
	Body   ListMessageBody        `json:"body,omitempty"`             // Body of the message.
	Footer *ListMessageFooter     `json:"footer,omitempty"`           // Footer of the message.
	Header *ListMessageHeader     `json:"header,omitempty"`           // Header of the message.
}

// ListMessageParams represents the parameters for creating a list message.
type ListMessageParams struct {
	ButtonText string `json:"-" validate:"required"` // Text of the button.
	BodyText   string `json:"-" validate:"required"` // Text of the body.
}

// NewListMessage creates a new list message with the given parameters.
func NewListMessage(params ListMessageParams) (*listMessage, error) {
	if err := internal.GetValidator().Struct(params); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}

	return &listMessage{
		Type: InteractiveMessageTypeList,
		Body: ListMessageBody{
			Text: params.BodyText,
		},
		Action: listMessageAction{
			ButtonText: params.ButtonText,
		},
	}, nil
}

// ListMessageApiPayload represents the API payload for the list message.
type ListMessageApiPayload struct {
	BaseMessagePayload
	Interactive listMessage `json:"interactive" validate:"required"` // Interactive message.
}

// AddSection adds a new section to the list message.
func (m *listMessage) AddSection(section *ListSection) {
	m.Action.Sections = append(m.Action.Sections, *section)
}

// SetBodyText sets the body text of the list message.
func (m *listMessage) SetBodyText(section *ListSection) {
	m.Body.Text = section.Title
}

func (m *listMessage) SetFooterText(footerText string) {
	m.Footer = &ListMessageFooter{
		Text: footerText,
	}
}

func (m *listMessage) SetHeaderText(headerText string) {
	m.Header = &ListMessageHeader{
		Type: "text", // only text is supported
		Text: headerText,
	}
}

// ToJson converts the list message to JSON.
func (m *listMessage) ToJson(configs ApiCompatibleJsonConverterConfigs) ([]byte, error) {
	if err := internal.GetValidator().Struct(configs); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}

	jsonData := ListMessageApiPayload{
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
