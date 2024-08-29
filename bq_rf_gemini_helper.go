package bqrfgemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/google/generative-ai-go/genai"
)

// textsToTexts processes multiple text inputs concurrently using the Gemini AI model
func textsToTexts(ctx context.Context, client *genai.Client, bqReq *BigQueryRequest) *BigQueryResponse {
	texts := make([]string, len(bqReq.Calls))
	wait := new(sync.WaitGroup)

	for i, call := range bqReq.Calls {
		wait.Add(1)

		go func(j int, promptInput, model string) {
			defer wait.Done()

			for {
				select {
				case <-ctx.Done():
					log.Printf("Got cancellation signal in Goroutine #%d", j)

					return
				default:
					log.Printf("Running in Goroutine #%d for input: %v", j, promptInput)

					input := promptRequest{
						PromptInput: promptInput,
						Model:       model,
					}
					text := textToText(ctx, client, input)
					texts[i] = generateText(text, &input)

					return
				}
			}
		}(i, fmt.Sprint(call[0]), fmt.Sprint(call[1]))
	}
	wait.Wait()

	bqResp := new(BigQueryResponse)
	bqResp.Replies = texts

	return bqResp
}

// textToText generates content for a single text input using the specified Gemini AI model
func textToText(ctx context.Context, client *genai.Client, input promptRequest) *genai.GenerateContentResponse {
	mdl := client.GenerativeModel(input.Model)

	resp, err := mdl.GenerateContent(ctx, genai.Text(input.PromptInput))
	if err != nil {
		log.Fatal(err)
	}

	return resp
}

// generateText extracts the generated text from the AI response and formats it as JSON
func generateText(resp *genai.GenerateContentResponse, input *promptRequest) string {
	var output string

	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				output += string(part.(genai.Text))
			}
		}
	}

	input.PromptOutput = output

	jsonInput, err := json.Marshal(input)
	if err != nil {
		log.Printf("Error marshaling input to JSON: %v", err)
		return ""
	}

	return string(jsonInput)
}
