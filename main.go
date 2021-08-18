package gcp_cost_notification

import (
	"context"
	"log"
	"time"

	"github.com/tatamiya/gcp-cost-notification/src/billing"
	"github.com/tatamiya/gcp-cost-notification/src/datetime"
	"github.com/tatamiya/gcp-cost-notification/src/db"
	"github.com/tatamiya/gcp-cost-notification/src/notification"
	"github.com/tatamiya/gcp-cost-notification/src/query"
	"github.com/tatamiya/gcp-cost-notification/src/utils"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

// The function called by Cloud Functions.
// Get GCP cost from BigQuery and
// send notification to Slack.
//
// The period of cost retrieval is
// from the first date of the month
// upto one day before the execution date.
// If the execution date is the first date of the month,
// the period is the previous month.
func CostNotifier(ctx context.Context, m PubSubMessage) error {
	tzConverter := datetime.NewTimeZoneConverter()
	currentDateTime := tzConverter.From(time.Now())

	BQClient := db.NewBQClient()

	slackClient := notification.NewSlackClient()

	message, err := mainProcess(currentDateTime, &BQClient, &slackClient)
	if err == nil {
		log.Println("Message was successfully sent to Slack!: ", message)
	} else {
		log.Println("Failed in sending message!: ", err.Error())
	}
	return err
}

type bqClientInterface interface {
	SendQuery(query string) ([]*db.QueryResult, *utils.CustomError)
}

type slackClientInterface interface {
	Send(messenger notification.Messenger) (string, *utils.CustomError)
}

func mainProcess(
	reportingDateTime time.Time,
	BQClient bqClientInterface,
	slackClient slackClientInterface,
) (string, error) {

	reportingPeriod := datetime.NewReportingPeriod(reportingDateTime)

	queryBuilder := query.NewQueryBuilder()
	query := queryBuilder.Build(reportingPeriod)

	queryResult, err := BQClient.SendQuery(query)
	if err != nil {
		log.Print(err)
		_, slackError := slackClient.Send(err)
		if slackError != nil {
			log.Println("Error notification to Slack also failed!: ", slackError.Error())
		}
		return "", err
	}

	invoice, err := billing.NewInvoice(&reportingPeriod, queryResult)
	if err != nil {
		log.Print(err)
		_, slackError := slackClient.Send(err)
		if slackError != nil {
			log.Println("Error notification to Slack also failed!: ", slackError.Error())
		}
		return "", err
	}

	sentMessage, err := slackClient.Send(invoice)
	if err != nil {
		log.Print(err)
		return "", err
	}

	return sentMessage, nil
}
