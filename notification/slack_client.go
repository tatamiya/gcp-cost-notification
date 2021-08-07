package notification

import (
	"os"

	"github.com/slack-go/slack"
)

type SlackClient struct {
	webhookURL string
}

func NewSlackClient() SlackClient {
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	return SlackClient{webhookURL: webhookURL}
}

func (c *SlackClient) Send(message string) error {
	msg := slack.WebhookMessage{
		Text: message,
	}
	err := slack.PostWebhook(c.webhookURL, &msg)
	return err
}

func sendMessageToSlack(webhookURL string, messageText string) error {
	msg := slack.WebhookMessage{
		Text: messageText,
	}
	err := slack.PostWebhook(webhookURL, &msg)
	return err
}
