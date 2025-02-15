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
		fmt.Println("Error validating WebhookManagerConfig:", err)
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
	fmt.Println(hubVerificationToken, hubChallenge)
	if hubMode == "subscribe" && hubVerificationToken == wh.secret {
		return c.String(200, hubChallenge)
	} else {
		return c.String(400, "invalid token")
	}
}

// PostRequestHandler handles POST requests to the webhook endpoint.
func (wh *WebhookManager) PostRequestHandler(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		c.String(400, "error reading request body")
	}

	var payload WhatsappApiNotificationPayloadSchemaType
	if err := json.Unmarshal(body, &payload); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		c.String(400, "Invalid JSON data")
	}

	if err := internal.GetValidator().Struct(payload); err != nil {
		fmt.Println("Error validating JSON:", err)
		c.String(400, "Invalid JSON data")
	}

	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			switch change.Field {
			case WebhookFieldEnumMessages:
				var messageValue MessagesValue
				valueBytes, err := json.Marshal(change.Value)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Error marshaling messages value")
				}
				if err := json.Unmarshal(valueBytes, &messageValue); err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid MessagesValue JSON: %v", err))
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
				var accountReviewValue AccountReviewUpdateValue
				valueBytes, err := json.Marshal(change.Value)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Error marshaling messages value")
				}
				if err := json.Unmarshal(valueBytes, &accountReviewValue); err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid accountReviewValue JSON: %v", err))
				}
				wh.handleAccountReviewSubscriptionEvents(accountReviewValue)
				if err != nil {
					fmt.Println("Error handling messages subscription events:", err)
					c.String(500, "Internal server error")
					return err
				}
			case WebhookFieldEnumAccountAlerts:
				var accountAlertValue AccountAlertsValue
				valueBytes, err := json.Marshal(change.Value)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Error marshaling messages value")
				}
				if err := json.Unmarshal(valueBytes, &accountAlertValue); err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid accountAlertValue JSON: %v", err))
				}
				wh.handleAccountAlertsSubscriptionEvents(accountAlertValue)
				if err != nil {
					fmt.Println("Error handling messages subscription events:", err)
					c.String(500, "Internal server error")
					return err
				}
			case WebhookFieldEnumAccountUpdate:
				var accountUpdate AccountUpdateValue
				valueBytes, err := json.Marshal(change.Value)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Error marshaling messages value")
				}
				if err := json.Unmarshal(valueBytes, &accountUpdate); err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid accountUpdate JSON: %v", err))
				}
				wh.handleAccountUpdateSubscriptionEvents(accountUpdate)
				if err != nil {
					fmt.Println("Error handling messages subscription events:", err)
					c.String(500, "Internal server error")
					return err
				}
			case WebhookFieldEnumTemplateCategoryUpdate:
				var templateCategoryUpdate TemplateCategoryUpdateValue
				valueBytes, err := json.Marshal(change.Value)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Error marshaling messages value")
				}
				if err := json.Unmarshal(valueBytes, &templateCategoryUpdate); err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid templateCategoryUpdate JSON: %v", err))
				}
				wh.handleTemplateCategoryUpdateSubscriptionEvents(templateCategoryUpdate)
				if err != nil {
					fmt.Println("Error handling messages subscription events:", err)
					c.String(500, "Internal server error")
					return err
				}
			case WebhookFieldEnumMessageTemplateQuality:
				var qualityUpdate TemplateQualityUpdateValue
				valueBytes, err := json.Marshal(change.Value)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Error marshaling messages value")
				}
				if err := json.Unmarshal(valueBytes, &qualityUpdate); err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid qualityUpdate JSON: %v", err))
				}
				wh.handleMessageTemplateQualitySubscriptionEvents(qualityUpdate)
				if err != nil {
					fmt.Println("Error handling messages subscription events:", err)
					c.String(500, "Internal server error")
					return err
				}
			case WebhookFieldEnumMessageTemplateStatus:
				var statusUpdate TemplateStatusUpdateValue
				valueBytes, err := json.Marshal(change.Value)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Error marshaling messages value")
				}
				if err := json.Unmarshal(valueBytes, &statusUpdate); err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid statusUpdate JSON: %v", err))
				}
				wh.handleMessageTemplateStatusSubscriptionEvents(statusUpdate)
				if err != nil {
					fmt.Println("Error handling messages subscription events:", err)
					c.String(500, "Internal server error")
					return err
				}
			case WebhookFieldEnumPhoneNumberName:
				var nameUpdate PhoneNumberNameUpdateValue
				valueBytes, err := json.Marshal(change.Value)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Error marshaling messages value")
				}
				if err := json.Unmarshal(valueBytes, &nameUpdate); err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid nameUpdate JSON: %v", err))
				}
				wh.handlePhoneNumberNameSubscriptionEvents(nameUpdate)
				if err != nil {
					fmt.Println("Error handling messages subscription events:", err)
					c.String(500, "Internal server error")
					return err
				}
			case WebhookFieldEnumPhoneNumberQuality:
				var qualityUpdate PhoneNumberQualityUpdateValue
				valueBytes, err := json.Marshal(change.Value)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Error marshaling messages value")
				}
				if err := json.Unmarshal(valueBytes, &qualityUpdate); err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid qualityUpdate JSON: %v", err))
				}
				wh.handlePhoneNumberQualitySubscriptionEvents(qualityUpdate)
				if err != nil {
					fmt.Println("Error handling messages subscription events:", err)
					c.String(500, "Internal server error")
					return err
				}
			case WebhookFieldEnumBusinessCapability:
				var capabilityUpdate BusinessCapabilityUpdateValue
				valueBytes, err := json.Marshal(change.Value)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Error marshaling messages value")
				}
				if err := json.Unmarshal(valueBytes, &capabilityUpdate); err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid capabilityUpdate JSON: %v", err))
				}
				wh.handleBusinessCapabilitySubscriptionEvents(capabilityUpdate)
				if err != nil {
					fmt.Println("Error handling messages subscription events:", err)
					c.String(500, "Internal server error")
					return err
				}
			case WebhookFieldEnumSecurity:
				var securityChange SecurityValue
				valueBytes, err := json.Marshal(change.Value)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Error marshaling messages value")
				}
				if err := json.Unmarshal(valueBytes, &securityChange); err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid securityChange JSON: %v", err))
				}
				wh.handleSecuritySubscriptionEvents(securityChange)
				if err != nil {
					fmt.Println("Error handling messages subscription events:", err)
					c.String(500, "Internal server error")
					return err
				}
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
		if err := server.Start("127.0.0.1:8080"); err != nil {
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
					}, status.Conversation.Id, status.RecipientId))
				}

			case string(MessageStatusRead):
				{
					wh.EventManager.Publish(events.MessageReadEventType, events.NewMessageReadEvent(events.BaseSystemEvent{
						Timestamp: status.Timestamp,
					}, status.Conversation.Id, status.RecipientId))
				}
			case string(MessageStatusSent):
				{
					wh.EventManager.Publish(events.MessageSentEventType, events.NewMessageSentEvent(events.BaseSystemEvent{
						Timestamp: status.Timestamp,
					}, status.Conversation.Id, status.RecipientId))
				}
			case string(MessageStatusFailed):
				{
					// ! TODO: check and properly emit the error event here.
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
				documentMessageComponent, err := components.NewVideoMessage(components.VideoMessageConfigs{
					Id:      message.Document.Id,
					Caption: message.Document.Caption,
				})

				if err != nil {
					// ! TODO: emit error event here
					fmt.Println("Error creating document message:", err)
					return err
				}

				wh.EventManager.Publish(events.DocumentMessageEventType, events.NewVideoMessageEvent(
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
				if message.Interactive.Type == "list" {
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
						message.Interactive.ButtonReply.ReplyId,
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
				wh.EventManager.Publish(events.OrderReceivedEventType, events.NewTextMessageEvent(
					baseMessageEvent,
					message.Text.Body),
				)
			}
		case NotificationMessageTypeSystem:
			{
				if message.System.Type == SystemNotificationTypeCustomerIdentityChanged {
					wh.EventManager.Publish(events.CustomerIdentityChangedEventType, events.CustomerIdentityChangedEvent{
						BaseSystemEvent: events.BaseSystemEvent{
							Timestamp: message.Timestamp,
						},
						Acknowledged:      message.Identity.Acknowledged,
						CreationTimestamp: message.Identity.CreatedTimestamp,
						Hash:              message.Identity.Hash,
					})
				} else {
					wh.EventManager.Publish(events.CustomerNumberChangedEventType, events.CustomerNumberChangedEvent{
						BaseSystemEvent: events.BaseSystemEvent{
							Timestamp: message.Timestamp,
						},
						NewWaId:           message.System.WaId,
						OldWaId:           message.System.Customer,
						ChangeDescription: message.System.Body,
					})
				}
			}
		case NotificationMessageTypeUnknown:
			{
				// ! TODO: handle error in the event and then emit it.
			}
		}
	}

	return nil
}

