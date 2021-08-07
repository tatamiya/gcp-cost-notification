package notification

type SlackClientStub struct {
	err error
}

func NewSlackClientStub(err error) SlackClientStub {
	var slackError error
	if err == nil {
		slackError = nil
	} else {
		slackError = NewSlackError("Failed", err)
	}
	return SlackClientStub{slackError}
}

func (c *SlackClientStub) Send(message string) error {
	return c.err
}
