package gcp_cost_notification

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBuildQuery(t *testing.T) {
	inputTableName := "sample_table_name"
	inputTimestamp := "2020-01-01T09:00:00Z"
	outputQuery := buildQuery(inputTableName, inputTimestamp)

	assert.EqualValues(t, true, strings.Contains(outputQuery, "SELECT"))
	assert.EqualValues(t, true, strings.Contains(outputQuery, inputTimestamp))
	assert.EqualValues(t, true, strings.Contains(outputQuery, inputTableName))
}

func TestSendQueryToBQ(t *testing.T) {
	projectID := os.Getenv("GCP_PROJECT")
	inputQuery := fmt.Sprintf("SELECT * FROM `%s.gcp_costs.test_cost_notiification`", projectID)

	actualOutput, err := sendQueryToBQ(inputQuery, projectID)
	assert.Nil(t, err)

	expectedOutput := []*QueryResult{
		{Service: "Total", Monthly: 100.0, Yesterday: 100.0},
		{Service: "BigQuery", Monthly: 90.0, Yesterday: 10.0},
	}
	assert.EqualValues(t, expectedOutput, actualOutput)
}

func TestCreateSingleMessageLine(t *testing.T) {
	sampleQueryResult := &QueryResult{
		Service: "Cloud SQL", Monthly: 1000.0, Yesterday: 400.0,
	}
	expectedLine := "Cloud SQL: ¥ 1,000 (¥ 400)"
	actualLine := sampleQueryResult.asMessageLine()

	assert.EqualValues(t, expectedLine, actualLine)
}

func TestCreateBillings(t *testing.T) {
	inputQueryResults := []*QueryResult{
		{Service: "Total", Monthly: 1000.07, Yesterday: 400.0},
		{Service: "Cloud SQL", Monthly: 1000.0, Yesterday: 400.0},
		{Service: "BigQuery", Monthly: 0.07, Yesterday: 0.0},
	}
	inputReportingPeriod := ReportingPeriod{
		From: time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
		To:   time.Date(2021, 5, 8, 0, 0, 0, 0, time.Local),
	}

	expectedBillings := &Billings{
		AggregationPeriod: AggregationPeriod{
			From: time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
			To:   time.Date(2021, 5, 8, 0, 0, 0, 0, time.Local),
		},
		Total: &QueryResult{Service: "Total", Monthly: 1000.07, Yesterday: 400.0},
		Services: []*QueryResult{
			{Service: "Cloud SQL", Monthly: 1000.0, Yesterday: 400.0},
			{Service: "BigQuery", Monthly: 0.07, Yesterday: 0.0},
		},
	}
	actualBillings, err := NewBillings(&inputReportingPeriod, inputQueryResults)

	assert.Nil(t, err)
	assert.EqualValues(t, expectedBillings, actualBillings)
}

func TestCreateBillingsFromEmptyQueryResults(t *testing.T) {
	inputQueryResults := []*QueryResult{}
	inputReportingPeriod := ReportingPeriod{
		From: time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
		To:   time.Date(2021, 5, 8, 0, 0, 0, 0, time.Local),
	}

	expectedBillings := &Billings{
		AggregationPeriod: AggregationPeriod{
			From: time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
			To:   time.Date(2021, 5, 8, 0, 0, 0, 0, time.Local),
		},
		Total:    &QueryResult{Service: "Total", Monthly: 0.00, Yesterday: 0.0},
		Services: []*QueryResult{},
	}
	actualBillings, err := NewBillings(&inputReportingPeriod, inputQueryResults)

	assert.Nil(t, err)
	assert.EqualValues(t, expectedBillings, actualBillings)
}

func TestCreateBillingsFromSingleElementQueryResult(t *testing.T) {
	inputQueryResults := []*QueryResult{
		{Service: "Total", Monthly: 0.07, Yesterday: 0.0},
	}
	inputReportingPeriod := ReportingPeriod{
		From: time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
		To:   time.Date(2021, 5, 8, 0, 0, 0, 0, time.Local),
	}

	expectedBillings := &Billings{
		AggregationPeriod: AggregationPeriod{
			From: time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
			To:   time.Date(2021, 5, 8, 0, 0, 0, 0, time.Local),
		},
		Total:    &QueryResult{Service: "Total", Monthly: 0.07, Yesterday: 0.0},
		Services: []*QueryResult{},
	}
	actualBillings, err := NewBillings(&inputReportingPeriod, inputQueryResults)

	assert.Nil(t, err)
	assert.EqualValues(t, expectedBillings, actualBillings)
}

