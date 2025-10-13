package manager

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"

    "github.com/wapikit/wapi.go/internal"
    "github.com/wapikit/wapi.go/internal/request_client"
)

// MessageTemplateStatus represents the status of a WhatsApp Business message template.
type MessageTemplateStatus string

const (
	MessageTemplateStatusApproved MessageTemplateStatus = "APPROVED"
	MessageTemplateStatusRejected MessageTemplateStatus = "REJECTED"
	MessageTemplateStatusPending  MessageTemplateStatus = "PENDING"
)

// MessageTemplateCategory represents the category of a WhatsApp Business message template.
type MessageTemplateCategory string

const (
	MessageTemplateCategoryUtility        MessageTemplateCategory = "UTILITY"
	MessageTemplateCategoryMarketing      MessageTemplateCategory = "MARKETING"
	MessageTemplateCategoryAuthentication MessageTemplateCategory = "AUTHENTICATION"
)

// WhatsAppBusinessMessageTemplateNode represents a WhatsApp Business message template.
type WhatsAppBusinessMessageTemplateNode struct {
	Id                         string                                    `json:"id,omitempty"`
	Category                   MessageTemplateCategory                   `json:"category,omitempty"`
	Components                 []WhatsAppBusinessHSMWhatsAppHSMComponent `json:"components,omitempty"`
	CorrectCategory            string                                    `json:"correct_category,omitempty"`
	CtaUrlLinkTrackingOptedOut bool                                      `json:"cta_url_link_tracking_opted_out,omitempty"`
	Language                   string                                    `json:"language,omitempty"`
	LibraryTemplateName        string                                    `json:"library_template_name,omitempty"`
	MessageSendTtlSeconds      int                                       `json:"message_send_ttl_seconds,omitempty"`
	Name                       string                                    `json:"name,omitempty"`
	PreviousCategory           string                                    `json:"previous_category,omitempty"`
	QualityScore               TemplateMessageQualityScore               `json:"quality_score,omitempty"`
	RejectedReason             string                                    `json:"rejected_reason,omitempty"`
	Status                     MessageTemplateStatus                     `json:"status,omitempty"`
}

// TemplateManager is responsible for managing WhatsApp Business message templates.
type TemplateManager struct {
	businessAccountId string
	apiAccessToken    string
	requester         *request_client.RequestClient
}

// TemplateManagerConfig represents the configuration for creating a new TemplateManager.
type TemplateManagerConfig struct {
	BusinessAccountId string
	ApiAccessToken    string
	Requester         *request_client.RequestClient
}

// NewTemplateManager creates a new TemplateManager with the given configuration.
func NewTemplateManager(config *TemplateManagerConfig) *TemplateManager {
	return &TemplateManager{
		apiAccessToken:    config.ApiAccessToken,
		businessAccountId: config.BusinessAccountId,
		requester:         config.Requester,
	}
}

// WhatsAppBusinessTemplatesFetchResponseEdge represents the response structure for fetching templates.
type WhatsAppBusinessTemplatesFetchResponseEdge struct {
	Data   []WhatsAppBusinessMessageTemplateNode      `json:"data,omitempty"`
	Paging internal.WhatsAppBusinessApiPaginationMeta `json:"paging,omitempty"`
}

// TemplateMessageComponentCard represents a card component in a message template.
type TemplateMessageComponentCard struct {
	// Add card-specific fields if needed.
	Components []WhatsAppBusinessHSMWhatsAppHSMComponent `json:"components,omitempty"`
}

// TemplateMessageButtonType represents the type of a button.
type TemplateMessageButtonType string

const (
	TemplateMessageButtonTypeQuickReply          TemplateMessageButtonType = "QUICK_REPLY"
	TemplateMessageButtonTypeUrl                 TemplateMessageButtonType = "URL"
	TemplateMessageButtonTypePhoneNumber         TemplateMessageButtonType = "PHONE_NUMBER"
	TemplateMessageButtonTypeCopyCode            TemplateMessageButtonType = "COPY_CODE"
	TemplateMessageButtonTypeCatalog             TemplateMessageButtonType = "CATALOG"
	TemplateMessageButtonTypeMultiProductMessage TemplateMessageButtonType = "MPM"
)

// TemplateMessageComponentButton represents a button component in a message template.
type TemplateMessageComponentButton struct {
	Type        TemplateMessageButtonType `json:"type,omitempty"`
	Text        string                    `json:"text,omitempty"`
	PhoneNumber string                    `json:"phone_number,omitempty"` // required when Type = PHONE_NUMBER
	Example     []string                  `json:"example,omitempty"`      // required when Type = URL and button uses a variable
	Url         string                    `json:"url,omitempty"`          // required when Type = URL
}

