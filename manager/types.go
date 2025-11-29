package manager

import (
	"github.com/wapikit/wapi.go/pkg/components"
)

type NotificationReasonEnum string

const (
	NotificationReasonMessage NotificationReasonEnum = "message"
)

type NotificationPayloadErrorSchemaType struct {
	Code      int    `json:"code"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	ErrorData struct {
		Details string `json:"details"`
	} `json:"error_data,omitempty"`
}

type NotificationPayloadMessageContextSchemaType struct {
	Forwarded           bool   `json:"forwarded,omitempty"`
	FrequentlyForwarded bool   `json:"frequently_forwarded,omitempty"`
	From                string `json:"from,omitempty"`
	Id                  string `json:"id"`
	ReferredProduct     struct {
		CatalogId         string `json:"catalog_id"`
		ProductRetailerId string `json:"product_retailer_id"`
	} `json:"referred_product,omitempty"`
}

// ReferralInfo represents referral data from Click to WhatsApp ads
type ReferralInfo struct {
	SourceUrl      string                           `json:"source_url"`
	SourceType     AdInteractionSourceTypeEnum      `json:"source_type"`
	SourceId       string                           `json:"source_id"`
	Headline       string                           `json:"headline"`
	Body           string                           `json:"body"`
	ImageUrl       string                           `json:"image_url,omitempty"`
	VideoUrl       string                           `json:"video_url,omitempty"`
	ThumbnailUrl   string                           `json:"thumbnail_url,omitempty"`
	CtwaClid       string                           `json:"ctwa_clid,omitempty"`
	MediaType      AdInteractionSourceMediaTypeEnum `json:"media_type"`
	WelcomeMessage struct {
		Text string `json:"text"`
	} `json:"welcome_message,omitempty"`
}

type NotificationPayloadTextMessageSchemaType struct {
	Text struct {
		Body string `json:"body"`
	} `json:"text,omitempty"`
}

type NotificationPayloadAudioMessageSchemaType struct {
	Audio struct {
		Id       string `json:"id,omitempty"`
		MIMEType string `json:"mime_type,omitempty"`
		SHA256   string `json:"sha256,omitempty"`
		Url      string `json:"url,omitempty"`
		Voice    bool   `json:"voice,omitempty"` // Is Voice Recording?
	} `json:"audio,omitempty"`
}

type NotificationPayloadImageMessageSchemaType struct {
	Image struct {
		Id       string `json:"id"`
		MIMEType string `json:"mime_type"`
		SHA256   string `json:"sha256"`
		Caption  string `json:"caption,omitempty"`
		Url      string `json:"url,omitempty"`
	} `json:"image,omitempty"`
}

type NotificationPayloadButtonMessageSchemaType struct {
	Button struct {
		Payload string `json:"payload"`
		Text    string `json:"text"`
	} `json:"button,omitempty"`
}

type NotificationPayloadDocumentMessageSchemaType struct {
	Document struct {
		Id       string `json:"id"`
		MIMEType string `json:"mime_type"`
		SHA256   string `json:"sha256"`
		Caption  string `json:"caption,omitempty"`
		Filename string `json:"filename,omitempty"`
		Link     string `json:"link,omitempty"`
	} `json:"document,omitempty"`
}

type NotificationPayloadOrderMessageSchemaType struct {
	// OrderText string `json:"text"`
	Order struct {
		CatalogId    string `json:"catalog_id"`
		ProductItems []struct {
			ProductRetailerId string  `json:"product_retailer_id"`
			Quantity          int     `json:"quantity"`
			ItemPrice         float64 `json:"item_price"`
			Currency          string  `json:"currency"`
		} `json:"product_items"`
		Text *string `json:"text,omitempty"`
	} `json:"order,omitempty"`
}

type NotificationPayloadStickerMessageSchemaType struct {
	Sticker struct {
		Id       string `json:"id"`
		MIMEType string `json:"mime_type"`
		SHA256   string `json:"sha256"`
		Animated bool   `json:"animated"`
		Url      string `json:"url,omitempty"`
	} `json:"sticker,omitempty"`
}

type NotificationPayloadSystemMessageSchemaType struct {
	System struct {
		Body string                     `json:"body"`
		Type SystemNotificationTypeEnum `json:"type"`
		WaId string                     `json:"wa_id"`
	} `json:"system,omitempty"`
}

type NotificationPayloadVideoMessageSchemaType struct {
	Video struct {
		Id       string `json:"id"`
		MIMEType string `json:"mime_type"`
		SHA256   string `json:"sha256"`
		Caption  string `json:"caption,omitempty"`
		Filename string `json:"filename,omitempty"`
		Url      string `json:"url,omitempty"`
	} `json:"video,omitempty"`
}

type NotificationPayloadReactionMessageSchemaType struct {
	Reaction struct {
		MessageId string `json:"message_id"`
		Emoji     string `json:"emoji"`
	} `json:"reaction,omitempty"`
}

type NotificationPayloadInteractionMessageSchemaType struct {
	Interactive struct {
		Type                                                  InteractiveNotificationTypeEnum `json:"type"`
		NotificationPayloadButtonInteractionMessageSchemaType `json:",inline,omitempty"`
		NotificationPayloadListInteractionMessageSchemaType   `json:",inline,omitempty"`
	} `json:"interactive,omitempty"`
}

type NotificationPayloadButtonInteractionMessageSchemaType struct {
	ButtonReply struct {
		Id    string `json:"id"`
		Title string `json:"title"`
	} `json:"button_reply,omitempty"`
}

type NotificationPayloadListInteractionMessageSchemaType struct {
	ListReply struct {
		Id          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"list_reply,omitempty"`
}

