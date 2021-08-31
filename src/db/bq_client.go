// db package implements objects to send a query to BigQuery.
package db

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/bigquery"
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
	Service   string  // GCP service name
	Monthly   float32 // Monthly cost
	Yesterday float32 // The cost in the day before
}

func (r *QueryResult) String() string {
	return fmt.Sprintf("{Service: %s, Monthly: %f, Yesterday: %f}", r.Service, r.Monthly, r.Yesterday)
}

// BQClient is an object to connect to BigQuery and send a query
// to retrieve the GCP cost.
type BQClient struct {
	client *bigquery.Client
}

// NewBQClient constructs a BQClient object.
// The GCP project name of BQ is fetched from the
// environmental variable.
func NewBQClient() BQClient {
	projectID := os.Getenv("GCP_PROJECT")
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}

	return BQClient{client: client}
}

// SendQuery receives a query as a string and send it to BQ
// to retrieve the GCP cost.
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
