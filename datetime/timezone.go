package datetime

import (
	"log"
	"os"
	"time"
)

type TimeZoneConverter struct {
	location *time.Location
}

func NewTimeZoneConverter() TimeZoneConverter {
	timeZone := os.Getenv("TIMEZONE")
	location, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Printf("Failed in setting timezone from environment variable '%s': %s", timeZone, err.Error())
		location = time.Local
		log.Printf("Runtime local timezone '%s' is set instead.", location.String())
	}
	return TimeZoneConverter{
		location: location,
	}
}

func (t *TimeZoneConverter) From(datetime time.Time) time.Time {
	return datetime.In(t.location)
}
