package datetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConvertFromUTC2JST(t *testing.T) {
	AsiaTokyo, _ := time.LoadLocation("Asia/Tokyo")

	inputDateTime := time.Date(2021, 5, 7, 23, 30, 0, 0, time.UTC)

	testConverter := TimeZoneConverter{
		location: AsiaTokyo,
	}

	expectedDateTime := time.Date(2021, 5, 8, 8, 30, 0, 0, AsiaTokyo)
	actualDateTime := testConverter.From(inputDateTime)

	assert.EqualValues(t, expectedDateTime, actualDateTime)

}

func TestConvertFromJST2UTC(t *testing.T) {
	AsiaTokyo, _ := time.LoadLocation("Asia/Tokyo")
	inputDateTime := time.Date(2021, 5, 8, 8, 30, 0, 0, AsiaTokyo)

	testConverter := TimeZoneConverter{
		location: time.UTC,
	}

	expectedDateTime := time.Date(2021, 5, 7, 23, 30, 0, 0, time.UTC)
	actualDateTime := testConverter.From(inputDateTime)

	assert.EqualValues(t, expectedDateTime, actualDateTime)
}