type NotificationPayloadLocationMessageSchemaType struct {
	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Name      string  `json:"name,omitempty"`
		Address   string  `json:"address,omitempty"`
	} `json:"location,omitempty"`
}

type NotificationPayloadUnsupportedMessageSchemaType struct {
	Unsupported struct {
		Type string `json:"type"`
	} `json:"unsupported,omitempty"`
}

type NotificationPayloadContactMessageSchemaType struct {
	Contacts []components.Contact `json:"contacts"`
}

type NotificationMessageTypeEnum string

const (
	NotificationMessageTypeText        NotificationMessageTypeEnum = "text"
	NotificationMessageTypeAudio       NotificationMessageTypeEnum = "audio"
	NotificationMessageTypeImage       NotificationMessageTypeEnum = "image"
	NotificationMessageTypeButton      NotificationMessageTypeEnum = "button"
	NotificationMessageTypeDocument    NotificationMessageTypeEnum = "document"
	NotificationMessageTypeOrder       NotificationMessageTypeEnum = "order"
	NotificationMessageTypeSticker     NotificationMessageTypeEnum = "sticker"
	NotificationMessageTypeSystem      NotificationMessageTypeEnum = "system"
	NotificationMessageTypeVideo       NotificationMessageTypeEnum = "video"
	NotificationMessageTypeReaction    NotificationMessageTypeEnum = "reaction"
	NotificationMessageTypeInteractive NotificationMessageTypeEnum = "interactive"
	NotificationMessageTypeUnknown     NotificationMessageTypeEnum = "unknown"
	NotificationMessageTypeLocation    NotificationMessageTypeEnum = "location"
	NotificationMessageTypeContacts    NotificationMessageTypeEnum = "contacts"
	NotificationMessageTypeUnsupported NotificationMessageTypeEnum = "unsupported"
)

type InteractiveNotificationTypeEnum string

const (
	NotificationTypeButtonReply InteractiveNotificationTypeEnum = "button_reply"
	NotificationTypeListReply   InteractiveNotificationTypeEnum = "list_reply"
)

type AdInteractionSourceTypeEnum string

const (
	AdInteractionSourceTypeUnknown AdInteractionSourceTypeEnum = "unknown"
	// Add other ad interaction source types
)

type AdInteractionSourceMediaTypeEnum string

const (
	AdInteractionSourceMediaTypeImage AdInteractionSourceMediaTypeEnum = "image"
	AdInteractionSourceMediaTypeVideo AdInteractionSourceMediaTypeEnum = "video"
	// Add other ad interaction source media types
)

type SystemNotificationTypeEnum string

const (
	SystemNotificationTypeCustomerPhoneNumberChange SystemNotificationTypeEnum = "user_changed_number"
	SystemNotificationTypeCustomerIdentityChanged   SystemNotificationTypeEnum = "customer_identity_changed"
)

type SenderContact struct {
	WaId            string  `json:"wa_id"`
	Profile         Profile `json:"profile"`
	IdentityKeyHash string  `json:"identity_key_hash,omitempty"`
}

type Profile struct {
	Name string `json:"name"`
}

