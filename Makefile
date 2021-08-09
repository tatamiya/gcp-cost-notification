REGION=asia-northeast1
TRIGGER_TOPIC=cost_notification

local-run:
	godotenv -f ./.env go test -run TestCostNotifier

test:
	godotenv -f ./.env go test ./... -cover -short

test-all:
	godotenv -f ./.env go test ./... -cover

deploy:
	gcloud functions deploy CostNotifier --env-vars-file env.yaml --trigger-topic $(TRIGGER_TOPIC) --region=$(REGION) --runtime=go113
