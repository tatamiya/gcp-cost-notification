package message

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/tatamiya/gcp-cost-notification/datetime"
	"github.com/tatamiya/gcp-cost-notification/db"
)

type AggregationPeriod struct {
	From time.Time
	To   time.Time
}

func (a *AggregationPeriod) String() string {
	return fmt.Sprintf("%d/%d ~ %d/%d", a.From.Month(), a.From.Day(), a.To.Month(), a.To.Day())
}

type Billings struct {
	AggregationPeriod AggregationPeriod
	Total             *db.QueryResult
	Services          []*db.QueryResult
}

func NewBillings(period *datetime.ReportingPeriod, queryResults []*db.QueryResult) (*Billings, error) {

	aggregationPeriod := AggregationPeriod{
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
			// TODO: Display queryResults in error message.
			log.Println("Unexpected query results: ", queryResults)
			return nil, fmt.Errorf("Unexpected query results! The results might not be correctly sorted!")
		}
		serviceCosts = queryResults[1:]
	}

	return &Billings{
		AggregationPeriod: aggregationPeriod,
		Total:             totalCost,
		Services:          serviceCosts,
	}, nil

}

func (b *Billings) detailLines() string {
	serviceCosts := b.Services
	var listOfLines []string
	for _, cost := range serviceCosts {
		listOfLines = append(listOfLines, cost.AsMessageLine())
	}
	return strings.Join(listOfLines, "\n")
}

func (b *Billings) AsNotification() string {

	notification := fmt.Sprintf("＜%s の GCP 利用料金＞ ※ () 内は前日分\n\n", &b.AggregationPeriod)
	notification += b.Total.AsMessageLine()

	if len(b.Services) > 0 {
		notification += "\n\n" + "----- 内訳 -----" + "\n"
		notification += b.detailLines()
	}

	return notification
}