type WhatsappApiNotificationPayloadSchemaType struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

type Entry struct {
	Id      string   `json:"id"`
	Changes []Change `json:"changes"`
	Time    int64    `json:"time,omitempty"`
}

type WebhookFieldEnum string

const (
	WebhookFieldEnumAccountAlerts                   WebhookFieldEnum = "account_alerts"
	WebhookFieldEnumMessages                        WebhookFieldEnum = "messages"
	WebhookFieldEnumSecurity                        WebhookFieldEnum = "security"
	WebhookFieldEnumAccountUpdate                   WebhookFieldEnum = "account_update"
	WebhookFieldEnumAccountReview                   WebhookFieldEnum = "account_review"
	WebhookFieldEnumBusinessCapability              WebhookFieldEnum = "business_capability"
	WebhookFieldEnumMessageTemplateQuality          WebhookFieldEnum = "message_template_quality"
	WebhookFieldEnumMessageTemplateStatus           WebhookFieldEnum = "message_template_status"
	WebhookFieldEnumPhoneNumberName                 WebhookFieldEnum = "phone_number_name"
	WebhookFieldEnumPhoneNumberQuality              WebhookFieldEnum = "phone_number_quality"
	WebhookFieldEnumTemplateCategoryUpdate          WebhookFieldEnum = "template_category"
	WebhookFieldEnumUserPreferences                 WebhookFieldEnum = "user_preferences"
	WebhookFieldEnumMessageTemplateComponentsUpdate WebhookFieldEnum = "message_template_components_update"
	WebhookFieldEnumPaymentConfigurationUpdate      WebhookFieldEnum = "payment_configuration_update"
	WebhookFieldEnumSmbAppStateSync                 WebhookFieldEnum = "smb_app_state_sync"
	WebhookFieldEnumSmbMessageEchoes                WebhookFieldEnum = "smb_message_echoes"
	WebhookFieldEnumHistory                         WebhookFieldEnum = "history"
)

type TemplateMessageStatusUpdateEventEnum string

const (
	TemplateMessageStatusUpdateEventEnumApproved            TemplateMessageStatusUpdateEventEnum = "APPROVED"
	TemplateMessageStatusUpdateEventEnumRejected            TemplateMessageStatusUpdateEventEnum = "REJECTED"
	TemplateMessageStatusUpdateEventEnumFlaggedForDisabling TemplateMessageStatusUpdateEventEnum = "FLAGGED"
	TemplateMessageStatusUpdateEventEnumPaused              TemplateMessageStatusUpdateEventEnum = "PAUSED"
	TemplateMessageStatusUpdateEventEnumPendingDeletion     TemplateMessageStatusUpdateEventEnum = "PENDING_DELETION"
)

type TemplateMessageStatusUpdateDisableInfo struct {
	DisableDate string `json:"disable_date"`
}

type TemplateMessageStatusUpdateOtherInfo struct {
	Title string `json:"title"`
}

type TemplateMessageRejectionReasonEnum string

const (
	TemplateMessageRejectionReasonEnumAbusiveContent    TemplateMessageRejectionReasonEnum = "ABUSIVE_CONTENT"
	TemplateMessageRejectionReasonEnumIncorrectCategory TemplateMessageRejectionReasonEnum = "INCORRECT_CATEGORY"
	TemplateMessageRejectionReasonEnumInvalidFormat     TemplateMessageRejectionReasonEnum = "INVALID_FORMAT"
	TemplateMessageRejectionReasonEnumNone              TemplateMessageRejectionReasonEnum = "NONE"
	TemplateMessageRejectionReasonEnumScam              TemplateMessageRejectionReasonEnum = "SCAM"
)

type TemplateStatusUpdateValue struct {
	Event                   TemplateMessageStatusUpdateEventEnum   `json:"event"`
	MessageTemplateId       int64                                  `json:"message_template_id"`
	MessageTemplateName     string                                 `json:"message_template_name"`
	MessageTemplateLanguage string                                 `json:"message_template_language"`
	Reason                  TemplateMessageRejectionReasonEnum     `json:"reason"`
	DisableInfo             TemplateMessageStatusUpdateDisableInfo `json:"disable_info,omitempty"`
	OtherInfo               TemplateMessageStatusUpdateOtherInfo   `json:"other_info,omitempty"`
}

