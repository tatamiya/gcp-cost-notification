# gcp-cost-notification

![](sample_image.png)

Send notification to Slack on GCP costs.

This system uses

- Cloud Billing
- BigQuery
- Cloud Scheduler
- Pub/Sub
- Google Cloud Functions
- Slack API

## Prerequisites

- [Export Cloud Billing data to BigQuery](https://cloud.google.com/billing/docs/how-to/export-data-bigquery-setup)
- [Create Pub/Sub topic](https://cloud.google.com/pubsub/docs/quickstart-console)
- [Set Cloud Scheduler](https://cloud.google.com/scheduler/docs/quickstart)

## Enviroment Variables

Enviroment variables to use in GCF runtime are set in `env.yaml` file.

(sample)
```yaml
GCP_PROJECT: <your GCP poject-id>
DATASET_NAME: <BQ dataset name>
TABLE_NAME: <BQ table name>
SLACK_WEBHOOK_URL: <slack webhook url>
FILE_DIRECTORY: "serverless_function_source_code/" # this should be fixed
TIMEZONE: <Your TimeZone. e.g. Asia/Tokyo>
```

## Deploy Command

```
gcloud functions deploy CostNotifier --env-vars-file env.yaml --trigger-topic <Pub/Sub topic name> --region=<region> --runtime=go113
```