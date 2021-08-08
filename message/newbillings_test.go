package message

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tatamiya/gcp-cost-notification/datetime"
	"github.com/tatamiya/gcp-cost-notification/db"
)

var InputReportingPeriod datetime.ReportingPeriod = datetime.ReportingPeriod{
	From: time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
	To:   time.Date(2021, 5, 8, 0, 0, 0, 0, time.Local),
}

func TestCreateBillingsCorrectly(t *testing.T) {
	inputReportingPeriod := InputReportingPeriod
	inputQueryResults := []*db.QueryResult{
		{Service: "Total", Monthly: 1000.07, Yesterday: 400.0},
		{Service: "Cloud SQL", Monthly: 1000.0, Yesterday: 400.0},
		{Service: "BigQuery", Monthly: 0.07, Yesterday: 0.0},
	}

	expectedBillings := &Billings{
		AggregationPeriod: AggregationPeriod{
			From: time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
			To:   time.Date(2021, 5, 8, 0, 0, 0, 0, time.Local),
		},
		Total: &db.QueryResult{Service: "Total", Monthly: 1000.07, Yesterday: 400.0},
		Services: []*db.QueryResult{
			{Service: "Cloud SQL", Monthly: 1000.0, Yesterday: 400.0},
			{Service: "BigQuery", Monthly: 0.07, Yesterday: 0.0},
		},
	}
	actualBillings, err := NewBillings(&inputReportingPeriod, inputQueryResults)

	assert.Nil(t, err)
	assert.EqualValues(t, expectedBillings, actualBillings)
}

func TestBillingsFromEmptyQueryResultHasZeroTotalCost(t *testing.T) {
	inputQueryResults := []*db.QueryResult{}

	actualBillings, err := NewBillings(&InputReportingPeriod, inputQueryResults)

	assert.Nil(t, err)
	assert.EqualValues(
		t,
		db.QueryResult{Service: "Total", Monthly: 0.00, Yesterday: 0.0},
		*actualBillings.Total,
	)
	assert.EqualValues(t, []*db.QueryResult{}, actualBillings.Services)
}

func TestBillingsFromSingleElementQueryResultHasEmptyServiceCosts(t *testing.T) {
	inputQueryResults := []*db.QueryResult{
		{Service: "Total", Monthly: 0.07, Yesterday: 0.0},
	}

	actualBillings, err := NewBillings(&InputReportingPeriod, inputQueryResults)

	assert.Nil(t, err)
	assert.EqualValues(
		t,
		db.QueryResult{Service: "Total", Monthly: 0.07, Yesterday: 0.0},
		*actualBillings.Total,
	)
	assert.EqualValues(t, []*db.QueryResult{}, actualBillings.Services)
}

func TestNewBillingsReturnErrorWhenQueryResultsUnexpectedlyOrderd(t *testing.T) {
	inputQueryResults := []*db.QueryResult{
		{Service: "Cloud SQL", Monthly: 1000.0, Yesterday: 400.0},
		{Service: "BigQuery", Monthly: 0.07, Yesterday: 0.0},
		{Service: "Total", Monthly: 1000.07, Yesterday: 400.0},
	}

	actualBillings, err := NewBillings(&InputReportingPeriod, inputQueryResults)

	assert.NotNil(t, err)
	assert.Nil(t, actualBillings)
	assert.EqualValues(t,
		"Error in Query Results Parser. Unexpected query results! The results might not be correctly sorted!: First element of the query results was Cloud SQL, not Total",
		err.Error(),
	)
}
