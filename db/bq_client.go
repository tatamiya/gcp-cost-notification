package db

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

type QueryResult struct {
	Service   string
	Monthly   float32
	Yesterday float32
}

type BQClientInterface interface {
	SendQuery(query string) ([]*QueryResult, error)
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

func (c *BQClient) SendQuery(query string) ([]*QueryResult, error) {
	var queryResults []*QueryResult

	q := c.client.Query(query)
	ctx := context.Background()
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
