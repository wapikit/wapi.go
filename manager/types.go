package manager

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

type NotificationPayloadTextMessageSchemaType struct {
	Text struct {
		Body string `json:"body"`
	} `json:"text,omitempty"`
	Referral struct {
		SourceUrl    string                           `json:"source_url"`
		SourceType   AdInteractionSourceTypeEnum      `json:"source_type"`
		SourceId     string                           `json:"source_id"`
		Headline     string                           `json:"headline"`
		Body         string                           `json:"body"`
		ImageUrl     string                           `json:"image_url,omitempty"`
		VideoUrl     string                           `json:"video_url,omitempty"`
		ThumbnailUrl string                           `json:"thumbnail_url"`
		CtwaCLId     string                           `json:"ctwa_clid"`
		MediaType    AdInteractionSourceMediaTypeEnum `json:"media_type"`
	} `json:"referral,omitempty"`
}

type NotificationPayloadAudioMessageSchemaType struct {
	Audio struct {
		Id       string `json:"id,omitempty"`
		MIMEType string `json:"mime_type,omitempty"`
		SHA256   string `json:"sha256,omitempty"`
	} `json:"audio,omitempty"`
}

type NotificationPayloadImageMessageSchemaType struct {
	Image struct {
		Id       string `json:"id"`
		MIMEType string `json:"mime_type"`
		SHA256   string `json:"sha256"`
		Caption  string `json:"caption,omitempty"`
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
	} `json:"document,omitempty"`
}

type NotificationPayloadOrderMessageSchemaType struct {
	// OrderText string `json:"text"`
	Order struct {
		CatalogId    string `json:"catalog_id"`
		ProductItems []struct {
			ProductRetailerId string `json:"product_retailer_id"`
			Quantity          string `json:"quantity"`
			ItemPrice         string `json:"item_price"`
			Currency          string `json:"currency"`
		} `json:"product_items"`
	} `json:"order,omitempty"`
}

type NotificationPayloadStickerMessageSchemaType struct {
	Sticker struct {
		Id       string `json:"id"`
		MIMEType string `json:"mime_type"`
		SHA256   string `json:"sha256"`
		Animated bool   `json:"animated"`
	} `json:"sticker,omitempty"`
}

type NotificationPayloadSystemMessageSchemaType struct {
	System struct {
		Identity string                     `json:"identity"`
		Body     string                     `json:"body"`
		Customer string                     `json:"customer"`
		Type     SystemNotificationTypeEnum `json:"type"`
		WaId     string                     `json:"wa_id"`
	} `json:"system,omitempty"`
	Identity struct {
		Acknowledged     string `json:"acknowledged"`
		CreatedTimestamp string `json:"created_timestamp"`
		Hash             string `json:"hash"`
	} `json:"identity,omitempty"`
}

