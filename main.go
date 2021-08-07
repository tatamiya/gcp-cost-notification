package gcp_cost_notification

import (
	"context"
	"log"
	"time"

	"github.com/tatamiya/gcp-cost-notification/db"
	"github.com/tatamiya/gcp-cost-notification/message"
	"github.com/tatamiya/gcp-cost-notification/notification"
	"github.com/tatamiya/gcp-cost-notification/query"
	reportingperiod "github.com/tatamiya/gcp-cost-notification/reporting-period"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

func CostNotifier(ctx context.Context, m PubSubMessage) error {
	timeZone := "Asia/Tokyo"
	currentTime := time.Now()
	reportingPeriod, err := reportingperiod.NewReportingPeriod(currentTime, timeZone)

	queryBuilder := query.NewQueryBuilder()
	query := queryBuilder.Build(reportingPeriod)

	BQClient := db.NewBQClient()
	costSummary, err := BQClient.SendQuery(query)
	if err != nil {
		log.Print(err)
		return err
	}

	billings, err := message.NewBillings(&reportingPeriod, costSummary)
	messageString := billings.AsNotification()

	slackClient := notification.NewSlackClient()
	err = slackClient.Send(messageString)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
