package notification

import (
	"os"

	"github.com/slack-go/slack"
	"github.com/tatamiya/gcp-cost-notification/utils"
)

func NewSlackError(message string, err error) *utils.CustomError {
	return &utils.CustomError{
		Process: "Slack Notification",
		Message: message,
		Err:     err,
	}
}

type SlackClientInterface interface {
	Send(message string) error
}

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
	if err != nil {
		return NewSlackError(
			"Could not send message!",
			err,
		)
	}
	return nil
}
