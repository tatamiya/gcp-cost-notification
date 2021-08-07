package notification

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlackPost(t *testing.T) {
	testClient := NewSlackClient()
	inputMessage := "test\nこれはテスト投稿です。"

	err := testClient.Send(inputMessage)
	assert.Nil(t, err)
}
