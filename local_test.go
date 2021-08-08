package gcp_cost_notification

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCostNotifier(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping")
	}
	m := PubSubMessage{
		Data: []byte(""),
	}
	err := CostNotifier(context.Background(), m)
	assert.Nil(t, err)
}