// NamedParamExample represents a single named parameter and its example value.
type NamedParamExample struct {
	ParamName string `json:"param_name"`
	Example   string `json:"example"`
}

// TemplateMessageComponentExample represents an example object for a template component.
type TemplateMessageComponentExample struct {
	HeaderHandle          *[]string            `json:"header_handle,omitempty"`            // for media headers (IMAGE, VIDEO, DOCUMENT)
	HeaderTextNamedParams *[]NamedParamExample `json:"header_text_named_params,omitempty"` // for named header params
	HeaderText            *[]string            `json:"header_text,omitempty"`              // for text headers (positional examples)
	BodyText              *[][]string          `json:"body_text,omitempty"`                // for body components (array of arrays for positional examples)
	BodyTextNamedParams   *[]NamedParamExample `json:"body_text_named_params,omitempty"`
}

// TemplateMessageLimitedTimeOfferParameter represents a limited time offer parameter.
type TemplateMessageLimitedTimeOfferParameter struct {
	Text          string `json:"text,omitempty"`
	HasExpiration bool   `json:"has_expiration,omitempty"`
}

// MessageTemplateComponentType represents the type of a template component.
type MessageTemplateComponentType string

const (
	MessageTemplateComponentTypeGreeting         MessageTemplateComponentType = "GREETING"
	MessageTemplateComponentTypeHeader           MessageTemplateComponentType = "HEADER"
	MessageTemplateComponentTypeBody             MessageTemplateComponentType = "BODY"
	MessageTemplateComponentTypeFooter           MessageTemplateComponentType = "FOOTER"
	MessageTemplateComponentTypeButtons          MessageTemplateComponentType = "BUTTONS"
	MessageTemplateComponentTypeCarousel         MessageTemplateComponentType = "CAROUSEL" // this appears in case of product caraousel
	MessageTemplateComponentTypeLimitedTimeOffer MessageTemplateComponentType = "LIMITED_TIME_OFFER"
)

// MessageTemplateComponentFormat represents the format of a template component.
type MessageTemplateComponentFormat string

const (
	MessageTemplateComponentFormatText     MessageTemplateComponentFormat = "TEXT"
	MessageTemplateComponentFormatImage    MessageTemplateComponentFormat = "IMAGE"
	MessageTemplateComponentFormatDocument MessageTemplateComponentFormat = "DOCUMENT"
	MessageTemplateComponentFormatVideo    MessageTemplateComponentFormat = "VIDEO"
	MessageTemplateComponentFormatLocation MessageTemplateComponentFormat = "LOCATION"
)

// WhatsAppBusinessHSMWhatsAppHSMComponent represents a component in a message template.
// Note: The "Type" field here now uses MessageTemplateComponentType.
type WhatsAppBusinessHSMWhatsAppHSMComponent struct {
	AddSecurityRecommendation bool                                      `json:"add_security_recommendation,omitempty"`
	Buttons                   []TemplateMessageComponentButton          `json:"buttons,omitempty"`
	Cards                     []TemplateMessageComponentCard            `json:"cards,omitempty"`
	CodeExpirationMinutes     int                                       `json:"code_expiration_minutes,omitempty"`
	Example                   *TemplateMessageComponentExample          `json:"example,omitempty"`
	Format                    MessageTemplateComponentFormat            `json:"format,omitempty"`
	LimitedTimeOffer          *TemplateMessageLimitedTimeOfferParameter `json:"limited_time_offer,omitempty"`
	Text                      string                                    `json:"text,omitempty"`
	Type                      MessageTemplateComponentType              `json:"type,omitempty"`
}

// TemplateMessageQualityScore represents the quality score of a template.
type TemplateMessageQualityScore struct {
	Date    int      `json:"date,omitempty"`
	Reasons []string `json:"reasons,omitempty"`
	Score   int      `json:"score,omitempty"`
}

// FetchAll fetches all WhatsApp Business message templates.
func (manager *TemplateManager) FetchAll() (*WhatsAppBusinessTemplatesFetchResponseEdge, error) {
	apiRequest := manager.requester.NewApiRequest(strings.Join([]string{manager.businessAccountId, "/", "message_templates"}, ""), http.MethodGet)

	fields := []string{
		"id", "category", "components", "correct_category", "cta_url_link_tracking_opted_out",
		"language", "library_template_name", "message_send_ttl_seconds", "name", "previous_category",
		"quality_score", "rejected_reason", "status", "sub_category",
	}

	for _, field := range fields {
		apiRequest.AddField(request_client.ApiRequestQueryParamField{
			Name:    field,
			Filters: map[string]string{},
		})
	}

	response, err := apiRequest.Execute()
	if err != nil {
		return nil, err
	}

	var responseToReturn WhatsAppBusinessTemplatesFetchResponseEdge
	json.Unmarshal([]byte(response), &responseToReturn)
	return &responseToReturn, nil
}

