// notification package implements an object to send a message to Slack
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

// An object implemented with Messenger interface
// can be converted into a notification message to send to Slack.
type Messenger interface {
	AsMessage() string
}

// SlackClient is an object to send a message to Slack
// via webhook URL.
type SlackClient struct {
	webhookURL string
}

// NewSlackClient constructs a SlackClient object.
// The webhook URL is fetched from the environment variable
// `SLACK_WEBHOOK_URL` in construction.
func NewSlackClient() SlackClient {
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	return SlackClient{webhookURL: webhookURL}
}

// Send method receives an object which can be converted into
// a notification message and sends it to Slack.
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
