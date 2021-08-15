package notification

import "github.com/tatamiya/gcp-cost-notification/utils"

type SlackClientStub struct {
	err *utils.CustomError
}

func NewSlackClientStub(err error) SlackClientStub {
	var slackError *utils.CustomError
	if err == nil {
		slackError = nil
	} else {
		slackError = NewSlackError("Failed", err)
	}
	return SlackClientStub{slackError}
}

func (c *SlackClientStub) Send(messenger Messenger) *utils.CustomError {
	return c.err
}
