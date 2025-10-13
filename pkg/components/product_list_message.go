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
    Sections          []ProductSection `json:"sections" validate:"required"` // minimum 1 and maximum 10
    CatalogId         string           `json:"catalog_id" validate:"required"`
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
    Header *ProductListMessageHeader `json:"header,omitempty"`
    Type   InteractiveMessageType    `json:"type" validate:"required"`
}

func (message *ProductListMessage) AddSection(section ProductSection) {
	fmt.Println("message", message, "section", section)
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

// SetProductRetailerId is deprecated. Product retailer IDs belong to each item.
// This method is kept for backward compatibility and is now a no-op.
func (message *ProductListMessage) SetProductRetailerId(productRetailerId string) {
    // no-op
}

func (message *ProductListMessage) SetFooter(text string) {
	message.Footer = &ProductListMessageFooter{
		Text: text,
	}
}

func (message *ProductListMessage) SetHeader(text string) {
    message.Header = &ProductListMessageHeader{
        Type: ProductListMessageHeaderTypeText,
        Text: text,
    }
}

// ProductListMessageParams represents the parameters for creating a product list message.
type ProductListMessageParams struct {
    CatalogId         string `validate:"required"`
    // Deprecated: action-level product_retailer_id. Use item-level IDs in sections.
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
            Sections:          params.Sections,
        },
    }, nil
}

// ToJson converts the product list message to JSON.
func (m *ProductListMessage) ToJson(configs ApiCompatibleJsonConverterConfigs) ([]byte, error) {
    if err := internal.GetValidator().Struct(configs); err != nil {
        return nil, fmt.Errorf("error validating configs: %v", err)
    }
    // Validate message structure and section/item limits
    if err := internal.GetValidator().Struct(m); err != nil {
        return nil, fmt.Errorf("error validating product list message: %v", err)
    }
    // Enforce Meta limits: ≤10 sections, ≤30 items per section
    if len(m.Action.Sections) == 0 || len(m.Action.Sections) > 10 {
        return nil, fmt.Errorf("product_list must contain between 1 and 10 sections")
    }
    for i, s := range m.Action.Sections {
        if len(s.Products) == 0 || len(s.Products) > 30 {
            return nil, fmt.Errorf("section %d must contain between 1 and 30 products", i)
        }
    }

	jsonData := ProductListMessageApiPayload{
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
