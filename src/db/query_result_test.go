package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorrectlyDisplayStringOfQueryResult(t *testing.T) {
	sampleQueryResult := &QueryResult{
		Service: "Cloud SQL", Monthly: 1000.0, Yesterday: 400.0,
	}
	expectedString := "{Service: Cloud SQL, Monthly: 1000.000000, Yesterday: 400.000000}"
	actualString := sampleQueryResult.String()

	assert.EqualValues(t, expectedString, actualString)

}
