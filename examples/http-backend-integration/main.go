package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	wapi "github.com/wapikit/wapi.go/pkg/client"
	"github.com/wapikit/wapi.go/pkg/components"
	"github.com/wapikit/wapi.go/pkg/events"
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
		if _, err := textMessageEvent.Reply(reply); err != nil {
			fmt.Println("error sending reply", err)
		}
	})

	getHandler := client.GetWebhookGetRequestHandler()
	postHandler := client.GetWebhookPostRequestHandler()

	server := echo.New()

	server.GET("/webhook", getHandler)
	server.POST("/webhook", postHandler)

	if err := server.Start(":8080"); err != nil {
		fmt.Println("server stopped", err)
	}

}
