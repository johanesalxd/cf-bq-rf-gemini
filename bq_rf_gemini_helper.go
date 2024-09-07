package bqrfgemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"cloud.google.com/go/vertexai/genai"
)

// textsToTexts processes multiple text inputs concurrently using the Gemini AI model
func textsToTexts(ctx context.Context, client *genai.Client, bqReq *BigQueryRequest) *BigQueryResponse {
	// Initialize a slice to store the processed texts
	texts := make([]string, len(bqReq.Calls))
	wait := new(sync.WaitGroup)

	// Process each call concurrently
	for i, call := range bqReq.Calls {
		wait.Add(1)

		go func(i int, call []interface{}) {
			defer wait.Done()

			select {
			case <-ctx.Done():
				log.Printf("Got cancellation signal in Goroutine #%d", i)

				return
			default:
				input := newPromptRequest()

				// Check if call has at least 3 elements
				if len(call) != 3 {
					log.Printf("Error in Goroutine #%d: call does not have enough elements", i)

					// Set error message in PromptOutput
					input.PromptOutput = json.RawMessage(`{"error": "Invalid input: expected at least 3 elements"}`)
					texts[i] = string(GenerateJSONResponse(input))

					return
				}

				log.Printf("Processing request in Goroutine #%d", i)

				input = PromptRequest{
					PromptInput: fmt.Sprint(call[0]),
					Model:       fmt.Sprint(call[1]),
					ModelConfig: ParseModelConfig(fmt.Sprint(call[2])),
				}
				texts[i] = textToText(ctx, client, &input)
			}
		}(i, call)
	}
	wait.Wait()

	// Prepare and return the BigQuery response
	bqResp := new(BigQueryResponse)
	bqResp.Replies = texts

	return bqResp
}

// textToText processes a single text input using the Gemini AI model
func textToText(ctx context.Context, client *genai.Client, input *PromptRequest) string {
	// Configure the generative model with input parameters
	mdl := client.GenerativeModel(input.Model)
	mdl.SetMaxOutputTokens(input.ModelConfig.MaxOutputTokens)
	mdl.SetTemperature(input.ModelConfig.Temperature)
	mdl.SetTopP(input.ModelConfig.TopP)
	mdl.SetTopK(input.ModelConfig.TopK)

	// Generate content using the model
	resp, err := mdl.GenerateContent(ctx, genai.Text(input.PromptInput))
	if err != nil {
		log.Printf("Error generating text for input: %v", err)
		input.PromptOutput = json.RawMessage(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
	} else {
		input.PromptOutput = GenerateJSONResponse(resp)
	}

	// Return the JSON representation of the entire PromptRequest
	return string(GenerateJSONResponse(input))
}

// GenerateJSONResponse converts the input to JSON format
func GenerateJSONResponse(input any) json.RawMessage {
	jsonInput, err := json.Marshal(input)
	if err != nil {
		log.Printf("Error marshaling input to JSON: %v", err)
		return json.RawMessage(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
	}

	return jsonInput
}

// ParseModelConfig converts a JSON string to ModelConfig struct
func ParseModelConfig(input string) ModelConfig {
	config := newModelConfig()

	// Attempt to unmarshal the input JSON into the config
	if err := json.Unmarshal([]byte(input), &config); err != nil {
		log.Printf("Default value used due to error unmarshaling model config: %v", err)
		return config
	}

	return config
}