type TemplateCategoryUpdateValue struct {
	MessageTemplateId       int64                   `json:"message_template_id"`
	MessageTemplateName     string                  `json:"message_template_name"`
	MessageTemplateLanguage string                  `json:"message_template_language"`
	PreviousCategory        MessageTemplateCategory `json:"previous_category"`
	NewCategory             MessageTemplateCategory `json:"new_category"`
	CorrectCategory         MessageTemplateCategory `json:"correct_category"`
}

type TemplateQualityUpdateValue struct {
	PreviousQualityScore    string `json:"previous_quality_score"`
	NewQualityScore         string `json:"new_quality_score"`
	MessageTemplateId       int64  `json:"message_template_id"`
	MessageTemplateName     string `json:"message_template_name"`
	MessageTemplateLanguage string `json:"message_template_language"`
}

type MessageTemplateComponentsUpdateValue struct {
	MessageTemplateId       int64  `json:"message_template_id"`
	MessageTemplateName     string `json:"message_template_name"`
	MessageTemplateLanguage string `json:"message_template_language"`
	MessageTemplateElement  string `json:"message_template_element"`          // Template body text
	MessageTemplateTitle    string `json:"message_template_title,omitempty"`  // Only if template has text header
	MessageTemplateFooter   string `json:"message_template_footer,omitempty"` // Only if template has footer
	MessageTemplateButtons  []struct {
		MessageTemplateButtonType        string `json:"message_template_button_type"`
		MessageTemplateButtonText        string `json:"message_template_button_text"`
		MessageTemplateButtonUrl         string `json:"message_template_button_url,omitempty"`          // Only for url buttons
		MessageTemplateButtonPhoneNumber string `json:"message_template_button_phone_number,omitempty"` // Only for phone number buttons
	} `json:"message_template_buttons,omitempty"` // Only if template has url or phone number button
}

type PhoneNumberNameUpdateValue struct {
	DisplayPhoneNumber    string `json:"display_phone_number"`
	Decision              string `json:"decision"`
	RequestedVerifiedName string `json:"requested_verified_name"`
	RejectionReason       string `json:"rejection_reason"`
}

type PhoneNumberQualityUpdateCurrentLimitEnum string

const (
	PhoneNumberQualityUpdateCurrentLimitEnumTier50        PhoneNumberQualityUpdateCurrentLimitEnum = "TIER_50"
	PhoneNumberQualityUpdateCurrentLimitEnumTier250       PhoneNumberQualityUpdateCurrentLimitEnum = "TIER_250"
	PhoneNumberQualityUpdateCurrentLimitEnumTier1K        PhoneNumberQualityUpdateCurrentLimitEnum = "TIER_1K"
	PhoneNumberQualityUpdateCurrentLimitEnumTier10K       PhoneNumberQualityUpdateCurrentLimitEnum = "TIER_10K"
	PhoneNumberQualityUpdateCurrentLimitEnumTier100K      PhoneNumberQualityUpdateCurrentLimitEnum = "TIER_100K"
	PhoneNumberQualityUpdateCurrentLimitEnumTierUnlimited PhoneNumberQualityUpdateCurrentLimitEnum = "TIER_UNLIMITED"
)

type PhoneNumberQualityUpdateValue struct {
	DisplayPhoneNumber string                                   `json:"display_phone_number"`
	Event              string                                   `json:"event"`
	CurrentLimit       PhoneNumberQualityUpdateCurrentLimitEnum `json:"current_limit"`
}

type AccountAlertSeverityEnum string

const (
	AccountAlertSeverityEnumCritical AccountAlertSeverityEnum = "CRITICAL"
	AccountAlertSeverityEnumWarning  AccountAlertSeverityEnum = "WARNING"
)

type AccountAlertsValue struct {
	EntityType string `json:"entity_type"`
	EntityId   string `json:"entity_id"`
	AlertInfo  struct {
		AlertSeverity    AccountAlertSeverityEnum `json:"alert_severity"`
		AlertStatus      string                   `json:"alert_status"`
		AlertType        string                   `json:"alert_type"`
		AlertDescription string                   `json:"alert_description"`
	} `json:"alert_info"`
}

type AccountUpdateEventEnum string

type AccountUpdateBanInfo struct {
	WabaBanState []string `json:"waba_ban_state"`
	WabaBanDate  string   `json:"waba_ban_date"`
}

type AccountUpdateRestrictionInfo struct {
	RestrictionType string `json:"restriction_type"`
	Expiration      string `json:"expiration"`
}

