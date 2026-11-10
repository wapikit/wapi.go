package manager

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/wapikit/wapi.go/internal"
	"github.com/wapikit/wapi.go/internal/request_client"
	"github.com/wapikit/wapi.go/pkg/components"
	"github.com/wapikit/wapi.go/pkg/events"
)

// WebhookManager represents a manager for handling webhooks.
type WebhookManager struct {
	secret       string
	appSecret    string
	path         string
	port         int
	EventManager *EventManager
	Requester    request_client.RequestClient
}

// WebhookManagerConfig represents the configuration options for creating a new WebhookManager.
type WebhookManagerConfig struct {
	Secret       string                       `validate:"required"`
	EventManager *EventManager                `validate:"required"`
	Requester    request_client.RequestClient `validate:"required"`
	Path         string
	Port         int
	// AppSecret is the Meta app secret. When set, POST requests whose
	// X-Hub-Signature-256 header does not match the HMAC-SHA256 of the raw
	// body are rejected with 401. When empty, signatures are not checked.
	AppSecret string
}

// NewWebhook creates a new WebhookManager with the given options.
// It returns nil if Secret or EventManager is missing.
func NewWebhook(options *WebhookManagerConfig) *WebhookManager {
	if err := internal.GetValidator().Struct(options); err != nil {
		return nil
	}
	return &WebhookManager{
		secret:       options.Secret,
		appSecret:    options.AppSecret,
		path:         options.Path,
		port:         options.Port,
		EventManager: options.EventManager,
		Requester:    options.Requester,
	}
}

// ErrInvalidWebhookSignature is returned by VerifyWebhookSignature when the
// X-Hub-Signature-256 header is missing, malformed or does not match.
var ErrInvalidWebhookSignature = errors.New("invalid webhook signature")

// VerifyWebhookSignature checks a Meta X-Hub-Signature-256 header value
// ("sha256=<hex>") against the HMAC-SHA256 of body keyed with appSecret.
func VerifyWebhookSignature(appSecret string, body []byte, signatureHeader string) error {
	hexSignature, ok := strings.CutPrefix(signatureHeader, "sha256=")
	if !ok {
		return ErrInvalidWebhookSignature
	}
	signature, err := hex.DecodeString(hexSignature)
	if err != nil {
		return ErrInvalidWebhookSignature
	}
	mac := hmac.New(sha256.New, []byte(appSecret))
	mac.Write(body)
	if !hmac.Equal(signature, mac.Sum(nil)) {
		return ErrInvalidWebhookSignature
	}
	return nil
}

// logSkippedChange records a webhook change whose value could not be decoded.
// The change is skipped so the rest of the batch is still delivered; failing
// the request would make Meta retry the whole payload and re-deliver events
// that were already published.
func logSkippedChange(field WebhookFieldEnum, err error) {
	slog.Error("wapi: skipping undecodable webhook change", "field", string(field), "error", err)
}

// logSkippedMessage records an inbound message that could not be turned into
// an event. Like logSkippedChange, it skips only that message.
func logSkippedMessage(message Message, err error) {
	slog.Error("wapi: skipping webhook message", "message_id", message.Id, "type", string(message.Type), "error", err)
}

// createEchoHttpServer creates a new instance of Echo HTTP server.
// This function is used in case the client has not provided any custom HTTP server.
func (wh *WebhookManager) createEchoHttpServer() *echo.Echo {
	e := echo.New()
	return e
}

// GetRequestHandler handles GET requests to the webhook endpoint.
func (wh *WebhookManager) GetRequestHandler(c echo.Context) error {
	hubVerificationToken := c.QueryParam("hub.verify_token")
	hubChallenge := c.QueryParam("hub.challenge")
	hubMode := c.QueryParam("hub.mode")
	if hubMode == "subscribe" && hubVerificationToken == wh.secret {
		return c.String(200, hubChallenge)
	} else {
		return c.String(400, "invalid token")
	}
}

// unmarshalWebhookValue is a generic helper to unmarshal webhook change values
func unmarshalWebhookValue[T any](value interface{}) (T, error) {
	var result T
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return result, fmt.Errorf("error marshaling webhook value: %w", err)
	}
	if err := json.Unmarshal(valueBytes, &result); err != nil {
		return result, fmt.Errorf("error unmarshaling webhook value: %w", err)
	}
	return result, nil
}

