// datetime package implements a timezone converter and
// an object to contain reporting period.
package datetime

import "time"

// ReportingPeriod contains the date period
// to aggregate and report the GCP cost.
type ReportingPeriod struct {
	TimeZone string
	From     time.Time
	To       time.Time
}

// NewReportingPeriod constructs the date period to report the GCP cost
// from the reporting datetime.
//
// The period is from the first day of the month to the one day before
// the reporting date. (e.g. 2021/8/30 -> 2021/8/1 ~ 2021/8/29)
//
// If the reporting date is the first day of a month,
// the period starts from the first day of the previous month.
// (e.g. 2021/8/1 -> 2021/7/1 ~ 2021/7/31)
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