// Fetch fetches a single WhatsApp Business message template by its ID.
func (manager *TemplateManager) Fetch(Id string) (*WhatsAppBusinessMessageTemplateNode, error) {
	apiRequest := manager.requester.NewApiRequest(strings.Join([]string{Id}, ""), http.MethodGet)
	fields := []string{
		"id", "category", "components", "correct_category", "cta_url_link_tracking_opted_out",
		"language", "library_template_name", "message_send_ttl_seconds", "name", "previous_category",
		"quality_score", "rejected_reason", "status", "sub_category",
	}
	for _, field := range fields {
		apiRequest.AddField(request_client.ApiRequestQueryParamField{
			Name:    field,
			Filters: map[string]string{},
		})
	}
	response, err := apiRequest.Execute()
	if err != nil {
		return nil, err
	}
	var responseToReturn WhatsAppBusinessMessageTemplateNode
	json.Unmarshal([]byte(response), &responseToReturn)
	return &responseToReturn, nil
}

// WhatsappMessageTemplateButtonCreateRequestBody represents the request body for creating a button.
type WhatsappMessageTemplateButtonCreateRequestBody struct {
	Type                 string   `json:"type,omitempty"`
	Text                 string   `json:"text,omitempty"`
	Url                  string   `json:"url,omitempty"`
	PhoneNumber          string   `json:"phone_number,omitempty"`
	Example              []string `json:"example,omitempty"` // For URL buttons with variables.
	FlowId               string   `json:"flow_id,omitempty"`
	ZeroTapTermsAccepted bool     `json:"zero_tap_terms_accepted,omitempty"`
}

// WhatsappMessageTemplateButtonCreateRequestBody alias used in component creation.
type WhatsappMessageTemplateButtonCreateRequestBodyAlias = WhatsappMessageTemplateButtonCreateRequestBody

// WhatsappMessageTemplateComponentCreateOrUpdateRequestBody represents the request body for creating/updating a component.
type WhatsappMessageTemplateComponentCreateOrUpdateRequestBody struct {
    Type    MessageTemplateComponentType                     `json:"type,omitempty"`
    Format  MessageTemplateComponentFormat                   `json:"format,omitempty"`
    Text    string                                           `json:"text,omitempty"`
    Buttons []WhatsappMessageTemplateButtonCreateRequestBody `json:"buttons,omitempty"`
    Example *TemplateMessageComponentExample                 `json:"example,omitempty"`
}

// validateCatalogAndMPMButtons ensures CATALOG and MPM buttons have required params
// and disallow unsupported fields.
func validateCatalogAndMPMButtons(components []WhatsappMessageTemplateComponentCreateOrUpdateRequestBody) error {
    for _, c := range components {
        if c.Type != MessageTemplateComponentTypeButtons {
            continue
        }
        for _, b := range c.Buttons {
            switch TemplateMessageButtonType(b.Type) {
            case TemplateMessageButtonTypeCatalog, TemplateMessageButtonTypeMultiProductMessage:
                if b.Text == "" {
                    return fmt.Errorf("template button of type %s requires non-empty text", b.Type)
                }
                // URL and phone_number are not relevant for catalog/mpm; ignore if present
            default:
                // no-op
            }
        }
    }
    return nil
}

// WhatsappMessageTemplateCreateRequestBody represents the request body for creating a message template.
type WhatsappMessageTemplateCreateRequestBody struct {
	AllowCategoryChange         bool                                                        `json:"allow_category_change,omitempty"`
	Category                    string                                                      `json:"category,omitempty" validate:"required"`
	Components                  []WhatsappMessageTemplateComponentCreateOrUpdateRequestBody `json:"components" validate:"required"`
	Name                        string                                                      `json:"name,omitempty" validate:"required"`
	Language                    string                                                      `json:"language" validate:"required"`
	LibraryTemplateName         string                                                      `json:"library_template_name,omitempty"`
	LibraryTemplateButtonInputs []WhatsappMessageTemplateButtonCreateRequestBody            `json:"library_template_button_inputs,omitempty"`
}

// AddComponent appends a component to the template creation request body.
func (body *WhatsappMessageTemplateCreateRequestBody) AddComponent(component WhatsappMessageTemplateComponentCreateOrUpdateRequestBody) {
	body.Components = append(body.Components, component)
}

// MessageTemplateCreationResponse represents the response after creating a template.
type MessageTemplateCreationResponse struct {
	Id       string                  `json:"id,omitempty"`
	Status   MessageTemplateStatus   `json:"status,omitempty"`
	Category MessageTemplateCategory `json:"category,omitempty"`
}