// PostRequestHandler handles POST requests to the webhook endpoint.
func (wh *WebhookManager) PostRequestHandler(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(400, "error reading request body")
	}

	if wh.appSecret != "" {
		if err := VerifyWebhookSignature(wh.appSecret, body, c.Request().Header.Get("X-Hub-Signature-256")); err != nil {
			return c.String(http.StatusUnauthorized, "invalid signature")
		}
	}

	var payload WhatsappApiNotificationPayloadSchemaType
	if err := json.Unmarshal(body, &payload); err != nil {
		return c.String(400, "Invalid JSON data")
	}

	if err := internal.GetValidator().Struct(payload); err != nil {
		return c.String(400, "Invalid JSON data")
	}

	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			switch change.Field {
			case WebhookFieldEnumMessages:
				messageValue, err := unmarshalWebhookValue[MessagesValue](change.Value)
				if err != nil {
					logSkippedChange(change.Field, err)
					continue
				}

				senderName := ""
				senderWaId := ""
				senderUserId := ""
				senderParentUserId := ""
				senderUsername := ""
				if len(messageValue.Contacts) > 0 {
					contact := messageValue.Contacts[0]
					senderName = contact.Profile.Name
					senderWaId = contact.WaId
					senderUserId = contact.UserId
					senderParentUserId = contact.ParentUserId
					senderUsername = contact.Profile.Username
				}

				wh.handleMessagesSubscriptionEvents(HandleMessageSubscriptionEventPayload{
					Messages: messageValue.Messages,
					Statuses: messageValue.Statuses,
					PhoneNumber: events.BusinessPhoneNumber{
						DisplayNumber: messageValue.Metadata.DisplayPhoneNumber,
						Id:            messageValue.Metadata.PhoneNumberId,
					},
					BusinessAccountId:  entry.Id,
					SenderName:         senderName,
					SenderWaId:         senderWaId,
					SenderUserId:       senderUserId,
					SenderParentUserId: senderParentUserId,
					SenderUsername:     senderUsername,
				})

			case WebhookFieldEnumAccountReview:
				accountReviewValue, err := unmarshalWebhookValue[AccountReviewUpdateValue](change.Value)
				if err != nil {
					logSkippedChange(change.Field, err)
					continue
				}
				wh.handleAccountReviewSubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, accountReviewValue)
			case WebhookFieldEnumAccountAlerts:
				accountAlertValue, err := unmarshalWebhookValue[AccountAlertsValue](change.Value)
				if err != nil {
					logSkippedChange(change.Field, err)
					continue
				}
				wh.handleAccountAlertsSubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, accountAlertValue)
			case WebhookFieldEnumAccountUpdate:
				accountUpdate, err := unmarshalWebhookValue[AccountUpdateValue](change.Value)
				if err != nil {
					logSkippedChange(change.Field, err)
					continue
				}
				wh.handleAccountUpdateSubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, accountUpdate)
			case WebhookFieldEnumTemplateCategoryUpdate:
				templateCategoryUpdate, err := unmarshalWebhookValue[TemplateCategoryUpdateValue](change.Value)
				if err != nil {
					logSkippedChange(change.Field, err)
					continue
				}
				wh.handleTemplateCategoryUpdateSubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, templateCategoryUpdate)
			case WebhookFieldEnumMessageTemplateQuality:
				qualityUpdate, err := unmarshalWebhookValue[TemplateQualityUpdateValue](change.Value)
				if err != nil {
					logSkippedChange(change.Field, err)
					continue
				}
				wh.handleMessageTemplateQualitySubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, qualityUpdate)
			case WebhookFieldEnumMessageTemplateStatus:
				statusUpdate, err := unmarshalWebhookValue[TemplateStatusUpdateValue](change.Value)
				if err != nil {
					logSkippedChange(change.Field, err)
					continue
				}
				wh.handleMessageTemplateStatusSubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, statusUpdate)
			case WebhookFieldEnumPhoneNumberName:
				nameUpdate, err := unmarshalWebhookValue[PhoneNumberNameUpdateValue](change.Value)
				if err != nil {
					logSkippedChange(change.Field, err)
					continue
				}
				wh.handlePhoneNumberNameSubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, nameUpdate)
			case WebhookFieldEnumPhoneNumberQuality:
				qualityUpdate, err := unmarshalWebhookValue[PhoneNumberQualityUpdateValue](change.Value)
				if err != nil {
					logSkippedChange(change.Field, err)
					continue
				}
				wh.handlePhoneNumberQualitySubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, qualityUpdate)
			case WebhookFieldEnumBusinessCapability:
				capabilityUpdate, err := unmarshalWebhookValue[BusinessCapabilityUpdateValue](change.Value)
				if err != nil {
					logSkippedChange(change.Field, err)
					continue
				}
				wh.handleBusinessCapabilitySubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, capabilityUpdate)
			case WebhookFieldEnumSecurity:
				securityChange, err := unmarshalWebhookValue[SecurityValue](change.Value)
				if err != nil {
					logSkippedChange(change.Field, err)
					continue
				}
				wh.handleSecuritySubscriptionEvents(securityChange)
			case WebhookFieldEnumUserPreferences:
				userPrefsValue, err := unmarshalWebhookValue[UserPreferencesValue](change.Value)
				if err != nil {
					logSkippedChange(change.Field, err)
					continue
				}
				wh.handleUserPreferencesSubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, userPrefsValue)
			case WebhookFieldEnumMessageTemplateComponentsUpdate:
				componentsValue, err := unmarshalWebhookValue[MessageTemplateComponentsUpdateValue](change.Value)
				if err != nil {
					logSkippedChange(change.Field, err)
					continue
				}
				wh.handleMessageTemplateComponentsUpdateSubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, componentsValue)
			case WebhookFieldEnumPaymentConfigurationUpdate:
				paymentConfigValue, err := unmarshalWebhookValue[PaymentConfigurationUpdateValue](change.Value)
				if err != nil {
					logSkippedChange(change.Field, err)
					continue
				}
				wh.handlePaymentConfigurationUpdateSubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, paymentConfigValue)
			case WebhookFieldEnumSmbAppStateSync:
				stateSyncValue, err := unmarshalWebhookValue[SmbAppStateSyncValue](change.Value)
				if err != nil {
					logSkippedChange(change.Field, err)
					continue
				}
				wh.handleSmbAppStateSyncSubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, stateSyncValue)
			case WebhookFieldEnumSmbMessageEchoes:
				messageEchoesValue, err := unmarshalWebhookValue[SmbMessageEchoesValue](change.Value)
				if err != nil {
					logSkippedChange(change.Field, err)
					continue
				}
				wh.handleSmbMessageEchoesSubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, messageEchoesValue)
			case WebhookFieldEnumHistory:
				historyValue, err := unmarshalWebhookValue[HistoryValue](change.Value)
				if err != nil {
					logSkippedChange(change.Field, err)
					continue
				}
				wh.handleHistorySubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, historyValue)
			case WebhookFieldEnumUserIdUpdate:
				userIdUpdate, err := unmarshalWebhookValue[UserIdUpdateValue](change.Value)
				if err != nil {
					logSkippedChange(change.Field, err)
					continue
				}
				wh.handleUserIdUpdateSubscriptionEvents(events.BaseSystemEvent{
					Timestamp: fmt.Sprint(entry.Time),
				}, entry.Id, userIdUpdate)
			case WebhookFieldEnumBusinessUsernameUpdates:
				usernameUpdate, err := unmarshalWebhookValue[BusinessUsernameUpdateValue](change.Value)
				if err != nil {
					logSkippedChange(change.Field, err)
					continue
				}
				wh.handleBusinessUsernameUpdateSubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, usernameUpdate)
			}
		}
	}

	return c.String(200, "Message received")
}

