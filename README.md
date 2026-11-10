<div align="center">
<br />
<p align="center">
<a href="https://golang.wapikit.com"><img src="https://mintlify.s3-us-west-1.amazonaws.com/softlancer/assets/banner.svg" alt="Wapi.go" /></a>
</p>
<br />

[![Go Reference](https://pkg.go.dev/badge/github.com/wapikit/wapi.go.svg)](https://pkg.go.dev/github.com/wapikit/wapi.go)
[![Test](https://github.com/wapikit/wapi.go/actions/workflows/test.yaml/badge.svg)](https://github.com/wapikit/wapi.go/actions/workflows/test.yaml)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](./LICENSE)

</div>

**wapi.go** is a Go SDK for the [WhatsApp Business Platform](https://developers.facebook.com/docs/whatsapp): the Cloud API (sending messages and receiving webhooks) and the Business Management API (business accounts, phone numbers, templates and catalogs).

📖 Full documentation: [golang.wapikit.com](https://golang.wapikit.com) · API reference: [pkg.go.dev](https://pkg.go.dev/github.com/wapikit/wapi.go)

> This SDK is not affiliated with or endorsed by Meta. It is an open source project to make building on the WhatsApp Cloud API easier. A Node.js version is available as [wapi.js](https://github.com/wapikit/wapi.js).

## Contents

- [Features](#features)
- [Installation](#installation)
- [Quick start](#quick-start)
- [Sending messages](#sending-messages)
- [Receiving webhooks](#receiving-webhooks)
- [Business management](#business-management)
- [Error handling](#error-handling)
- [Configuration](#configuration)
- [Packages](#packages)
- [Examples](#examples)
- [Contributing](#contributing)

## Features

- **Messages:** send text, media (image, video, audio, document, sticker), location, contacts, reactions, templates and interactive messages (reply buttons, lists, CTA URL buttons, product, product list and catalog messages, location and contact requests).
- **Recipients:** send to a phone number or to a business-scoped user id (BSUID).
- **Read receipts:** mark messages as read, optionally with a typing indicator.
- **Media:** upload media, resolve a media id to its download URL, and delete media.
- **Webhooks:** verify the subscription, optionally check `X-Hub-Signature-256`, and get typed events for incoming messages, delivery statuses (including pricing) and account notifications.
- **Integration:** run the built-in webhook server, or mount the handlers in your own [Echo](https://echo.labstack.com) app.
- **Business Management API:** business account details and analytics, phone numbers (registration, verification, QR codes, usernames), message templates and product catalogs, with pagination helpers.
- **Errors:** Graph API failures come back as a typed `GraphAPIError` with the HTTP status, Meta's error code and the raw body.
- **Timeouts:** requests share one HTTP client with sensible timeouts. The client, API version and base URL can all be replaced.

## Installation

Requires Go 1.21 or later.

```sh
go get github.com/wapikit/wapi.go
```

## Quick start

You need a Meta app with the WhatsApp product, a business account id, a phone number id and an access token. The [setup guide](https://golang.wapikit.com) walks through getting them.

This bot replies "Hello, from wapi.go" to every text message it receives:

```go
package main

import (
	"log"

	wapi "github.com/wapikit/wapi.go/pkg/client"
	"github.com/wapikit/wapi.go/pkg/components"
	"github.com/wapikit/wapi.go/pkg/events"
)

func main() {
	client := wapi.New(&wapi.ClientConfig{
		BusinessAccountId: "<WHATSAPP_BUSINESS_ACCOUNT_ID>",
		ApiAccessToken:    "<ACCESS_TOKEN>",
		WebhookSecret:     "<WEBHOOK_VERIFY_TOKEN>", // the verify token you set in the Meta app dashboard
		WebhookPath:       "/webhook",
		WebhookServerPort: 8080,
	})

	client.On(events.TextMessageEventType, func(event events.BaseEvent) {
		message := event.(*events.TextMessageEvent)

		reply, err := components.NewTextMessage(components.TextMessageConfigs{
			Text: "Hello, from wapi.go",
		})
		if err != nil {
			log.Println("build reply:", err)
			return
		}
		if _, err := message.Reply(reply); err != nil {
			log.Println("send reply:", err)
		}
	})

	// Starts the webhook server on 127.0.0.1:8080 and blocks until Ctrl+C.
	client.Initiate()
}
```

The built-in server listens on `127.0.0.1` only, so put a reverse proxy or a tunnel such as ngrok in front of it. Then point your Meta app's webhook at `https://<your-domain>/webhook` and subscribe to the `messages` field.

## Sending messages

Messages are sent through a messaging client for one of your business phone numbers. Every message type has a constructor in `pkg/components` that validates its input.

```go
messaging := client.NewMessagingClient("<PHONE_NUMBER_ID>")

text, err := components.NewTextMessage(components.TextMessageConfigs{
	Text:         "Your order has shipped: https://example.com/track/123",
	AllowPreview: true,
})
if err != nil {
	return err
}

resp, err := messaging.Message.Send(text, "15551234567")
if err != nil {
	return err
}
log.Println("sent message", resp.Messages[0].ID)
```

### Media

```go
image, err := components.NewImageMessage(components.ImageMessageConfigs{
	Link:    "https://example.com/cat.png", // or Id: a media id returned by UploadMedia
	Caption: "Look at this",
})
if err != nil {
	return err
}
_, err = messaging.Message.Send(image, "15551234567")
```

Upload a file once and send it by id:

```go
file, err := os.Open("invoice.pdf")
if err != nil {
	return err
}
defer file.Close()

mediaId, err := messaging.Media.UploadMedia("<PHONE_NUMBER_ID>", file, "invoice.pdf", "application/pdf")
```

`Media.GetMediaUrlById(id)` resolves the media id from an incoming message into a download URL. `Media.DeleteMedia(id)` deletes uploaded media.

### Templates

Outside the 24-hour customer service window, only approved templates can be sent:

```go
template, err := components.NewTemplateMessage(&components.TemplateMessageConfigs{
	Name:     "order_shipped",
	Language: "en_US",
})
if err != nil {
	return err
}

name := "Ada"
template.AddBody(components.TemplateMessageComponentBodyType{
	Type: components.TemplateMessageComponentTypeBody,
	Parameters: []components.TemplateMessageParameter{
		components.TemplateMessageBodyAndHeaderParameter{
			Type: components.TemplateMessageParameterTypeText,
			Text: &name,
		},
	},
})

_, err = messaging.Message.Send(template, "15551234567")
```

`AddHeader`, `AddButton` and `AddCarousel` add the other template components. `Message.SendMarketingMessage` sends through the Marketing Messages API, for business accounts onboarded to it.

### Interactive messages

```go
buttons, err := components.NewQuickReplyButtonMessage("How was your delivery?")
if err != nil {
	return err
}
if err := buttons.AddButton("good", "Great"); err != nil {
	return err
}
if err := buttons.AddButton("bad", "Not great"); err != nil {
	return err
}
_, err = messaging.Message.Send(buttons, "15551234567")
```

Lists (`NewListMessage`), CTA URL buttons (`NewCtaMessage`), products (`NewProductMessage`, `NewProductListMessage`, `NewCatalogMessage`) and location and contact requests (`NewLocationRequestMessage`, `NewContactRequestMessage`) follow the same pattern.

### Replies, reactions and read receipts

```go
// Reply in-thread to a message you received:
messaging.Message.Reply(text, "15551234567", "<INCOMING_MESSAGE_ID>")

// Or, inside an event handler:
message.Reply(text)
message.React("👍")

// Mark a message as read, with or without a typing indicator:
messaging.Message.ReadMessageOnly("<INCOMING_MESSAGE_ID>")
messaging.Message.ReadMessageWithTyping("<INCOMING_MESSAGE_ID>")
```

### Sending to a business-scoped user id

Users who chat through a WhatsApp username may only be known by their BSUID. Send to them with a message target:

```go
import "github.com/wapikit/wapi.go/manager"

resp, err := messaging.Message.SendToTarget(text, manager.NewBSUIDTarget("<BSUID>"))
```

`manager.NewPhoneTarget` and `manager.NewParentBSUIDTarget` build the other target types. `ReplyToTarget` and `SendMarketingMessageToTarget` are the target-based versions of `Reply` and `SendMarketingMessage`.

## Receiving webhooks

### Subscribing to events

Register a handler per event type with `client.On`. Each handler runs in its own goroutine and receives every event of its type, one at a time. Type-assert the event to its concrete pointer type:

```go
client.On(events.ImageMessageEventType, func(event events.BaseEvent) {
	image := event.(*events.ImageMessageEvent)
	url, err := messaging.Media.GetMediaUrlById(image.MediaId)
	// download url with your access token ...
})

client.On(events.MessageFailedEventType, func(event events.BaseEvent) {
	failed := event.(*events.MessageFailedEvent)
	log.Printf("message %s to %s failed: %d %s", failed.MessageId, failed.SentTo, failed.ErrorCode, failed.ErrorMessage)
})
```

| Category | Event types |
| --- | --- |
| Incoming messages | `TextMessageEventType`, `ImageMessageEventType`, `VideoMessageEventType`, `AudioMessageEventType`, `DocumentMessageEventType`, `StickerMessageEventType`, `LocationMessageEventType`, `ContactMessageEventType`, `ReactionMessageEventType` |
| Interactions | `QuickReplyMessageEventType`, `ReplyButtonInteractionEventType`, `ListInteractionMessageEventType`, `OrderReceivedEventType`, `AdInteractionEventType` |
| Delivery statuses | `MessageSentEventType`, `MessageDeliveredEventType`, `MessageReadEventType`, `MessageFailedEventType`, `MessageUndeliveredEventType` |
| Account and identity | `AccountAlertsEventType`, `AccountUpdateEventType`, `UserPreferencesEventType`, `UserIdUpdateEventType`, `BusinessUsernameUpdateEventType`, `CustomerNumberChangedEventType` |
| Templates and other | `MessageTemplateComponentsUpdateEventType`, `PaymentConfigurationUpdateEventType`, `HistoryEventType`, `SmbAppStateSyncEventType`, `SmbMessageEchoesEventType`, `ReadyEventType` |

> **Note:** template status, quality and category updates, account review, business capability, phone number name and quality, and security notifications are currently all delivered under `AccountAlertsEventType`. Switch on the concrete type inside that handler (for example `*events.MessageTemplateStatusUpdateEvent` or `*events.PhoneNumberQualityUpdateEvent`). Their dedicated event types are defined but are not published yet.

Status events carry Meta's pricing block in `Pricing` when it is present. Message events carry the sender's identity: `From`, `SenderName`, and the BSUID fields `UserId`, `ParentUserId` and `Username`.

If one message or change in a webhook delivery can't be parsed, it is logged with `log/slog` and skipped, and the rest of the batch is still delivered. A panic in your handler is recovered and logged, so it can't take down the process.

### Using your own server

Instead of `client.Initiate()`, mount the handlers on an existing Echo app:

```go
e := echo.New()
e.GET("/webhook", client.GetWebhookGetRequestHandler())   // subscription verification
e.POST("/webhook", client.GetWebhookPostRequestHandler()) // event delivery
e.Logger.Fatal(e.Start(":8080"))
```

### Verifying signatures

Meta signs every webhook delivery with your app secret. Set `WebhookAppSecret` and the POST handler rejects any request whose `X-Hub-Signature-256` header doesn't match the raw body, returning 401:

```go
client := wapi.New(&wapi.ClientConfig{
	// ...
	WebhookAppSecret: os.Getenv("META_APP_SECRET"),
})
```

If you need to check the signature yourself (for example in middleware that changes the body), use `wapi.VerifyWebhookSignature(appSecret, body, header)`.

## Business management

`client.Business` wraps the Business Management API for the configured business account:

```go
account, err := client.Business.Fetch()

phone, err := client.Business.PhoneNumber.Fetch("<PHONE_NUMBER_ID>")

template, err := client.Business.Template.Fetch("<TEMPLATE_ID>")

// Page through templates:
page, err := client.Business.Template.FetchAllWithPagination(&manager.PaginationInput{Limit: 50})
for err == nil && page.Paging.HasNextPage() {
	page, err = client.Business.Template.FetchAllWithPagination(&manager.PaginationInput{
		Limit: 50,
		After: page.Paging.GetNextCursor(),
	})
}

// Or fetch every page at once:
templates, err := client.Business.Template.GetAllTemplates()

catalogs, err := client.Business.Catalog.GetAllCatalogsWithPagination()
products, err := client.Business.Catalog.GetAllCatalogProducts("<CATALOG_ID>")
```

Also available:

- **Phone numbers:** `Create`, `RequestVerificationCode`, `VerifyCode`, QR codes (`GenerateQrCode`, `GetAllQrCodes`, `GetQrCodeById`, `UpdateQrCode`, `DeleteQrCode`) and usernames (`GetUsername`, `SetUsername`, `DeleteUsername`).
- **Templates:** `Create`, `Update` and `MigrateFromOtherBusinessAccount`.
- **Analytics:** `client.Business.FetchAnalytics` and `client.Business.ConversationAnalytics`.
- **Registration:** `messaging.Register(pin)` and `messaging.Deregister()` register a phone number with the Cloud API and remove it again.

## Error handling

When the Graph API responds with a non-2xx status, the error is a `*wapi.GraphAPIError`. It carries the HTTP status, Meta's error code, subcode, message and `fbtrace_id`, and the raw response body:

```go
resp, err := messaging.Message.Send(text, "15551234567")
if apiErr, ok := wapi.AsGraphAPIError(err); ok {
	log.Printf("graph error: status=%d code=%d trace=%s: %s",
		apiErr.StatusCode, apiErr.Code, apiErr.FBTraceID, apiErr.Message)
	if apiErr.IsRetryable() { // 429 or 5xx
		// back off and retry
	}
}
```

For message sends, the parsed response is returned along with the error when Meta includes one, so `resp.Error` gives you Meta's `error_data.details` too. Network failures and timeouts are returned as the standard `net/http` errors.

## Configuration

Every field below `WebhookServerPort` is optional:

```go
client := wapi.New(&wapi.ClientConfig{
	BusinessAccountId: "...",
	ApiAccessToken:    "...",
	WebhookSecret:     "...", // webhook verify token (required)
	WebhookPath:       "/webhook",
	WebhookServerPort: 8080,

	WebhookAppSecret: "...",                                 // check X-Hub-Signature-256 on deliveries
	HTTPClient:       &http.Client{Timeout: 2 * time.Minute}, // replace the default HTTP client
	ApiVersion:       "v21.0",                                // Graph API version (default "v20.0")
	GraphBaseURL:     "http://127.0.0.1:9999",                // point at a test server
})
```

- **Default HTTP client:** shared by all clients. It bounds connecting, the TLS handshake and waiting for Meta's response headers (60s), but sets no overall timeout, so large media uploads are not cut off.
- **Bounding calls:** manager methods don't take a `context` yet. To cap how long any call can take, pass an `HTTPClient` with a `Timeout`.
- **Testing:** `GraphBaseURL` makes it easy to test your integration against an `httptest.Server`.
- **Logging:** the library logs through `log/slog`'s default logger. Call `slog.SetDefault` to route its logs into your own logger.

## Packages

| Package | Import path | Contents |
| --- | --- | --- |
| `wapi` | `github.com/wapikit/wapi.go/pkg/client` | `New`, `ClientConfig`, webhook handlers, `GraphAPIError` |
| `components` | `github.com/wapikit/wapi.go/pkg/components` | Message builders for every message type |
| `events` | `github.com/wapikit/wapi.go/pkg/events` | Event types and event structs |
| `manager` | `github.com/wapikit/wapi.go/manager` | Message, media, template, phone number and catalog managers, response types, message targets |
| `business` | `github.com/wapikit/wapi.go/pkg/business` | Business account client (reached through `client.Business`) |
| `messaging` | `github.com/wapikit/wapi.go/pkg/messaging` | Per-phone-number messaging client (from `client.NewMessagingClient`) |

## Examples

- [`examples/chat-bot`](./examples/chat-bot): a bot that answers keywords with each message type, using the built-in server.
- [`examples/http-backend-integration`](./examples/http-backend-integration): the webhook handlers mounted on your own Echo server.

## References

- [Cloud API: messages](https://developers.facebook.com/docs/whatsapp/cloud-api/reference/messages)
- [Webhooks: payloads](https://developers.facebook.com/docs/whatsapp/cloud-api/webhooks/components)
- [Business Management API](https://developers.facebook.com/docs/whatsapp/business-management-api)

## Contributing

Contributions are welcome. See [CONTRIBUTING.md](./CONTRIBUTING.md) for the workflow and [COMMIT_CONVENTION.md](./COMMIT_CONVENTION.md) for commit messages; releases are cut automatically from Conventional Commits. Before opening a PR, run:

```sh
gofmt -l . && go vet ./... && go test -race ./...
```

Report bugs and ideas on [GitHub issues](https://github.com/wapikit/wapi.go/issues).

## License

Distributed under the Apache 2.0 License. See [LICENSE](./LICENSE).

## Contact

- Email: sarthak@softlancer.co
- [Twitter](https://twitter.com/wapikit) | [LinkedIn](https://www.linkedin.com/in/wapikit)
