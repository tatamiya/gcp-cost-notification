package gcp_cost_notification

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCostNotifierAll(t *testing.T) {
	m := PubSubMessage{
		Data: []byte(""),
	}
	err := CostNotifier(context.Background(), m)
	assert.Nil(t, err)
}