// ListenToEvents starts the built-in webhook server on 127.0.0.1:<port> and
// blocks until SIGINT, then shuts the server down.
func (wh *WebhookManager) ListenToEvents() {
	server := wh.createEchoHttpServer()
	server.GET(wh.path, wh.GetRequestHandler)
	server.POST(wh.path, wh.PostRequestHandler)

	addr := fmt.Sprintf("127.0.0.1:%d", wh.port)
	go func() {
		if err := server.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("wapi: webhook server stopped", "addr", addr, "error", err)
		}
	}()
	slog.Info("wapi: listening for webhook events", "addr", addr, "path", wh.path)

	_ = wh.EventManager.Publish(events.ReadyEventType, events.NewReadyEvent())
	// Wait for an interrupt signal (e.g., Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	defer signal.Stop(quit)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("wapi: webhook server shutdown failed", "error", err)
	}
}

type HandleMessageSubscriptionEventPayload struct {
	Messages          []Message                  `json:"messages"`
	Statuses          []Status                   `json:"statuses"`
	PhoneNumber       events.BusinessPhoneNumber `json:"phone_number_id"`     // * this is the phone number to which this event has bee sent to
	BusinessAccountId string                     `json:"business_account_id"` // * business account id to which this event has been sent to
	SenderName        string                     `json:"sender_name"`
	// Sender identity from contacts[0] (BSUID/username rollout). Additive; empty
	// for phone-only contacts.
	SenderWaId         string `json:"sender_wa_id"`
	SenderUserId       string `json:"sender_user_id"`
	SenderParentUserId string `json:"sender_parent_user_id"`
	SenderUsername     string `json:"sender_username"`
}

