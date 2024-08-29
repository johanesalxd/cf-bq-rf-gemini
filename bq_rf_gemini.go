package bqrfgemini

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// BQRFGemini handles HTTP requests for the BigQuery Remote Function using Gemini AI
func BQRFGemini(w http.ResponseWriter, r *http.Request) {
	bqReq := new(BigQueryRequest)
	if err := json.NewDecoder(r.Body).Decode(bqReq); err != nil {
		SendError(w, err, http.StatusBadRequest)

		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer func() {
		log.Print("Done, Goroutines closed")
		cancel()
	}()

	// TODO: Replace the API key with an environment variable for better security
	// TODO: Change from Google AI to Vertex AI
	// TODO: Update client initialization to use Vertex AI
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		SendError(w, err, http.StatusInternalServerError)

		return
	}
	defer client.Close()

	bqResp := textsToTexts(ctx, client, bqReq)
	SendSuccess(w, bqResp)
}
