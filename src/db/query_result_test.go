package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSingleMessageLine(t *testing.T) {
	sampleQueryResult := &QueryResult{
		Service: "Cloud SQL", Monthly: 1000.0, Yesterday: 400.0,
	}
	expectedLine := "Cloud SQL: ¥ 1,000 (¥ 400)"
	actualLine := sampleQueryResult.AsMessageLine()

	assert.EqualValues(t, expectedLine, actualLine)
}

func TestCorrectlyDisplayStringOfQueryResult(t *testing.T) {
	sampleQueryResult := &QueryResult{
		Service: "Cloud SQL", Monthly: 1000.0, Yesterday: 400.0,
	}
	expectedString := "{Service: Cloud SQL, Monthly: 1000.000000, Yesterday: 400.000000}"
	actualString := sampleQueryResult.String()

	assert.EqualValues(t, expectedString, actualString)

}
