package gcp_cost_notification

import (
	"context"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/stretchr/testify/assert"
)

func TestCostNotifier(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping")
	}
	m := pubsub.Message{}
	err := CostNotifier(context.Background(), m)
	assert.Nil(t, err)
}
