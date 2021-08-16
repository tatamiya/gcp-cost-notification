package db

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/dustin/go-humanize"
	"github.com/tatamiya/gcp-cost-notification/src/utils"
	"google.golang.org/api/iterator"
)

func NewQueryError(message string, err error) *utils.CustomError {
	return &utils.CustomError{
		Process: "Query Execution",
		Message: message,
		Err:     err,
	}
}

type QueryResult struct {
	Service   string
	Monthly   float32
	Yesterday float32
}

func (r *QueryResult) String() string {
	return fmt.Sprintf("{Service: %s, Monthly: %f, Yesterday: %f}", r.Service, r.Monthly, r.Yesterday)
}

func (r *QueryResult) AsMessageLine() string {
	service := r.Service
	monthly := humanize.CommafWithDigits(float64(r.Monthly), 2)
	yesterday := humanize.CommafWithDigits(float64(r.Yesterday), 2)

	return fmt.Sprintf("%s: ¥ %s (¥ %s)", service, monthly, yesterday)
}

type BQClient struct {
	client *bigquery.Client
}

func NewBQClient() BQClient {
	projectID := os.Getenv("GCP_PROJECT")
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}

	return BQClient{client: client}
}

func (c *BQClient) SendQuery(query string) ([]*QueryResult, *utils.CustomError) {
	var queryResults []*QueryResult

	q := c.client.Query(query)
	ctx := context.Background()
	it, err := q.Read(ctx)
	if err != nil {
		return queryResults, NewQueryError("Failed in executing query", err)
	}

	for {
		var result QueryResult
		err := it.Next(&result)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []*QueryResult{}, NewQueryError("Failed in parsing query results", err)
		}
		queryResults = append(queryResults, &result)
	}

	return queryResults, nil
}
