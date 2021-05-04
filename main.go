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

func buildQuery(tableName string, executionTimestamp string) string {

	noEscapeTimestamp := template.HTML(executionTimestamp)
	params := QueryParameters{tableName, noEscapeTimestamp}
	var buf bytes.Buffer
	t := template.Must(template.ParseFiles("query.sql"))
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

func createSingleReportLine(cost *QueryResult) string {
	service := cost.Service
	monthly := humanize.CommafWithDigits(float64(cost.Monthly), 2)
	yesterday := humanize.CommafWithDigits(float64(cost.Yesterday), 2)

	return fmt.Sprintf("\n%s: ¥ %s (¥ %s)", service, monthly, yesterday)
}

func createNotificationString(costSummary []*QueryResult, executionTimestamp time.Time) string {

	oneDayBefore := executionTimestamp.AddDate(0, 0, -1)
	month := oneDayBefore.Month()
	day := oneDayBefore.Day()

	output := fmt.Sprintf("＜%d/1 ~ %d/%d の GCP 利用料金＞ ※ () 内は前日分\n", month, month, day)

	firstLine := costSummary[0]
	if firstLine.Service != "Total" {
		return "Something Wrong!"
	}
	output += createSingleReportLine(firstLine)
	if len(costSummary) < 1 {
		return output
	}

	output += "\n\n----- 明細 -----"

	for _, detail := range costSummary[1:] {
		output += createSingleReportLine(detail)
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
	currentTimestamp := time.Now()

	projectID := os.Getenv("GCP_PROJECT")
	datasetName := os.Getenv("DATASET_NAME")
	tableName := os.Getenv("TABLE_NAME")

	fullTableName := fmt.Sprintf("%s.%s.%s", projectID, datasetName, tableName)

	timestampString := currentTimestamp.Format(time.RFC3339)
	query := buildQuery(fullTableName, timestampString)
	costSummary, err := sendQueryToBQ(query, projectID)
	if err != nil {
		log.Print(err)
		return err
	}

	messageString := createNotificationString(costSummary, currentTimestamp)
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	err = sendMessageToSlack(webhookURL, messageString)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
