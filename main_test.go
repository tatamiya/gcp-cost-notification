package gcp_cost_notification

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tatamiya/gcp-cost-notification/db"
	"github.com/tatamiya/gcp-cost-notification/notification"
)

func TestRunWholeProcessCorrectly(t *testing.T) {

	reportingDateTime := time.Date(2021, 8, 7, 20, 15, 0, 0, time.Local)
	inputQueryResults := []*db.QueryResult{
		{Service: "Total", Monthly: 1000.07, Yesterday: 400.0},
		{Service: "Cloud SQL", Monthly: 1000.0, Yesterday: 400.0},
		{Service: "BigQuery", Monthly: 0.07, Yesterday: 0.0},
	}
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
