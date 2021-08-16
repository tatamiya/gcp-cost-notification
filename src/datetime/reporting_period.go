package datetime

import "time"

type ReportingPeriod struct {
	TimeZone string
	From     time.Time
	To       time.Time
}

func NewReportingPeriod(reportingDateTime time.Time) ReportingPeriod {
	location := reportingDateTime.Location()
	oneDayBefore := reportingDateTime.AddDate(0, 0, -1)

	year := oneDayBefore.Year()
	month := oneDayBefore.Month()
	day := oneDayBefore.Day()
	return ReportingPeriod{
		TimeZone: location.String(),
		From:     time.Date(year, month, 1, 0, 0, 0, 0, location),
		To:       time.Date(year, month, day, 0, 0, 0, 0, location),
	}
}
