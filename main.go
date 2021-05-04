package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

type QueryParameters struct {
	TableName string
}

type QueryResult struct {
	Service   string
	Monthly   float32
	Yesterday float32
}

func buildQuery(tableName string) string {

	params := QueryParameters{tableName}
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
