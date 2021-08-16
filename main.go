package gcp_cost_notification

import (
	"context"
	"log"
	"time"

	"github.com/tatamiya/gcp-cost-notification/billing"
	"github.com/tatamiya/gcp-cost-notification/datetime"
	"github.com/tatamiya/gcp-cost-notification/db"
	"github.com/tatamiya/gcp-cost-notification/notification"
	"github.com/tatamiya/gcp-cost-notification/query"
	"github.com/tatamiya/gcp-cost-notification/utils"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

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

type BQClientInterface interface {
	SendQuery(query string) ([]*db.QueryResult, *utils.CustomError)
}

type SlackClientInterface interface {
	Send(messenger notification.Messenger) (string, *utils.CustomError)
}

func mainProcess(
	reportingDateTime time.Time,
	BQClient BQClientInterface,
	slackClient SlackClientInterface,
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

	billings, err := billing.NewBillings(&reportingPeriod, queryResult)
	if err != nil {
		log.Print(err)
		_, slackError := slackClient.Send(err)
		if slackError != nil {
			log.Println("Error notification to Slack also failed!: ", slackError.Error())
		}
		return "", err
	}

	sentMessage, err := slackClient.Send(billings)
	if err != nil {
		log.Print(err)
		return "", err
	}

	return sentMessage, nil
}
