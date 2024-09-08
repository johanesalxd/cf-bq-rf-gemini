BigQuery Remote Function (BQ RF) with Gemini
-----------------------------
Details TBA

# How to run
## Run locally
```
FUNCTION_TARGET=BQRFGemini PROJECT_ID=YOUR_PROJECT_ID LOCATION=YOUR_LOCATION go run cmd/main.go
```

## Run locally with Pack and Docker
```
pack build --builder=gcr.io/buildpacks/builder cf-bq-rf-gemini

gcloud auth application-default login

ADC=~/.config/gcloud/application_default_credentials.json && \
docker run -p8080:8080 \
-e PROJECT_ID=YOUR_PROJECT_ID \
-e LOCATION=YOUR_LOCATION \
-e GOOGLE_APPLICATION_CREDENTIALS=/tmp/keys/secret.json \
-v ${ADC}:/tmp/keys/secret.json \
cf-bq-rf-gemini
```

## Test locally (accept [BQ RF request contract](https://cloud.google.com/bigquery/docs/remote-functions#input_format))
```
curl -m 60 -X POST localhost:8080 \
-H "Content-Type: application/json" \
-d '{
  "requestId": "",
  "caller": "",
  "sessionUser": "",
  "userDefinedContext": {},
  "calls": [
    ["what is bigquery", "gemini-1.5-flash-001", "{\"temperature\":0.2,\"maxOutputTokens\":8000,\"topP\":0.8,\"topK\":40}"],
    ["default model config", "gemini-1.5-pro-001", "{\"temperature\":0.2,\"maxOutputTokens\":8000"],
    ["error model", "gemini-1.0-proo", "{\"topP\":0.8,\"topK\":40}"],
    ["missing element", "gemini-1.0-pro-002"]
  ]
}'
```

## Run on Cloud Function
```
gcloud functions deploy cf-bq-rf-gemini \
    --gen2 \
    --concurrency=8 \
    --cpu=1 \
    --memory=512Mi \
    --runtime=go122 \
    --region=us-central1 \
    --source=. \
    --entry-point=BQRFGemini \
    --trigger-http \
    --allow-unauthenticated \
    --env-vars-file=.env.yaml
```

## Run on Cloud Run
[![Run on Google Cloud](https://deploy.cloud.run/button.svg)](https://deploy.cloud.run)

# Additional notes
TBA

## Related links
* https://cloud.google.com/bigquery/docs/remote-functions
* https://cloud.google.com/functions/docs/concepts/go-runtime
* https://cloud.google.com/docs/buildpacks/build-function
