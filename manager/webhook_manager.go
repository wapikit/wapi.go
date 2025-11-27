package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
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
	path         string
	port         int
	EventManager EventManager
	Requester    request_client.RequestClient
}

// WebhookManagerConfig represents the configuration options for creating a new WebhookManager.
type WebhookManagerConfig struct {
	Secret       string                       `validate:"required"`
	EventManager EventManager                 `validate:"required"`
	Requester    request_client.RequestClient `validate:"required"`
	Path         string
	Port         int
}

// NewWebhook creates a new WebhookManager with the given options.
func NewWebhook(options *WebhookManagerConfig) *WebhookManager {
	if err := internal.GetValidator().Struct(options); err != nil {
		return nil
	}
	return &WebhookManager{
		secret:       options.Secret,
		path:         options.Path,
		port:         options.Port,
		EventManager: options.EventManager,
		Requester:    options.Requester,
	}
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
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid messages webhook: %v", err))
				}

				senderName := ""
				if len(messageValue.Contacts) > 0 {
					senderName = messageValue.Contacts[0].Profile.Name
				}

				err = wh.handleMessagesSubscriptionEvents(HandleMessageSubscriptionEventPayload{
					Messages: messageValue.Messages,
					Statuses: messageValue.Statuses,
					PhoneNumber: events.BusinessPhoneNumber{
						DisplayNumber: messageValue.Metadata.DisplayPhoneNumber,
						Id:            messageValue.Metadata.PhoneNumberId,
					},
					BusinessAccountId: entry.Id,
					SenderName:        senderName,
				})

				if err != nil {
					fmt.Println("Error handling messages subscription events:", err)
					c.String(500, "Internal server error")
					return err
				}
			case WebhookFieldEnumAccountReview:
				accountReviewValue, err := unmarshalWebhookValue[AccountReviewUpdateValue](change.Value)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid account_review webhook: %v", err))
				}
				err = wh.handleAccountReviewSubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, accountReviewValue)
				if err != nil {
					fmt.Println("Error handling account_review webhook:", err)
					c.String(500, "Internal server error")
					return err
				}
			case WebhookFieldEnumAccountAlerts:
				accountAlertValue, err := unmarshalWebhookValue[AccountAlertsValue](change.Value)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid account_alerts webhook: %v", err))
				}
				err = wh.handleAccountAlertsSubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, accountAlertValue)
				if err != nil {
					fmt.Println("Error handling account_alerts webhook:", err)
					c.String(500, "Internal server error")
					return err
				}
			case WebhookFieldEnumAccountUpdate:
				accountUpdate, err := unmarshalWebhookValue[AccountUpdateValue](change.Value)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid account_update webhook: %v", err))
				}
				wh.handleAccountUpdateSubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, accountUpdate)
			case WebhookFieldEnumTemplateCategoryUpdate:
				templateCategoryUpdate, err := unmarshalWebhookValue[TemplateCategoryUpdateValue](change.Value)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid template_category webhook: %v", err))
				}
				wh.handleTemplateCategoryUpdateSubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, templateCategoryUpdate)
				if err != nil {
					fmt.Println("Error handling template_category webhook:", err)
					c.String(500, "Internal server error")
					return err
				}
			case WebhookFieldEnumMessageTemplateQuality:
				qualityUpdate, err := unmarshalWebhookValue[TemplateQualityUpdateValue](change.Value)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid message_template_quality webhook: %v", err))
				}
				wh.handleMessageTemplateQualitySubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, qualityUpdate)
				if err != nil {
					fmt.Println("Error handling message_template_quality webhook:", err)
					c.String(500, "Internal server error")
					return err
				}
			case WebhookFieldEnumMessageTemplateStatus:
				statusUpdate, err := unmarshalWebhookValue[TemplateStatusUpdateValue](change.Value)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid message_template_status webhook: %v", err))
				}
				wh.handleMessageTemplateStatusSubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, statusUpdate)
				if err != nil {
					fmt.Println("Error handling message_template_status webhook:", err)
					c.String(500, "Internal server error")
					return err
				}
			case WebhookFieldEnumPhoneNumberName:
				nameUpdate, err := unmarshalWebhookValue[PhoneNumberNameUpdateValue](change.Value)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid phone_number_name webhook: %v", err))
				}
				wh.handlePhoneNumberNameSubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, nameUpdate)
				if err != nil {
					fmt.Println("Error handling phone_number_name webhook:", err)
					c.String(500, "Internal server error")
					return err
				}
			case WebhookFieldEnumPhoneNumberQuality:
				qualityUpdate, err := unmarshalWebhookValue[PhoneNumberQualityUpdateValue](change.Value)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid phone_number_quality webhook: %v", err))
				}
				wh.handlePhoneNumberQualitySubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, qualityUpdate)
				if err != nil {
					fmt.Println("Error handling phone_number_quality webhook:", err)
					c.String(500, "Internal server error")
					return err
				}
			case WebhookFieldEnumBusinessCapability:
				capabilityUpdate, err := unmarshalWebhookValue[BusinessCapabilityUpdateValue](change.Value)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid business_capability webhook: %v", err))
				}
				wh.handleBusinessCapabilitySubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, capabilityUpdate)
				if err != nil {
					fmt.Println("Error handling business_capability webhook:", err)
					c.String(500, "Internal server error")
					return err
				}
			case WebhookFieldEnumSecurity:
				securityChange, err := unmarshalWebhookValue[SecurityValue](change.Value)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid security webhook: %v", err))
				}
				wh.handleSecuritySubscriptionEvents(securityChange)
			case WebhookFieldEnumUserPreferences:
				userPrefsValue, err := unmarshalWebhookValue[UserPreferencesValue](change.Value)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid user_preferences webhook: %v", err))
				}
				wh.handleUserPreferencesSubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, userPrefsValue)
			case WebhookFieldEnumMessageTemplateComponentsUpdate:
				componentsValue, err := unmarshalWebhookValue[MessageTemplateComponentsUpdateValue](change.Value)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid message_template_components_update webhook: %v", err))
				}
				wh.handleMessageTemplateComponentsUpdateSubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, componentsValue)
			case WebhookFieldEnumPaymentConfigurationUpdate:
				paymentConfigValue, err := unmarshalWebhookValue[PaymentConfigurationUpdateValue](change.Value)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid payment_configuration_update webhook: %v", err))
				}
				wh.handlePaymentConfigurationUpdateSubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, paymentConfigValue)
			case WebhookFieldEnumSmbAppStateSync:
				stateSyncValue, err := unmarshalWebhookValue[SmbAppStateSyncValue](change.Value)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid smb_app_state_sync webhook: %v", err))
				}
				wh.handleSmbAppStateSyncSubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, stateSyncValue)
			case WebhookFieldEnumSmbMessageEchoes:
				messageEchoesValue, err := unmarshalWebhookValue[SmbMessageEchoesValue](change.Value)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid smb_message_echoes webhook: %v", err))
				}
				wh.handleSmbMessageEchoesSubscriptionEvents(events.BaseBusinessAccountEvent{
					BusinessAccountId: entry.Id,
					Timestamp:         fmt.Sprint(entry.Time),
				}, messageEchoesValue)
			}
		}
	}

	c.String(200, "Message received")
	return nil
}

