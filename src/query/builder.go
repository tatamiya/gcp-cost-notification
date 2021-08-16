package query

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/tatamiya/gcp-cost-notification/src/datetime"
)

type QueryBuilder struct {
	tableID      string
	templatePath string
}

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