func (wh *WebhookManager) handleAccountAlertsSubscriptionEvents(value AccountAlertsValue) error {
	wh.EventManager.Publish(events.AccountAlertsEventType, events.AccountAlertEvent{})
	return nil
}

func (wh *WebhookManager) handleSecuritySubscriptionEvents(value SecurityValue) {
	wh.EventManager.Publish(events.AccountAlertsEventType, events.SecurityEvent{})

}

func (wh *WebhookManager) handleAccountUpdateSubscriptionEvents(value AccountUpdateValue) {
	wh.EventManager.Publish(events.AccountAlertsEventType, events.AccountUpdateEvent{})

}

func (wh *WebhookManager) handleAccountReviewSubscriptionEvents(value AccountReviewUpdateValue) {
	wh.EventManager.Publish(events.AccountAlertsEventType, events.AccountReviewUpdateEvent{})

}

func (wh *WebhookManager) handleBusinessCapabilitySubscriptionEvents(value BusinessCapabilityUpdateValue) {
	wh.EventManager.Publish(events.AccountAlertsEventType, events.BusinessCapabilityUpdateEvent{})

}

func (wh *WebhookManager) handleMessageTemplateQualitySubscriptionEvents(value TemplateQualityUpdateValue) {
	wh.EventManager.Publish(events.AccountAlertsEventType, events.MessageTemplateQualityUpdateEvent{})

}

func (wh *WebhookManager) handleMessageTemplateStatusSubscriptionEvents(value TemplateStatusUpdateValue) {
	wh.EventManager.Publish(events.AccountAlertsEventType, events.MessageTemplateStatusUpdateEvent{})

}

func (wh *WebhookManager) handlePhoneNumberNameSubscriptionEvents(value PhoneNumberNameUpdateValue) {
	wh.EventManager.Publish(events.AccountAlertsEventType, events.PhoneNumberNameUpdateEvent{})

}

func (wh *WebhookManager) handlePhoneNumberQualitySubscriptionEvents(value PhoneNumberQualityUpdateValue) {
	wh.EventManager.Publish(events.AccountAlertsEventType, events.PhoneNumberQualityUpdateEvent{})

}

func (wh *WebhookManager) handleTemplateCategoryUpdateSubscriptionEvents(value TemplateCategoryUpdateValue) {
	wh.EventManager.Publish(events.AccountAlertsEventType, events.TemplateCategoryUpdateEvent{})

}
