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
	clientInterface := clientPool.Get()
	defer func() {
		if clientInterface != nil {
			clientPool.Put(clientInterface)
		}

		log.Print("Client returned to pool")
	}()

	// Type assert the clientInterface to *genai.Client
	client := clientInterface.(*genai.Client)
	log.Print("Client retrieved from pool")

	// Process the request using textsToTexts function
	bqResp := textsToTexts(ctx, client, bqReq)

	// Send the successful response
	SendSuccess(w, bqResp)
}
