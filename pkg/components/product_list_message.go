package components

import (
	"encoding/json"
	"fmt"

	"github.com/wapikit/wapi.go/internal"
)

type Product struct {
	RetailerId string `json:"product_retailer_id" validate:"required"`
}

func (p *Product) SetRetailerId(id string) {
	p.RetailerId = id
}

type ProductSection struct {
	Title    string    `json:"title" validate:"required"`
	Products []Product `json:"product_items" validate:"required"`
}

func (ps *ProductSection) SetTitle(title string) {
	ps.Title = title
}

func (ps *ProductSection) AddProduct(product Product) {
	ps.Products = append(ps.Products, product)
}

type ProductListMessageAction struct {
	Sections  []ProductSection `json:"sections" validate:"required"` // minimum 1 and maximum 10
	CatalogId string           `json:"catalog_id" validate:"required"`
	// ProductRetailerId is not part of Meta's product_list action and is not
	// sent. It is kept so existing code compiles.
	ProductRetailerId string `json:"-"`
}

func (a *ProductListMessageAction) AddSection(section ProductSection) {
	a.Sections = append(a.Sections, section)
}

type ProductListMessageBody struct {
	Text string `json:"text" validate:"required"`
}

type ProductListMessageFooter struct {
	Text string `json:"text" validate:"required"`
}

type ProductListMessageHeaderType string

const (
	ProductListMessageHeaderTypeText ProductListMessageHeaderType = "text"
)

// ! TODO: support more header types
type ProductListMessageHeader struct {
	Type ProductListMessageHeaderType `json:"type" validate:"required"`
	Text string                       `json:"text" validate:"required"`
}

// ProductListMessage represents a product list message.
type ProductListMessage struct {
	Action ProductListMessageAction  `json:"action" validate:"required"`
	Body   ProductListMessageBody    `json:"body" validate:"required"`
	Footer *ProductListMessageFooter `json:"footer,omitempty"`
	Header ProductListMessageHeader  `json:"header,omitempty"`
	Type   InteractiveMessageType    `json:"type" validate:"required"`
}

// MarshalJSON leaves out the header when SetHeader was never called, instead
// of sending an empty header object.
func (message ProductListMessage) MarshalJSON() ([]byte, error) {
	type plain ProductListMessage
	aux := struct {
		plain
		Header *ProductListMessageHeader `json:"header,omitempty"`
	}{plain: plain(message)}
	if message.Header != (ProductListMessageHeader{}) {
		header := message.Header
		aux.Header = &header
	}
	return json.Marshal(aux)
}

func (message *ProductListMessage) AddSection(section ProductSection) {
	message.Action.Sections = append(message.Action.Sections, section)
}

func (message *ProductListMessage) SetBody(text string) {
	message.Body = ProductListMessageBody{
		Text: text,
	}
}

func (message *ProductListMessage) SetCatalogId(catalogId string) {
	message.Action.CatalogId = catalogId
}

func (message *ProductListMessage) SetProductRetailerId(productRetailerId string) {
	message.Action.ProductRetailerId = productRetailerId
}

func (message *ProductListMessage) SetFooter(text string) {
	message.Footer = &ProductListMessageFooter{
		Text: text,
	}
}

func (message *ProductListMessage) SetHeader(text string) {
	message.Header = ProductListMessageHeader{
		Type: ProductListMessageHeaderTypeText,
		Text: text,
	}
}

// ProductListMessageParams represents the parameters for creating a product list message.
type ProductListMessageParams struct {
	CatalogId string `validate:"required"`
	// ProductRetailerId is optional and not sent; Meta's product_list action
	// has no such field.
	ProductRetailerId string
	BodyText          string `validate:"required"`
	Sections          []ProductSection
}

// ProductListMessageApiPayload represents the API payload for a product list message.
type ProductListMessageApiPayload struct {
	BaseMessagePayload
	Interactive ProductListMessage `json:"interactive" validate:"required"`
}

// NewProductListMessage creates a new product list message.
func NewProductListMessage(params ProductListMessageParams) (*ProductListMessage, error) {
	if err := internal.GetValidator().Struct(params); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}

	return &ProductListMessage{
		Type: InteractiveMessageTypeProductList,
		Body: ProductListMessageBody{
			Text: params.BodyText,
		},
		Action: ProductListMessageAction{
			CatalogId:         params.CatalogId,
			ProductRetailerId: params.ProductRetailerId,
			Sections:          params.Sections,
		},
	}, nil
}

// ToJson converts the product list message to JSON.
func (m *ProductListMessage) ToJson(configs ApiCompatibleJsonConverterConfigs) ([]byte, error) {
	if err := internal.GetValidator().Struct(configs); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}

	jsonData := ProductListMessageApiPayload{
		BaseMessagePayload: NewBaseMessagePayload(configs, MessageTypeInteractive),
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