type AccountUpdateViolationInfo struct {
	ViolationType string `json:"violation_type"`
}

const (
	AccountUpdateEventEnumVerifiedAccount    AccountUpdateEventEnum = "VERIFIED_ACCOUNT"
	AccountUpdateEventEnumDisabledAccount    AccountUpdateEventEnum = "DISABLED_UPDATE"
	AccountUpdateEventEnumAccountViolation   AccountUpdateEventEnum = "ACCOUNT_VIOLATION"
	AccountUpdateEventEnumAccountRestriction AccountUpdateEventEnum = "ACCOUNT_RESTRICTION"
	AccountUpdateEventEnumAccountDeleted     AccountUpdateEventEnum = "ACCOUNT_DELETED"
	AccountUpdateEventEnumPartnerRemoved     AccountUpdateEventEnum = "PARTNER_REMOVED"
)

type AccountUpdateWabaInfo struct {
	WabaId                     string   `json:"waba_id"`
	OwnerBusinessId            string   `json:"owner_business_id"`
	AdAccountLinked            string   `json:"ad_account_linked,omitempty"`             // Only for AD_ACCOUNT_LINKED event
	PartnerAppId               string   `json:"partner_app_id,omitempty"`                // Only for PARTNER_APP_INSTALLED, PARTNER_APP_UNINSTALLED
	SolutionId                 string   `json:"solution_id,omitempty"`                   // Only if customer onboarded via multi-partner solution
	SolutionPartnerBusinessIds []string `json:"solution_partner_business_ids,omitempty"` // Only if customer onboarded via multi-partner solution
}

type AccountUpdateAuthInternationalRateEligibility struct {
	ExceptionCountries []struct {
		CountryCode string `json:"country_code"`
		StartTime   int64  `json:"start_time"`
	} `json:"exception_countries,omitempty"`
	StartTime int64 `json:"start_time"`
}

type AccountUpdatePartnerClientCertificationInfo struct {
	ClientBusinessId string   `json:"client_business_id"`
	Status           string   `json:"status"`
	RejectionReasons []string `json:"rejection_reasons,omitempty"`
}

type AccountUpdateValue struct {
	PhoneNumber                      string                                         `json:"phone_number,omitempty"`
	Event                            AccountUpdateEventEnum                         `json:"event"`
	Country                          string                                         `json:"country,omitempty"`                             // Only for BUSINESS_PRIMARY_LOCATION_COUNTRY_UPDATE
	WabaInfo                         *AccountUpdateWabaInfo                         `json:"waba_info,omitempty"`                           // For various events
	ViolationInfo                    *AccountUpdateViolationInfo                    `json:"violation_info,omitempty"`                      // Only for ACCOUNT_VIOLATION
	AuthInternationalRateEligibility *AccountUpdateAuthInternationalRateEligibility `json:"auth_international_rate_eligibility,omitempty"` // Only for AUTH_INTL_PRICE_ELIGIBILITY_UPDATE
	BanInfo                          *AccountUpdateBanInfo                          `json:"ban_info,omitempty"`                            // Only for DISABLED_UPDATE
	RestrictionInfo                  *AccountUpdateRestrictionInfo                  `json:"restriction_info,omitempty"`                    // Only for ACCOUNT_RESTRICTION
	PartnerClientCertificationInfo   *AccountUpdatePartnerClientCertificationInfo   `json:"partner_client_certification_info,omitempty"`   // Only for PARTNER_CLIENT_CERTIFICATION_STATUS_UPDATE
}

type AccountReviewUpdateValue struct {
	Decision string `json:"decision"`
}

type BusinessCapabilityUpdateValue struct {
	MaxDailyConversationPerPhone int `json:"max_daily_conversation_per_phone"`
	MaxPhoneNumbersPerBusiness   int `json:"max_phone_numbers_per_business"`
}

type SecurityValue struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	Event              string `json:"event"`
	Requester          string `json:"requester"`
}

type UserPreferencesValue struct {
	UserPreferences []struct {
		WaId      string `json:"wa_id"`
		Detail    string `json:"detail"`
		Category  string `json:"category"` // e.g., "marketing_messages"
		Value     string `json:"value"`    // Preference value
		Timestamp int64  `json:"timestamp"`
	} `json:"user_preferences"`
}

