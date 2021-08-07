package query

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"time"

	reportingperiod "github.com/tatamiya/gcp-cost-notification/reporting-period"
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
		templatePath: "./" + fileDir + "query.sql",
	}
}

func (b *QueryBuilder) Build(period reportingperiod.ReportingPeriod) string {

	reportingToTimestamp := period.To.Format(time.RFC3339)
	reportingDateTo := template.HTML(reportingToTimestamp)

	reportingFromTimestamp := period.From.Format(time.RFC3339)
	reportingDateFrom := template.HTML(reportingFromTimestamp)

	params := struct {
		TableName         string
		ReportingDateFrom template.HTML
		ReportingDateTo   template.HTML
	}{
		TableName:         b.tableID,
		ReportingDateFrom: reportingDateFrom,
		ReportingDateTo:   reportingDateTo,
	}
	var buf bytes.Buffer
	t := template.Must(template.ParseFiles(b.templatePath))
	t.Execute(&buf, params)

	return buf.String()
}