func (wh *WebhookManager) handleMessagesSubscriptionEvents(payload HandleMessageSubscriptionEventPayload) {
	// consider the field here too, because we will be supporting more events
	if len(payload.Statuses) > 0 {
		for _, status := range payload.Statuses {
			// Meta's pricing block — present on billable status updates. It is
			// surfaced on every status event so consumers can pick whichever
			// status they bill on; deduplicate downstream by message id, since
			// the block may repeat across sent/delivered/read for one message.
			pricing := buildPricingInfo(status.Pricing)

			switch status.Status {
			case string(MessageStatusDelivered):
				{
					ev := events.NewMessageDeliveredEvent(events.BaseSystemEvent{
						Timestamp: status.Timestamp,
					}, status.Id, status.RecipientId)
					ev.Pricing = pricing
					ev.RecipientUserId = status.RecipientUserId
					ev.RecipientParentUserId = status.RecipientParentUserId
					_ = wh.EventManager.Publish(events.MessageDeliveredEventType, ev)
				}

			case string(MessageStatusRead):
				{
					ev := events.NewMessageReadEvent(events.BaseSystemEvent{
						Timestamp: status.Timestamp,
					}, status.Id, status.RecipientId)
					ev.Pricing = pricing
					ev.RecipientUserId = status.RecipientUserId
					ev.RecipientParentUserId = status.RecipientParentUserId
					_ = wh.EventManager.Publish(events.MessageReadEventType, ev)
				}
			case string(MessageStatusSent):
				{
					ev := events.NewMessageSentEvent(events.BaseSystemEvent{
						Timestamp: status.Timestamp,
					}, status.Id, status.RecipientId)
					ev.Pricing = pricing
					ev.RecipientUserId = status.RecipientUserId
					ev.RecipientParentUserId = status.RecipientParentUserId
					_ = wh.EventManager.Publish(events.MessageSentEventType, ev)
				}
			case string(MessageStatusFailed):
				{
					failedReason := ""
					errorCode := 0
					errorMessage := ""
					if len(status.Errors) > 0 {
						for _, err := range status.Errors {
							failedReason = err.Title
							errorCode = err.Code
							errorMessage = err.Message
							break
						}
					}

					_ = wh.EventManager.Publish(events.MessageFailedEventType, events.NewMessageFailedEvent(events.BaseSystemEvent{
						Timestamp: status.Timestamp,
					}, status.Id, status.RecipientId, failedReason, errorCode, errorMessage))
				}
			case string(MessageStatusUnDelivered):
				{
					undeliveredReason := ""
					errorCode := 0
					errorMessage := ""
					if len(status.Errors) > 0 {
						for _, err := range status.Errors {
							undeliveredReason = err.Title
							errorCode = err.Code
							errorMessage = err.Message
							break
						}
					}

					_ = wh.EventManager.Publish(events.MessageUndeliveredEventType, events.NewMessageUndeliveredEvent(events.BaseSystemEvent{
						Timestamp: status.Timestamp,
					}, status.Id, status.RecipientId, undeliveredReason, errorCode, errorMessage))
				}
			}

		}
	}

	for _, message := range payload.Messages {
		var repliedTo string
		if message.Context.Id != "" {
			repliedTo = message.Context.Id
		}

		baseMessageEvent := events.NewBaseMessageEvent(events.BaseMessageEventParams{
			BusinessAccountId: payload.BusinessAccountId,
			MessageId:         message.Id,
			PhoneNumber:       payload.PhoneNumber,
			Timestamp:         message.Timestamp,
			From:              message.From,
			SenderName:        payload.SenderName,
			IsForwarded:       message.Context.Forwarded,
			Context: events.MessageContext{
				RepliedToMessageId: repliedTo,
			},
			Requester: wh.Requester,
			// Identity fields (BSUID/username rollout). Sender-contact-level
			// fields come from contacts[0]; the from_* fields come from the
			// message itself. All additive/optional.
			WaId:             payload.SenderWaId,
			UserId:           payload.SenderUserId,
			ParentUserId:     payload.SenderParentUserId,
			Username:         payload.SenderUsername,
			FromUserId:       message.FromUserId,
			FromParentUserId: message.FromParentUserId,
		})

		if message.Referral != nil {
			mediaUrl := message.Referral.ImageUrl
			if message.Referral.VideoUrl != "" {
				mediaUrl = message.Referral.VideoUrl
			}
			adSource := events.AdSource{
				Url:          message.Referral.SourceUrl,
				Id:           message.Referral.SourceId,
				Type:         events.AdInteractionSourceType(message.Referral.SourceType),
				Title:        message.Referral.Headline,
				Description:  message.Referral.Body,
				MediaUrl:     mediaUrl,
				MediaType:    events.AdInteractionSourceMediaType(message.Referral.MediaType),
				ThumbnailUrl: message.Referral.ThumbnailUrl,
				CtwaClid:     message.Referral.CtwaClid,
			}
			welcomeText := message.Referral.WelcomeMessage.Text
			if welcomeText == "" && message.Type == NotificationMessageTypeText {
				welcomeText = message.Text.Body
			}
			_ = wh.EventManager.Publish(events.AdInteractionEventType, events.NewAdInteractionEvent(
				baseMessageEvent,
				adSource,
				welcomeText,
			))
		}

		switch message.Type {
		case NotificationMessageTypeText:
			{
				_ = wh.EventManager.Publish(events.TextMessageEventType, events.NewTextMessageEvent(
					baseMessageEvent,
					message.Text.Body),
				)
			}
		case NotificationMessageTypeImage:
			{
				imageMessageComponent, err := components.NewImageMessage(components.ImageMessageConfigs{
					Id:      message.Image.Id,
					Link:    message.Image.Url,
					Caption: message.Image.Caption,
				})

				if err != nil {
					logSkippedMessage(message, err)
					continue
				}

				_ = wh.EventManager.Publish(events.ImageMessageEventType, events.NewImageMessageEvent(
					baseMessageEvent,
					*imageMessageComponent,
					message.Image.MIMEType, message.Image.SHA256, message.Image.Id),
				)
			}
		case NotificationMessageTypeAudio:
			{
				audioMessageComponent, err := components.NewAudioMessage(components.AudioMessageConfigs{
					Id:   message.Audio.Id,
					Link: message.Audio.Url,
				})

				if err != nil {
					logSkippedMessage(message, err)
					continue
				}

				_ = wh.EventManager.Publish(events.AudioMessageEventType, events.NewAudioMessageEvent(
					baseMessageEvent,
					*audioMessageComponent,
					message.Audio.MIMEType, message.Audio.SHA256, message.Audio.Id),
				)
			}
		case NotificationMessageTypeVideo:
			{
				videoMessageComponent, err := components.NewVideoMessage(components.VideoMessageConfigs{
					Id:      message.Video.Id,
					Caption: message.Video.Caption,
					Link:    message.Video.Url,
				})

				if err != nil {
					logSkippedMessage(message, err)
					continue
				}

				_ = wh.EventManager.Publish(events.VideoMessageEventType, events.NewVideoMessageEvent(
					baseMessageEvent,
					*videoMessageComponent,
					message.Video.MIMEType, message.Video.SHA256, message.Video.Id),
				)
			}
		case NotificationMessageTypeDocument:
			{
				documentMessageComponent, err := components.NewDocumentMessage(components.DocumentMessageConfigs{
					Id:       message.Document.Id,
					Caption:  &message.Document.Caption,
					FileName: message.Document.Filename,
					Link:     message.Document.Link,
				})

				if err != nil {
					logSkippedMessage(message, err)
					continue
				}

				_ = wh.EventManager.Publish(events.DocumentMessageEventType, events.NewDocumentMessageEvent(
					baseMessageEvent,
					*documentMessageComponent,
					message.Document.Id, message.Document.SHA256, message.Document.MIMEType),
				)
			}
		case NotificationMessageTypeLocation:
			{
				locationMessageComponent, err := components.NewLocationMessage(message.Location.Latitude, message.Location.Longitude)

				if err != nil {
					logSkippedMessage(message, err)
					continue
				}

				_ = wh.EventManager.Publish(events.LocationMessageEventType, events.NewLocationMessageEvent(
					baseMessageEvent,
					*locationMessageComponent),
				)
			}
		case NotificationMessageTypeContacts:
			{
				contactMessageComponent, err := components.NewContactMessage(message.Contacts)
				if err != nil {
					logSkippedMessage(message, err)
					continue
				}
				_ = wh.EventManager.Publish(events.ContactMessageEventType, events.NewContactsMessageEvent(
					baseMessageEvent,
					*contactMessageComponent,
				))
			}
		case NotificationMessageTypeSticker:
			{
				stickerMessageComponent, err := components.NewStickerMessage(&components.StickerMessageConfigs{
					Id:         message.Sticker.Id,
					Link:       message.Sticker.Url,
					IsAnimated: message.Sticker.Animated,
				})

				if err != nil {
					logSkippedMessage(message, err)
					continue
				}

				_ = wh.EventManager.Publish(events.StickerMessageEventType, events.NewStickerMessageEvent(
					baseMessageEvent,
					*stickerMessageComponent,
					message.Sticker.Id, message.Sticker.SHA256, message.Sticker.MIMEType),
				)
			}
		case NotificationMessageTypeButton:
			{
				_ = wh.EventManager.Publish(events.QuickReplyMessageEventType, events.NewQuickReplyButtonInteractionEvent(
					baseMessageEvent,
					message.Button.Text,
					message.Button.Payload,
				))
			}
		case NotificationMessageTypeInteractive:
			{
				if message.Interactive.Type == "list_reply" {
					_ = wh.EventManager.Publish(events.ListInteractionMessageEventType, events.NewListInteractionEvent(
						baseMessageEvent,
						message.Interactive.ListReply.Title,
						message.Interactive.ListReply.Id,
						message.Interactive.ListReply.Description,
					))
				} else {
					_ = wh.EventManager.Publish(events.ReplyButtonInteractionEventType, events.NewReplyButtonInteractionEvent(
						baseMessageEvent,
						message.Interactive.ButtonReply.Title,
						message.Interactive.ButtonReply.Id,
					))
				}

			}
		case NotificationMessageTypeReaction:
			{
				reactionMessageComponent, err := components.NewReactionMessage(components.ReactionMessageParams{
					MessageId: message.Reaction.MessageId,
					Emoji:     message.Reaction.Emoji,
				})

				if err != nil {
					logSkippedMessage(message, err)
					continue
				}

				_ = wh.EventManager.Publish(events.ReactionMessageEventType, events.NewReactionMessageEvent(
					baseMessageEvent,
					*reactionMessageComponent,
				))
			}
		case NotificationMessageTypeOrder:
			{

				productItems := make([]components.ProductItem, len(message.Order.ProductItems))
				for i, item := range message.Order.ProductItems {
					productItems[i] = components.ProductItem{
						Currency:          item.Currency,
						ItemPrice:         item.ItemPrice,
						ProductRetailerID: item.ProductRetailerId,
						Quantity:          item.Quantity,
					}
				}

				_ = wh.EventManager.Publish(events.OrderReceivedEventType, events.NewOrderEvent(
					baseMessageEvent,
					components.Order{
						CatalogID:    message.Order.CatalogId,
						ProductItems: productItems,
						Text:         message.Order.Text,
					},
				))
			}
		case NotificationMessageTypeSystem:
			{
				// According to official WhatsApp docs, system messages only have: body, wa_id, and type
				// The user_changed_number type is the primary system message type
				if message.System.Type == SystemNotificationTypeCustomerPhoneNumberChange {
					_ = wh.EventManager.Publish(events.CustomerNumberChangedEventType, events.CustomerNumberChangedEvent{
						BaseSystemEvent: events.BaseSystemEvent{
							Timestamp: message.Timestamp,
						},
						NewWaId:           message.System.WaId,
						OldWaId:           message.From, // The old number is in the 'from' field
						ChangeDescription: message.System.Body,
					})
				}
				// Note: customer_identity_changed is no longer supported in the current webhook structure
			}
		case NotificationMessageTypeUnknown:
			{
				// ! TODO: handle error in the event and then emit it.
			}
		}
	}

}

