// query package implements a object build a query to retrieve
// GCP cost from BigQuery.
package query

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/tatamiya/gcp-cost-notification/src/datetime"
)

// QueryBuilder is an object to build a query from a template.
type QueryBuilder struct {
	tableID      string
	templatePath string
}

// NewQueryBuilder constructs QueryBuilder.
//
// Four environment variables are needed in construction.
//
// `GCP_PROJECT`, `DATASET_NAME`, `TABLE_NAME` ... identify the table to retrieve the cost from.
//
// `FILE_DIRECTORY` ... the directory name where a query template file `template.sql` is.
// On Cloud Functions, it must be `serverless_function_source_code/`.
func NewQueryBuilder() QueryBuilder {

	projectID := os.Getenv("GCP_PROJECT")
	datasetName := os.Getenv("DATASET_NAME")
	tableName := os.Getenv("TABLE_NAME")
	tableID := fmt.Sprintf("%s.%s.%s", projectID, datasetName, tableName)

	fileDir := os.Getenv("FILE_DIRECTORY")

	return QueryBuilder{
		tableID:      tableID,
		templatePath: "./" + fileDir + "src/query/template.sql",
	}
}

// Build method renders a query tamplate with the cost aggregation period to report and BQ table ID.
func (b *QueryBuilder) Build(period datetime.ReportingPeriod) string {

	reportingToTimestamp := period.To.Format(time.RFC3339)
	reportingDateTo := template.HTML(reportingToTimestamp)

	reportingFromTimestamp := period.From.Format(time.RFC3339)
	reportingDateFrom := template.HTML(reportingFromTimestamp)

	params := struct {
		TableName         string
		TimeZone          string
		ReportingDateFrom template.HTML
		ReportingDateTo   template.HTML
	}{
		TableName:         b.tableID,
		TimeZone:          period.TimeZone,
		ReportingDateFrom: reportingDateFrom,
		ReportingDateTo:   reportingDateTo,
	}
	var buf bytes.Buffer
	t := template.Must(template.ParseFiles(b.templatePath))
	t.Execute(&buf, params)

	return buf.String()
}
