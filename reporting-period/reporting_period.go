package reportingperiod

import "time"

type ReportingPeriod struct {
	TimeZone string
	From     time.Time
	To       time.Time
}

func NewReportingPeriod(reportingDateTime time.Time, timeZone string) (ReportingPeriod, error) {
	tz := timeZone
	location, err := time.LoadLocation(tz)
	if err != nil {
		location = reportingDateTime.Location()
		tz = location.String()
	}
	localizedDateTime := reportingDateTime.In(location)
	oneDayBefore := localizedDateTime.AddDate(0, 0, -1)

	year := oneDayBefore.Year()
	month := oneDayBefore.Month()
	day := oneDayBefore.Day()
	return ReportingPeriod{
		TimeZone: tz,
		From:     time.Date(year, month, 1, 0, 0, 0, 0, location),
		To:       time.Date(year, month, day, 0, 0, 0, 0, location),
	}, err
}
