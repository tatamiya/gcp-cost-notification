package notification

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type messengerStub struct {
	message string
}

func (m *messengerStub) AsMessage() string {
	return m.message
}

func TestSlackPost(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping")
	}
	testClient := NewSlackClient()
	testMessenger := messengerStub{
		message: "test\nこれはテスト投稿です。",
	}

	sentMessage, err := testClient.Send(&testMessenger)
	assert.Nil(t, err)
	assert.EqualValues(t, "test\nこれはテスト投稿です。", sentMessage)
}
