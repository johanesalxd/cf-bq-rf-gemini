package bqrfgemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"cloud.google.com/go/vertexai/genai"
)

var (
	clientPool *sync.Pool
	initOnce   sync.Once
)

// SendError sends an error response with the given error message and HTTP status code
func SendError(w http.ResponseWriter, err error, code int) {
	bqResp := new(BigQueryResponse)
	bqResp.ErrorMessage = fmt.Sprintf("Got error with details: %v", err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(bqResp)
}

// SendSuccess sends a successful response with the given BigQueryResponse
func SendSuccess(w http.ResponseWriter, bqResp *BigQueryResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bqResp)
}

func initializePool() {
	clientPool = &sync.Pool{
		New: func() interface{} {
			client, err := genai.NewClient(context.Background(), os.Getenv("PROJECT_ID"), os.Getenv("LOCATION"))
			if err != nil {
				log.Fatalf("Failed to create client: %v", err)
			}
			log.Print("Client created")

			return client
		},
	}
}
