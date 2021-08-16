package notification

import "github.com/tatamiya/gcp-cost-notification/src/utils"

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

func (c *SlackClientStub) Send(messenger Messenger) (string, *utils.CustomError) {
	return messenger.AsMessage(), c.err
}
