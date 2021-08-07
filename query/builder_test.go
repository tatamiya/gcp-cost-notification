package query

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	reportingperiod "github.com/tatamiya/gcp-cost-notification/reporting-period"
)

func TestRenderQueryFromTemplateCorrectly(t *testing.T) {
	inputTableID := "sample_project.sample_dataset.sample_table"

	builder := QueryBuilder{
		tableID:      inputTableID,
		templatePath: "./template.sql",
	}

	AsiaTokyo, _ := time.LoadLocation("Asia/Tokyo")
	inputReportingPeriod := reportingperiod.ReportingPeriod{
		TimeZone: "Asia/Tokyo",
		From:     time.Date(2021, 5, 1, 0, 0, 0, 0, AsiaTokyo),
		To:       time.Date(2021, 5, 7, 0, 0, 0, 0, AsiaTokyo),
	}
	outputQuery := builder.Build(inputReportingPeriod)

	assert.EqualValues(t, true, strings.Contains(outputQuery, "SELECT"))
	assert.EqualValues(t, true, strings.Contains(outputQuery, "Asia/Tokyo"))
	assert.EqualValues(t, true, strings.Contains(outputQuery, "2021-05-07T00:00:00+09:00"))
	assert.EqualValues(t, true, strings.Contains(outputQuery, "2021-05-07T00:00:00+09:00"))
	assert.EqualValues(t, true, strings.Contains(outputQuery, inputTableID))
}
