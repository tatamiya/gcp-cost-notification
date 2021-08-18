package billing

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/tatamiya/gcp-cost-notification/src/datetime"
	"github.com/tatamiya/gcp-cost-notification/src/db"
	"github.com/tatamiya/gcp-cost-notification/src/utils"
)

func newResultValidationError(message string, err error) *utils.CustomError {
	return &utils.CustomError{
		Process: "Query Results Validation",
		Message: message,
		Err:     err,
	}
}

// The date period to aggregate the GCP cost.
type BillingPeriod struct {
	From time.Time
	To   time.Time
}

// Display the period in the "MM/DD ~ MM/DD" format.
func (a *BillingPeriod) String() string {
	return fmt.Sprintf("%d/%d ~ %d/%d", a.From.Month(), a.From.Day(), a.To.Month(), a.To.Day())
}

// Invoice contains the data of the cost aggregation period,
// the total cost, and costs for each service.
type Invoice struct {
	BillingPeriod BillingPeriod
	Total         *db.QueryResult
	Services      []*db.QueryResult
}

// NewInvoice constructs a new Invoice from cost reporting period and BigQuery Results.
//
// The first element of the BQ results should be the total cost.
// If it is not, an error is returned.
//
// If the BQ result is empty, the new Invoice has 0 total cost and empty service costs.
func NewInvoice(period *datetime.ReportingPeriod, queryResults []*db.QueryResult) (*Invoice, *utils.CustomError) {

	billingPeriod := BillingPeriod{
		From: period.From,
		To:   period.To,
	}

	var totalCost *db.QueryResult
	var serviceCosts []*db.QueryResult

	if len(queryResults) == 0 {
		totalCost = &db.QueryResult{Service: "Total", Monthly: 0.00, Yesterday: 0.00}
		serviceCosts = []*db.QueryResult{}
	} else {
		totalCost = queryResults[0]
		if totalCost.Service != "Total" {
			log.Printf("Unexpected query results: %v", queryResults)
			return nil, newResultValidationError(
				"Unexpected query results! The results might not be correctly sorted!",
				fmt.Errorf("First element of the query results was %s, not Total", totalCost.Service),
			)
		}
		serviceCosts = queryResults[1:]
	}

	return &Invoice{
		BillingPeriod: billingPeriod,
		Total:         totalCost,
		Services:      serviceCosts,
	}, nil

}

func (b *Invoice) details() string {
	serviceCosts := b.Services
	var listOfLines []string
	for _, cost := range serviceCosts {
		listOfLines = append(listOfLines, cost.AsMessageLine())
	}
	return strings.Join(listOfLines, "\n")
}

// AsMessage creates a notification message of GCP costs from Invoice.
func (b *Invoice) AsMessage() string {

	message := fmt.Sprintf("＜%s の GCP 利用料金＞ ※ () 内は前日分\n\n", &b.BillingPeriod)
	message += b.Total.AsMessageLine()

	if len(b.Services) > 0 {
		message += "\n\n" + "----- 内訳 -----" + "\n"
		message += b.details()
	}

	return message
}
