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

type SlackClientInterface interface {
	Send(messenger notification.Messenger) *utils.CustomError
}

func mainProcess(
	reportingDateTime time.Time,
	BQClient db.BQClientInterface,
	slackClient SlackClientInterface,
) (string, error) {

	reportingPeriod := datetime.NewReportingPeriod(reportingDateTime)

	queryBuilder := query.NewQueryBuilder()
	query := queryBuilder.Build(reportingPeriod)

	queryResult, err := BQClient.SendQuery(query)
	if err != nil {
		log.Print(err)
		slackError := slackClient.Send(err)
		if slackError != nil {
			log.Println("Error notification to Slack also failed!: ", slackError.Error())
		}
		return "", err
	}

	billings, err := message.NewBillings(&reportingPeriod, queryResult)
	if err != nil {
		log.Print(err)
		slackError := slackClient.Send(err)
		if slackError != nil {
			log.Println("Error notification to Slack also failed!: ", slackError.Error())
		}
		return "", err
	}

	err = slackClient.Send(billings)
	if err != nil {
		log.Print(err)
		return "", err
	}

	return billings.AsMessage(), nil
}