// ListenToEvents starts listening to events and handles incoming requests.
func (wh *WebhookManager) ListenToEvents() {
	fmt.Println("Listening to events")
	server := wh.createEchoHttpServer()
	server.GET(wh.path, wh.GetRequestHandler)
	server.POST(wh.path, wh.PostRequestHandler)

	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf("127.0.0.1:%d", wh.port)
		if err := server.Start(addr); err != nil {
			return
		}
	}()

	wh.EventManager.Publish(events.ReadyEventType, events.NewReadyEvent())
	// Wait for an interrupt signal (e.g., Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt) // Capture SIGINT (Ctrl+C)
	<-quit                            // Wait for the signal

	// Gracefully shut down the server (optional)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err) // Handle shutdown errors gracefully
	}
}

type HandleMessageSubscriptionEventPayload struct {
	Messages          []Message                  `json:"messages"`
	Statuses          []Status                   `json:"statuses"`
	PhoneNumber       events.BusinessPhoneNumber `json:"phone_number_id"`     // * this is the phone number to which this event has bee sent to
	BusinessAccountId string                     `json:"business_account_id"` // * business account id to which this event has been sent to
	SenderName        string                     `json:"sender_name"`
}

func (wh *WebhookManager) handleMessagesSubscriptionEvents(payload HandleMessageSubscriptionEventPayload) error {
	// consider the field here too, because we will be supporting more events
	if len(payload.Statuses) > 0 {
		for _, status := range payload.Statuses {
			switch status.Status {
			case string(MessageStatusDelivered):
				{
					wh.EventManager.Publish(events.MessageDeliveredEventType, events.NewMessageDeliveredEvent(events.BaseSystemEvent{
						Timestamp: status.Timestamp,
					}, status.Id, status.RecipientId))
				}

			case string(MessageStatusRead):
				{
					wh.EventManager.Publish(events.MessageReadEventType, events.NewMessageReadEvent(events.BaseSystemEvent{
						Timestamp: status.Timestamp,
					}, status.Id, status.RecipientId))
				}
			case string(MessageStatusSent):
				{
					wh.EventManager.Publish(events.MessageSentEventType, events.NewMessageSentEvent(events.BaseSystemEvent{
						Timestamp: status.Timestamp,
					}, status.Id, status.RecipientId))
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

					wh.EventManager.Publish(events.MessageFailedEventType, events.NewMessageFailedEvent(events.BaseSystemEvent{
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

					wh.EventManager.Publish(events.MessageUndeliveredEventType, events.NewMessageUndeliveredEvent(events.BaseSystemEvent{
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
		})

		switch message.Type {
		case NotificationMessageTypeText:
			{
				wh.EventManager.Publish(events.TextMessageEventType, events.NewTextMessageEvent(
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
					// ! TODO: emit error event here
					fmt.Println("Error creating image message:", err)
					return err
				}

				wh.EventManager.Publish(events.ImageMessageEventType, events.NewImageMessageEvent(
					baseMessageEvent,
					*imageMessageComponent,
					message.Image.MIMEType, message.Image.SHA256, message.Image.Id),
				)
			}
		case NotificationMessageTypeAudio:
			{

				audioMessageComponent, err := components.NewAudioMessage(components.AudioMessageConfigs{
					Id: message.Audio.Id,
				})

				if err != nil {
					// ! TODO: emit error event here
					fmt.Println("Error creating audio message:", err)
					return err
				}

				wh.EventManager.Publish(events.AudioMessageEventType, events.NewAudioMessageEvent(
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
				})

				if err != nil {
					// ! TODO: emit error event here
					fmt.Println("Error creating Video message:", err)
					return err
				}

				wh.EventManager.Publish(events.VideoMessageEventType, events.NewVideoMessageEvent(
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
				})

				if err != nil {
					// ! TODO: emit error event here
					fmt.Println("Error creating document message:", err)
					return err
				}

				wh.EventManager.Publish(events.DocumentMessageEventType, events.NewDocumentMessageEvent(
					baseMessageEvent,
					*documentMessageComponent,
					message.Document.MIMEType, message.Document.SHA256, message.Document.Id),
				)
			}
		case NotificationMessageTypeLocation:
			{
				locationMessageComponent, err := components.NewLocationMessage(message.Location.Latitude, message.Location.Longitude)

				if err != nil {
					// ! TODO: emit error event here
					fmt.Println("Error creating location message:", err)
					return err
				}

				wh.EventManager.Publish(events.LocationMessageEventType, events.NewLocationMessageEvent(
					baseMessageEvent,
					*locationMessageComponent),
				)
			}
		case NotificationMessageTypeContacts:
			{
				contactMessageComponent, _ := components.NewContactMessage([]components.Contact{})
				// ! TODO: add the contact here to the contact message component
				wh.EventManager.Publish(events.ContactMessageEventType, events.NewContactsMessageEvent(
					baseMessageEvent,
					*contactMessageComponent,
				))
			}
		case NotificationMessageTypeSticker:
			{

				stickerMessageComponent, err := components.NewStickerMessage(&components.StickerMessageConfigs{
					Id: message.Sticker.Id,
				})

				if err != nil {
					// ! TODO: emit error event here
					fmt.Println("Error creating Sticker message:", err)
					return err
				}

				wh.EventManager.Publish(events.StickerMessageEventType, events.NewStickerMessageEvent(
					baseMessageEvent,
					*stickerMessageComponent,
					message.Sticker.MIMEType, message.Sticker.SHA256, message.Sticker.Id),
				)

			}
		case NotificationMessageTypeButton:
			{
				wh.EventManager.Publish(events.QuickReplyMessageEventType, events.NewQuickReplyButtonInteractionEvent(
					baseMessageEvent,
					message.Button.Text,
					message.Button.Payload,
				))
			}
		case NotificationMessageTypeInteractive:
			{
				if message.Interactive.Type == "list_reply" {
					wh.EventManager.Publish(events.ListInteractionMessageEventType, events.NewListInteractionEvent(
						baseMessageEvent,
						message.Interactive.ListReply.Title,
						message.Interactive.ListReply.Id,
						message.Interactive.ListReply.Description,
					))
				} else {
					wh.EventManager.Publish(events.ReplyButtonInteractionEventType, events.NewReplyButtonInteractionEvent(
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
					// ! TODO: emit error event here
					fmt.Println("Error creating location message:", err)
					return err
				}

				wh.EventManager.Publish(events.ReactionMessageEventType, events.NewReactionMessageEvent(
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

				wh.EventManager.Publish(events.OrderReceivedEventType, events.NewOrderEvent(
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
					wh.EventManager.Publish(events.CustomerNumberChangedEventType, events.CustomerNumberChangedEvent{
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

	return nil
}

func (wh *WebhookManager) handleAccountAlertsSubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value AccountAlertsValue) error {
	wh.EventManager.Publish(events.AccountAlertsEventType, events.NewAccountAlertEvent(
		&baseEvent,
		value.EntityType,
		value.EntityId,
		events.AccountAlertSeverityEnum(value.AlertInfo.AlertSeverity),
		events.AccountAlertStatusEnum(value.AlertInfo.AlertStatus),
		value.AlertInfo.AlertType,
		value.AlertInfo.AlertDescription,
	))
	return nil
}

func (wh *WebhookManager) handleSecuritySubscriptionEvents(value SecurityValue) {
	wh.EventManager.Publish(events.AccountAlertsEventType, events.SecurityEvent{})
}

func (wh *WebhookManager) handleAccountUpdateSubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value AccountUpdateValue) {
	wh.EventManager.Publish(events.AccountAlertsEventType, events.NewAccountUpdateEvent(
		&baseEvent,
		events.AccountUpdateEventEnum(value.Event),
		value.PhoneNumber,
	))

}

func (wh *WebhookManager) handleAccountReviewSubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value AccountReviewUpdateValue) error {
	wh.EventManager.Publish(events.AccountAlertsEventType, events.NewAccountReviewUpdateEvent(
		&baseEvent,
		events.AccountReviewUpdateEventEnum(value.Decision),
	))
	return nil

}

func (wh *WebhookManager) handleBusinessCapabilitySubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value BusinessCapabilityUpdateValue) error {
	wh.EventManager.Publish(events.AccountAlertsEventType, events.NewBusinessCapabilityUpdateEvent(
		&baseEvent,
		int64(value.MaxDailyConversationPerPhone),
		int64(value.MaxPhoneNumbersPerBusiness),
	))
	return nil

}

func (wh *WebhookManager) handleMessageTemplateQualitySubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value TemplateQualityUpdateValue) error {
	wh.EventManager.Publish(events.AccountAlertsEventType, events.NewMessageTemplateQualityUpdateEvent(
		&baseEvent,
		events.MessageTemplateQualityUpdateQualityScoreEnum(value.PreviousQualityScore),
		events.MessageTemplateQualityUpdateQualityScoreEnum(value.NewQualityScore),
		value.MessageTemplateId,
		value.MessageTemplateName,
		value.MessageTemplateLanguage,
	))

	return nil

}

func (wh *WebhookManager) handleMessageTemplateStatusSubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value TemplateStatusUpdateValue) error {
	wh.EventManager.Publish(events.AccountAlertsEventType, events.NewMessageTemplateStatusUpdateEvent(
		&baseEvent,
		events.MessageTemplateStatusUpdateEventEnum(value.Event),
		value.MessageTemplateId,
		value.MessageTemplateName,
		value.MessageTemplateLanguage,
		events.MessageTemplateStatusUpdateReason(value.Reason),
	))
	return nil

}

func (wh *WebhookManager) handlePhoneNumberNameSubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value PhoneNumberNameUpdateValue) error {
	wh.EventManager.Publish(events.AccountAlertsEventType, events.NewPhoneNumberNameUpdateEvent(
		&baseEvent,
		value.DisplayPhoneNumber,
		value.RequestedVerifiedName,
		value.Decision,
		&value.RejectionReason,
	))
	return nil
}

func (wh *WebhookManager) handlePhoneNumberQualitySubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value PhoneNumberQualityUpdateValue) error {
	wh.EventManager.Publish(events.AccountAlertsEventType, events.NewPhoneNumberQualityUpdateEvent(
		&baseEvent,
		value.DisplayPhoneNumber,
		events.PhoneNumberUpdateEventEnum(value.Event),
		events.PhoneNumberQualityUpdateCurrentLimitEnum(value.CurrentLimit),
	))
	return nil
}

func (wh *WebhookManager) handleTemplateCategoryUpdateSubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value TemplateCategoryUpdateValue) error {
	wh.EventManager.Publish(events.AccountAlertsEventType, events.NewMessageTemplateCategoryUpdateEvent(
		&baseEvent,
		value.MessageTemplateId,
		value.MessageTemplateName,
		value.MessageTemplateLanguage,
		events.MessageTemplateCategoryEnum(value.PreviousCategory),
		events.MessageTemplateCategoryEnum(value.NewCategory),
	))
	return nil
}

func (wh *WebhookManager) handleUserPreferencesSubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value UserPreferencesValue) {
	// TODO: Create proper event type for user preferences in events package
	// For now, we'll publish a generic event with the user preferences data
	// SDK users can subscribe to this event type when it's added to the events package
	_ = baseEvent
	_ = value
	// wh.EventManager.Publish(events.UserPreferencesEventType, events.NewUserPreferencesEvent(&baseEvent, value))
}

func (wh *WebhookManager) handleMessageTemplateComponentsUpdateSubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value MessageTemplateComponentsUpdateValue) {
	// TODO: Create proper event type for message template components update in events package
	// For now, we'll publish a generic event with the components update data
	// SDK users can subscribe to this event type when it's added to the events package
	_ = baseEvent
	_ = value
	// wh.EventManager.Publish(events.MessageTemplateComponentsUpdateEventType, events.NewMessageTemplateComponentsUpdateEvent(&baseEvent, value))
}

func (wh *WebhookManager) handlePaymentConfigurationUpdateSubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value PaymentConfigurationUpdateValue) {
	// TODO: Create proper event type for payment configuration update in events package
	// For now, we'll publish a generic event with the payment configuration data
	// SDK users can subscribe to this event type when it's added to the events package
	_ = baseEvent
	_ = value
	// wh.EventManager.Publish(events.PaymentConfigurationUpdateEventType, events.NewPaymentConfigurationUpdateEvent(&baseEvent, value))
}

func (wh *WebhookManager) handleSmbAppStateSyncSubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value SmbAppStateSyncValue) {
	// TODO: Create proper event type for SMB app state sync in events package
	// For now, we'll publish a generic event with the state sync data
	// SDK users can subscribe to this event type when it's added to the events package
	_ = baseEvent
	_ = value
	// wh.EventManager.Publish(events.SmbAppStateSyncEventType, events.NewSmbAppStateSyncEvent(&baseEvent, value))
}

func (wh *WebhookManager) handleSmbMessageEchoesSubscriptionEvents(baseEvent events.BaseBusinessAccountEvent, value SmbMessageEchoesValue) {
	// TODO: Create proper event type for SMB message echoes in events package
	// For now, we'll publish a generic event with the message echoes data
	// SDK users can subscribe to this event type when it's added to the events package
	_ = baseEvent
	_ = value
	// wh.EventManager.Publish(events.SmbMessageEchoesEventType, events.NewSmbMessageEchoesEvent(&baseEvent, value))
}