func (wh *WebhookManager) handleAccountAlertsSubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value AccountAlertsValue) {
	_ = wh.EventManager.Publish(events.AccountAlertsEventType, events.NewAccountAlertEvent(
		&baseEvent,
		value.EntityType,
		value.EntityId,
		events.AccountAlertSeverityEnum(value.AlertInfo.AlertSeverity),
		events.AccountAlertStatusEnum(value.AlertInfo.AlertStatus),
		value.AlertInfo.AlertType,
		value.AlertInfo.AlertDescription,
	))
}

func (wh *WebhookManager) handleSecuritySubscriptionEvents(value SecurityValue) {
	_ = wh.EventManager.Publish(events.AccountAlertsEventType, events.SecurityEvent{})
}

func (wh *WebhookManager) handleAccountUpdateSubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value AccountUpdateValue) {
	var wabaInfo *events.WabaInfo
	if value.WabaInfo != nil {
		wabaInfo = &events.WabaInfo{
			WabaId:                     value.WabaInfo.WabaId,
			OwnerBusinessId:            value.WabaInfo.OwnerBusinessId,
			AdAccountLinked:            value.WabaInfo.AdAccountLinked,
			PartnerAppId:               value.WabaInfo.PartnerAppId,
			SolutionId:                 value.WabaInfo.SolutionId,
			SolutionPartnerBusinessIds: value.WabaInfo.SolutionPartnerBusinessIds,
		}
	}

	var violationInfo *events.ViolationInfo
	if value.ViolationInfo != nil {
		violationInfo = &events.ViolationInfo{
			ViolationType: value.ViolationInfo.ViolationType,
		}
	}

	var authIntlRateEligibility *events.AuthInternationalRateEligibility
	if value.AuthInternationalRateEligibility != nil {
		exceptionCountries := make([]events.ExceptionCountry, len(value.AuthInternationalRateEligibility.ExceptionCountries))
		for i, ec := range value.AuthInternationalRateEligibility.ExceptionCountries {
			exceptionCountries[i] = events.ExceptionCountry{
				CountryCode: ec.CountryCode,
				StartTime:   ec.StartTime,
			}
		}
		authIntlRateEligibility = &events.AuthInternationalRateEligibility{
			ExceptionCountries: exceptionCountries,
			StartTime:          value.AuthInternationalRateEligibility.StartTime,
		}
	}

	var banInfo *events.BanInfo
	if value.BanInfo != nil {
		banInfo = &events.BanInfo{
			WabaBanState: value.BanInfo.WabaBanState,
			WabaBanDate:  value.BanInfo.WabaBanDate,
		}
	}

	var restrictionInfo *events.RestrictionInfo
	if value.RestrictionInfo != nil {
		restrictionInfo = &events.RestrictionInfo{
			RestrictionType: value.RestrictionInfo.RestrictionType,
			Expiration:      value.RestrictionInfo.Expiration,
		}
	}

	var partnerClientCertInfo *events.PartnerClientCertificationInfo
	if value.PartnerClientCertificationInfo != nil {
		partnerClientCertInfo = &events.PartnerClientCertificationInfo{
			ClientBusinessId: value.PartnerClientCertificationInfo.ClientBusinessId,
			Status:           value.PartnerClientCertificationInfo.Status,
			RejectionReasons: value.PartnerClientCertificationInfo.RejectionReasons,
		}
	}

	_ = wh.EventManager.Publish(events.AccountUpdateEventType, events.NewAccountUpdateEvent(
		&baseEvent,
		events.AccountUpdateEventEnum(value.Event),
		value.PhoneNumber,
		value.Country,
		wabaInfo,
		violationInfo,
		authIntlRateEligibility,
		banInfo,
		restrictionInfo,
		partnerClientCertInfo,
	))
}