type PaymentConfigurationUpdateValue struct {
	ConfigurationName string `json:"configuration_name"`
	ProviderName      string `json:"provider_name"`
	ProviderMid       string `json:"provider_mid"`
	Status            string `json:"status"`
	CreatedTimestamp  int64  `json:"created_timestamp"`
	UpdatedTimestamp  int64  `json:"updated_timestamp"`
}

type SmbAppStateSyncValue struct {
	MessagingProduct string `json:"messaging_product"`
	Metadata         struct {
		DisplayPhoneNumber string `json:"display_phone_number"`
		PhoneNumberId      string `json:"phone_number_id"`
	} `json:"metadata"`
	StateSync []struct {
		Type    string `json:"type"` // e.g., "contact"
		Contact struct {
			FullName    string `json:"full_name,omitempty"`
			FirstName   string `json:"first_name,omitempty"`
			PhoneNumber string `json:"phone_number,omitempty"`
		} `json:"contact,omitempty"`
		Action   string `json:"action"` // e.g., "add", "remove", "update"
		Metadata struct {
			Timestamp string `json:"timestamp"`
		} `json:"metadata"`
	} `json:"state_sync"`
}

type SmbMessageEchoesValue struct {
	MessagingProduct string `json:"messaging_product"`
	Metadata         struct {
		DisplayPhoneNumber string `json:"display_phone_number"`
		PhoneNumberId      string `json:"phone_number_id"`
	} `json:"metadata"`
	MessageEchoes []struct {
		From      string `json:"from"`
		To        string `json:"to"`
		Id        string `json:"id"`
		Timestamp string `json:"timestamp"`
		Type      string `json:"type"`
		// Message contents are dynamic based on type
		// Using map for flexibility
		MessageContent map[string]interface{} `json:"-"` // Will be populated from the type-specific field
	} `json:"message_echoes"`
}

type HistoryValue struct {
	MessagingProduct string `json:"messaging_product"`
	Metadata         struct {
		DisplayPhoneNumber string `json:"display_phone_number"`
		PhoneNumberId      string `json:"phone_number_id"`
	} `json:"metadata"`
	History []struct {
		Metadata struct {
			Phase      int `json:"phase"`
			ChunkOrder int `json:"chunk_order"`
			Progress   int `json:"progress"`
		} `json:"metadata"`
		Threads []struct {
			Id       string `json:"id"` // WhatsApp user phone number
			Messages []struct {
				From           string `json:"from"`
				To             string `json:"to,omitempty"` // Only included if SMB message echo
				Id             string `json:"id"`
				Timestamp      string `json:"timestamp"`
				Type           string `json:"type"`
				HistoryContext struct {
					Status string `json:"status"` // Message status
				} `json:"history_context"`
				// Message contents are dynamic based on type
				// Would need to be unmarshaled based on Type field
			} `json:"messages"`
		} `json:"threads"`
	} `json:"history"`
}

type Change struct {
	Value interface{}      `json:"value"`
	Field WebhookFieldEnum `json:"field"`
}

type MessagesValue struct {
	MessagingProduct string          `json:"messaging_product"`
	Metadata         Metadata        `json:"metadata"`
	Contacts         []SenderContact `json:"contacts,omitempty"`
	Statuses         []Status        `json:"statuses,omitempty"`
	Messages         []Message       `json:"messages,omitempty"`
	Errors           []Error         `json:"errors,omitempty"`
}

type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberId      string `json:"phone_number_id"`
}

type Status struct {
	Id                       string       `json:"id"`
	Conversation             Conversation `json:"conversation,omitempty"`
	Errors                   []Error      `json:"errors,omitempty"`
	Status                   string       `json:"status"`
	Timestamp                string       `json:"timestamp"`
	RecipientId              string       `json:"recipient_id"`
	RecipientType            string       `json:"recipient_type,omitempty"`              // Only included if message sent to a group
	RecipientParticipantId   string       `json:"recipient_participant_id,omitempty"`    // Only included if message sent to a group
	RecipientIdentityKeyHash string       `json:"recipient_identity_key_hash,omitempty"` // Only included if identity change check enabled
	BizOpaqueCallbackData    string       `json:"biz_opaque_callback_data,omitempty"`    // Only included if message sent with biz_opaque_callback_data
	Pricing                  Pricing      `json:"pricing,omitempty"`
}

type Conversation struct {
	Id                  string `json:"id"`
	ExpirationTimestamp string `json:"expiration_timestamp,omitempty"`
	Origin              Origin `json:"origin,omitempty"`
}

