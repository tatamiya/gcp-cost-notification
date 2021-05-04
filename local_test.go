package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCostNotifier(t *testing.T) {
	m := PubSubMessage{
		Data: []byte(""),
	}
	err := CostNotifier(context.Background(), m)
	assert.Nil(t, err)
}
