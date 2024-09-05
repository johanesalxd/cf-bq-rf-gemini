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

		go func(i int, promptInput, model string) {
			defer wait.Done()

			select {
			case <-ctx.Done():
				log.Printf("Got cancellation signal in Goroutine #%d", i)
				return
			default:
				//TODO: remove promptInput for less verbose logging
				log.Printf("Running in Goroutine #%d for input: %v", i, promptInput)

				input := promptRequest{
					PromptInput: promptInput,
					Model:       model,
				}

				// Process the input and store the result
				texts[i] = textToText(ctx, client, &input)
			}
		}(i, fmt.Sprint(call[0]), fmt.Sprint(call[1]))
	}
	wait.Wait()

	// Prepare and return the BigQuery response
	bqResp := new(BigQueryResponse)
	bqResp.Replies = texts

	return bqResp
}

// textToText processes a single text input using the Gemini AI model
func textToText(ctx context.Context, client *genai.Client, input *promptRequest) string {
	// Get the generative model
	mdl := client.GenerativeModel(input.Model)

	// Generate content using the model
	resp, err := mdl.GenerateContent(ctx, genai.Text(input.PromptInput))
	if err != nil {
		log.Printf("Error generating text for input: %v", err)
		input.PromptOutput = err.Error()

		return generateText(input)
	}

	// Extract the generated text from the response
	var output string

	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if text, ok := part.(genai.Text); ok {
					output += string(text)
				}
			}
		}
	}

	input.PromptOutput = output

	return generateText(input)
}

// generateText converts the promptRequest to JSON format
func generateText(input *promptRequest) string {
	jsonInput, err := json.Marshal(input)
	if err != nil {
		log.Printf("Error marshaling input to JSON: %v", err)
		return fmt.Sprintf(`{"error": "Failed to marshal input: %v"}`, err)
	}

	return string(jsonInput)
}