// Create sends a creation request for a message template.
func (manager *TemplateManager) Create(body WhatsappMessageTemplateCreateRequestBody) (*MessageTemplateCreationResponse, error) {
    // Pre-validate catalog and multi-product message buttons
    if err := validateCatalogAndMPMButtons(body.Components); err != nil {
        return nil, err
    }
    apiRequest := manager.requester.NewApiRequest(strings.Join([]string{manager.businessAccountId, "/", "message_templates"}, ""), http.MethodPost)
    jsonBody, err := json.Marshal(body)
    if err != nil {
        return nil, err
    }
	apiRequest.SetBody(string(jsonBody))
	response, err := apiRequest.Execute()
	if err != nil {
		return nil, err
	}

	var responseToReturn MessageTemplateCreationResponse
	json.Unmarshal([]byte(response), &responseToReturn)
	return &responseToReturn, nil
}

// WhatsAppBusinessAccountMessageTemplateUpdateRequestBody represents the request body for updating a template.
type WhatsAppBusinessAccountMessageTemplateUpdateRequestBody struct {
	Components            []WhatsappMessageTemplateComponentCreateOrUpdateRequestBody `json:"components,omitempty"`
	Category              string                                                      `json:"category,omitempty"`
	MessageSendTtlSeconds int                                                         `json:"message_send_ttl_seconds,omitempty"`
}

// Update sends an update request for a template.
func (manager *TemplateManager) Update(templateId string, updates WhatsAppBusinessAccountMessageTemplateUpdateRequestBody) (*MessageTemplateCreationResponse, error) {
    // Pre-validate catalog and multi-product message buttons
    if len(updates.Components) > 0 {
        if err := validateCatalogAndMPMButtons(updates.Components); err != nil {
            return nil, err
        }
    }
    apiRequest := manager.requester.NewApiRequest(strings.Join([]string{templateId}, ""), http.MethodPost)
    jsonBody, err := json.Marshal(updates)
    if err != nil {
        return nil, err
    }
	apiRequest.SetBody(string(jsonBody))
	response, err := apiRequest.Execute()
	if err != nil {
		return nil, err
	}

	var responseToReturn MessageTemplateCreationResponse
	json.Unmarshal([]byte(response), &responseToReturn)
	return &responseToReturn, nil
}

// WhatsAppBusinessAccountMessageTemplateDeleteRequestBody represents the request body for deleting a template.
type WhatsAppBusinessAccountMessageTemplateDeleteRequestBody struct {
	HsmId string `json:"hsm_id,omitempty"`
	Name  string `json:"name,omitempty"`
}

// Delete dissociates a template (delete implementation to be added as needed).
func (tm *TemplateManager) Delete(id string) {
	// Implement deletion logic here as per API requirements.
}

// WhatsAppBusinessAccountMessageTemplatePreviewButton represents a preview button.
type WhatsAppBusinessAccountMessageTemplatePreviewButton struct {
	AutoFillText string `json:"auto_fill_text,omitempty"`
	Text         string `json:"text,omitempty"`
}

// TemplateMessagePreviewNode represents a preview node.
type TemplateMessagePreviewNode struct {
	Body     string                                                `json:"body,omitempty"`
	Buttons  []WhatsAppBusinessAccountMessageTemplatePreviewButton `json:"buttons,omitempty"`
	Footer   string                                                `json:"footer,omitempty"`
	Header   string                                                `json:"header,omitempty"`
	Language string                                                `json:"language,omitempty"`
}

// TemplateMessagePreviewEdge represents the preview response.
type TemplateMessagePreviewEdge struct {
	Data   []TemplateMessagePreviewNode               `json:"data,omitempty"`
	Paging internal.WhatsAppBusinessApiPaginationMeta `json:"paging,omitempty"`
}

// TemplateMigrationResponse represents the migration response.
type TemplateMigrationResponse struct {
	MigratedTemplates []string          `json:"migrated_templates,omitempty"`
	FailedTemplates   map[string]string `json:"failed_templates,omitempty"`
}

// MigrateFromOtherBusinessAccount migrates templates from another business account.
func (manager *TemplateManager) MigrateFromOtherBusinessAccount(sourcePageNumber int, sourceWabaId int) (*TemplateMigrationResponse, error) {
	apiRequest := manager.requester.NewApiRequest(strings.Join([]string{manager.businessAccountId, "migrate_message_templates"}, "/"), http.MethodGet)
	apiRequest.AddQueryParam("page_number", string(sourcePageNumber))
	apiRequest.AddQueryParam("source_waba_id", string(sourceWabaId))
	response, err := apiRequest.Execute()
	if err != nil {
		return nil, err
	}
	var responseToReturn TemplateMigrationResponse
	json.Unmarshal([]byte(response), &responseToReturn)
	return &responseToReturn, nil
}