type NotificationPayloadVideoMessageSchemaType struct {
	Video struct {
		Id       string `json:"id"`
		MIMEType string `json:"mime_type"`
		SHA256   string `json:"sha256"`
		Caption  string `json:"caption,omitempty"`
		Filename string `json:"filename,omitempty"`
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
		ReplyId string `json:"reply_id"`
		Title   string `json:"title"`
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

type NotificationPayloadContactMessageSchemaType struct {
	Contacts []Contact `json:"contacts"`
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

type Contact struct {
	WaId    string  `json:"wa_id"`
	Profile Profile `json:"profile"`
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
}

type WebhookFieldEnum string

const (
	WebhookFieldEnumAccountAlerts          WebhookFieldEnum = "account_alerts"
	WebhookFieldEnumMessages               WebhookFieldEnum = "messages"
	WebhookFieldEnumSecurity               WebhookFieldEnum = "security"
	WebhookFieldEnumAccountUpdate          WebhookFieldEnum = "account_update"
	WebhookFieldEnumAccountReview          WebhookFieldEnum = "account_review"
	WebhookFieldEnumBusinessCapability     WebhookFieldEnum = "business_capability"
	WebhookFieldEnumMessageTemplateQuality WebhookFieldEnum = "message_template_quality"
	WebhookFieldEnumMessageTemplateStatus  WebhookFieldEnum = "message_template_status"
	WebhookFieldEnumPhoneNumberName        WebhookFieldEnum = "phone_number_name"
	WebhookFieldEnumPhoneNumberQuality     WebhookFieldEnum = "phone_number_quality"
	WebhookFieldEnumTemplateCategoryUpdate WebhookFieldEnum = "template_category"
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
	MessageTemplateId       string                                 `json:"message_template_id"`
	MessageTemplateName     string                                 `json:"message_template_name"`
	MessageTemplateLanguage string                                 `json:"message_template_language"`
	Reason                  TemplateMessageRejectionReasonEnum     `json:"reason"`
	DisableInfo             TemplateMessageStatusUpdateDisableInfo `json:"disable_info,omitempty"`
	OtherInfo               TemplateMessageStatusUpdateOtherInfo   `json:"other_info,omitempty"`
}

type TemplateCategoryUpdateValue struct {
	MessageTemplateId       string                  `json:"message_template_id"`
	MessageTemplateName     string                  `json:"message_template_name"`
	MessageTemplateLanguage string                  `json:"message_template_language"`
	PreviousCategory        MessageTemplateCategory `json:"previous_category"`
	NewCategory             MessageTemplateCategory `json:"new_category"`
	CorrectCategory         MessageTemplateCategory `json:"correct_category"`
}

type TemplateQualityUpdateValue struct {
	PreviousQualityScore    string `json:"previous_quality_score"`
	NewQualityScore         string `json:"new_quality_score"`
	MessageTemplateId       string `json:"message_template_id"`
	MessageTemplateName     string `json:"message_template_name"`
	MessageTemplateLanguage string `json:"message_template_language"`
}

type PhoneNumberNameUpdateValue struct {
	DisplayPhoneNumber    string `json:"display_phone_number"`
	Decision              string `json:"decision"`
	RequestedVerifiedName string `json:"requested_verified_name"`
	RejectionReason       string `json:"rejection_reason"`
}

type PhoneNumberQualityUpdateValue struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	Event              string `json:"event"`
	CurrentLimit       string `json:"current_limit"`
}

type AccountAlertSeverityEnum string

const (
	AccountAlertSeverityEnumCritical AccountAlertSeverityEnum = "CRITICAL"
	AccountAlertSeverityEnumWarning  AccountAlertSeverityEnum = "WARNING"
)

type AccountAlertsValue struct {
	EntityType       string                   `json:"entity_type"`
	EntityId         string                   `json:"entity_id"`
	AlertSeverity    AccountAlertSeverityEnum `json:"alert_severity"`
	AlertStatus      string                   `json:"alert_status"`
	AlertType        string                   `json:"alert_type"`
	AlertDescription string                   `json:"alert_description"`
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

type AccountUpdateValue struct {
	PhoneNumber string                 `json:"phone_number,omitempty"`
	Event       AccountUpdateEventEnum `json:"event"`
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

type Change struct {
	Value interface{}      `json:"value"`
	Field WebhookFieldEnum `json:"field"`
}

type MessagesValue struct {
	MessagingProduct string    `json:"messaging_product"`
	Metadata         Metadata  `json:"metadata"`
	Contacts         []Contact `json:"contacts,omitempty"`
	Statuses         []Status  `json:"statuses,omitempty"`
	Messages         []Message `json:"messages,omitempty"`
	Errors           []Error   `json:"errors,omitempty"`
}

type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberId      string `json:"phone_number_id"`
}

type Status struct {
	Conversation Conversation `json:"conversation,omitempty"`
	Errors       []Error      `json:"errors,omitempty"`
	Status       string       `json:"status"`
	Timestamp    string       `json:"timestamp"`
	RecipientId  string       `json:"recipient_id"`
	Pricing      Pricing      `json:"pricing,omitempty"`
}

type Conversation struct {
	Id     string `json:"id"`
	Origin Origin `json:"origin,omitempty"`
}

type Origin struct {
	Type                MessageStatusCategoryEnum `json:"type"`
	ExpirationTimestamp string                    `json:"expiration_timestamp,omitempty"`
}

type Pricing struct {
	PricingModel string                    `json:"pricing_model"`
	Category     MessageStatusCategoryEnum `json:"category"`
}

type Message struct {
	Id                                              string                                      `json:"id"`
	From                                            string                                      `json:"from"`
	Timestamp                                       string                                      `json:"timestamp"`
	Type                                            NotificationMessageTypeEnum                 `json:"type"`
	Context                                         NotificationPayloadMessageContextSchemaType `json:"context"`
	Errors                                          []Error                                     `json:",inline"`
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
}

type Error struct {
	// Add fields for error details
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
