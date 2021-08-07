package notification

type SlackClientStub struct {
	err error
}

func NewSlackClientStub(err error) SlackClientStub {
	return SlackClientStub{err}
}

func (c *SlackClientStub) Send(message string) error {
	return c.err
}
