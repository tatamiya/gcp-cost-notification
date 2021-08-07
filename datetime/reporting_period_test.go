package datetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBuildReportingPeriodCorrectly(t *testing.T) {
	inputDateTime := time.Date(2021, 5, 8, 8, 30, 0, 0, time.Local)

	expectedReportingPeriod := ReportingPeriod{
		TimeZone: time.Local.String(),
		From:     time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
		To:       time.Date(2021, 5, 7, 0, 0, 0, 0, time.Local),
	}
	actualReportingPeriod := NewReportingPeriod(inputDateTime)

	assert.EqualValues(t, expectedReportingPeriod, actualReportingPeriod)
}

func TestBuildReportingPeriodOnFirstDayOfMonthCorrectly(t *testing.T) {
	inputDateTime := time.Date(2021, 5, 1, 8, 30, 0, 0, time.Local)

	expectedReportingPeriod := ReportingPeriod{
		TimeZone: time.Local.String(),
		From:     time.Date(2021, 4, 1, 0, 0, 0, 0, time.Local),
		To:       time.Date(2021, 4, 30, 0, 0, 0, 0, time.Local),
	}
	actualReportingPeriod := NewReportingPeriod(inputDateTime)

	assert.EqualValues(t, expectedReportingPeriod, actualReportingPeriod)
}

func TestReportingPeriodPreservesTimezone(t *testing.T) {
	AsiaTokyo, _ := time.LoadLocation("Asia/Tokyo")
	inputDateTime := time.Date(2021, 5, 8, 8, 00, 0, 0, AsiaTokyo)

	expectedReportingPeriod := ReportingPeriod{
		TimeZone: "Asia/Tokyo",
		From:     time.Date(2021, 5, 1, 0, 0, 0, 0, AsiaTokyo),
		To:       time.Date(2021, 5, 7, 0, 0, 0, 0, AsiaTokyo),
	}
	actualReportingPeriod := NewReportingPeriod(inputDateTime)

	assert.EqualValues(t, expectedReportingPeriod, actualReportingPeriod)
}