func (wh *WebhookManager) handleAccountReviewSubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value AccountReviewUpdateValue) {
	_ = wh.EventManager.Publish(events.AccountAlertsEventType, events.NewAccountReviewUpdateEvent(
		&baseEvent,
		events.AccountReviewUpdateEventEnum(value.Decision),
	))

}

func (wh *WebhookManager) handleBusinessCapabilitySubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value BusinessCapabilityUpdateValue) {
	_ = wh.EventManager.Publish(events.AccountAlertsEventType, events.NewBusinessCapabilityUpdateEvent(
		&baseEvent,
		int64(value.MaxDailyConversationPerPhone),
		int64(value.MaxPhoneNumbersPerBusiness),
	))

}

func (wh *WebhookManager) handleMessageTemplateQualitySubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value TemplateQualityUpdateValue) {
	_ = wh.EventManager.Publish(events.AccountAlertsEventType, events.NewMessageTemplateQualityUpdateEvent(
		&baseEvent,
		events.MessageTemplateQualityUpdateQualityScoreEnum(value.PreviousQualityScore),
		events.MessageTemplateQualityUpdateQualityScoreEnum(value.NewQualityScore),
		value.MessageTemplateId,
		value.MessageTemplateName,
		value.MessageTemplateLanguage,
	))
}

