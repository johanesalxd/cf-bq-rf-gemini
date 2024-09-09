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
	wg := new(sync.WaitGroup)
	semaphore := make(chan struct{}, concurrencyLimit)

	for i, call := range bqReq.Calls {
		select {
		case <-ctx.Done():
			log.Printf("Context cancelled before starting goroutine #%d", i)
			texts[i] = string(GenerateJSONResponse(&PromptRequest{
				PromptOutput: json.RawMessage(`{"error": "Request cancelled"}`),
			}))

			continue
		default:
			wg.Add(1)

			// Process each call concurrently
			go func(i int, call []interface{}) {
				// Acquire semaphore
				semaphore <- struct{}{}
				defer func() {
					// Release semaphore
					<-semaphore
					wg.Done()
				}()
				log.Printf("Processing request in Goroutine #%d", i)

				// Check if call has 3 elements
				if len(call) != 3 {
					log.Printf("Error in Goroutine #%d: call does not have enough elements", i)
					texts[i] = string(GenerateJSONResponse(&PromptRequest{
						PromptOutput: json.RawMessage(`{"error": "Invalid input: expected 3 elements"}`),
					}))

					return
				}

				// Update the input from the call slice
				input := newPromptRequest()
				input.PromptInput = fmt.Sprint(call[0])
				input.Model = fmt.Sprint(call[1])
				input.ModelConfig = ParseModelConfig(fmt.Sprint(call[2]))
				input.PromptOutput = textToText(ctx, client, &input)

				texts[i] = string(GenerateJSONResponse(input))
			}(i, call)
		}
	}
	wg.Wait()

	// Prepare and return the BigQuery response
	return &BigQueryResponse{
		Replies: texts,
	}
}

// Generates content based on the provided input
func textToText(ctx context.Context, client *genai.Client, input *PromptRequest) json.RawMessage {
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
		return json.RawMessage(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
	}

	return GenerateJSONResponse(resp)
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
