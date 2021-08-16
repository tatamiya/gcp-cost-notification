package notification

import (
	"os"

	"github.com/slack-go/slack"
	"github.com/tatamiya/gcp-cost-notification/src/utils"
)

func NewSlackError(message string, err error) *utils.CustomError {
	return &utils.CustomError{
		Process: "Slack Notification",
		Message: message,
		Err:     err,
	}
}

type Messenger interface {
	AsMessage() string
}

type SlackClient struct {
	webhookURL string
}

func NewSlackClient() SlackClient {
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	return SlackClient{webhookURL: webhookURL}
}

func (c *SlackClient) Send(messenger Messenger) (string, *utils.CustomError) {
	message := messenger.AsMessage()
	msg := slack.WebhookMessage{
		Text: message,
	}
	err := slack.PostWebhook(c.webhookURL, &msg)
	if err != nil {
		return "", NewSlackError(
			"Could not send message!",
			err,
		)
	}
	return message, nil
}
