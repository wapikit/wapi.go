package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	wapi "github.com/gTahidi/wapi.go/pkg/client"
	"github.com/gTahidi/wapi.go/pkg/components"
	"github.com/gTahidi/wapi.go/pkg/events"
)

func main() {

	client := wapi.New(&wapi.ClientConfig{
		ApiAccessToken:    "",
		BusinessAccountId: "",
		WebhookPath:       "/webhook",
		WebhookSecret:     "1234567890",
		WebhookServerPort: 8080,
	})

	client.On(events.TextMessageEventType, func(event events.BaseEvent) {
		textMessageEvent := event.(*events.TextMessageEvent)
		reply, err := components.NewTextMessage(components.TextMessageConfigs{
			Text: "Hello, from wapi.go",
		})
		if err != nil {
			fmt.Println("error creating text message", err)
			return
		}
		textMessageEvent.Reply(reply)
	})

	getHandler := client.GetWebhookGetRequestHandler()
	postHandler := client.GetWebhookPostRequestHandler()

	server := echo.New()

	server.GET("/webhook", getHandler)
	server.POST("/webhook", postHandler)

	server.Start(":8080")

}
