package bqrfgemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/vertexai/genai"
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

	client := clientPool.Get().(*genai.Client)
	if client == nil {
		SendError(w, fmt.Errorf("failed to get client from pool"), http.StatusInternalServerError)

		return
	}
	defer clientPool.Put(client)

	log.Print("Client retrieved from pool")

	bqResp := textsToTexts(ctx, client, bqReq)
	SendSuccess(w, bqResp)
}
