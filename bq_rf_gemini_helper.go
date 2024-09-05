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
	texts := make([]string, len(bqReq.Calls))
	wait := new(sync.WaitGroup)

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

				text := textToText(ctx, client, &input)
				texts[i] = generateText(text, &input)
			}
		}(i, fmt.Sprint(call[0]), fmt.Sprint(call[1]))
	}
	wait.Wait()

	bqResp := new(BigQueryResponse)
	bqResp.Replies = texts

	return bqResp
}

// textToText generates content for a single text input using the specified Gemini AI model
func textToText(ctx context.Context, client *genai.Client, input *promptRequest) *genai.GenerateContentResponse {
	mdl := client.GenerativeModel(input.Model)

	resp, err := mdl.GenerateContent(ctx, genai.Text(input.PromptInput))
	if err != nil {
		input.PromptOutput = err.Error()

		return nil
	}

	return resp
}

// generateText extracts the generated text from the AI response and formats it as JSON
func generateText(resp *genai.GenerateContentResponse, input *promptRequest) string {
	if resp == nil {
		log.Printf("Error: Received nil response: %v", input.PromptOutput)
	} else {
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
	}

	jsonInput, err := json.Marshal(input)
	if err != nil {
		log.Printf("Error marshaling input to JSON: %v", err)
		return ""
	}

	return string(jsonInput)
}
