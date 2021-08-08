package db

import (
	"fmt"
	"os"
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

func TestSendQueryToBQ(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping")
	}
	projectID := os.Getenv("GCP_PROJECT")
	inputQuery := fmt.Sprintf("SELECT * FROM `%s.gcp_costs.test_cost_notiification`", projectID)

	testClient := NewBQClient()

	actualOutput, err := testClient.SendQuery(inputQuery)
	assert.Nil(t, err)

	expectedOutput := []*QueryResult{
		{Service: "Total", Monthly: 100.0, Yesterday: 100.0},
		{Service: "BigQuery", Monthly: 90.0, Yesterday: 10.0},
	}
	assert.EqualValues(t, expectedOutput, actualOutput)
}
