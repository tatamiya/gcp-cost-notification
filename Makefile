REGION=asia-northeast1
TRIGGER_TOPIC=cost_notification

local-run:
	godotenv -f ./.env go test -run TestCostNotifier

test-all:
	godotenv -f ./.env go test ./... -cover

test-short:
	godotenv -f ./.env go test ./... -cover -short

deploy:
	gcloud functions deploy CostNotifier --env-vars-file env.yaml --trigger-topic $(TRIGGER_TOPIC) --region=$(REGION) --runtime=go113
