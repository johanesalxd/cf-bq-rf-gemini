package bqrfgemini

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// TODO: Replace the API key with an environment variable for better security
// TODO: Change from Google AI to Vertex AI
// TODO: Update client initialization to use Vertex AI
var clientPool = &sync.Pool{
	New: func() interface{} {
		client, err := genai.NewClient(context.Background(), option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
		if err != nil {
			log.Printf("Error creating new client: %v", err)
			return nil
		}

		log.Print("Client created")

		return client
	},
}

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

	client := clientPool.Get().(*genai.Client)
	if client == nil {
		SendError(w, errors.New("failed to get client from pool"), http.StatusInternalServerError)

		return
	}
	defer clientPool.Put(client)

	log.Print("Client retrieved from pool")

	bqResp := textsToTexts(ctx, client, bqReq)
	SendSuccess(w, bqResp)
}