func (wh *WebhookManager) handleMessageTemplateStatusSubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value TemplateStatusUpdateValue) {
	_ = wh.EventManager.Publish(events.AccountAlertsEventType, events.NewMessageTemplateStatusUpdateEvent(
		&baseEvent,
		events.MessageTemplateStatusUpdateEventEnum(value.Event),
		value.MessageTemplateId,
		value.MessageTemplateName,
		value.MessageTemplateLanguage,
		events.MessageTemplateStatusUpdateReason(value.Reason),
	))

}

func (wh *WebhookManager) handlePhoneNumberNameSubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value PhoneNumberNameUpdateValue) {
	_ = wh.EventManager.Publish(events.AccountAlertsEventType, events.NewPhoneNumberNameUpdateEvent(
		&baseEvent,
		value.DisplayPhoneNumber,
		value.RequestedVerifiedName,
		value.Decision,
		&value.RejectionReason,
	))
}

func (wh *WebhookManager) handlePhoneNumberQualitySubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value PhoneNumberQualityUpdateValue) {
	_ = wh.EventManager.Publish(events.AccountAlertsEventType, events.NewPhoneNumberQualityUpdateEvent(
		&baseEvent,
		value.DisplayPhoneNumber,
		events.PhoneNumberUpdateEventEnum(value.Event),
		events.PhoneNumberQualityUpdateCurrentLimitEnum(value.CurrentLimit),
	))
}

func (wh *WebhookManager) handleTemplateCategoryUpdateSubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value TemplateCategoryUpdateValue) {
	_ = wh.EventManager.Publish(events.AccountAlertsEventType, events.NewMessageTemplateCategoryUpdateEvent(
		&baseEvent,
		value.MessageTemplateId,
		value.MessageTemplateName,
		value.MessageTemplateLanguage,
		events.MessageTemplateCategoryEnum(value.PreviousCategory),
		events.MessageTemplateCategoryEnum(value.NewCategory),
	))
}

func (wh *WebhookManager) handleUserPreferencesSubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value UserPreferencesValue) {
	// Convert webhook value to event preferences
	prefs := make([]events.UserPreference, len(value.UserPreferences))
	for i, p := range value.UserPreferences {
		prefs[i] = events.UserPreference{
			WaId:      p.WaId,
			Detail:    p.Detail,
			Category:  p.Category,
			Value:     p.Value,
			Timestamp: p.Timestamp,
		}
	}
	_ = wh.EventManager.Publish(events.UserPreferencesEventType, events.NewUserPreferencesEvent(&baseEvent, prefs))
}

func (wh *WebhookManager) handleMessageTemplateComponentsUpdateSubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value MessageTemplateComponentsUpdateValue) {
	// Convert webhook buttons to event buttons
	buttons := make([]events.MessageTemplateButton, len(value.MessageTemplateButtons))
	for i, b := range value.MessageTemplateButtons {
		buttons[i] = events.MessageTemplateButton{
			Type:        b.MessageTemplateButtonType,
			Text:        b.MessageTemplateButtonText,
			Url:         b.MessageTemplateButtonUrl,
			PhoneNumber: b.MessageTemplateButtonPhoneNumber,
		}
	}
	_ = wh.EventManager.Publish(events.MessageTemplateComponentsUpdateEventType, events.NewMessageTemplateComponentsUpdateEvent(
		&baseEvent,
		value.MessageTemplateId,
		value.MessageTemplateName,
		value.MessageTemplateLanguage,
		value.MessageTemplateElement,
		value.MessageTemplateTitle,
		value.MessageTemplateFooter,
		buttons,
	))
}

func (wh *WebhookManager) handlePaymentConfigurationUpdateSubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value PaymentConfigurationUpdateValue) {
	_ = wh.EventManager.Publish(events.PaymentConfigurationUpdateEventType, events.NewPaymentConfigurationUpdateEvent(
		&baseEvent,
		value.ConfigurationName,
		value.ProviderName,
		value.ProviderMid,
		value.Status,
		value.CreatedTimestamp,
		value.UpdatedTimestamp,
	))
}

