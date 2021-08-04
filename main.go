package gcp_cost_notification

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	humanize "github.com/dustin/go-humanize"
	"github.com/slack-go/slack"
	"google.golang.org/api/iterator"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

type QueryParameters struct {
	TableName          string
	ExecutionTimestamp template.HTML
}

type QueryResult struct {
	Service   string
	Monthly   float32
	Yesterday float32
}

func (r *QueryResult) asMessageLine() string {
	service := r.Service
	monthly := humanize.CommafWithDigits(float64(r.Monthly), 2)
	yesterday := humanize.CommafWithDigits(float64(r.Yesterday), 2)

	return fmt.Sprintf("\n%s: ¥ %s (¥ %s)", service, monthly, yesterday)
}

type ReportingPeriod struct {
	From time.Time
	To   time.Time
}

type AggregationPeriod struct {
	From time.Time
	To   time.Time
}

type Billings struct {
	AggregationPeriod AggregationPeriod
	Total             *QueryResult
	Services          []*QueryResult
}

func NewBillings(period *ReportingPeriod, queryResults []*QueryResult) (*Billings, error) {

	aggregationPeriod := AggregationPeriod{
		From: period.From,
		To:   period.To,
	}

	var totalCost *QueryResult
	var serviceCosts []*QueryResult

	if len(queryResults) == 0 {
		totalCost = &QueryResult{Service: "Total", Monthly: 0.00, Yesterday: 0.00}
		serviceCosts = []*QueryResult{}
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

func (b *Billings) headline() string {
	from := b.AggregationPeriod.From
	to := b.AggregationPeriod.To

	monthFrom := from.Month()
	dayFrom := from.Day()

	monthTo := to.Month()
	dayTo := to.Day()

	return fmt.Sprintf("＜%d/%d ~ %d/%d の GCP 利用料金＞ ※ () 内は前日分", monthFrom, dayFrom, monthTo, dayTo)
}

func (b *Billings) detailLines() string {
	serviceCosts := b.Services
	if len(serviceCosts) == 0 {
		return ""
	}
	output := "----- 内訳 -----"
	for _, cost := range serviceCosts {
		output += cost.asMessageLine()
	}
	return output
}

func buildQuery(tableName string, executionTimestamp string) string {

	fileDir := os.Getenv("FILE_DIRECTORY")

	noEscapeTimestamp := template.HTML(executionTimestamp)
	params := QueryParameters{tableName, noEscapeTimestamp}
	var buf bytes.Buffer
	t := template.Must(template.ParseFiles("./" + fileDir + "query.sql"))
	t.Execute(&buf, params)

	return buf.String()
}

func sendQueryToBQ(query string, projectID string) ([]*QueryResult, error) {
	var queryResults []*QueryResult

	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return queryResults, fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	q := client.Query(query)
	it, err := q.Read(ctx)
	if err != nil {
		return queryResults, fmt.Errorf("client.Query: %v", err)
	}

	for {
		var result QueryResult
		err := it.Next(&result)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []*QueryResult{}, err
		}
		queryResults = append(queryResults, &result)
	}

	return queryResults, nil
}

func createNotificationString(costSummary []*QueryResult, executionTime time.Time) string {

	location, _ := time.LoadLocation("Asia/Tokyo")
	localizedTime := executionTime.In(location)
	oneDayBefore := localizedTime.AddDate(0, 0, -1)
	month := oneDayBefore.Month()
	day := oneDayBefore.Day()

	output := fmt.Sprintf("＜%d/1 ~ %d/%d の GCP 利用料金＞ ※ () 内は前日分\n", month, month, day)

	firstLine := costSummary[0]
	if firstLine.Service != "Total" {
		return "Something Wrong!"
	}
	output += firstLine.asMessageLine()
	if len(costSummary) < 1 {
		return output
	}

	output += "\n\n----- 内訳 -----"

	for _, detail := range costSummary[1:] {
		output += detail.asMessageLine()
	}
	return output
}

func sendMessageToSlack(webhookURL string, messageText string) error {
	msg := slack.WebhookMessage{
		Text: messageText,
	}
	err := slack.PostWebhook(webhookURL, &msg)
	return err
}

func CostNotifier(ctx context.Context, m PubSubMessage) error {
	currentTime := time.Now()

	projectID := os.Getenv("GCP_PROJECT")
	datasetName := os.Getenv("DATASET_NAME")
	tableName := os.Getenv("TABLE_NAME")

	fullTableName := fmt.Sprintf("%s.%s.%s", projectID, datasetName, tableName)

	timestampString := currentTime.Format(time.RFC3339)
	query := buildQuery(fullTableName, timestampString)
	costSummary, err := sendQueryToBQ(query, projectID)
	if err != nil {
		log.Print(err)
		return err
	}

	messageString := createNotificationString(costSummary, currentTime)
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	err = sendMessageToSlack(webhookURL, messageString)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
