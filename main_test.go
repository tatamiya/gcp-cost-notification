package gcp_cost_notification

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tatamiya/gcp-cost-notification/db"
	"github.com/tatamiya/gcp-cost-notification/notification"
)

var InputReportingDateTime time.Time = time.Date(2021, 8, 7, 20, 15, 0, 0, time.Local)
var InputQueryResults []*db.QueryResult = []*db.QueryResult{
	{Service: "Total", Monthly: 1000.07, Yesterday: 400.0},
	{Service: "Cloud SQL", Monthly: 1000.0, Yesterday: 400.0},
	{Service: "BigQuery", Monthly: 0.07, Yesterday: 0.0},
}

func TestRunWholeProcessCorrectly(t *testing.T) {

	reportingDateTime := InputReportingDateTime
	inputQueryResults := InputQueryResults
	BQClientStub := db.NewBQClientStub(inputQueryResults, nil)
	SlackClientStub := notification.NewSlackClientStub(nil)

	expectedMessage :=
		`＜8/1 ~ 8/6 の GCP 利用料金＞ ※ () 内は前日分

Total: ¥ 1,000.07 (¥ 400)

----- 内訳 -----
Cloud SQL: ¥ 1,000 (¥ 400)
BigQuery: ¥ 0.07 (¥ 0)`

	actualMessage, err := mainProcess(reportingDateTime, &BQClientStub, &SlackClientStub)

	assert.Nil(t, err)
	assert.EqualValues(t, expectedMessage, actualMessage)
}

func TestDisplayPreviousMonthConstsOnFirstDayOfMonth(t *testing.T) {

	reportingDateTime := time.Date(2021, 8, 1, 0, 0, 0, 0, time.Local)

	BQClientStub := db.NewBQClientStub(InputQueryResults, nil)
	SlackClientStub := notification.NewSlackClientStub(nil)

	actualMessage, err := mainProcess(reportingDateTime, &BQClientStub, &SlackClientStub)

	assert.Nil(t, err)
	assert.True(t, strings.Contains(actualMessage, "＜7/1 ~ 7/31 の GCP 利用料金＞"), actualMessage)
}

func TestNotDisplayServiceCostsWhenQueryResultHasNoServiceCosts(t *testing.T) {

	inputQueryResults := []*db.QueryResult{
		{Service: "Total", Monthly: 1000.07, Yesterday: 400.0},
	}
	BQClientStub := db.NewBQClientStub(inputQueryResults, nil)
	SlackClientStub := notification.NewSlackClientStub(nil)

	expectedMessage :=
		"＜8/1 ~ 8/6 の GCP 利用料金＞ ※ () 内は前日分\n\nTotal: ¥ 1,000.07 (¥ 400)"

	actualMessage, err := mainProcess(InputReportingDateTime, &BQClientStub, &SlackClientStub)

	assert.Nil(t, err)
	assert.EqualValues(t, expectedMessage, actualMessage)
}

func TestDisplayZeroTotalCostWhenQueryResultIsEmpty(t *testing.T) {

	inputQueryResults := []*db.QueryResult{}
	BQClientStub := db.NewBQClientStub(inputQueryResults, nil)
	SlackClientStub := notification.NewSlackClientStub(nil)

	expectedMessage :=
		"＜8/1 ~ 8/6 の GCP 利用料金＞ ※ () 内は前日分\n\nTotal: ¥ 0 (¥ 0)"

	actualMessage, err := mainProcess(InputReportingDateTime, &BQClientStub, &SlackClientStub)

	assert.Nil(t, err)
	assert.EqualValues(t, expectedMessage, actualMessage)
}

func TestReturnErrorWhenQueryResultIsUnexpectedlySorted(t *testing.T) {

	inputQueryResults := []*db.QueryResult{
		{Service: "Cloud SQL", Monthly: 1000.0, Yesterday: 400.0},
		{Service: "BigQuery", Monthly: 0.07, Yesterday: 0.0},
		{Service: "Total", Monthly: 1000.07, Yesterday: 400.0},
	}
	BQClientStub := db.NewBQClientStub(inputQueryResults, nil)
	SlackClientStub := notification.NewSlackClientStub(nil)

	actualMessage, err := mainProcess(InputReportingDateTime, &BQClientStub, &SlackClientStub)

	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "Error in Query Results Validation."), err)
	assert.EqualValues(t, "", actualMessage)
}

func TestReturnErrorWhenBQFailed(t *testing.T) {

	BQClientStub := db.NewBQClientStub([]*db.QueryResult{}, fmt.Errorf("Something Happened!"))
	SlackClientStub := notification.NewSlackClientStub(nil)

	actualMessage, err := mainProcess(InputReportingDateTime, &BQClientStub, &SlackClientStub)

	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "Error in Query Execution."), err)
	assert.EqualValues(t, "", actualMessage)
}

func TestReturnErrorWhenSlackNotificationFailed(t *testing.T) {

	BQClientStub := db.NewBQClientStub(InputQueryResults, nil)
	SlackClientStub := notification.NewSlackClientStub(fmt.Errorf("Something Happened!"))

	actualMessage, err := mainProcess(InputReportingDateTime, &BQClientStub, &SlackClientStub)

	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "Error in Slack Notification."), err)
	assert.EqualValues(t, "", actualMessage)
}