func (wh *WebhookManager) handleSmbAppStateSyncSubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value SmbAppStateSyncValue) {
	// Convert webhook state sync to event state sync
	stateSync := make([]events.StateSyncItem, len(value.StateSync))
	for i, s := range value.StateSync {
		stateSync[i] = events.StateSyncItem{
			Type: s.Type,
			Contact: events.ContactSync{
				FullName:    s.Contact.FullName,
				FirstName:   s.Contact.FirstName,
				PhoneNumber: s.Contact.PhoneNumber,
			},
			Action:    s.Action,
			Timestamp: s.Metadata.Timestamp,
		}
	}
	_ = wh.EventManager.Publish(events.SmbAppStateSyncEventType, events.NewSmbAppStateSyncEvent(
		&baseEvent,
		value.MessagingProduct,
		value.Metadata.DisplayPhoneNumber,
		value.Metadata.PhoneNumberId,
		stateSync,
	))
}

func (wh *WebhookManager) handleSmbMessageEchoesSubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value SmbMessageEchoesValue) {
	// Convert webhook message echoes to event message echoes
	echoes := make([]events.MessageEcho, len(value.MessageEchoes))
	for i, e := range value.MessageEchoes {
		echoes[i] = events.MessageEcho{
			From:      e.From,
			To:        e.To,
			Id:        e.Id,
			Timestamp: e.Timestamp,
			Type:      e.Type,
		}
	}
	_ = wh.EventManager.Publish(events.SmbMessageEchoesEventType, events.NewSmbMessageEchoesEvent(
		&baseEvent,
		value.MessagingProduct,
		value.Metadata.DisplayPhoneNumber,
		value.Metadata.PhoneNumberId,
		echoes,
	))
}

func (wh *WebhookManager) handleHistorySubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value HistoryValue) {
	// Convert webhook history to event history chunks
	chunks := make([]events.HistoryChunk, len(value.History))
	for i, h := range value.History {
		threads := make([]events.HistoryThread, len(h.Threads))
		for j, t := range h.Threads {
			messages := make([]events.HistoryMessage, len(t.Messages))
			for k, m := range t.Messages {
				messages[k] = events.HistoryMessage{
					From:          m.From,
					To:            m.To,
					Id:            m.Id,
					Timestamp:     m.Timestamp,
					Type:          m.Type,
					MessageStatus: m.HistoryContext.Status,
				}
			}
			threads[j] = events.HistoryThread{
				Id:       t.Id,
				Messages: messages,
			}
		}
		chunks[i] = events.HistoryChunk{
			Phase:      h.Metadata.Phase,
			ChunkOrder: h.Metadata.ChunkOrder,
			Progress:   h.Metadata.Progress,
			Threads:    threads,
		}
	}
	_ = wh.EventManager.Publish(events.HistoryEventType, events.NewHistoryEvent(
		&baseEvent,
		value.MessagingProduct,
		value.Metadata.DisplayPhoneNumber,
		value.Metadata.PhoneNumberId,
		chunks,
	))
}

func (wh *WebhookManager) handleUserIdUpdateSubscriptionEvents(baseEvent events.BaseSystemEvent, businessAccountId string, value UserIdUpdateValue) {
	// Prefer explicit old/new; fall back to user_id when Meta only sends the
	// current id. Raw values preserved for downstream identity bridging.
	newUserId := value.NewUserId
	if newUserId == "" {
		newUserId = value.UserId
	}
	_ = wh.EventManager.Publish(events.UserIdUpdateEventType, events.NewUserIdUpdateEvent(
		baseEvent,
		businessAccountId,
		value.WaId,
		value.OldUserId,
		newUserId,
		value.ParentUserId,
	))
}

func (wh *WebhookManager) handleBusinessUsernameUpdateSubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value BusinessUsernameUpdateValue) {
	_ = wh.EventManager.Publish(events.BusinessUsernameUpdateEventType, events.NewBusinessUsernameUpdateEvent(
		&baseEvent,
		value.Metadata.PhoneNumberId,
		value.Username,
		value.PreviousUsername,
		value.Event,
	))
}

// buildPricingInfo converts the internal Pricing struct (parsed from
// Meta's webhook payload) into the public events.MessagePricingInfo
// surface. Returns nil for the empty case so downstream handlers can
// distinguish 'pricing block absent' from 'pricing block present with
// Billable=false'.
func buildPricingInfo(p Pricing) *events.MessagePricingInfo {
	// Treat truly empty Pricing (no fields populated) as nil. Meta omits
	// the entire pricing block on non-billable internal status pings; we
	// don't want every downstream to have to null-check three empty
	// strings.
	if !p.Billable && p.PricingModel == "" && string(p.Category) == "" {
		return nil
	}
	return &events.MessagePricingInfo{
		Billable:     p.Billable,
		PricingModel: p.PricingModel,
		Category:     string(p.Category),
	}
}
