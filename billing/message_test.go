package billing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tatamiya/gcp-cost-notification/db"
)

func TestAggregationPeriodIntoStringCorrectly(t *testing.T) {
	period := AggregationPeriod{
		From: time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
		To:   time.Date(2021, 5, 8, 0, 0, 0, 0, time.Local),
	}
	assert.EqualValues(t, "5/1 ~ 5/8", period.String())
}

func TestCreateDetailLinesCorrectly(t *testing.T) {
	inputBillings := &Billings{
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
	expectedDetailLines := "Cloud SQL: ¥ 1,000 (¥ 400)\nBigQuery: ¥ 0.07 (¥ 0)"

	actualDetailLines := inputBillings.detailLines()
	assert.EqualValues(t, expectedDetailLines, actualDetailLines)
}

func TestCreateBlankDetailLineWhenServiceCostIsEmpty(t *testing.T) {
	inputBillings := &Billings{
		AggregationPeriod: AggregationPeriod{
			From: time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
			To:   time.Date(2021, 5, 8, 0, 0, 0, 0, time.Local),
		},
		Total:    &db.QueryResult{Service: "Total", Monthly: 1000.07, Yesterday: 400.0},
		Services: []*db.QueryResult{},
	}
	expectedDetailLines := ""

	actualDetailLines := inputBillings.detailLines()
	assert.EqualValues(t, expectedDetailLines, actualDetailLines)
}

func TestCreateNotificationFromBillings(t *testing.T) {
	inputBillings := &Billings{
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

	expectedNotification :=
		`＜5/1 ~ 5/8 の GCP 利用料金＞ ※ () 内は前日分

Total: ¥ 1,000.07 (¥ 400)

----- 内訳 -----
Cloud SQL: ¥ 1,000 (¥ 400)
BigQuery: ¥ 0.07 (¥ 0)`

	actualNotification := inputBillings.AsMessage()
	assert.EqualValues(t, expectedNotification, actualNotification)
}

func TestCreateNotificationFromBillingsWithNoServiceCosts(t *testing.T) {
	inputBillings := &Billings{
		AggregationPeriod: AggregationPeriod{
			From: time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
			To:   time.Date(2021, 5, 8, 0, 0, 0, 0, time.Local),
		},
		Total:    &db.QueryResult{Service: "Total", Monthly: 0.00, Yesterday: 0.00},
		Services: []*db.QueryResult{},
	}

	expectedNotification :=
		`＜5/1 ~ 5/8 の GCP 利用料金＞ ※ () 内は前日分

Total: ¥ 0 (¥ 0)`

	actualNotification := inputBillings.AsMessage()
	assert.EqualValues(t, expectedNotification, actualNotification)
}
