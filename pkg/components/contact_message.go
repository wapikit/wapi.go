package components

import (
	"encoding/json"
	"fmt"

	"github.com/gTahidi/wapi.go/internal"
)

type AddressType string

const (
	HomeAddress AddressType = "HOME"
	WorkAddress AddressType = "WORK"
)

type UrlType string

const (
	HomeUrl AddressType = "HOME"
	WorkUrl AddressType = "WORK"
)

type EmailType string

const (
	HomeEmail EmailType = "HOME"
	WorkEmail EmailType = "WORK"
)

type PhoneType string

const (
	CellPhone   PhoneType = "CELL"
	MainPhone   PhoneType = "MAIN"
	IphonePhone PhoneType = "IPHONE"
	HomePhone   PhoneType = "HOME"
	WorkPhone   PhoneType = "WORK"
)

type ContactAddress struct {
	Street      string      `json:"street,omitempty"`
	City        string      `json:"city,omitempty"`
	State       string      `json:"state,omitempty"`
	Zip         string      `json:"zip,omitempty"`
	Country     string      `json:"country,omitempty"`
	CountryCode string      `json:"countryCode,omitempty"`
	Type        AddressType `json:"type" validate:"required"`
}

type ContactName struct {
	FormattedName string `json:"formatted_name" validate:"required"`
	FirstName     string `json:"first_name,omitempty"`
	LastName      string `json:"last_name,omitempty"`
	MiddleName    string `json:"middle_name,omitempty"`
	Suffix        string `json:"suffix,omitempty"`
	Prefix        string `json:"prefix,omitempty"`
}

type ContactOrg struct {
	Company    string `json:"company,omitempty"`
	Title      string `json:"title,omitempty"`
	Department string `json:"department,omitempty"`
}

type ContactEmail struct {
	Email string    `json:"email,omitempty"`
	Type  EmailType `json:"type,omitempty"`
}

type ContactPhone struct {
	Phone string    `json:"phone,omitempty"`
	WaId  string    `json:"wa_id,omitempty"`
	Type  PhoneType `json:"type" validate:"required"`
}

type ContactUrl struct {
	Url  string  `json:"url" validate:"required"`
	Type UrlType `json:"type" validate:"required"`
}

type Contact struct {
	Name      ContactName      `json:"name" validate:"required"`
	Org       ContactOrg       `json:"org,omitempty"`
	Addresses []ContactAddress `json:"addresses,omitempty"`
	Urls      []ContactUrl     `json:"urls,omitempty"`
	Emails    []ContactEmail   `json:"emails,omitempty"`
	Phones    []ContactPhone   `json:"phones,omitempty"`
	Birthday  string           `json:"birthday,omitempty"`
}

func NewContact(name ContactName) *Contact {
	return &Contact{
		Name: name,
	}
}

func (contact *ContactMessage) AddContact(params Contact) {
	contact.Contacts = append(contact.Contacts, params)
}

func (contact *Contact) SetFirstName(firstName string) {
	contact.Name.FirstName = firstName
}

func (contact *Contact) SetLastName(lastName string) {
	contact.Name.LastName = lastName
}

func (contact *Contact) SetMiddleName(middleName string) {
	contact.Name.MiddleName = middleName
}

func (contact *Contact) SetOrg(params ContactOrg) {
	contact.Org = params
}

func (contact *Contact) AddPhone(params ContactPhone) {
	contact.Phones = append(contact.Phones, params)
}

func (contact *Contact) AddEmail(params ContactEmail) {
	contact.Emails = append(contact.Emails, params)
}

func (contact *Contact) AddUrl(params ContactUrl) {
	contact.Urls = append(contact.Urls, params)
}

// ! TODO: add regex check here in the params
func (contact *Contact) SetBirthday(params string) {
	contact.Birthday = params
}

type ContactMessage struct {
	Contacts []Contact `json:"contacts" validate:"required"`
}

type ContactMessageConfigs struct {
	Name string `json:"name" validate:"required"`
}

type ContactMessageApiPayload struct {
	BaseMessagePayload
	Contacts []Contact `json:"contacts" validate:"required"`
}

func NewContactMessage(configs []Contact) (*ContactMessage, error) {
	return &ContactMessage{
		Contacts: configs,
	}, nil
}

func (m *ContactMessage) ToJson(configs ApiCompatibleJsonConverterConfigs) ([]byte, error) {
	if err := internal.GetValidator().Struct(configs); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}
	jsonData := ContactMessageApiPayload{
		BaseMessagePayload: NewBaseMessagePayload(configs.SendToPhoneNumber, MessageTypeContact),
		Contacts:           m.Contacts,
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
