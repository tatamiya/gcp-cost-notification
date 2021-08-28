package billing

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ExampleBillingPeriod_String() {
	period := BillingPeriod{
		From: time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
		To:   time.Date(2021, 5, 8, 0, 0, 0, 0, time.Local),
	}
	fmt.Println(period.String())
	// Output: 5/1 ~ 5/8

}

func TestCreateSingleMessageLine(t *testing.T) {
	sampleCost := &Cost{
		Service: "Cloud SQL", Monthly: 1000.0, Yesterday: 400.0,
	}
	expectedLine := "Cloud SQL: ¥ 1,000 (¥ 400)"
	actualLine := sampleCost.asMessageLine()

	assert.EqualValues(t, expectedLine, actualLine)
}

func TestDescribeDetailsCorrectly(t *testing.T) {
	inputInvoice := &Invoice{
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
	expectedDetailLines := "Cloud SQL: ¥ 1,000 (¥ 400)\nBigQuery: ¥ 0.07 (¥ 0)"

	actualDetailLines := inputInvoice.details()
	assert.EqualValues(t, expectedDetailLines, actualDetailLines)
}

func TestShowNoDetailWhenServiceCostIsEmpty(t *testing.T) {
	inputInvoice := &Invoice{
		BillingPeriod: BillingPeriod{
			From: time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
			To:   time.Date(2021, 5, 8, 0, 0, 0, 0, time.Local),
		},
		Total:    &Cost{Service: "Total", Monthly: 1000.07, Yesterday: 400.0},
		Services: []*Cost{},
	}
	expectedDetailLines := ""

	actualDetailLines := inputInvoice.details()
	assert.EqualValues(t, expectedDetailLines, actualDetailLines)
}

func TestCreateMessageFromInvoice(t *testing.T) {
	inputInvoice := &Invoice{
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

	expectedMessage :=
		`＜5/1 ~ 5/8 の GCP 利用料金＞ ※ () 内は前日分

Total: ¥ 1,000.07 (¥ 400)

----- 内訳 -----
Cloud SQL: ¥ 1,000 (¥ 400)
BigQuery: ¥ 0.07 (¥ 0)`

	actualMessage := inputInvoice.AsMessage()
	assert.EqualValues(t, expectedMessage, actualMessage)
}

func TestCreateMessageFromInvoiceWithNoServiceCosts(t *testing.T) {
	inputInvoice := &Invoice{
		BillingPeriod: BillingPeriod{
			From: time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
			To:   time.Date(2021, 5, 8, 0, 0, 0, 0, time.Local),
		},
		Total:    &Cost{Service: "Total", Monthly: 0.00, Yesterday: 0.00},
		Services: []*Cost{},
	}

	expectedMessage :=
		`＜5/1 ~ 5/8 の GCP 利用料金＞ ※ () 内は前日分

Total: ¥ 0 (¥ 0)`

	actualMessage := inputInvoice.AsMessage()
	assert.EqualValues(t, expectedMessage, actualMessage)
}
