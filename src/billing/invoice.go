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

func NewResultValidationError(message string, err error) *utils.CustomError {
	return &utils.CustomError{
		Process: "Query Results Validation",
		Message: message,
		Err:     err,
	}
}

type BillingPeriod struct {
	From time.Time
	To   time.Time
}

func (a *BillingPeriod) String() string {
	return fmt.Sprintf("%d/%d ~ %d/%d", a.From.Month(), a.From.Day(), a.To.Month(), a.To.Day())
}

type Invoice struct {
	BillingPeriod BillingPeriod
	Total         *db.QueryResult
	Services      []*db.QueryResult
}

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
			return nil, NewResultValidationError(
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

func (b *Invoice) AsMessage() string {

	message := fmt.Sprintf("＜%s の GCP 利用料金＞ ※ () 内は前日分\n\n", &b.BillingPeriod)
	message += b.Total.AsMessageLine()

	if len(b.Services) > 0 {
		message += "\n\n" + "----- 内訳 -----" + "\n"
		message += b.details()
	}

	return message
}
