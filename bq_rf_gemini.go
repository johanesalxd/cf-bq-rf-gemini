package bqrfgemini

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"cloud.google.com/go/vertexai/genai"
)

// BQRFGemini handles HTTP requests for the BigQuery Remote Function using Gemini AI
func BQRFGemini(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming BigQuery request
	bqReq := new(BigQueryRequest)
	if err := json.NewDecoder(r.Body).Decode(bqReq); err != nil {
		SendError(w, err, http.StatusBadRequest)

		return
	}

	// Create a cancellable context
	ctx, cancel := context.WithCancel(r.Context())
	defer func() {
		cancel()

		log.Print("Done, Goroutines closed")
	}()

	// Get a client from the pool
	client := clientPool.Get().(*genai.Client)
	defer func() {
		if client != nil {
			clientPool.Put(client)
			log.Print("Client returned to pool")
		}
	}()
	log.Print("Client retrieved from pool")

	// Process the request using textsToTexts function
	bqResp := textsToTexts(ctx, client, bqReq)

	// Send the successful response
	SendSuccess(w, bqResp)
}
