package datetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

func TestNewReportingPeriodReturnErrorForInvalidTimeZone(t *testing.T) {
	inputDateTime := time.Date(2021, 5, 7, 8, 30, 0, 0, time.Local)

	expectedReportingPeriod := ReportingPeriod{
		TimeZone: time.Local.String(),
		From:     time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
		To:       time.Date(2021, 5, 6, 0, 0, 0, 0, time.Local),
	}

	actualReportingPeriod, err := NewReportingPeriod(inputDateTime, "Invalid/TimeZone")

	assert.EqualValues(t, expectedReportingPeriod, actualReportingPeriod)

	assert.NotNil(t, err)
	assert.EqualValues(t, "unknown time zone Invalid/TimeZone", err.Error())
}
