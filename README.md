BigQuery Remote Function (BQ RF) with Gemini
-----------------------------
Details TBA

# How to run
## Run locally
```
FUNCTION_TARGET=BQRFGemini GEMINI_API_KEY=YOUR_GEMINI_API_KEY go run cmd/main.go
```

## Run locally with Pack and Docker
```
pack build --builder=gcr.io/buildpacks/builder cf-bq-rf-gemini

gcloud auth application-default login

ADC=~/.config/gcloud/application_default_credentials.json && \
docker run -p8080:8080 \
-e GEMINI_API_KEY=YOUR_GEMINI_API_KEY \
-e GOOGLE_APPLICATION_CREDENTIALS=/tmp/keys/secret.json \
-v ${ADC}:/tmp/keys/secret.json \
cf-bq-rf-gemini
```

## Test locally (accept BQ RF [request contract](https://cloud.google.com/bigquery/docs/remote-functions#input_format))
Notes: *model* should be the same for the same call
```
curl -m 60 -X POST localhost:8080 \
-H "Content-Type: application/json" \
-d '{
  "requestId": "",
  "caller": "",
  "sessionUser": "",
  "userDefinedContext": {},
  "calls": [
    ["prompt_1", "model"],
    ["prompt_2", "model"],
    ["prompt_3", "model"]
  ]
  }'
```

## Run on Cloud Function
```
gcloud functions deploy cf-bq-rf-gemini \
    --gen2 \
    --concurrency=8 \
    --runtime=go122 \
    --region=us-central1 \
    --source=. \
    --entry-point=BQRFGemini \
    --trigger-http \
    --allow-unauthenticated
```

## Run on Cloud Run
[![Run on Google Cloud](https://deploy.cloud.run/button.svg)](https://deploy.cloud.run)

# Additional notes
Details TBA
## Related links
* https://cloud.google.com/bigquery/docs/remote-functions
* https://cloud.google.com/functions/docs/concepts/go-runtime
* https://cloud.google.com/docs/buildpacks/build-function
