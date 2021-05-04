package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	humanize "github.com/dustin/go-humanize"
	"google.golang.org/api/iterator"
)

type QueryParameters struct {
	TableName          string
	ExecutionTimestamp string
}

type QueryResult struct {
	Service   string
	Monthly   float32
	Yesterday float32
}

func buildQuery(tableName string, executionTimestamp string) string {

	params := QueryParameters{tableName, executionTimestamp}
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

func createNotificationString(costSummary []*QueryResult) string {

	output := "＜5/1 ~ 7 の GCP 利用料金＞\n() 内は前日分"

	firstLine := costSummary[0]
	if firstLine.Service != "Total" {
		return "Something Wrong!"
	}
	output += createSingleReportLine(firstLine)
	if len(costSummary) < 1 {
		return output
	}

	output += "\n以下サービス別"

	for _, detail := range costSummary[1:] {
		output += createSingleReportLine(detail)
	}
	return output
}

func CostNotifier() {
	currentTimestamp := time.Now()

	projectID := os.Getenv("GCP_PROJECT")
	datasetName := os.Getenv("DATASET_NAME")
	tableName := os.Getenv("TABLE_NAME")

	fullTableName := fmt.Sprintf("%s.%s.%s", projectID, datasetName, tableName)

	timestampString := currentTimestamp.Format(time.RFC3339)
	query := buildQuery(fullTableName, timestampString)
	_, _ = sendQueryToBQ(query, projectID)
}
