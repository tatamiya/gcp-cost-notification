package main

import (
	"bytes"
	"html/template"
)

type QueryParameters struct {
	TableName string
}

func buildQuery(tableName string) string {

	params := QueryParameters{tableName}
	var buf bytes.Buffer
	t := template.Must(template.ParseFiles("query.sql"))
	t.Execute(&buf, params)

	return buf.String()
}
