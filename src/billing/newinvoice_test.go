package billing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tatamiya/gcp-cost-notification/src/datetime"
	"github.com/tatamiya/gcp-cost-notification/src/db"
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

	expectedInvoice := &Invoice{
		BillingPeriod: BillingPeriod{
			From: time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
			To:   time.Date(2021, 5, 8, 0, 0, 0, 0, time.Local),
		},
		Total: &Cost{Service: "Total", Monthly: 1000.07, Yesterday: 400.0},
		Services: []*Cost{
			{Service: "Cloud SQL", Monthly: 1000.0, Yesterday: 400.0},
			{Service: "BigQuery", Monthly: 0.07, Yesterday: 0.0},
		},
	}
	actualInvoice, err := NewInvoice(&inputReportingPeriod, inputQueryResults)

	assert.Nil(t, err)
	assert.EqualValues(t, expectedInvoice, actualInvoice)
}

func TestBillingsFromEmptyQueryResultHasZeroTotalCost(t *testing.T) {
	inputQueryResults := []*db.QueryResult{}

	actualInvoice, err := NewInvoice(&InputReportingPeriod, inputQueryResults)

	assert.Nil(t, err)
	assert.EqualValues(
		t,
		Cost{Service: "Total", Monthly: 0.00, Yesterday: 0.0},
		*actualInvoice.Total,
	)
	assert.EqualValues(t, []*Cost{}, actualInvoice.Services)
}

func TestBillingsFromSingleElementQueryResultHasEmptyServiceCosts(t *testing.T) {
	inputQueryResults := []*db.QueryResult{
		{Service: "Total", Monthly: 0.07, Yesterday: 0.0},
	}

	actualInvoice, err := NewInvoice(&InputReportingPeriod, inputQueryResults)

	assert.Nil(t, err)
	assert.EqualValues(
		t,
		Cost{Service: "Total", Monthly: 0.07, Yesterday: 0.0},
		*actualInvoice.Total,
	)
	assert.EqualValues(t, []*Cost{}, actualInvoice.Services)
}

func TestNewBillingsReturnErrorWhenQueryResultsUnexpectedlyOrderd(t *testing.T) {
	inputQueryResults := []*db.QueryResult{
		{Service: "Cloud SQL", Monthly: 1000.0, Yesterday: 400.0},
		{Service: "BigQuery", Monthly: 0.07, Yesterday: 0.0},
		{Service: "Total", Monthly: 1000.07, Yesterday: 400.0},
	}

	actualInvoice, err := NewInvoice(&InputReportingPeriod, inputQueryResults)

	assert.NotNil(t, err)
	assert.Nil(t, actualInvoice)
	assert.EqualValues(t,
		"Error in Query Results Validation. Unexpected query results! The results might not be correctly sorted!: First element of the query results was Cloud SQL, not Total",
		err.Error(),
	)
}