func TestBillingNotCreatedFromUnsortedQueryResults(t *testing.T) {
	inputQueryResults := []*QueryResult{
		{Service: "Cloud SQL", Monthly: 1000.0, Yesterday: 400.0},
		{Service: "BigQuery", Monthly: 0.07, Yesterday: 0.0},
		{Service: "Total", Monthly: 1000.07, Yesterday: 400.0},
	}
	inputReportingPeriod := ReportingPeriod{
		From: time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
		To:   time.Date(2021, 5, 8, 0, 0, 0, 0, time.Local),
	}

	actualBillings, err := NewBillings(&inputReportingPeriod, inputQueryResults)

	assert.NotNil(t, err)
	assert.Nil(t, actualBillings)
	assert.EqualValues(t, "Unexpected query results! The results might not be correctly sorted!", err.Error())
}

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
		Total: &QueryResult{Service: "Total", Monthly: 1000.07, Yesterday: 400.0},
		Services: []*QueryResult{
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
		Total:    &QueryResult{Service: "Total", Monthly: 1000.07, Yesterday: 400.0},
		Services: []*QueryResult{},
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
		Total: &QueryResult{Service: "Total", Monthly: 1000.07, Yesterday: 400.0},
		Services: []*QueryResult{
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

	actualNotification := inputBillings.AsNotification()
	assert.EqualValues(t, expectedNotification, actualNotification)
}

func TestCreateNotificationFromBillingsWithNoServiceCosts(t *testing.T) {
	inputBillings := &Billings{
		AggregationPeriod: AggregationPeriod{
			From: time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
			To:   time.Date(2021, 5, 8, 0, 0, 0, 0, time.Local),
		},
		Total:    &QueryResult{Service: "Total", Monthly: 0.00, Yesterday: 0.00},
		Services: []*QueryResult{},
	}

	expectedNotification :=
		`＜5/1 ~ 5/8 の GCP 利用料金＞ ※ () 内は前日分

Total: ¥ 0 (¥ 0)`

	actualNotification := inputBillings.AsNotification()
	assert.EqualValues(t, expectedNotification, actualNotification)
}

func TestCreateNotificationString(t *testing.T) {
	inputCostSummary := []*QueryResult{
		{Service: "Total", Monthly: 1000.07, Yesterday: 400.0},
		{Service: "Cloud SQL", Monthly: 1000.0, Yesterday: 400.0},
		{Service: "BigQuery", Monthly: 0.07, Yesterday: 0.0},
	}
	executionTimestamp := time.Date(2021, 5, 8, 8, 30, 0, 0, time.Local)

	expectedOutput :=
		`＜5/1 ~ 5/7 の GCP 利用料金＞ ※ () 内は前日分

Total: ¥ 1,000.07 (¥ 400)

----- 内訳 -----
Cloud SQL: ¥ 1,000 (¥ 400)
BigQuery: ¥ 0.07 (¥ 0)`

	actualOutput := createNotificationString(inputCostSummary, executionTimestamp)
	assert.EqualValues(t, expectedOutput, actualOutput)
}

func TestCreateNotificationOnFirstDayOfMonth(t *testing.T) {
	inputCostSummary := []*QueryResult{
		{Service: "Total", Monthly: 1000.07, Yesterday: 400.0},
		{Service: "Cloud SQL", Monthly: 1000.0, Yesterday: 400.0},
		{Service: "BigQuery", Monthly: 0.07, Yesterday: 0.0},
	}
	AsiaTokyo, _ := time.LoadLocation("Asia/Tokyo")
	executionTimestamp := time.Date(2021, 5, 1, 8, 30, 0, 0, AsiaTokyo)

	expectedFirstLineOfOutput := "＜4/1 ~ 4/30 の GCP 利用料金＞ ※ () 内は前日分"

	actualOutput := createNotificationString(inputCostSummary, executionTimestamp)
	actualFirstLineOfOutput := strings.Split(actualOutput, "\n")[0]
	assert.EqualValues(t, expectedFirstLineOfOutput, actualFirstLineOfOutput)
}

func TestCreateNotificationInJST(t *testing.T) {
	inputCostSummary := []*QueryResult{
		{Service: "Total", Monthly: 1000.07, Yesterday: 400.0},
		{Service: "Cloud SQL", Monthly: 1000.0, Yesterday: 400.0},
		{Service: "BigQuery", Monthly: 0.07, Yesterday: 0.0},
	}
	// 2021-05-08 in JST
	executionTimestamp := time.Date(2021, 5, 7, 23, 00, 0, 0, time.UTC)

	expectedFirstLineOfOutput := "＜5/1 ~ 5/7 の GCP 利用料金＞ ※ () 内は前日分"

	actualOutput := createNotificationString(inputCostSummary, executionTimestamp)
	actualFirstLineOfOutput := strings.Split(actualOutput, "\n")[0]
	assert.EqualValues(t, expectedFirstLineOfOutput, actualFirstLineOfOutput)
}
func TestSlackPost(t *testing.T) {
	inputURL := os.Getenv("SLACK_WEBHOOK_URL")
	inputMessage := "test\nこれはテスト投稿です。"

	err := sendMessageToSlack(inputURL, inputMessage)
	assert.Nil(t, err)
}

func TestBuildReportingPeriodCorrectly(t *testing.T) {
	AsiaTokyo, _ := time.LoadLocation("Asia/Tokyo")
	inputDateTime := time.Date(2021, 5, 8, 8, 30, 0, 0, AsiaTokyo)

	expectedReportingPeriod := ReportingPeriod{
		TimeZone: "Asia/Tokyo",
		From:     time.Date(2021, 5, 1, 0, 0, 0, 0, AsiaTokyo),
		To:       time.Date(2021, 5, 7, 0, 0, 0, 0, AsiaTokyo),
	}
	actualReportingPeriod, err := NewReportingPeriod(inputDateTime, "Asia/Tokyo")

	assert.EqualValues(t, expectedReportingPeriod, actualReportingPeriod)
	assert.Nil(t, err)
}

func TestBuildReportingPeriodOnFirstDayOfMonthCorrectly(t *testing.T) {
	AsiaTokyo, _ := time.LoadLocation("Asia/Tokyo")
	inputDateTime := time.Date(2021, 5, 1, 8, 30, 0, 0, AsiaTokyo)

	expectedReportingPeriod := ReportingPeriod{
		TimeZone: "Asia/Tokyo",
		From:     time.Date(2021, 4, 1, 0, 0, 0, 0, AsiaTokyo),
		To:       time.Date(2021, 4, 30, 0, 0, 0, 0, AsiaTokyo),
	}
	actualReportingPeriod, err := NewReportingPeriod(inputDateTime, "Asia/Tokyo")

	assert.EqualValues(t, expectedReportingPeriod, actualReportingPeriod)
	assert.Nil(t, err)
}

func TestBuildReportingPeriodFromJSTToUTCCorrectly(t *testing.T) {
	// 2021-05-08 in JST
	inputDateTime := time.Date(2021, 5, 7, 23, 00, 0, 0, time.UTC)

	AsiaTokyo, _ := time.LoadLocation("Asia/Tokyo")
	expectedReportingPeriod := ReportingPeriod{
		TimeZone: "Asia/Tokyo",
		From:     time.Date(2021, 5, 1, 0, 0, 0, 0, AsiaTokyo),
		To:       time.Date(2021, 5, 7, 0, 0, 0, 0, AsiaTokyo),
	}
	actualReportingPeriod, err := NewReportingPeriod(inputDateTime, "Asia/Tokyo")

	assert.EqualValues(t, expectedReportingPeriod, actualReportingPeriod)
	assert.Nil(t, err)
}

func TestBuildReportingPeriodFromUTCToJSTCorrectly(t *testing.T) {
	// 2021-05-06 in UTC
	AsiaTokyo, _ := time.LoadLocation("Asia/Tokyo")
	inputDateTime := time.Date(2021, 5, 7, 8, 30, 0, 0, AsiaTokyo)

	utc := time.UTC
	expectedReportingPeriod := ReportingPeriod{
		TimeZone: "UTC",
		From:     time.Date(2021, 5, 1, 0, 0, 0, 0, utc),
		To:       time.Date(2021, 5, 5, 0, 0, 0, 0, utc),
	}
	actualReportingPeriod, err := NewReportingPeriod(inputDateTime, "UTC")

	assert.EqualValues(t, expectedReportingPeriod, actualReportingPeriod)
	assert.Nil(t, err)
}
