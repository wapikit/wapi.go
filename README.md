<div align="center">
<br />
<p align="center">
<a href="https://wapijs.co"><img src="https://mintlify.s3-us-west-1.amazonaws.com/softlancer/assets/banner.svg" alt="@wapijs/Wapi.go" /></a>
</p>
<br />
</div>

## 📌 Status

Beta Version - This SDK is not stable right now. It is currently in beta version. Report issues [here](https://github.com/sarthakjdev/wapi.go/issues).

## 📖 About

Wapi.go is a Golang SDK, that supports WhatsApp API products i.e., Business Management API and Cloud API
to build WhatsApp applications easily.

## ✨ Features

- Single Client Model
- Send Messages with the least configuration
- Event Listener for Notifications (support both User and System Notifications)
- Upload Media to WhatsApp servers
- Reply and React to incoming messages.

## 💻 Installation

This assumes you already have a working Go environment, if not please see
[this page](https://golang.org/doc/install) first.

`go get` _will always pull the latest tagged release from the master branch._

```sh
go get github.com/sarthakjdev/wapi.go
```

> Note: This SDK is not affiliated with the official WhatsApp Cloud API or does not act as any official solution provided the the Meta Inclusive Private Limited, this is just a open source SDK built for developers to support them in building whatsapp cloud api based chat bots easily.

## 🚀 Usage

You can check out the example WhatsApp bot here. [Example Chatbot](./example-chat-bot/)

# Example Usage

Import the package into your project.
This repository has three packages exported:

- github.com/sarthakjdev/wapi.go/components
- github.com/sarthakjdev/wapi.go/wapi/wapi
- github.com/sarthakjdev/wapi.go/wapi/business
- github.com/sarthakjdev/wapi.go/wapi/events

```go
import "github.com/sarthakjdev/wapi.go/wapi/wapi"
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

## 🔗 References

- **Message Structures**: Refer to the WhatsApp Docs [here](https://developers.facebook.com/docs/whatsapp/cloud-api/reference/messages).

- **Notification Payloads**: Details can be found [here](https://developers.facebook.com/docs/whatsapp/cloud-api/webhooks/components).

## 🤝 Contribution Guidelines

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

For detailed guidelines, check [Contributing.md](./CONTRIBUTING.md).

# TODOs'

- Handle errors gracefully
- Handle template and interactive messages gracefully
- Add support for more interactive messaged types like address input.
- 

## 📜 License

Distributed under the AGPL 3.0 License. View [LICENSE](./LICENSE).

## 📞 Contact

- [Sarthak Jain](https://sarthakjdev.com)
- Email: sarthak@softlancer.co
- [Twitter](https://twitter.com/sarthakjdev) | [LinkedIn](https://www.linkedin.com/in/sarthakjdev)

Note: This SDK is part of an open-source product-building initiative by [Softlancer](https://github.com/softlancerhq), and this repository will soon be moved under the same organization.
