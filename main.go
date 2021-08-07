package gcp_cost_notification

import (
	"context"
	"log"
	"time"

	"github.com/tatamiya/gcp-cost-notification/datetime"
	"github.com/tatamiya/gcp-cost-notification/db"
	"github.com/tatamiya/gcp-cost-notification/message"
	"github.com/tatamiya/gcp-cost-notification/notification"
	"github.com/tatamiya/gcp-cost-notification/query"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

func CostNotifier(ctx context.Context, m PubSubMessage) error {
	tzConverter := datetime.NewTimeZoneConverter()
	currentDateTime := tzConverter.From(time.Now())

	BQClient := db.NewBQClient()

	slackClient := notification.NewSlackClient()

	return mainProcess(currentDateTime, &BQClient, &slackClient)
}

func mainProcess(
	reportingDateTime time.Time,
	BQClient db.BQClientInterface,
	slackClient notification.SlackClientInterface,
) error {

	reportingPeriod := datetime.NewReportingPeriod(reportingDateTime)

	queryBuilder := query.NewQueryBuilder()
	query := queryBuilder.Build(reportingPeriod)

	costSummary, err := BQClient.SendQuery(query)
	if err != nil {
		log.Print(err)
		return err
	}

	billings, err := message.NewBillings(&reportingPeriod, costSummary)
	messageString := billings.AsNotification()

	err = slackClient.Send(messageString)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
