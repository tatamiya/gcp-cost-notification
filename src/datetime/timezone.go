package datetime

import (
	"log"
	"os"
	"time"
)

// TimeZoneConverter is an object to convert the timezone of an Time object
// to that of the designated location.
type TimeZoneConverter struct {
	location *time.Location
}

// NewTimeZoneConverter constructs a TimeZoneConverter to
// convert the timezone into the location designated in the
// global enviroment variable.
func NewTimeZoneConverter() TimeZoneConverter {
	timeZone := os.Getenv("TIMEZONE")
	location, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Printf("Failed in setting timezone from environment variable '%s': %s", timeZone, err.Error())
		location = time.Local
		log.Printf("Runtime local timezone '%s' is set instead.", location)
	}
	return TimeZoneConverter{
		location: location,
	}
}

// From method converts the timezone of an input Time.
func (t *TimeZoneConverter) From(datetime time.Time) time.Time {
	return datetime.In(t.location)
}