type Origin struct {
	Type                MessageStatusCategoryEnum `json:"type"`
	ExpirationTimestamp string                    `json:"expiration_timestamp,omitempty"`
}

type Pricing struct {
	Billable     bool                      `json:"billable"`
	PricingModel string                    `json:"pricing_model"`
	Category     MessageStatusCategoryEnum `json:"category"`
}

type Message struct {
	Id                                              string                                      `json:"id"`
	From                                            string                                      `json:"from"`
	Timestamp                                       string                                      `json:"timestamp"`
	Type                                            NotificationMessageTypeEnum                 `json:"type"`
	GroupId                                         string                                      `json:"group_id,omitempty"`
	Context                                         NotificationPayloadMessageContextSchemaType `json:"context,omitempty"`
	Errors                                          []Error                                     `json:"errors,omitempty"`
	Referral                                        *ReferralInfo                               `json:"referral,omitempty"`
	NotificationPayloadTextMessageSchemaType        `json:",inline"`
	NotificationPayloadAudioMessageSchemaType       `json:",inline"`
	NotificationPayloadImageMessageSchemaType       `json:",inline"`
	NotificationPayloadButtonMessageSchemaType      `json:",inline"`
	NotificationPayloadDocumentMessageSchemaType    `json:",inline"`
	NotificationPayloadOrderMessageSchemaType       `json:",inline"`
	NotificationPayloadStickerMessageSchemaType     `json:",inline"`
	NotificationPayloadSystemMessageSchemaType      `json:",inline"`
	NotificationPayloadVideoMessageSchemaType       `json:",inline"`
	NotificationPayloadReactionMessageSchemaType    `json:",inline"`
	NotificationPayloadLocationMessageSchemaType    `json:",inline"`
	NotificationPayloadContactMessageSchemaType     `json:",inline"`
	NotificationPayloadInteractionMessageSchemaType `json:",inline"`
	NotificationPayloadUnsupportedMessageSchemaType `json:",inline"`
}

type Error struct {
	Code      int    `json:"code"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	Href      string `json:"href"`
	ErrorData struct {
		Details string `json:"details"`
	} `json:"error_data"`
}

type MessageStatusCategoryEnum string

const (
	MessageStatusCategorySent MessageStatusCategoryEnum = "sent"
)

type MessageStatusEnum string

const (
	MessageStatusDelivered   MessageStatusEnum = "delivered"
	MessageStatusRead        MessageStatusEnum = "read"
	MessageStatusUnDelivered MessageStatusEnum = "undelivered"
	MessageStatusFailed      MessageStatusEnum = "failed"
	MessageStatusSent        MessageStatusEnum = "sent"
)

// PaginationCursors represents the before/after cursors used in cursor-based pagination
type PaginationCursors struct {
	Before string `json:"before,omitempty"`
	After  string `json:"after,omitempty"`
}

// PaginationDetails contains the pagination metadata returned by WhatsApp API
type PaginationDetails struct {
	Cursors  PaginationCursors `json:"cursors,omitempty"`
	Next     string            `json:"next,omitempty"`
	Previous string            `json:"previous,omitempty"`
}

// PaginationInput represents the pagination parameters that can be passed to API calls
type PaginationInput struct {
	Limit  int    `json:"limit,omitempty"`  // Number of results per page
	After  string `json:"after,omitempty"`  // Cursor for next page
	Before string `json:"before,omitempty"` // Cursor for previous page
}

// PaginatedResponse is a generic wrapper for paginated responses
type PaginatedResponse[T any] struct {
	Data   []T               `json:"data"`
	Paging PaginationDetails `json:"paging,omitempty"`
}

// HasNextPage checks if there's a next page available
func (pd *PaginationDetails) HasNextPage() bool {
	return pd.Cursors.After != "" || pd.Next != ""
}

// HasPreviousPage checks if there's a previous page available
func (pd *PaginationDetails) HasPreviousPage() bool {
	return pd.Cursors.Before != "" || pd.Previous != ""
}

// GetNextCursor returns the cursor for the next page
func (pd *PaginationDetails) GetNextCursor() string {
	return pd.Cursors.After
}

// GetPreviousCursor returns the cursor for the previous page
func (pd *PaginationDetails) GetPreviousCursor() string {
	return pd.Cursors.Before
}
