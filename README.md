<div align="center">
<br />
<p align="center">
<a href="https://wapijs.co"><img src="https://mintlify.s3-us-west-1.amazonaws.com/softlancer/assets/banner.svg" alt="Wapi.go" /></a>
</p>
<br />
</div>

Visit the documentation of the SDK [here](https://golang.wapikit.com)

## Status

Beta Version - This SDK is not stable right now. It is currently in beta version. Report issues [here](https://github.com/wapikit/wapi.go/issues).

This SDK is part of a technical suite built to support the WhatsApp Business Application Development ecosystem. This SDK also has a Node.js version, you can check it out [here](https://wapikit/wapi.js/js).

## About

Wapi.go is a Golang SDK, that supports WhatsApp API products i.e., Business Management API and Cloud API
to build WhatsApp applications easily.
This SDK supports managing WhatsApp business accounts, and managing phone numbers for a WhatsApp business account which includes creating, verifying, and registering a phone number to use for messaging via cloud API and deregistering a phone number. This SDK also supports the management of message templates which includes creating new templates or updating/deleting the existing ones.

You can listen to the incoming webhook events via the inbuilt standalone HTTP server built using echo and also you can integrate the SDK within your existing backend applications using the getters for the route handlers of the webhook server.

## Features

- Supports multiple phone number messaging clients.
- Supports WhatsApp
- Send Messages with the least configuration
- Event Listener for Notifications (support both User and System Notifications) and all the events provided by the WhatsApp Business Platform.
- Manage Phone numbers and messaging templates right from the SDK programmatically.
- Upload Media to WhatsApp servers
- Reply and React to incoming messages.

## Installation

This assumes you already have a working Go environment, if not please see
[this page](https://golang.org/doc/install) first.

`go get` _will always pull the latest tagged release from the master branch._

```sh
go get github.com/wapikit/wapi.go
```

> Note: This SDK is not affiliated with the official WhatsApp Cloud API or does not act as any official solution provided the the Meta Inclusive Private Limited, this is just a open source SDK built for developers to support them in building whatsapp cloud api based chat bots easily.

## Usage

You can check out the example WhatsApp bot here. [Example Chatbot](./example-chat-bot/)

# Example Usage

Import the package into your project.
This repository has three packages exported:

- github.com/wapikit/wapi.go/components
- github.com/wapikit/wapi.go/wapi/wapi
- github.com/wapikit/wapi.go/wapi/business
- github.com/wapikit/wapi.go/wapi/events

```go
import "github.com/wapikit/wapi.go/wapi/wapi"
```

Construct a new Wapi Client to access the managers in order to send messages and listen to incoming notifications.

```go
whatsappClient, err := wapi.New(wapi.ClientConfig{
		PhoneNumberId:     "",
		ApiAccessToken: "",
		BusinessAccountId: "",
		WebhookPath:       "/webhook",
		WebhookSecret:     "",
		WebhookServerPort: 8080,
	})
```

## References

- **Message Structures**: Refer to the WhatsApp Docs [here](https://developers.facebook.com/docs/whatsapp/cloud-api/reference/messages).

- **Notification Payloads**: Details can be found [here](https://developers.facebook.com/docs/whatsapp/cloud-api/webhooks/components).

## Contribution Guidelines

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

For detailed guidelines, check [Contributing.md](./CONTRIBUTING.md).

# TODOs'

- Handle errors gracefully
- Handle template and interactive messages gracefully
- Add support for more interactive message types like address input.

## License

Distributed under the Apache 2.0 License. View [LICENSE](./LICENSE).

## Contact

- Email: sarthak@softlancer.co
- [Twitter](https://twitter.com/wapikit) | [LinkedIn](https://www.linkedin.com/in/wapikit)
